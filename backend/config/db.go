package config

import (
	"fmt"
	"log"
	"strings"

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

	dsn := configureSQLiteDSN(env.DBDSN)
	client, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
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

	// SQLite is sensitive to concurrent writes, so keep one open connection
	// and allow a short wait before surfacing "database is locked" errors.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	if err := client.Exec("PRAGMA journal_mode = WAL;").Error; err != nil {
		log.Printf("warning: failed to enable WAL mode: %v", err)
	}
	if err := client.Exec("PRAGMA busy_timeout = 5000;").Error; err != nil {
		log.Printf("warning: failed to set busy timeout: %v", err)
	}

	log.Printf("connected to %s database (%s)", env.DBDriver, fmt.Sprintf("dsn=%s", dsn))
	return &DB{Client: client}
}

func configureSQLiteDSN(raw string) string {
	if strings.Contains(raw, "?") {
		return raw + "&_busy_timeout=5000&_journal_mode=WAL"
	}

	return raw + "?_busy_timeout=5000&_journal_mode=WAL"
}
