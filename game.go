package main

import (
	"flag"
	"fmt"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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

func NewGame() *Game {

	var seed int64
	var size int
	var tileset string
	tilesets := ParseTilesetsFromJSON()
	flag.Int64Var(&seed, "seed", time.Now().UTC().Unix(), "random seed for the generation, default is time.Now().UTC().Unix()")
	flag.IntVar(&size, "size", 32, "size of the board")
	flag.StringVar(&tileset, "tileset", tilesets[0].TilesetName, fmt.Sprintf("specifies the set of tiles used to generate the image, possible values: [%v]", strings.Join(GetTilesetsNames(tilesets), ", ")))
	flag.Parse()

	if !IsTilesetValid(tileset, tilesets) {
		fmt.Fprintln(os.Stderr, "invalid tileset option")
		flag.PrintDefaults()
		os.Exit(1)
	}

	rg := rand.New(rand.NewSource(seed))

	fmt.Printf("     === using ===\nrandom seed:\t%v\nboard size:\t%v\ntileset:\t%v\n", seed, size, tileset)

	tiles := CreateTileset(tileset, tilesets)

	grid := make([]*Cell, size*size)
	for i := 0; i < size*size; i++ {
		grid[i] = &Cell{collapsed: false, possibleTiles: tiles}
	}

	board := &Board{grid: grid, width: size, height: size}

	game := &Game{tiles: tiles, board: board, size: size, randGenerator: rg}

	return game
}
