package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
)

type StoryOption struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Story struct {
	Title   string        `json:"title"`
	Story   []string      `json:"story"`
	Options []StoryOption `json:"options"`
}

func StoryHandler(s map[string]Story, t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		arc := r.URL.Path[1:]
		if arc == "" {
			arc = "intro"
		}
		story, ok := s[arc]
		if !ok {
			http.NotFound(w, r)
		}
		t.Execute(w, story)
	}
}

func main() {
	// Load stories.json
	f, err := os.Open("stories.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var stories map[string]Story
	json.NewDecoder(f).Decode(&stories)

	// Load story.html
	pageTemplate, err := template.ParseFiles("story.html")
	if err != nil {
		panic(err)
	}

	// Start HTTP server
	ipAddr := ":8081"
	fmt.Printf("Starting the server on %s\n", ipAddr)
	storyHandler := StoryHandler(stories, pageTemplate)
	err = http.ListenAndServe(ipAddr, storyHandler)
	if err != nil {
		panic(fmt.Sprintf("Error starting the server on %s: %s", ipAddr, err))
	}
}
