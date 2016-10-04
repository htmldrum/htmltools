package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		outline(url)
	}
}

func outline(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	forEachNode(doc, startElement, endElement)

	return nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}


var depth int

func startElement(n *html.Node) {
	if isElement(n) {
		if inlineElement(n) {
			fmt.Printf("%*s<%s/>\n", depth*2, "", n.Data)
		} else {
			fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
			if inlineData(n){
				depth++
				fmt.Printf("%*s%s\n", depth*2, "", n.FirstChild.Data)
				depth--
			}
		}
		depth++
	}
}

func endElement(n *html.Node) {
	if isElement(n) {
		depth--
		if !inlineElement(n) {
			fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
		}
	}
}

func isElement(n *html.Node) bool {
	return n.Type == html.ElementNode
}
func isRoot(n *html.Node) bool {
	return n.FirstChild == nil
}

func inlineElement(n *html.Node) bool {
	return isRoot(n) ||
		n.Data == "br" ||
		n.Data == "img"
}

func inlineData(n *html.Node) bool {
	return n.Data == "script" ||
		n.Data == "link" ||
		n.Data == "style" ||
		n.Data == "code" ||
		n.Data == "h1" ||
		n.Data == "tt";
}
