package testutil

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func CreateAPIGateway(config aws.Config, name string) (*apigateway.CreateRestApiOutput, error) {
	apiClient := apigateway.NewFromConfig(config)

	apiInput := &apigateway.CreateRestApiInput{
		Name:    aws.String(name),
		Version: aws.String("v1"),
		EndpointConfiguration: &types.EndpointConfiguration{
			Types: []types.EndpointType{types.EndpointTypeRegional},
		},
	}

	result, err := apiClient.CreateRestApi(context.TODO(), apiInput)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func CleanupAPIGateways(config aws.Config) error {
	apiClient := apigateway.NewFromConfig(config)

	listInput := &apigateway.GetRestApisInput{}

	listResult, err := apiClient.GetRestApis(context.TODO(), listInput)
	if err != nil {
		return err
	}

	var errs []error

	for _, api := range listResult.Items {
		deleteInput := &apigateway.DeleteRestApiInput{
			RestApiId: api.Id,
		}

		_, err := apiClient.DeleteRestApi(context.TODO(), deleteInput)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
