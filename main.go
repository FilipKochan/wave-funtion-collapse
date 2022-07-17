package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := NewGame()
	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("Wave Function Collapse")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
