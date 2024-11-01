package apigateway

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func ListAPIGateways(config aws.Config) ([]types.RestApi, error) {
	apiClient := apigateway.NewFromConfig(config)
	input := &apigateway.GetRestApisInput{}

	result, err := apiClient.GetRestApis(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}
