package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

var url, ident string
var git Git
var priceFile PriceFile

func main() {
	parseOpts()
	updatePrice()
}

func parseOpts() {
	var repo string
	flag.StringVar(&repo, "repo", "", "path to Git workdir containing price-db file")
	flag.StringVar(&url, "url", "", "URL to Avanza page from which to fetch data")
	flag.StringVar(&ident, "ident", "", "price identifier use in price-db")
	flag.Parse()

	repo = strings.TrimSpace(repo)
	url = strings.TrimSpace(url)
	ident = strings.TrimSpace(ident)
	if repo == "" {
		log.Fatal("Need repo")
	}
	if url == "" {
		log.Fatal("Need URL")
	}
	if ident == "" {
		log.Fatal("Need ident")
	}
	git = Git{repo}
	priceFile = PriceFile{strings.TrimRight(repo, "/") + "/price-db"}
	if !git.RepoValid() {
		log.Fatal("Not a repo")
	}
}

func updatePrice() {
	theirPrice := priceFromURL(url)
	theirStr := fmt.Sprintf("P %s %s %s", theirPrice.date, theirPrice.amount, theirPrice.currency)
	ourStr := priceFile.lastLine()
	if theirStr != ourStr {
		priceFile.addLine(theirStr)
		git.Add("price-db")
		git.Commit("Update " + theirPrice.date)
		git.Push()
	}
}