package testutil

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

type BaseTestSuite struct {
	suite.Suite

	Ctx                 context.Context
	Config              *aws.Config
	LocalStackContainer *localstack.LocalStackContainer
}

func (suite *BaseTestSuite) SetupSuite() {
	suite.Ctx = context.Background()

	cfg, lsc, err := CreateLocalStackContainer(suite.Ctx)
	suite.Require().NoError(err, "failed to set up LocalStack container")

	suite.Config = cfg
	suite.LocalStackContainer = lsc
}

func (suite *BaseTestSuite) TearDownSuite() {
	err := CleanupAPIGateways(*suite.Config)
	suite.Require().NoError(err, "failed to clean up API gateways")

	err = suite.LocalStackContainer.Terminate(suite.Ctx)
	suite.Require().NoError(err, "failed to terminate LocalStack container")
}

func (suite *BaseTestSuite) TearDownTest() {
	err := CleanupAPIGateways(*suite.Config)
	suite.Require().NoError(err, "failed to clean up API gateways after test")
}
