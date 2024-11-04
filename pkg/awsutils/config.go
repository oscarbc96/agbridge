package awsutils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadConfigFor(profile, region string) (*aws.Config, error) {
	var options []func(*config.LoadOptions) error

	if profile != "" {
		options = append(options, config.WithSharedConfigProfile(profile))
	}

	if region != "" {
		options = append(options, config.WithRegion(region))
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config (profile: %s, region: %s), %w", profile, region, err)
	}

	return &cfg, nil
}
