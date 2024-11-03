package proxy

import (
	"fmt"
	"os"

	"github.com/samber/lo"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/oscarbc96/agbridge/pkg/apigateway"
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

func (c *Config) Validate() (map[string]Handler, error) {
	var (
		awsCfg *aws.Config
		err    error
		result = make(map[string]Handler)
	)

	for _, gw := range c.Gateways {
		awsCfg, err = apigateway.LoadConfigFor(gw.ProfileName, gw.Region)
		if err != nil {
			return nil, fmt.Errorf("couldn't load AWS Config for profile %s: %w", gw.ProfileName, err)
		}

		resources, err := apigateway.DescribeAPIGateway(*awsCfg, gw.RestAPIID)
		if err != nil {
			return nil, fmt.Errorf("couldn't describe API Gateway for RestAPIID %s: %w", gw.RestAPIID, err)
		}

		for _, resource := range resources {
			if _, exists := result[*resource.Path]; exists {
				return nil, fmt.Errorf("duplicate path %s found in the configuration for Rest API ID %s", *resource.Path, gw.RestAPIID)
			}

			if resource.ResourceMethods != nil {
				result[*resource.Path] = Handler{
					ResourceID: *resource.Id,
					RestAPIID:  gw.RestAPIID,
					Methods:    lo.Keys(resource.ResourceMethods),
					Config:     *awsCfg,
				}
			}
		}
	}

	return result, nil
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open Config file: %w", err)
	}
	defer file.Close()

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
