package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Loading .env", "err", err)
	}

	timeToExecute := os.Getenv("TIME")
	s := gocron.NewScheduler(time.UTC)
	_, err = s.Every(1).Days().At(timeToExecute).Do(scrapeThenMail)
	if err != nil {
		slog.Error("Specifying func to call", "err", err)
	}

	s.StartBlocking()
}

func scrapeThenMail() {
	_, err := scrape() // TODO - handle entries
	if err != nil {
		slog.Error("Scraping", "err", err)
		return
	}

	err = sendEmail()
	if err != nil {
		slog.Error("Sending email", "err", err)
		return
	}

	slog.Info("Success.")
}

type entry struct {
	title string
	link  string
}

func scrape() ([]entry, error) {
	var entries []entry

	url := os.Getenv("SCRAPE_TARGET_URL")
	entrySelector := os.Getenv("ENTRY_QUERY_SELECTOR")
	titleSelector := os.Getenv("TITLE_QUERY_SELECTOR")
	linkSelector := os.Getenv("LINK_QUERY_SELECTOR")

	c := colly.NewCollector()
	c.OnHTML(entrySelector, func(elem *colly.HTMLElement) {
		t := strings.TrimSpace(elem.ChildText(titleSelector))
		l := strings.TrimSpace(elem.ChildAttr(linkSelector, "href"))

		if t == "" || l == "" {
			return
		}

		entries = append(entries, entry{title: t, link: l})
	})

	err := c.Visit(url)
	if err != nil {
		return entries, fmt.Errorf("visiting target url: %w", err)
	}

	return entries, nil
}

// TODO
func sendEmail() error {
	return nil
}
