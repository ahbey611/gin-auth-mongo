package jwt

import (
	"encoding/base64"
	"errors"
	"regexp"
	"time"

	crand "crypto/rand"

	"gin-auth-mongo/graph/model"
	"gin-auth-mongo/models"
	"gin-auth-mongo/repositories"

	"gin-auth-mongo/utils/consts"

	"gin-auth-mongo/utils/jwkmanager"

	"github.com/gin-gonic/gin"
	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
)

// generate random string for refresh token
func GenerateRefreshToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := crand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// pick a random signing key from the ./private/keys.json file
func pickRandomSigningKey() (jose.JSONWebKey, jose.SigningKey, error) {

	// load from memory
	privateJWK, err := jwkmanager.GetRandomJWK()
	if err != nil {
		return jose.JSONWebKey{}, jose.SigningKey{}, err
	}

	key := jose.SigningKey{Algorithm: jose.EdDSA, Key: privateJWK}
	return privateJWK, key, nil
}

// generate public and private claims
func generateClaims(user *models.User, privateJWK jose.JSONWebKey, key jose.SigningKey) (*jwt.Builder, time.Time, error) {

	var signerOptions = jose.SignerOptions{}
	signerOptions.WithType("JWT")
	signerOptions.WithHeader("kid", privateJWK.KeyID)

	rsaSigner, err := jose.NewSigner(key, &signerOptions)
	if err != nil {
		return nil, time.Time{}, err
	}

	builder := jwt.Signed(rsaSigner)

	// public claims
	issuedAt := time.Now()
	publicClaims := jwt.Claims{
		Issuer:  consts.JWT_ISSUER,
		Subject: user.ID.Hex(),
		// Audience:
		IssuedAt: jwt.NewNumericDate(issuedAt),
		Expiry:   jwt.NewNumericDate(issuedAt.Add(time.Duration(consts.JWT_ACCESS_TOKEN_EXPIRY) * time.Minute)),
	}

	// private claims
	privateClaims := map[string]interface{}{
		"email": user.Email,
		// YOU CAN ADD MORE PRIVATE CLAIMS HERE
	}

	builder = builder.Claims(publicClaims).Claims(privateClaims)

	return &builder, issuedAt, nil
}

// generate access and refresh token, then save to database
func GenerateToken(user *models.User, device string, builder jwt.Builder, issuedAt time.Time) (*model.Token, error) {

	accessToken, err := builder.CompactSerialize()
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(32)
	if err != nil {
		return nil, err
	}

	err = repositories.CreateRefreshToken(user.ID.Hex(), refreshToken, issuedAt.AddDate(0, 0, consts.JWT_REFRESH_TOKEN_EXPIRY), device)

	if err != nil {
		return nil, err
	}

	token := &model.Token{
		UserID:             user.ID.Hex(),
		AccessToken:        accessToken,
		AccessTokenExpiry:  issuedAt.Add(time.Duration(consts.JWT_ACCESS_TOKEN_EXPIRY) * time.Minute).Format(consts.DATETIME_NANO_FORMAT),
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: issuedAt.Add(time.Duration(consts.JWT_REFRESH_TOKEN_EXPIRY) * time.Hour).Format(consts.DATETIME_NANO_FORMAT),
		Device:             device,
	}

	return token, nil
}

// use to generate access token and refresh token
func HandleLogin(user *models.User, device string) (*model.Token, error) {

	privateJWK, key, err := pickRandomSigningKey()
	if err != nil {
		return nil, err
	}

	builder, issuedAt, err := generateClaims(user, privateJWK, key)
	if err != nil {
		return nil, err
	}

	return GenerateToken(user, device, *builder, issuedAt)
}

func GenerateNewAccessToken(builder jwt.Builder, issuedAt time.Time) (*model.AccessToken, error) {

	accessToken, err := builder.CompactSerialize()
	if err != nil {
		return nil, err
	}

	accessTokenExpiry := issuedAt.Add(time.Duration(consts.JWT_ACCESS_TOKEN_EXPIRY) * time.Minute).Format(consts.DATETIME_NANO_FORMAT)

	return &model.AccessToken{
		AccessToken:       accessToken,
		AccessTokenExpiry: accessTokenExpiry,
	}, nil
}

// use to refresh access token
func HandleRefreshToken(user *models.User) (*model.AccessToken, error) {

	privateJWK, key, err := pickRandomSigningKey()
	if err != nil {
		return nil, err
	}

	builder, issuedAt, err := generateClaims(user, privateJWK, key)
	if err != nil {
		return nil, err
	}

	return GenerateNewAccessToken(*builder, issuedAt)
}

// parse part

// parse the jwt token
func ParseToken(token string) (*jwt.JSONWebToken, error) {
	return jwt.ParseSigned(token)
}

// get and parse the jwt token from the authorization header
func GetTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")

	if len(authHeader) <= 0 || authHeader == "" {
		return "", errors.New("empty authorization header")
	}

	regex := regexp.MustCompile(`^Bearer (\S+)$`)
	matches := regex.FindStringSubmatch(authHeader)

	if len(matches) <= 0 {
		return "", errors.New("invalid authorization header")
	}

	token := matches[1]

	return token, nil
}

// parse the jwt token and return the claims
func ParseJWTClaims(token string) (map[string]interface{}, error) {
	parsedJWT, err := ParseToken(token)
	if err != nil {
		return nil, err
	}
	// check if the token has a header
	if len(parsedJWT.Headers) <= 0 {
		return nil, errors.New("no headers found")
	}

	// get public jwks
	publicJWKs, err := jwkmanager.GetPublicJWKs()
	if err != nil {
		return nil, err
	}

	// create jwks
	JWKs := jose.JSONWebKeySet{
		Keys: publicJWKs,
	}
	publicJWK := JWKs.Key(parsedJWT.Headers[0].KeyID)

	// check if the public jwk is valid
	if len(publicJWK) <= 0 {
		return nil, errors.New("invalid public jwk")
	}

	// extract the claims
	allClaims := make(map[string]interface{})
	if err := parsedJWT.Claims(publicJWK[0].Key, &allClaims); err != nil {
		return nil, err
	}

	return allClaims, nil
}
