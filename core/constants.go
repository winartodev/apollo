package core

import "errors"

const (
	API = "/api"
	V1  = "/v1"

	AccessProtected = "/protected"
	AccessInternal  = "/internal"

	JwtAccessTokenSecretKey  = "JWT_ACCESS_TOKEN_SECRET_KEY"
	JwtRefreshTokenSecretKey = "JWT_REFRESH_TOKEN_SECRET_KEY"
	JoblessApiKey            = "JOBLESS_API_KEY"

	DefaultOffset int64 = 0
	DefaultPage   int64 = 1
	DefaultLimit  int64 = 10
	MaxLimit      int64 = 100
	DefaultSort         = "asc"
	DefaultOrder        = "id"
)

var (
	ErrRefreshTokenExists = errors.New("refresh token already exists")
)
