package main

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// rand.Seed(42)
	rand.Seed(69)
	game := NewGame(8)
	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("Wave Function Collapse")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
