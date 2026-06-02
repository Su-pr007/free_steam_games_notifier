package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	chatsTable = "chats"
	gamesTable = "games"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSL      string
}

func (cfg *Config) Debug() {
	fmt.Println(cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSL)
}

func ConnectDb(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host,
			cfg.Port,
			cfg.Username,
			cfg.Password,
			cfg.DBName,
			cfg.SSL,
		),
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
