package parser

import (
	"freeSteamGamesParser/internal/types"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func HasTargetBlock(html string) (bool, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return false, err
	}

	found := false

	doc.Find(".search_result_row").Each(func(i int, s *goquery.Selection) {
		found = true
	})

	return found, nil
}

func GetGamesLinks(html string) ([]types.Game, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return []types.Game{}, err
	}

	var games []types.Game

	doc.Find(".search_result_row").Each(func(i int, s *goquery.Selection) {
		link, foundHref := s.Attr("href")
		if foundHref != true {
			return
		}

		name := s.Find(".title").First().Text()
		if name == "" {
			return
		}

		game := types.Game{
			Name: s.Text(),
			Link: link,
		}
		games = append(games, game)
	})

	return games, nil
}
