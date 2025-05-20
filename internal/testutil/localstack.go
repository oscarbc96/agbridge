package testutil

import (
	"context"
	"net"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

func CreateLocalStackContainer(ctx context.Context) (*aws.Config, *localstack.LocalStackContainer, error) {
	lsContainer, err := localstack.Run(
		ctx,
		"localstack/localstack:4.4.0",
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Env: map[string]string{"SERVICES": "apigateway"},
			},
		}),
	)
	if err != nil {
		return nil, nil, err
	}

	host, err := lsContainer.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	port, err := lsContainer.MappedPort(ctx, "4566/tcp")
	if err != nil {
		return nil, nil, err
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-1"))
	if err != nil {
		return nil, nil, err
	}

	awsCfg.BaseEndpoint = aws.String("http://" + net.JoinHostPort(host, port.Port()))

	return &awsCfg, lsContainer, nil
}
