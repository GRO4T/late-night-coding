package main

import (
	"io"

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
			if err != nil {
				return err
			}
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

	tagName, _ := p.Tokenizer.TagName()

	if string(tagName) != "a" {
		return nil
	}

	for {
		attrName, attrValue, moreAttr := p.Tokenizer.TagAttr()
		if string(attrName) == "href" {
			p.CurrentLink = &Link{Href: string(attrValue), Text: ""}
			break
		} else if !moreAttr {
			break
		}
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
		if err == io.EOF {
			return p.Links, nil
		} else if err != nil {
			return nil, err
		}
	}
}
