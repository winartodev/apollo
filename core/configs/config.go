package configs

import (
	"github.com/winartodev/apollo/core"
	"github.com/winartodev/apollo/core/helpers"
	"log"
	"os"
	"strconv"
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

type SMTP struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Sender   string `yaml:"sender"`
	Password string `yaml:"password"`
}
type Twilio struct {
	AccountSid string `yaml:"sid"`
	AuthToken  string `yaml:"authToken"`
	PhoneNum   string `yaml:"phoneNumber"`
}

type OTP struct {
	Enable bool `yaml:"enable"`
}

type Config struct {
	App struct {
		Name string `yaml:"name"`
		Port struct {
			HTTP int `yaml:"http"`
		} `yaml:"port"`
	} `yaml:"app"`

	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`

	OTP    OTP    `yaml:"otp"`
	Auth   Auth   `yaml:"auth"`
	SMTP   SMTP   `yaml:"smtp"`
	Twilio Twilio `yaml:"twilio"`
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

	if err = SaveToEnv(config); err != nil {
		return nil, err
	}

	return config, nil
}

func SaveToEnv(config *Config) (err error) {
	err = os.Setenv(core.JwtAccessTokenSecretKey, config.Auth.JWT.AccessToken.SecretKey)
	if err != nil {
		return err
	}

	err = os.Setenv(core.JwtRefreshTokenSecretKey, config.Auth.JWT.RefreshToken.SecretKey)
	if err != nil {
		return err
	}

	err = os.Setenv(core.ApolloAPIKey, config.Auth.APIKey)
	if err != nil {
		return err
	}

	err = os.Setenv(core.EnvSMTPHost, config.SMTP.Host)
	if err != nil {
		return err
	}

	err = os.Setenv(core.EnvSMTPPort, strconv.Itoa(config.SMTP.Port))
	if err != nil {
		return err
	}

	err = os.Setenv(core.EnvSMTPSender, config.SMTP.Sender)
	if err != nil {
		return err
	}

	err = os.Setenv(core.EnvSMTPPassword, config.SMTP.Password)
	if err != nil {
		return err
	}

	return nil
}
