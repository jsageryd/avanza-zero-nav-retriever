package main

import (
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strings"
)

type Price struct {
	date     string
	amount   string
	currency string
}

func priceFromURL(url string) Price {
	node := theirSite(url)
	return priceFromNode(node)
}

func theirSite(url string) *html.Node {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func priceFromNode(node *html.Node) Price {
	productNode := nodeByAttr("itemtype", "http://schema.org/Product", node)

	// Date
	reviewNode := nodeByAttr("itemtype", "http://schema.org/Review", productNode)
	dateNode := nodeByAttr("itemprop", "datePublished", reviewNode)
	date := strings.TrimSpace(dateNode.FirstChild.Data)
	if date == "" {
		log.Fatal("Cannot find their date")
	}

	// Amount and currency
	offerNode := nodeByAttr("itemtype", "http://schema.org/Offer", productNode)
	priceNode := nodeByAttr("itemprop", "price", offerNode)
	currencyNode := nodeByAttr("itemprop", "priceCurrency", offerNode)
	amount := strings.TrimSpace(attrValue("content", priceNode))
	if amount == "" {
		log.Fatal("Cannot find their amount")
	}
	currency := strings.TrimSpace(attrValue("content", currencyNode))
	if amount == "" {
		log.Fatal("Cannot find their currency")
	}

	return Price{
		date:     date,
		amount:   amount,
		currency: currency,
	}
}

func nodeByAttr(attrKey string, attrVal string, node *html.Node) *html.Node {
	if node.Type == html.ElementNode {
		if hasAttr(attrKey, attrVal, node) {
			return node
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		res := nodeByAttr(attrKey, attrVal, c)
		if res != nil {
			return res
		}
	}
	return nil
}

func hasAttr(attrKey string, attrVal string, node *html.Node) bool {
	for _, attr := range node.Attr {
		if (attr.Key == attrKey) && (attr.Val == attrVal) {
			return true
		}
	}
	return false
}

func attrValue(attrKey string, node *html.Node) string {
	for _, attr := range node.Attr {
		if attr.Key == attrKey {
			return attr.Val
		}
	}
	return ""
}
