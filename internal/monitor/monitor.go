package monitor

import (
	"freeSteamGamesParser/internal/types"
	"log"
	"time"

	"freeSteamGamesParser/internal/fetcher"
	"freeSteamGamesParser/internal/notifier"
	"freeSteamGamesParser/internal/parser"
)

func Run(url string, token string, chatID string) int {
	for {
		// TODO: подключение к боту

		html, err := fetcher.Fetch(url)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Minute)
			continue
		}

		gamesList, err := parser.GetGamesLinks(html)
		if err != nil {
			log.Println(err)
			continue
		}

		if len(gamesList) > 0 {
			err := notifier.SendMessage(
				token,
				chatID,
				formatText(gamesList),
			)

			if err != nil {
				return 0
			}

			break
		}
	}

	return 1
}

func formatText(gamesList []types.Game) string {
	result := "Бесплатные игры: \n"

	for _, game := range gamesList {
		result += game.Name + " " + game.Link + "\n"
	}

	return result
}
