package utils

import (
	"errors"
	"html"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	xhtml "golang.org/x/net/html"
)

// LoadStringFromFile is just a helper function to load a string from a file. Primarily used for testing
func LoadStringFromFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(strings.TrimSuffix(html.UnescapeString(string(b)), "\n")), nil
}

func StringFromSelector(selector, source string) (string, error) {
	str := ""
	doc, err := QueryDocFromstring(source)
	if err != nil {
		return str, err
	}
	strContainer := doc.Find(selector)
	if strContainer.Length() == 0 {
		return str, errors.New("could not find name container")
	}

	str = strContainer.Text()

	return str, nil
}

func QueryDocFromstring(source string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(source))
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func isRealNode(source string, n *xhtml.Node) bool {
	if n == nil {
		return false
	}
	return strings.Contains(source, n.Data)
}

func firstRealNode(source string, node *xhtml.Node) *xhtml.Node {
	if node == nil {
		return nil
	}
	if isRealNode(source, node) {
		return node
	}
	currentChildNode := node.FirstChild
	for currentChildNode != nil {
		childNode := firstRealNode(source, node.FirstChild)
		if childNode != nil {
			return childNode
		}
		currentChildNode = currentChildNode.NextSibling
	}

	return nil
}

func firstRealNodeInList(source string, nodes []*xhtml.Node) *xhtml.Node {
	var realNode *xhtml.Node
	for _, n := range nodes {
		realNode = firstRealNode(source, n)
		if realNode != nil {
			break
		}
	}
	return realNode
}

func findRootTextNodes(source string, doc *goquery.Document) []*xhtml.Node {
	textNodes := make([]*xhtml.Node, 0)

	allNodes := doc.Children().Children()

	realNode := firstRealNodeInList(source, allNodes.Nodes)

	currentNode := realNode
	for currentNode != nil {
		if currentNode.Type == xhtml.TextNode {
			textNodes = append(textNodes, currentNode)
		}
		currentNode = currentNode.NextSibling
	}

	return textNodes
}

func RemoveChildTextNodes(source string) (string, error) {
	doc, err := QueryDocFromstring(source)
	if err != nil {
		return source, err
	}

	textNodes := findRootTextNodes(source, doc)
	if textNodes == nil {
		return source, nil
	}

	textNodeSelection := doc.FindNodes(textNodes...)
	if textNodeSelection == nil {
		return source, errors.New("could not find text nodes in document")
	}

	textNodeSelection.Remove()

	replaced, err := doc.Html()
	if err != nil {
		return source, err
	}

	replaced = strings.ReplaceAll(replaced, "<html><head></head><body>", "")
	replaced = strings.ReplaceAll(replaced, "</body></html>", "")

	return replaced, nil
}

func RemoveNodesFromString(source, selector string) (string, error) {
	doc, err := QueryDocFromstring(source)
	if err != nil {
		return source, err
	}

	nodes := doc.Find(selector)
	nodes.Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	html, err := doc.Html()
	if err != nil {
		return source, err
	}

	html = strings.ReplaceAll(html, "<html><head></head><body>", "")
	html = strings.ReplaceAll(html, "</body></html>", "")

	return html, nil
}
