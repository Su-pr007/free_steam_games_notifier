package cron

import (
	"log/slog"
	"steamGamesSales/internal/entity"
	"steamGamesSales/internal/lib/parser"
	repository "steamGamesSales/internal/repository/postgres"
	"time"
)

type Cron struct {
	log               *slog.Logger
	repo              *repository.Repository
	cronEvents        map[string]int64
	newFreeGamesEvent chan []entity.Game
}

func NewCron(log *slog.Logger, repo *repository.Repository, newFreeGamesEvent chan []entity.Game) *Cron {
	return &Cron{
		log:  log,
		repo: repo,
		cronEvents: map[string]int64{
			"ParseGames": 0,
		},
		newFreeGamesEvent: newFreeGamesEvent,
	}
}

func (c *Cron) Start() {
	for {
		for cronName, lastTrigger := range c.cronEvents {
			if lastTrigger > time.Now().Unix() {
				continue
			}

			if cronName == "ParseGames" {
				go c.ParseGames()
			}
			c.cronEvents[cronName] = time.Now().Add(8 * time.Hour).Unix()
		}

		time.Sleep(1 * time.Hour)
	}
}

func (c *Cron) ParseGames() {
	gamesParser := parser.NewGamesParser(c.log, c.repo, c.newFreeGamesEvent)
	gamesParser.Run()
}
