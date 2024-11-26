package main

import (
	"encoding/xml"
	"flag"
	"os"
)

func main() {
	rootUrl := flag.String("u", "", "Root URL to build sitemap")
	maxDepth := flag.Int("d", 3, "Max depth to crawl")
	flag.Parse()
	if *rootUrl == "" {
		panic("Root URL is required")
	}

	sitemap, err := BuildSitemap(*rootUrl, *maxDepth)
	if err != nil {
		panic(err)
	}

	sitemapXml, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("sitemap.xml", sitemapXml, 0644)
	if err != nil {
		panic(err)
	}
}
