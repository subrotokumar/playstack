package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"gitlab.com/subrotokumar/glitchr/internal/utility"
)

type IdentityProvider struct {
	CognitoClient *cognitoidentityprovider.Client
	ClientId      string
	ClientSecret  string
}

func NewIndentityProvider() IdentityProvider {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		log.Fatal(err)
	}
	cognitoClient := cognitoidentityprovider.NewFromConfig(sdkConfig)
	return IdentityProvider{
		CognitoClient: cognitoClient,
		ClientId:      os.Getenv("COGNITO_CLIENT_ID"),
		ClientSecret:  os.Getenv("COGNITO_CLIENT_SECRET"),
	}
}

func (actor *IdentityProvider) SignUp(ctx context.Context, name, userEmail, password string) (bool, error) {
	confirmed := false

	secretHash := utility.GetSecretHash(userEmail, actor.ClientId, actor.ClientId)
	output, err := actor.CognitoClient.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(actor.ClientId),
		Password:   aws.String(password),
		Username:   aws.String(userEmail),
		SecretHash: aws.String(secretHash),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(userEmail)},
			{Name: aws.String("name"), Value: aws.String(name)},
		},
	})
	if err != nil {
		var invalidPassword *types.InvalidPasswordException
		if errors.As(err, &invalidPassword) {
			log.Println(*invalidPassword.Message)
		} else {
			log.Printf("Couldn't sign up user %v. Here's why: %v\n", userEmail, err)
		}
	} else {
		confirmed = output.UserConfirmed
	}
	return confirmed, err
}
