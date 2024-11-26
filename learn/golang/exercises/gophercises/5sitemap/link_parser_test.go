package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

var testParseParams = []struct {
	testPageFilename string
	expectedLinks    []Link
}{
	{
		"test_data/ex1.html",
		[]Link{
			{Href: "/other-page", Text: "A link to another page"},
		},
	},
	{
		"test_data/ex2.html",
		[]Link{
			{
				Href: "https://www.twitter.com/joncalhoun",
				Text: "\n      Check me out on twitter\n      \n    ",
			},
			{
				Href: "https://github.com/gophercises",
				Text: "\n      Gophercises is on Github!\n    ",
			},
		},
	},
	{
		"test_data/ex3.html",
		[]Link{
			{
				Href: "#",
				Text: "Login ",
			},
			{
				Href: "/lost",
				Text: "Lost? Need help?",
			},
			{
				Href: "https://twitter.com/marcusolsson",
				Text: "@marcusolsson",
			},
		},
	},
	{
		"test_data/ex4.html",
		[]Link{
			{
				Href: "/dog-cat",
				Text: "dog cat ",
			},
		},
	},
}

func TestParse(t *testing.T) {
	for _, testParams := range testParseParams {
		// Arrange
		f, err := os.Open(testParams.testPageFilename)
		if err != nil {
			t.Errorf("Error opening file: %v", err)
		}
		defer f.Close()
		tokenizer := html.NewTokenizer(f)

		// Act
		parser := LinkParser{Tokenizer: tokenizer}
		links, err := parser.Parse()

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, testParams.expectedLinks, links)
	}
}
