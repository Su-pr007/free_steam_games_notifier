package main

import (
	"bufio"
	"fmt"
	"os"
	"steamGamesSales/internal/config"
	repository "steamGamesSales/internal/repository/postgres"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
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

	args := os.Args[1:]

	if len(args) > 0 {
		handleArgs(args, reader, m)
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

func handleArgs(args []string, reader *bufio.Reader, m *migrate.Migrate) {
	if args[0] == "down" {
		fmt.Print("Вы уверены? (y/N)")
		sure, err2 := reader.ReadString('\n')
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		if sure == "y" {
			err := m.Down()
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("Успех")
		}

		return
	}

	if args[0] == "stepback" {
		err := m.Steps(-1)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Последняя миграция успешно отменена")

		return
	}
}
