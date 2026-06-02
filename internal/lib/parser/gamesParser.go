package parser

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"steamGamesSales/internal/entity"
	repository "steamGamesSales/internal/repository/postgres"
	"strings"

	"golang.org/x/net/html"
)

const (
	steamSearchUri = "https://store.steampowered.com/search/results/"
	paginationSize = 100
)

var steamSearchParams = map[string]string{
	"query":    "",
	"start":    "0",
	"count":    string(rune(paginationSize)),
	"maxprice": "free",
	"specials": "1",
	"infinite": "1",
}

type GamesParser struct {
	log               *slog.Logger
	repo              *repository.Repository
	newFreeGamesEvent chan []*entity.Game
}

func NewGamesParser(log *slog.Logger, repo *repository.Repository, newFreeGamesEvent chan []*entity.Game) *GamesParser {
	return &GamesParser{
		log:               log,
		repo:              repo,
		newFreeGamesEvent: newFreeGamesEvent,
	}
}

func (gamesParser *GamesParser) Run() {
	start := 0
	newGames := make([]*entity.Game, 0)

	for {
		responseBody := gamesParser.Request(start)

		gamesOnPage := gamesParser.parseHtml(responseBody)

		gamesParser.log.Info(fmt.Sprintf(fmt.Sprintf("Free Games found: %d", len(gamesOnPage))))

		newGamesOnPage := gamesParser.filterNewGames(gamesOnPage)
		if len(newGamesOnPage) == 0 {
			gamesParser.log.Info("No new games found")
			break
		}
		newGames = slices.Concat(newGames, newGamesOnPage)

		for _, newGame := range newGamesOnPage {
			err := gamesParser.repo.Game.Create(newGame)
			if err == nil {
				continue
			}
		}

		gamesParser.log.Info(fmt.Sprintf("New Games found: %d", len(newGamesOnPage)))

		start = start + paginationSize

		if start >= responseBody.TotalCount {
			break
		}
	}

	gamesParser.sendNotification(newGames)
}

func (gamesParser *GamesParser) filterNewGames(games []entity.Game) (result []*entity.Game) {
	for _, game := range games {
		addedGame, _ := gamesParser.repo.Game.Get(game.Id)
		if addedGame == nil {
			result = append(result, &game)
		}
	}

	return result
}

type ResponseData struct {
	Success     int    `json:"success"`
	ResultsHtml string `json:"results_html"`
	TotalCount  int    `json:"total_count"`
	Start       int    `json:"start"`
}

func (gamesParser *GamesParser) Request(page int) (responseData ResponseData) {
	req, err := http.NewRequest("GET", steamSearchUri, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)

		return ResponseData{}
	}

	q := req.URL.Query()
	for key, value := range steamSearchParams {
		q.Add(key, value)
	}
	q.Set("start", fmt.Sprintf("%d", page))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)

		return ResponseData{}
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)

		return ResponseData{}
	}

	return responseData
}

func (gamesParser *GamesParser) parseHtml(body ResponseData) []entity.Game {
	doc, _ := html.Parse(strings.NewReader(body.ResultsHtml))

	var result []entity.Game

	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "search_result_row") {
					game := parseGameFromTagA(node)
					if game == nil {
						return
					}

					result = append(result, *game)
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return result
}

func (gamesParser *GamesParser) sendNotification(games []*entity.Game) {
	if len(games) == 0 {
		return
	}

	gamesParser.newFreeGamesEvent <- games
}

func parseGameFromTagA(node *html.Node) *entity.Game {
	// id
	id := getAttr(node, "data-ds-appid")

	// link
	link := getAttr(node, "href")

	// name
	nameBlock := findChildByClass(node, "title")
	if nameBlock == nil {
		return nil
	}
	name := getText(nameBlock)

	// image
	imageContainer := findChildByClass(node, "search_capsule")
	imageBlock := imageContainer.FirstChild
	if imageBlock.Type != html.ElementNode || imageBlock.Data != "img" {
		return nil
	}
	image := getAttr(imageBlock, "src")

	return &entity.Game{
		Id:        id,
		Name:      name,
		Link:      link,
		ImageLink: image,
	}
}

// Извлечение текста из узла
func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += getText(c)
	}
	return result
}

func getAttr(node *html.Node, attributeName string) string {
	if node.Type == html.TextNode {
		return ""
	}

	for _, attr := range node.Attr {
		if attr.Key == attributeName {
			return attr.Val
		}
	}

	return ""
}
func findChildByClass(node *html.Node, class string) *html.Node {
	var f func(*html.Node) *html.Node
	f = func(node *html.Node) *html.Node {
		if node.Type == html.ElementNode {
			for _, attr := range node.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, class) {
					return node
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			result := f(c)
			if result == nil {
				continue
			}

			return result
		}

		return nil
	}

	return f(node)
}
