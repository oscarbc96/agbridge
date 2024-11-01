package apigateway

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/oscarbc96/agbridge/internal/testutil"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

type APIGatewayTestSuite struct {
	suite.Suite

	Ctx                 context.Context //nolint:containedctx
	Config              *aws.Config
	LocalStackContainer *localstack.LocalStackContainer
}

func (suite *APIGatewayTestSuite) SetupSuite() {
	suite.Ctx = context.Background()

	cfg, lsc, err := testutil.CreateLocalStackContainer(suite.Ctx)
	if err != nil {
		os.Exit(1)
	}

	suite.Config = cfg
	suite.LocalStackContainer = lsc
}

func (suite *APIGatewayTestSuite) TearDownSuite() {
	err := testutil.CleanupAPIGateways(*suite.Config)
	if err != nil {
		os.Exit(1)
	}

	err = suite.LocalStackContainer.Terminate(suite.Ctx)
	if err != nil {
		os.Exit(1)
	}

	suite.Ctx.Done()
}

func (suite *APIGatewayTestSuite) SetupTest() {
	_, err := testutil.CreateAPIGateway(*suite.Config, "test")
	if err != nil {
		os.Exit(1)
	}
}

func (suite *APIGatewayTestSuite) TearDownTest() {
	err := testutil.CleanupAPIGateways(*suite.Config)
	if err != nil {
		os.Exit(1)
	}
}

func (suite *APIGatewayTestSuite) TestListAPIGateways() {
	apgws, err := ListAPIGateways(*suite.Config)

	suite.NoError(err) //nolint:testifylint
	suite.Len(apgws, 1)
}

func TestApiGatewayTestSuite(t *testing.T) { //nolint:paralleltest
	suite.Run(t, new(APIGatewayTestSuite))
}
