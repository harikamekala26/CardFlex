package config

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	Client *gorm.DB
}

func ConnectSQL(env Env) *DB {
	if env.DBDriver != "sqlite" {
		log.Fatalf("unsupported DB_DRIVER %q (supported: sqlite)", env.DBDriver)
	}

	client, err := gorm.Open(sqlite.Open(env.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to sql database: %v", err)
	}

	sqlDB, err := client.DB()
	if err != nil {
		log.Fatalf("failed to obtain sql handle: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping sql database: %v", err)
	}

	log.Printf("connected to %s database (%s)", env.DBDriver, fmt.Sprintf("dsn=%s", env.DBDSN))
	return &DB{Client: client}
}
