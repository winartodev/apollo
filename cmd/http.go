package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core/configs"
	"github.com/winartodev/apollo/core/routes"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := configs.NewConfig()
	if err != nil {
		panic(err)
	}

	db, err := cfg.Database.NewConnection()
	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		if err := configs.CloseDB(db); err != nil {
			log.Fatalf("db.Close() error: %v", err)
		} else {
			log.Println("DB gracefully stopped")
		}
	}(db)

	redisClient, err := cfg.Redis.NewRedis()
	if err != nil {
		panic(err)
	}
	defer func(client *redis.Client) {
		if err := configs.CloseRedis(client); err != nil {
			log.Fatalf("redis.Close() error: %v", err)
		} else {
			log.Println("Redis gracefully stopped")
		}
	}(redisClient)

	autoMigration, err := configs.NewAutoMigration(cfg.Database.Name, db)
	if err != nil {
		panic(err)
	}

	if autoMigration == nil {
		panic("autoMigration is nil")
	}

	if err := autoMigration.Start(); err != nil {
		panic(err)
	}

	smtpClient, err := configs.NewSMTPClient(cfg.SMTP)
	if err != nil {
		panic(err)
	}

	twilioClient := configs.NewTwilioClient(cfg.Twilio)

	app := fiber.New(fiber.Config{
		AppName: cfg.App.Name,
	})

	repository := routes.NewRepository(routes.RepositoryDependency{DB: db, Redis: redisClient})
	controller := routes.NewController(routes.ControllerDependency{
		Repository: repository,
		SMTPClient: smtpClient,
		Twilio:     twilioClient})
	handler := routes.NewHandler(routes.HandlerDependency{Controller: controller})

	if err = routes.RegisterHandler(app, handler); err != nil {
		panic(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		if err := app.Listen(fmt.Sprintf(":%v", cfg.App.Port.HTTP)); err != nil {
			serverErrors <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case sig := <-shutdown:
		log.Printf("Received signal: %v. Initiating shutdown...", sig)
	case err = <-serverErrors:
		log.Fatalf("Server error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("app.Shutdown error: %v", err)
	} else {
		log.Println("app gracefully stopped")
	}

	log.Println("Application shutdown complete")
}
