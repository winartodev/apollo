package configs

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	postgresDNS = "postgres://%s:%s@%s:%s/%s?sslmode=%s"
)

type Database struct {
	Driver          string `yaml:"driver"`
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	Name            string `yaml:"name"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	SSLMode         string `yaml:"sslMode"`
	MaxOpenConn     int    `yaml:"defaultMaxConn"`
	MaxIdleConn     int    `yaml:"defaultIdleConn"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
	ConnMaxIdleTime int    `yaml:"connMaxIdleTime"`
}

func (d *Database) NewConnection() (*sql.DB, error) {
	dsn := fmt.Sprintf(postgresDNS, d.Username, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
	db, err := sql.Open(d.Driver, dsn)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(errorFailedOpenConnection, err))
	}

	db.SetMaxOpenConns(d.MaxOpenConn)
	db.SetMaxIdleConns(d.MaxIdleConn)
	db.SetConnMaxIdleTime(time.Duration(d.ConnMaxIdleTime) * time.Minute)
	db.SetConnMaxLifetime(time.Duration(d.ConnMaxLifetime) * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, errors.New(fmt.Sprintf(errorFailedPingConnection, err))
	}

	return db, nil
}

func CloseDB(db *sql.DB) error {
	stats := db.Stats()
	if stats.OpenConnections > 0 {
		return errors.New(fmt.Sprintf(errorDatabaseHasOpenConnections, stats.OpenConnections))
	}

	if err := db.Close(); err != nil {
		return errors.New(fmt.Sprintf(errorFailedCloseDatabase, err))
	}

	return nil
}
