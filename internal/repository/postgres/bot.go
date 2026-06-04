package repository

import (
	"fmt"
	"steamGamesSales/internal/entity"
	"time"

	"github.com/jmoiron/sqlx"
)

type ChatPostgres struct {
	db *sqlx.DB
}

type GamePostgres struct {
	db *sqlx.DB
}

func NewChatPostgres(db *sqlx.DB) *ChatPostgres {
	return &ChatPostgres{
		db: db,
	}
}

func NewGamePostgres(db *sqlx.DB) *GamePostgres {
	return &GamePostgres{
		db: db,
	}
}

// Chats

func (r *ChatPostgres) Create(chat *entity.Chat) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (id, name, username, subscribed, language_code, date_add) VALUES ($1, $2, $3, $4, $5, $6) returning id`,
		chatsTable,
	)

	_, err := r.db.Exec(query, chat.Id, chat.Name, chat.Username, chat.Subscribed, chat.LanguageCode, chat.DateAdd)
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatPostgres) Get(id string) (*entity.Chat, error) {
	query := fmt.Sprintf(`SELECT id, name, username, subscribed, language_code, date_add FROM %s WHERE id = $1`, chatsTable)

	var user entity.Chat

	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ChatPostgres) Update(chat *entity.Chat) error {
	query := fmt.Sprintf(
		`UPDATE %s SET name=$1, subscribed=$2 WHERE id = $3`,
		chatsTable,
	)

	_, err := r.db.Exec(
		query,
		chat.Name,
		chat.Subscribed,
		chat.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatPostgres) GetSubscribed() ([]entity.Chat, error) {
	query := fmt.Sprintf(`SELECT id, name, username, language_code, date_add FROM %s WHERE subscribed = $1`, chatsTable)

	var chats []entity.Chat

	err := r.db.Select(&chats, query, true)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

// Games

func (r *GamePostgres) Create(game entity.Game) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (id, name, link, image_link, date_add) VALUES ($1, $2, $3, $4, $5)",
		gamesTable,
	)

	_, err := r.db.Exec(query, game.Id, game.Name, game.Link, game.ImageLink, time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

func (r *GamePostgres) Get(id string) (*entity.Game, error) {
	query := fmt.Sprintf(`SELECT id, name, link, image_link, date_add FROM %s WHERE id = $1`, gamesTable)

	var game entity.Game

	err := r.db.Get(&game, query, id)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *GamePostgres) Update(game *entity.Game) error {
	query := fmt.Sprintf(
		`UPDATE %s SET name=$1, link=$2, image_link=$3 WHERE id = $4`,
		gamesTable,
	)

	_, err := r.db.Exec(
		query,
		game.Name,
		game.Link,
		game.ImageLink,
		game.Id,
	)
	if err != nil {
		return err
	}

	return nil
}
