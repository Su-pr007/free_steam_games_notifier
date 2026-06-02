package main

import (
	"fmt"
	"steamGamesSales/internal/config"
	repository "steamGamesSales/internal/repository/postgres"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustLoadConfig()

	db, err := repository.ConnectDb(repository.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSL:      cfg.DB.SSL,
	})
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		fmt.Println("Ошибка создания экземпляра драйвера базы данных:", err)
		return
	}

	// Путь к директории с миграциями
	migrationsDir := "./sql/migrations"

	// Создание экземпляра объекта Migrate
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationsDir), "postgres", driver)
	if err != nil {
		fmt.Println("Ошибка создания экземпляра объекта Migrate:", err)
		return
	}

	// Применение всех миграций
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Println("Ошибка применения миграций:", err)
		return
	}

	fmt.Println("Все миграции успешно применены")
}
