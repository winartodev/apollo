package configs

import (
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/helpers"
	"log"
	"os"
)

const (
	development = "development"
	staging     = "staging"
	production  = "production"

	developmentConfigPath = "core/files/apollo.dev.yaml"
)

type Auth struct {
	APIKey string `yaml:"apiKey"`
	JWT    JWT    `yaml:"jwt"`
}

type JWT struct {
	AccessToken  AccessToken  `yaml:"accessToken"`
	RefreshToken RefreshToken `yaml:"refreshToken"`
}

type AccessToken struct {
	SecretKey string `yaml:"secretKey"`
}

type RefreshToken struct {
	SecretKey string `yaml:"secretKey"`
}

type Config struct {
	App struct {
		Name string `yaml:"name"`
		Port struct {
			HTTP int `yaml:"http"`
		} `yaml:"port"`
	} `yaml:"app"`

	Database Database `yaml:"database"`

	Auth Auth `yaml:"auth"`
}

func NewConfig() (*Config, error) {
	var config *Config
	var env = os.Getenv("ENV")
	var filePath string

	switch env {
	case development:
		filePath = developmentConfigPath
	case staging:
		filePath = developmentConfigPath
	case production:
		filePath = developmentConfigPath
	default:
		log.Printf("WARNING: Using development configuration for unknown environment %s\n", env)
		filePath = developmentConfigPath
	}

	err := helpers.ReadYAMLFile(filePath, &config)
	if err != nil {
		return nil, err
	}

	err = os.Setenv(core.JwtAccessTokenSecretKey, config.Auth.JWT.AccessToken.SecretKey)
	if err != nil {
		return nil, err
	}

	err = os.Setenv(core.JwtRefreshTokenSecretKey, config.Auth.JWT.RefreshToken.SecretKey)
	if err != nil {
		return nil, err
	}

	err = os.Setenv(core.JoblessApiKey, config.Auth.APIKey)
	if err != nil {
		return nil, err
	}

	return config, nil
}
