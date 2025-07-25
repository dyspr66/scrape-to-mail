# Scrape to Mail

Scrapes data from pages, then mails the data to you at a certain time every day.

## How to Use

Make sure you already [have Go installed](https://go.dev/doc/install).

1. Create a file named `.env` in this directory. In it, set up the [config](#config)
2. cd into this directory, then run `go run .`

### Config

Sample .env config

```
TIME=

SCRAPE_TARGET_URL=
ENTRY_QUERY_SELECTOR=
TITLE_QUERY_SELECTOR=
LINK_QUERY_SELECTOR=

MAILFROM=
MAILPASS=
MAILTO=
MAILHOST=
MAILPORT=
```

Config Guide:

-   TIME: the time the program will scrape and mail data to you
-   SCRAPE_TARGET_URL: the url of the page you want to scrape
-   ENTRY_QUERY_SELECTOR: the query selector for a single entry
-   TITLE_QUERY_SELECTOR: the query selector for an entry title (which should be within the entry)
-   LINK_QUERY_SELECTOR: the query selector for an entry link (which should be within the entry)
-   MAILFROM: the email you'll mail the data from
-   MAILPASS: the password to the email in MAILFROM
-   MAILTO: the email you'll receive the scraped data from
-   MAILHOST: the host for the email in MAILFROM
-   MAILPORT: the port for MAILHOST

## More on Scraping

This program is best for pages whose scraped data can be represented as a list of _entries_, each with a _title_ and _link_.

For example, a news website with a bunch of articles. Each article would be an _entry_, each article's title would be the entry _title_, and the link to the article would be the entry _link_.
