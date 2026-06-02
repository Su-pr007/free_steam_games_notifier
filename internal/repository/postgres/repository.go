package repository

import (
	"steamGamesSales/internal/entity"

	"github.com/jmoiron/sqlx"
)

type Chat interface {
	Create(chat *entity.Chat) error
	Get(string) (*entity.Chat, error)
	Update(chat *entity.Chat) error
	GetSubscribed() ([]entity.Chat, error)
}

type Game interface {
	Create(game *entity.Game) error
	Get(string) (*entity.Game, error)
	Update(game *entity.Game) error
}

type Repository struct {
	Chat
	Game
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Chat: NewChatPostgres(db),
		Game: NewGamePostgres(db),
	}
}
