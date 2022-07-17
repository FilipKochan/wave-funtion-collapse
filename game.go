package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	tiles         []*Tile
	board         *Board
	size          int
	randGenerator *rand.Rand
}

func (g *Game) Update() error {
	if g.board.IsFull() {
		return nil
	}

	nextCell := g.board.GetCellWithLeastEntropy(g.randGenerator)
	if nextCell == nil {
		fmt.Println("out of options, can't continue")
		return nil
	}

	nextCell.Collapse(g.randGenerator)

	g.board.UpdateEntropies()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			cell := g.board.CellAt(i, j)
			tile := cell.tile
			if !cell.collapsed {
				continue
			}
			image := &tile.image
			imageWidth := image.Bounds().Size().X
			screenWidth := screen.Bounds().Size().X
			wantedImageWidth := screenWidth / g.size
			ratio := float64(wantedImageWidth) / float64(imageWidth)
			geom := ebiten.GeoM{}
			geom.Scale(ratio, ratio)

			rot := cell.tile.rotation

			xShift := 0
			if rot == 1 || rot == 2 {
				xShift = 1
			}

			yShift := 0
			if rot > 1 {
				yShift = 1
			}

			geom.Rotate(float64(rot) * math.Pi / 2)
			geom.Translate(float64(wantedImageWidth*i)+float64(xShift)*float64(wantedImageWidth), float64(wantedImageWidth*j)+float64(yShift)*float64(wantedImageWidth))
			screen.DrawImage(image, &ebiten.DrawImageOptions{GeoM: geom})
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 800
}

func NewGame(size int) *Game {
	var tiles []*Tile
	for i, path := range []string{"tiles/pipe/blank.png", "tiles/pipe/down.png"} {
		eImage, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Fatal(err.Error())
		}

		t := &Tile{image: *eImage}
		tiles = append(tiles, t)
		if i == 0 {
			t.sides = Sides{
				top:    Side{SideEmpty, SideEmpty, SideEmpty},
				right:  Side{SideEmpty, SideEmpty, SideEmpty},
				bottom: Side{SideEmpty, SideEmpty, SideEmpty},
				left:   Side{SideEmpty, SideEmpty, SideEmpty},
			}
		}
		if i == 1 {
			t.sides = Sides{
				top:    Side{SideEmpty, SideEmpty, SideEmpty},
				right:  Side{SideEmpty, SidePipe, SideEmpty},
				bottom: Side{SideEmpty, SidePipe, SideEmpty},
				left:   Side{SideEmpty, SidePipe, SideEmpty},
			}
			tiles = append(tiles, t.Rotated(1), t.Rotated(2), t.Rotated(3))
		}
	}

	grid := make([]*Cell, size*size)
	for i := 0; i < size*size; i++ {
		grid[i] = &Cell{collapsed: false, possibleTiles: tiles}
	}

	board := &Board{grid: grid, width: size, height: size}

	rg := rand.New(rand.NewSource(69420))

	game := &Game{tiles: tiles, board: board, size: size, randGenerator: rg}

	return game
}
