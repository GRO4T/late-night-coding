package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

type LinkParser struct {
	Tokenizer   *html.Tokenizer
	CurrentLink *Link
	Links       []Link
}

func (p *LinkParser) ParseLinkText() error {
	for {
		token, err := p.GetToken()
		if err != nil {
			return err
		}

		switch token {
		case html.StartTagToken:
			err = p.ParseLinkText()
		case html.TextToken:
			p.CurrentLink.Text += string(p.Tokenizer.Text())
		case html.EndTagToken:
			return nil
		}
	}
}

func (p *LinkParser) TryParseLink(token html.TokenType) error {
	if token != html.StartTagToken {
		return nil
	}

	tagName, hasAttr := p.Tokenizer.TagName()

	if string(tagName) != "a" {
		return nil
	}

	for attrName, attrValue, moreAttr := p.Tokenizer.TagAttr(); hasAttr; hasAttr = moreAttr {
		if string(attrName) != "href" {
			continue
		}
		p.CurrentLink = &Link{Href: string(attrValue), Text: ""}
		break
	}

	err := p.ParseLinkText()
	if err != nil {
		return err
	}

	p.Links = append(p.Links, *p.CurrentLink)
	return nil
}

func (p *LinkParser) GetToken() (html.TokenType, error) {
	token := p.Tokenizer.Next()
	if token == html.ErrorToken {
		return token, p.Tokenizer.Err()
	}
	return token, nil
}

func (p *LinkParser) Parse() ([]Link, error) {
	p.Links = []Link{}
	p.CurrentLink = nil

	for {
		token, err := p.GetToken()
		if err == io.EOF {
			return p.Links, nil
		} else if err != nil {
			return nil, err
		}

		err = p.TryParseLink(token)
		if err != nil {
			return nil, err
		}
	}
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	tokenizer := html.NewTokenizer(f)
	parser := LinkParser{Tokenizer: tokenizer}
	links, err := parser.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Println(links)
}
