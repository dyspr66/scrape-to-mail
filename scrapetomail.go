package main

import (
	"fmt"
	"log/slog"
	"net/smtp"
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
		return
	}

	timeToExecute := os.Getenv("TIME")
	s := gocron.NewScheduler(time.UTC)
	_, err = s.Every(1).Days().At(timeToExecute).Do(scrapeThenMail)
	if err != nil {
		slog.Error("Specifying func to call", "err", err)
		return
	}

	s.StartBlocking()
}

func scrapeThenMail() {
	e, err := scrape()
	if err != nil {
		slog.Error("Scraping", "err", err)
		return
	}

	msg, err := prepareEntries(e)
	if err != nil {
		slog.Error("Preparing entires", "err", err)
		return
	}

	err = sendEmail(msg)
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

func prepareEntries(entries []entry) (string, error) {
	// TODO - reading this and customizability is a pain.
	message := "<!DOCTYPE html>\n<html lang=\"en\">\n<body>\n<h1>Entries</h1>\n<ul>\n"
	for _, entry := range entries {
		message += fmt.Sprintf("<li><a href=\"%s\">%s</a></li>\n", entry.link, entry.title)
	}
	message += "</ul>\n</body>\n</html>"
	return message, nil
}

func sendEmail(message string) error {
	MailFrom := os.Getenv("MAILFROM")
	MailPass := os.Getenv("MAILPASS")
	MailTo := os.Getenv("MAILTO")
	MailHost := os.Getenv("MAILHOST")
	MailPort := os.Getenv("MAILPORT")

	addr := MailHost + ":" + MailPort
	auth := smtp.PlainAuth("", MailFrom, MailPass, MailHost)
	to := []string{MailTo}

	// TODO - does this really need to be written out here?
	// Adding email headers for rendering html
	message = fmt.Sprintf("To: %s\nFrom: %s\nSubject: %s\nMIME-version: 1.0;\nContent-Type: text/html; charset=UTF-8;\n", MailTo, MailFrom, "Entries") + message

	err := smtp.SendMail(addr, auth, MailFrom, to, []byte(message))
	if err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	return nil
}
