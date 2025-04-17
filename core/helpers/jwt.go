package helpers

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/modules/application"
	userEntity "github.com/winartodev/apollo/modules/user/entities"
	"os"
	"strings"
	"time"
)

const (
	errorUnexpectedSigningMethod = "unexpected signing method: %v"
)

var (
	errorMissingSecretKey = errors.New("missing secret key")
	errorInvalidToken     = errors.New("invalid token")
	errorTokenExpired     = errors.New("token is expired")
)

type JWTClaims struct {
	ID       int64               `json:"id,omitempty"`
	Username string              `json:"username,omitempty"`
	Email    string              `json:"email,omitempty"`
	Access   *application.Access `json:"access,omitempty"`
	jwt.StandardClaims
}

type accessToken struct {
	SecretKey []byte
}

type refreshToken struct {
	SecretKey []byte
}

type JWT struct {
	AccessToken  accessToken
	RefreshToken refreshToken
}

type JWTResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewJWT() (*JWT, error) {
	atSecret := os.Getenv(core.JwtAccessTokenSecretKey)
	if atSecret == "" {
		return nil, errors.New("access token secret key is empty")
	}

	rtSecret := os.Getenv(core.JwtRefreshTokenSecretKey)
	if rtSecret == "" {
		return nil, errors.New("refresh token secret key is empty")
	}

	return &JWT{
		AccessToken: accessToken{
			SecretKey: []byte(atSecret),
		},
		RefreshToken: refreshToken{
			SecretKey: []byte(rtSecret),
		},
	}, nil
}

func (j *JWT) GenerateToken(user *userEntity.User, access *application.Access) (result *JWTResponse, err error) {
	if user == nil {
		return nil, errors.New("user not found")
	}

	if !isSecretKeyExists(j.AccessToken.SecretKey) && !isSecretKeyExists(j.RefreshToken.SecretKey) {
		return nil, errorMissingSecretKey
	}

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Access:   access,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	})

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		ID:     user.ID,
		Access: access,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	})

	newAccessTokenString, err := newAccessToken.SignedString(j.AccessToken.SecretKey)
	if err != nil {
		return nil, err
	}

	newRefreshTokenString, err := newRefreshToken.SignedString(j.RefreshToken.SecretKey)
	if err != nil {
		return nil, err
	}

	return &JWTResponse{
		AccessToken:  newAccessTokenString,
		RefreshToken: newRefreshTokenString,
	}, nil
}

func (j *JWT) VerifyToken(secretKey []byte, tokenString string) (result map[string]interface{}, isValid bool, err error) {
	if !isSecretKeyExists(secretKey) {
		return nil, false, errorMissingSecretKey
	}

	token, err := j.ParseToken(secretKey, tokenString)
	if err != nil {
		return nil, false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		result = claims
		exp := claims["exp"].(float64)
		expirationTime := time.Unix(int64(exp), 0)
		if time.Now().After(expirationTime) {
			return nil, false, errorTokenExpired
		}
	} else {
		return nil, false, errorInvalidToken
	}

	return result, true, nil
}

func (j *JWT) ParseToken(secretKey []byte, tokenString string) (result *jwt.Token, err error) {
	if !isSecretKeyExists(secretKey) {
		return nil, errorMissingSecretKey
	}

	const prefix = "Bearer "
	if strings.HasPrefix(tokenString, prefix) {
		tokenString = tokenString[len(prefix):]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf(errorUnexpectedSigningMethod, token.Header["alg"])
		}

		return secretKey, nil
	})
	if err != nil {
		return token, err
	}

	if !token.Valid {
		return token, errorInvalidToken
	}

	return token, err
}

func isSecretKeyExists(secretKey []byte) bool {
	return secretKey != nil && len(secretKey) > 0
}
