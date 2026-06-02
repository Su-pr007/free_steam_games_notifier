package entity

import "github.com/go-playground/validator/v10"

type Game struct {
	Id        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Link      string `json:"link" db:"link"`
	ImageLink string `json:"image_link" db:"image_link"`
	DateAdd   int    `json:"date_add" db:"date_add"`
}

func (u *Game) Validate() error {
	validate := validator.New()

	if err := validate.Struct(u); err != nil {
		return err
	}

	return nil
}
