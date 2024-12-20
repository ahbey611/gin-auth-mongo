package auth

import (
	"encoding/json"
	"errors"
	"gin-auth-mongo/graph/model"
	"gin-auth-mongo/repositories"
	"gin-auth-mongo/utils/jwt"
	"os"

	"github.com/square/go-jose/v3"
)

func RefreshToken(token string) (*model.AccessToken, error) {

	// check if the token is valid
	refreshToken, err := repositories.GetRefreshTokenByToken(token)
	if err != nil || refreshToken == nil {
		return nil, errors.New("invalid refresh token1")
	}

	// check if the user exists
	user, err := repositories.GetUserByID(refreshToken.UserID.Hex())
	if err != nil || user == nil {
		return nil, errors.New("invalid refresh token2")
	}

	// generate new token
	accessToken, err := jwt.HandleRefreshToken(user)
	if err != nil {
		return nil, errors.New("invalid refresh token3")
	}

	return accessToken, nil
}

func GetTokenInfo(token string) (map[string]interface{}, error) {
	// parse the token
	parsedJWT, err := jwt.ParseToken(token)
	if err != nil {
		return nil, err
	}

	// check if the token has a header
	if len(parsedJWT.Headers) <= 0 {
		return nil, errors.New("invalid token")
	}

	// load public key
	publicJWKs := make([]jose.JSONWebKey, 0)
	data, err := os.ReadFile(".public/keys.json")
	if err != nil {
		return nil, err
	}

	// parse public keys
	if err := json.Unmarshal(data, &publicJWKs); err != nil {
		return nil, err
	}

	JWKs := jose.JSONWebKeySet{
		Keys: publicJWKs,
	}
	publicJWK := JWKs.Key(parsedJWT.Headers[0].KeyID)

	if len(publicJWK) <= 0 {
		return nil, errors.New("invalid token")
	}

	// Extract the claims
	allClaims := make(map[string]interface{})
	if err := parsedJWT.Claims(publicJWK[0].Key, &allClaims); err != nil {
		return nil, err
	}

	// log.Println(allClaims)

	return allClaims, nil
}

func DeleteRefreshTokenByToken(token string) error {
	return repositories.DeleteRefreshTokenByToken(token)
}

func DeleteRefreshTokenByUserID(userID string) error {
	return repositories.DeleteRefreshTokenByUserID(userID)
}
