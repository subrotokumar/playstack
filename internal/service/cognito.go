package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"

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

func NewIndentityProvider(region, clientId, clientSecret string) IdentityProvider {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		log.Fatal(err)
	}
	cognitoClient := cognitoidentityprovider.NewFromConfig(sdkConfig)
	slog.Info("Cognito", "Id", clientId, "Secret", clientSecret)
	return IdentityProvider{
		CognitoClient: cognitoClient,
		ClientId:      clientId,
		ClientSecret:  clientSecret,
	}
}

func (idp *IdentityProvider) SignUp(
	ctx context.Context,
	name, email, password string,
) (bool, string, error) {
	secretHash := utility.GetSecretHash(email, idp.ClientId, idp.ClientSecret)

	out, err := idp.CognitoClient.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(idp.ClientId),
		Username:   aws.String(email),
		Password:   aws.String(password),
		SecretHash: aws.String(secretHash),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(email)},
			{Name: aws.String("name"), Value: aws.String(name)},
		},
	})
	if err != nil {
		var invalidPassword *types.InvalidPasswordException
		if errors.As(err, &invalidPassword) {
			log.Println(*invalidPassword.Message)
		} else {
			log.Printf("Couldn't sign up user %v. Here's why: %v\n", email, err)
		}
	}

	return out.UserConfirmed, aws.ToString(out.UserSub), nil
}

func (idp *IdentityProvider) ConfirmSignUp(
	ctx context.Context,
	email, otp string,
) error {
	secretHash := utility.GetSecretHash(email, idp.ClientId, idp.ClientSecret)

	_, err := idp.CognitoClient.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(idp.ClientId),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(otp),
		SecretHash:       aws.String(secretHash),
	})

	return err
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	IdToken      string
}

func (idp *IdentityProvider) Login(
	ctx context.Context,
	email, password string,
) (*AuthTokens, error) {

	secretHash := utility.GetSecretHash(email, idp.ClientId, idp.ClientSecret)

	out, err := idp.CognitoClient.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		ClientId: aws.String(idp.ClientId),
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME":    email,
			"PASSWORD":    password,
			"SECRET_HASH": secretHash,
		},
	})
	if err != nil {
		return nil, err
	}

	if out.AuthenticationResult == nil {
		return nil, errors.New("empty authentication result")
	}

	return &AuthTokens{
		AccessToken:  aws.ToString(out.AuthenticationResult.AccessToken),
		RefreshToken: aws.ToString(out.AuthenticationResult.RefreshToken),
		IdToken:      aws.ToString(out.AuthenticationResult.IdToken),
	}, nil
}

func (idp *IdentityProvider) RefreshAccessToken(
	ctx context.Context,
	username, refreshToken string,
) (string, error) {
	secretHash := utility.GetSecretHash(username, idp.ClientId, idp.ClientSecret)

	out, err := idp.CognitoClient.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		ClientId: aws.String(idp.ClientId),
		AuthFlow: types.AuthFlowTypeRefreshTokenAuth,
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
			"SECRET_HASH":   secretHash,
		},
	})
	if err != nil {
		return "", err
	}

	if out.AuthenticationResult == nil {
		return "", errors.New("empty authentication result")
	}

	return aws.ToString(out.AuthenticationResult.AccessToken), nil
}

func (idp *IdentityProvider) ChangePassword(
	ctx context.Context,
	accessToken, previousPassword, proposedPassword string,
) error {
	_, err := idp.CognitoClient.ChangePassword(ctx, &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(accessToken),
		PreviousPassword: aws.String(previousPassword),
		ProposedPassword: aws.String(proposedPassword),
	})
	if err != nil {
		return err
	}

	return nil
}
