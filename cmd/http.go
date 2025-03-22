package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core/configs"
	"github.com/winartodev/apollo/core/routes"
	"log"
	"os"
	"os/signal"
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
		err := db.Close()
		if err != nil {
			log.Fatalf("db.Close() error: %v", err)
		}
	}(db)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

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

	app := fiber.New(fiber.Config{
		AppName: cfg.App.Name,
	})

	repository := routes.NewRepository(routes.RepositoryDependency{DB: db})
	controller := routes.NewController(routes.ControllerDependency{Repository: repository})
	handler := routes.NewHandler(routes.HandlerDependency{Controller: controller})

	if err = routes.RegisterHandler(app, handler); err != nil {
		panic(err)
	}

	go func() {
		if err := app.Listen(fmt.Sprintf(":%v", cfg.App.Port.HTTP)); err != nil {
			log.Fatalf("server.ListenAndServe: %v", err)
		}
	}()

	<-shutdown
	log.Println("Shutting down server...")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		log.Fatalf("app.Shutdown error: %v", err)
	} else {
		log.Println("app gracefully stopped")
	}

	if err := configs.CloseDB(db); err != nil {
		log.Fatalf("db.Close() error: %v", err)
	} else {
		log.Println("DB gracefully stopped")
	}

	log.Println("Application shutdown complete")
}
