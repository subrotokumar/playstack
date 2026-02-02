package idp

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (idp *IdentityProvider) SignUp(
	ctx context.Context,
	name, email, password string,
) (bool, string, error) {
	secretHash := GetSecretHash(email, idp.ClientId, idp.ClientSecret)

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
		var userExists *types.UsernameExistsException
		if errors.As(err, &invalidPassword) {
			return false, "", fmt.Errorf("invalid password")
		} else if errors.As(err, &userExists) {
			return false, "", fmt.Errorf("user already exists")
		}
	}

	return out.UserConfirmed, aws.ToString(out.UserSub), nil
}

func (idp *IdentityProvider) ConfirmSignUp(
	ctx context.Context,
	email, otp string,
) error {
	secretHash := GetSecretHash(email, idp.ClientId, idp.ClientSecret)

	_, err := idp.CognitoClient.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(idp.ClientId),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(otp),
		SecretHash:       aws.String(secretHash),
	})

	return err
}

func (idp *IdentityProvider) ResendOTP(
	ctx context.Context,
	email string,
) error {
	secretHash := GetSecretHash(email, idp.ClientId, idp.ClientSecret)
	_, err := idp.CognitoClient.ResendConfirmationCode(ctx, &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId:   aws.String(idp.ClientId),
		Username:   aws.String(email),
		SecretHash: aws.String(secretHash),
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

	secretHash := GetSecretHash(email, idp.ClientId, idp.ClientSecret)

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
	secretHash := GetSecretHash(username, idp.ClientId, idp.ClientSecret)

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
