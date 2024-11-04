package awsutils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func DescribeAPIGateway(config aws.Config, apiID string) ([]types.Resource, error) {
	client := apigateway.NewFromConfig(config)
	input := &apigateway.GetResourcesInput{
		RestApiId: aws.String(apiID),
	}

	var result []types.Resource
	ctx := context.TODO()
	paginator := apigateway.NewGetResourcesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve resources: %w", err)
		}
		result = append(result, page.Items...)
	}

	return result, nil
}
