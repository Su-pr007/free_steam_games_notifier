package entity

import (
	"github.com/go-playground/validator/v10"
)

type Chat struct {
	Id           int64  `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Username     string `json:"username" db:"username"`
	Subscribed   bool   `json:"subscribed" db:"subscribed"`
	LanguageCode string `json:"language_code" db:"language_code"`
	DateAdd      int    `json:"date_add" db:"date_add"`
}

func (u *Chat) Validate() error {
	validate := validator.New()

	if err := validate.Struct(u); err != nil {
		return err
	}

	return nil
}
