package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	their_site := their_site()
	their_date := their_date(their_site)
	their_price := their_price(their_site)
	our_last_price := our_last_price()
	fmt.Println(our_last_price)
	fmt.Printf("P %s ZERO %s\n", their_date, their_price)
}

func our_last_price() string {
	price_db := fmt.Sprintf("%s/Documents/ledger/prices/price-db", os.Getenv("HOME"))
	return price_db
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
