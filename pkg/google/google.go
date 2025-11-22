package google

import (
	"context"
	"errors"

	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// VerifyGoogleToken verifies a Google OAuth token and returns the token info
func VerifyGoogleToken(tokenString string, googleClientID string) (*oauth2.Tokeninfo, error) {
	// Create OAuth2 service
	ctx := context.Background()
	oauth2Service, err := oauth2.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, errors.New("failed to create OAuth2 service")
	}

	// Verify the token
	tokenInfo, err := oauth2Service.Tokeninfo().IdToken(tokenString).Do()
	if err != nil {
		return nil, errors.New("invalid Google token")
	}

	// Verify the audience (client ID)
	if tokenInfo.Audience != googleClientID {
		return nil, errors.New("token audience does not match")
	}

	// Verify the token is issued by Google
	if tokenInfo.IssuedTo != googleClientID {
		return nil, errors.New("token was not issued to this client")
	}

	return tokenInfo, nil
}
