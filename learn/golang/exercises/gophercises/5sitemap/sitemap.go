package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Url struct {
	Loc string `xml:"loc"`
}

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Url     []Url    `xml:"url"`
}

func FetchPageContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func FilterByDomain(links []Link, domain string) []Link {
	filteredLinks := []Link{}
	for _, link := range links {
		if link.Href != "" && (link.Href[0] == '/' || strings.Contains(link.Href, domain)) {
			filteredLinks = append(filteredLinks, link)
		}
	}
	return filteredLinks
}

func GetLinks(url string) ([]Link, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []Link{}, err
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	parser := LinkParser{Tokenizer: tokenizer}
	links, err := parser.Parse()
	if err != nil {
		return []Link{}, err
	}

	links = FilterByDomain(links, url)
	return links, nil
}

func BuildSitemap(rootUrl string, maxDepth int) (Sitemap, error) {
	depth := 0
	url := rootUrl
	queue, err := GetLinks(url)
	if err != nil {
		return Sitemap{}, err
	}
	visited := map[Link]bool{}
	for _, link := range queue {
		visited[link] = true
	}
	counter := len(queue)

	for len(queue) > 0 {
		url, queue = queue[0].Href, queue[1:]
		if url[0] == '/' {
			url = rootUrl + url
		}
		fmt.Println(url)
		links, err := GetLinks(url)
		if err != nil {
			return Sitemap{}, nil
		}
		for _, link := range links {
			if _, ok := visited[link]; !ok {
				visited[link] = true
				queue = append(queue, link)
			}
		}
		counter--
		if counter == 0 {
			counter = len(queue)
			depth++
			if depth >= maxDepth {
				break
			}
		}
	}

	urlset := []Url{}
	for link := range visited {
		urlset = append(urlset, Url{Loc: link.Href})
	}
	return Sitemap{Url: urlset}, nil
}
