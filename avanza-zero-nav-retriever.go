package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	their_site := their_site()
	their_date := their_date(their_site)
	their_price := their_price(their_site)
	their_price_line := fmt.Sprintf("P %s ZERO %s", their_date, their_price)
	our_price_line := our_price_line()
	if their_price_line != our_price_line {
		fmt.Println("Ours:   ", our_price_line)
		fmt.Println("Theirs: ", their_price_line)
		append_price_line(their_price_line)
		commit_change(their_date)
	}
}

func commit_change(their_date string) {
	git_commit := exec.Command("git", "commit", "-am", fmt.Sprintf("Update %s", their_date))
	price_db_dir := filepath.Dir(price_db())
	git_commit.Dir = price_db_dir
	git_push := exec.Command("git", "push")
	git_push.Dir = price_db_dir
	err := git_commit.Run()
	if err != nil {
		log.Fatal("git_commit ", err)
	}
	err = git_push.Run()
	if err != nil {
		log.Fatal("git_push ", err)
	}
}

func append_price_line(price_line string) {
	file, err := os.OpenFile(price_db(), os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(price_line + "\n")
}

func last_line(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}
	return line
}

func price_db() string {
	return fmt.Sprintf("%s/Documents/ledger/prices/price-db", os.Getenv("HOME"))
}

func our_price_line() string {
	return last_line(price_db())
}

func their_price(node *html.Node) string {
	product_node := node_by_attr("itemtype", "http://schema.org/Product", node)
	offer_node := node_by_attr("itemtype", "http://schema.org/Offer", product_node)
	price_node := node_by_attr("itemprop", "price", offer_node)
	currency_node := node_by_attr("itemprop", "priceCurrency", offer_node)
	price := attr_value("content", price_node)
	currency := attr_value("content", currency_node)
	return fmt.Sprintf("%s %s", price, currency)
}

func their_date(node *html.Node) string {
	product_node := node_by_attr("itemtype", "http://schema.org/Product", node)
	review_node := node_by_attr("itemtype", "http://schema.org/Review", product_node)
	date_node := node_by_attr("itemprop", "datePublished", review_node)
	return strings.TrimSpace(date_node.FirstChild.Data)
}

func their_site() *html.Node {
	uri := "https://www.avanza.se/fonder/om-fonden.html/41567/avanza-zero"
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func node_by_attr(attr_key string, attr_val string, node *html.Node) *html.Node {
	if node.Type == html.ElementNode {
		if has_attr(attr_key, attr_val, node) {
			return node
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		res := node_by_attr(attr_key, attr_val, c)
		if res != nil {
			return res
		}
	}
	return nil
}

func has_attr(attr_key string, attr_val string, node *html.Node) bool {
	for _, attr := range node.Attr {
		if (attr.Key == attr_key) && (attr.Val == attr_val) {
			return true
		}
	}
	return false
}

func attr_value(attr_key string, node *html.Node) string {
	for _, attr := range node.Attr {
		if attr.Key == attr_key {
			return attr.Val
		}
	}
	return ""
}
