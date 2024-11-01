package apigateway

import (
	"context"
	"log"
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
		log.Fatalf("failed to set up LocalStack container: %v", err)
	}

	suite.Config = cfg
	suite.LocalStackContainer = lsc
}

func (suite *APIGatewayTestSuite) TearDownSuite() {
	if err := testutil.CleanupAPIGateways(*suite.Config); err != nil {
		log.Printf("failed to clean up API gateways: %v", err)
	}

	if err := suite.LocalStackContainer.Terminate(suite.Ctx); err != nil {
		log.Printf("failed to terminate LocalStack container: %v", err)
	}

	suite.Ctx.Done()
}

func (suite *APIGatewayTestSuite) SetupTest() {
	if _, err := testutil.CreateAPIGateway(*suite.Config, "test"); err != nil {
		log.Fatalf("failed to create test API gateway: %v", err)
	}
}

func (suite *APIGatewayTestSuite) TearDownTest() {
	if err := testutil.CleanupAPIGateways(*suite.Config); err != nil {
		log.Printf("failed to clean up API gateways after test: %v", err)
	}
}

func (suite *APIGatewayTestSuite) TestListAPIGateways() {
	apigws, err := ListAPIGateways(*suite.Config)

	suite.Require().NoError(err, "expected no error listing API gateways")
	suite.Len(apigws, 1, "expected exactly one API gateway")
}

func TestAPIGatewayTestSuite(t *testing.T) { //nolint:paralleltest
	suite.Run(t, new(APIGatewayTestSuite))
}
