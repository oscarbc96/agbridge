package apigateway

import (
	"testing"

	"github.com/oscarbc96/agbridge/internal/testutil"
	"github.com/stretchr/testify/suite"
)

type APIGatewayTestSuite struct {
	testutil.BaseTestSuite
}

func (suite *APIGatewayTestSuite) SetupTest() {
	_, err := testutil.CreateAPIGateway(*suite.Config, "test")
	suite.Require().NoError(err, "failed to create test API gateway")
}

func (suite *APIGatewayTestSuite) TestListAPIGateways() {
	apigws, err := ListAPIGateways(*suite.Config)

	suite.Require().NoError(err, "expected no error listing API gateways")
	suite.Len(apigws, 1, "expected exactly one API gateway")
}

func TestAPIGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(APIGatewayTestSuite))
}
