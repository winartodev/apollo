package core

import "errors"

const (
	API = "/api"
	V1  = "/v1"

	AccessProtected = "/protected"
	AccessInternal  = "/internal"

	JwtAccessTokenSecretKey  = "JWT_ACCESS_TOKEN_SECRET_KEY"
	JwtRefreshTokenSecretKey = "JWT_REFRESH_TOKEN_SECRET_KEY"
	ApolloAPIKey             = "APOLLO_API_KEY"

	EnvSMTPHost     = "SMTP_HOST"
	EnvSMTPPort     = "SMTP_PORT"
	EnvSMTPSender   = "SMTP_SENDER"
	EnvSMTPPassword = "SMTP_PASSWORD"

	DefaultOffset int64 = 0
	DefaultPage   int64 = 1
	DefaultLimit  int64 = 10
	MaxLimit      int64 = 100
	DefaultSort         = "asc"
	DefaultOrder        = "id"

	OSWindows = "windows"
)

var (
	ErrRefreshTokenExists = errors.New("refresh token already exists")
)
