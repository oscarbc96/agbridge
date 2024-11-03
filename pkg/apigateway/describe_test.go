package apigateway

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"

	"github.com/oscarbc96/agbridge/internal/testutil"
	"github.com/stretchr/testify/suite"
)

type APIGatewayTestSuite struct {
	testutil.BaseTestSuite

	ApiGateway *apigateway.CreateRestApiOutput
}

func (suite *APIGatewayTestSuite) SetupTest() {
	apigw, err := testutil.CreateAPIGateway(*suite.Config, "test")
	suite.Require().NoError(err, "failed to create test API gateway")
	suite.ApiGateway = apigw
}

func (suite *APIGatewayTestSuite) TestDescribeAPIGateway() {
	apigws, err := DescribeAPIGateway(*suite.Config, *suite.ApiGateway.Id)

	suite.Require().NoError(err, "expected no error Describing API gateways")
	suite.Len(apigws, 1, "expected exactly one API gateway")
	suite.Equal(*suite.ApiGateway.RootResourceId, *apigws[0].Id, "expected API Gateway ID")
}

func TestAPIGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(APIGatewayTestSuite))
}
