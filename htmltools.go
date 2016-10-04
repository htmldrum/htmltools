package main

import (
	"fmt"
	"os"
	"sort"
	"golang.org/x/net/html"
)

type VisibleNode struct {
	NodeType html.NodeType
	Contents string
}


// PrintContents prints the contents of all user-visible text nodes in a node tree
func main(){
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "htmltools: %v\n", err)
	}
	res := visitPrintContents(nil, doc)

	for _, s := range res {
		fmt.Printf("S): %v\n", s.Contents)
	}
	fmt.Printf("Total: %d\n", len(res))

}

func visitPrintContents(elems []VisibleNode, n *html.Node) []VisibleNode {
	if n.Type == html.TextNode && !blank(n.Data) {
		elems = append(elems, VisibleNode{n.Type, n.Data})
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if n.Data == "script" || n.Data == "style" {
				break
			}
		}
		elems = visitPrintContents(elems, c)
	}

	return elems
}

func blank(s string) bool {
	for _, v := range s {
		if v != ' ' && v != '\n' && v != '\t' {
			return false
		}
	}
	return true
}

// CountElems prints a map of DOM element types to their frequency in a node tree
func CountElems() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "htmltools: %v\n", err)
	}
	res := visitCountElems(nil, doc)
	summary := countRes(res)

	var alpha []string

	for k, _ := range summary {
		alpha = append(alpha, k)
	}
	sort.Strings(alpha)
	fmt.Printf("Rank\tType\tOccurs\n")
	fmt.Printf("==================\n")
	for v, k := range alpha {
		fmt.Printf("%v)\t%v\t%d\n", v, k, summary[k])
	}

}

func visitCountElems(elems []string, n *html.Node) []string {
	if n.Type == html.ElementNode {
		elems = append(elems, n.Data)
	}

	for c:= n.FirstChild; c != nil; c = c.NextSibling {
		elems = visitCountElems(elems, c)
	}

	return elems

}

func countRes(elems []string) map[string]int {
	uniq := make(map[string]int)
	for _, v := range elems {
		uniq[v]++
	}
	return uniq
}
