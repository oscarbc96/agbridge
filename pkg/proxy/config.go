package proxy

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
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
		awsCfg *aws.Config
		err    error
		result = make(map[*regexp.Regexp]Handler)
		seen   = make(map[string]struct{})
	)

	for _, gw := range c.Gateways {
		awsCfg, err = awsutils.LoadConfigFor(gw.ProfileName, gw.Region)
		if err != nil {
			return nil, fmt.Errorf("couldn't load AWS Config for profile %s: %w", gw.ProfileName, err)
		}

		resources, err := awsutils.DescribeAPIGateway(*awsCfg, gw.RestAPIID)
		if err != nil {
			return nil, fmt.Errorf("couldn't describe API Gateway for RestAPIID %s: %w", gw.RestAPIID, err)
		}

		for _, resource := range resources {
			if resource.ResourceMethods == nil {
				continue
			}

			path := *resource.Path

			if _, ok := seen[path]; ok {
				return nil, fmt.Errorf("duplicate path %s found in the configuration for Rest API ID %s", path, gw.RestAPIID)
			}
			seen[path] = struct{}{}

			regexPattern, err := convertPathToRegex(path)
			if err != nil {
				return nil, fmt.Errorf("invalid path %s: %w", path, err)
			}

			result[regexPattern] = Handler{
				Path:       path,
				ResourceID: *resource.Id,
				RestAPIID:  gw.RestAPIID,
				Methods:    lo.Keys(resource.ResourceMethods),
				Config:     *awsCfg,
			}
		}
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

func NewConfig(restAPIID, profileName, region string) *Config {
	return &Config{
		Gateways: []GatewayConfig{
			{
				RestAPIID:   restAPIID,
				ProfileName: profileName,
				Region:      region,
			},
		},
	}
}
