package proxy

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/oscarbc96/agbridge/pkg/awsutils"
	"github.com/oscarbc96/agbridge/pkg/log"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type GatewayConfig struct {
	RestAPIID   string `yaml:"rest_api_id"`
	ProfileName string `yaml:"profile_name"`
	Region      string `yaml:"region"`
	StageName   string `yaml:"stage_name"`
}

type Config struct {
	Gateways []GatewayConfig `yaml:"gateways"`
}

func convertPathToRegex(path string) (*regexp.Regexp, error) {
	// Replace `{param}` with `[^/]+`
	re := regexp.MustCompile(`\{[^/]+\}`)
	pattern := "^" + re.ReplaceAllString(path, `[^/]+`) + "$"
	return regexp.Compile(pattern)
}

func (c *Config) Validate() (map[*regexp.Regexp]Handler, error) {
	var (
		mu     sync.Mutex
		wg     sync.WaitGroup
		result = make(map[*regexp.Regexp]Handler)
		seen   = make(map[string]struct{})
		errCh  = make(chan error, len(c.Gateways))
	)

	for _, gw := range c.Gateways {
		wg.Add(1)
		go func(gw GatewayConfig) {
			defer wg.Done()

			awsCfg, err := awsutils.LoadConfigFor(gw.ProfileName, gw.Region)
			if err != nil {
				errCh <- fmt.Errorf("couldn't load AWS Config for profile %s: %w", gw.ProfileName, err)
				return
			}

			var stage *apigateway.GetStageOutput
			var stageVariables map[string]string
			if gw.StageName != "" {
				stage, err = awsutils.DescribeStage(*awsCfg, gw.RestAPIID, gw.StageName)
				if err != nil {
					errCh <- fmt.Errorf("couldn't describe stage with name %s: %w", gw.StageName, err)
					return
				}
				stageVariables = stage.Variables
			}

			resources, err := awsutils.DescribeAPIGateway(*awsCfg, gw.RestAPIID)
			if err != nil {
				errCh <- fmt.Errorf("couldn't describe API Gateway for RestAPIID %s: %w", gw.RestAPIID, err)
				return
			}

			for _, resource := range resources {
				if resource.ResourceMethods == nil {
					continue
				}

				path := *resource.Path
				stagePath := path
				if stage != nil {
					stagePath = fmt.Sprintf("/%s%s", *stage.StageName, path)
				}

				regexPattern, err := convertPathToRegex(stagePath)
				if err != nil {
					errCh <- fmt.Errorf("invalid path %s: %w", stagePath, err)
					return
				}

				mu.Lock()
				if _, ok := seen[stagePath]; ok {
					mu.Unlock()
					errCh <- fmt.Errorf("duplicate path %s found in the configuration for Rest API ID %s", stagePath, gw.RestAPIID)
					return
				}
				seen[stagePath] = struct{}{}
				result[regexPattern] = Handler{
					StagePath:      stagePath,
					Path:           path,
					ResourceID:     *resource.Id,
					RestAPIID:      gw.RestAPIID,
					Methods:        lo.Keys(resource.ResourceMethods),
					Config:         *awsCfg,
					StageVariables: stageVariables,
				}
				mu.Unlock()
			}
		}(gw)
	}

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return nil, err // return first error
	}

	return result, nil
}

func LoadConfig(fs afero.Fs, filename string) (*Config, error) {
	file, err := fs.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open Config file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Fatal("Failed to close config file", log.Err(cerr), log.String("file", filename))
		}
	}()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse Config file: %w", err)
	}

	return &config, nil
}

func NewConfig(restAPIID, profileName, region, stageName string) *Config {
	return &Config{
		Gateways: []GatewayConfig{
			{
				RestAPIID:   restAPIID,
				ProfileName: profileName,
				Region:      region,
				StageName:   stageName,
			},
		},
	}
}
