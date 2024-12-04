package main

import (
	"fmt"

	"github.com/GRO4T/deck-demo/deck"
)

func main() {
	fmt.Println("Hello World")
	myDeck := deck.New()
	deck.Shuffle(myDeck)
	fmt.Println(myDeck)
}
