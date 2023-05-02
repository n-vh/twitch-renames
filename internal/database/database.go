package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/n-vh/twitch-renames/internal/models/renames"
	"github.com/n-vh/twitch-renames/internal/models/users"
	"github.com/n-vh/twitch-renames/internal/models/workers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Renames *renames.Model
	Users   *users.Model
	Workers *workers.Model
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Schema   string
}

func Connect(config *Config) {
	uri := makeUri(config.User, config.Password, config.Host, config.Port, config.Database, config.Schema)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger:                 customLogger(),
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Panic("Unable to connect to Postgres", err)
	}

	log.Println("Connected to Postgres!")

	db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", config.Schema))

	Renames = renames.InitModel(db)
	Users = users.InitModel(db)
	Workers = workers.InitModel(db)
}

func makeUri(user, password, host, port, database, schema string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?search_path=%s", user, password, host, port, database, schema)
}

func customLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}
