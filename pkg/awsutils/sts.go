package awsutils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetAccountDetails(config aws.Config) (string, string, error) {
	stsClient := sts.NewFromConfig(config)

	callerIdentity, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", "", fmt.Errorf("failed to get caller identity: %w", err)
	}

	return *callerIdentity.Account, *callerIdentity.Arn, nil
}
