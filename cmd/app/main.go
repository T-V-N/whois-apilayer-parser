package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/T-V-N/whois-api-parser/internal/config"
	"github.com/T-V-N/whois-api-parser/internal/db"
)

type Result struct {
	Result string `json:"result"`
}

func main() {
	cfg, err := config.Init()

	if err != nil {
		log.Panic(err)
	}

	dbConn, err := db.Init(cfg)

	if err != nil {
		log.Panic(err)
	}

	linkChannel := make(chan string)

	go proccessLink(context.Background(), linkChannel, cfg, dbConn)
	go proccessLink(context.Background(), linkChannel, cfg, dbConn)
	go proccessLink(context.Background(), linkChannel, cfg, dbConn)
	go proccessLink(context.Background(), linkChannel, cfg, dbConn)
	go proccessLink(context.Background(), linkChannel, cfg, dbConn)
	go proccessLink(context.Background(), linkChannel, cfg, dbConn)
	go proccessLink(context.Background(), linkChannel, cfg, dbConn)
	go proccessLink(context.Background(), linkChannel, cfg, dbConn)
	go proccessLink(context.Background(), linkChannel, cfg, dbConn)
	go proccessLink(context.Background(), linkChannel, cfg, dbConn)

	fetchLink(dbConn, linkChannel)
}

func fetchLink(db *db.DBStorage, c chan string) {
	for {
		links, _ := db.GetUnproccessedDomains(context.Background())

		if len(links) == 0 {
			close(c)
			return
		}

		for _, link := range links {
			c <- link
		}
	}
}

func proccessLink(ctx context.Context, c chan string, cfg *config.Config, db *db.DBStorage) {
	for {
		var link string
		link, ok := <-c

		if !ok {
			return
		}

		url, err := url.Parse(strings.TrimSuffix(link, "\r"))
		if err != nil {
			fmt.Print(err)
		}

		parts := strings.Split(url.Host, ".")
		if len(parts) < 2 {
			continue
		}
		domain := parts[len(parts)-2] + "." + parts[len(parts)-1]

		whoisResult, err := fetchWhois(cfg, domain)
		if err != nil {
			fmt.Print(err)
		}

		err = db.UpdateDomainAvailability(ctx, link, whoisResult)
		if err != nil {
			fmt.Print(err)
		}
	}

}

func fetchWhois(cfg *config.Config, domain string) (string, error) {
	url := "https://api.apilayer.com/whois/check?domain=" + domain
	fmt.Println(domain)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", cfg.APIKey)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)

	if res.Body != nil {
		defer res.Body.Close()
	}

	var available Result

	if err := json.NewDecoder(res.Body).Decode(&available); err != nil {
		return "", err
	}

	return available.Result, nil
}
