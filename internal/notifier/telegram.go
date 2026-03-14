package notifier

import (
	"fmt"
	"net/http"
	"net/url"
)

func SendMessage(token, chatID, text string) error {
	api := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		token,
	)

	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", text)

	_, err := http.PostForm(api, data)

	return err
}
