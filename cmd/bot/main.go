package main

import (
	"freeSteamGamesParser/internal/monitor"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}

	url, found := os.LookupEnv("PAGE_LINK")
	if found == false {
		log.Fatal("Не найден env PAGE_LINK")
	}

	token, found := os.LookupEnv("TOKEN")
	if found == false {
		log.Fatal("Не найден env TOKEN")
	}

	chatID, found := os.LookupEnv("BOT_ID")
	if found == false {
		log.Fatal("Не найден env BOT_ID")
	}

	monitor.Run(url, token, chatID)
	log.Print("Успэх")
}
