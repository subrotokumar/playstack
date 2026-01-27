package idp

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"gitlab.com/subrotokumar/playstack/libs/core"
)

type IdentityProvider struct {
	CognitoClient *cognitoidentityprovider.Client
	ClientId      string
	ClientSecret  string
}

func NewIndentityProvider(region, clientId, clientSecret string) IdentityProvider {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		core.LogFatal(err.Error())
	}
	cognitoClient := cognitoidentityprovider.NewFromConfig(sdkConfig)
	return IdentityProvider{
		CognitoClient: cognitoClient,
		ClientId:      clientId,
		ClientSecret:  clientSecret,
	}
}
