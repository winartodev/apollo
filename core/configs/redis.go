package configs

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

func (r *Redis) NewRedis() (client *redis.Client, err error) {
	addr := fmt.Sprintf("%s:%s", r.Host, r.Port)
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: r.Password,
		DB:       0,
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CloseRedis(client *redis.Client) error {
	err := client.Close()
	if err != nil {
		return err
	}
	return nil
}
