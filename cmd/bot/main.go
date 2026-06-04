package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"steamGamesSales/internal/config"
	"steamGamesSales/internal/entity"
	"steamGamesSales/internal/lib/cron"
	"steamGamesSales/internal/lib/logger"
	"steamGamesSales/internal/lib/tgBot"
	repository "steamGamesSales/internal/repository/postgres"
	"syscall"

	"github.com/jmoiron/sqlx"
)

type TgBot struct {
	logger            *slog.Logger
	botToken          string
	repo              *repository.Repository
	newFreeGamesEvent chan []entity.Game
}

func (initData TgBot) initService() *tgBot.BotData {
	return tgBot.NewBot(initData.logger, initData.botToken, initData.repo, initData.newFreeGamesEvent).InitBot()
}

type Cron struct {
	logger            *slog.Logger
	repo              *repository.Repository
	newFreeGamesEvent chan []entity.Game
}

func (initData Cron) initService() *cron.Cron {
	return cron.NewCron(initData.logger, initData.repo, initData.newFreeGamesEvent)
}

func main() {
	cfg := config.MustLoadConfig()

	log := initLogger()

	db := getDb(cfg)

	repo := repository.NewRepository(db)

	newFreeGamesChannel := make(chan []entity.Game)

	botService := TgBot{
		logger:            log,
		botToken:          cfg.TgBot.Token,
		repo:              repo,
		newFreeGamesEvent: newFreeGamesChannel,
	}.initService()
	cronService := Cron{
		logger:            log,
		repo:              repo,
		newFreeGamesEvent: newFreeGamesChannel,
	}.initService()

	go botService.StartListening()
	go botService.StartWaitingForGames()
	go cronService.Start()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
}

func getDb(cfg *config.Config) *sqlx.DB {
	db, err := repository.ConnectDb(repository.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSL:      cfg.DB.SSL,
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to db: %s", err))
	}

	return db
}

func initLogger() *slog.Logger {
	opts := logger.PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	log := slog.New(logger.NewPrettyHandler(os.Stdout, opts))

	return log
}
