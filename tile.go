package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tile struct {
	image    ebiten.Image
	sides    Sides
	rotation int
}

type TileOption struct {
	ImageName  string `json:"imageName"`
	Sides      Sides  `json:"sides"`
	Rotate     bool   `json:"rotate"`
	RotateOnce bool   `json:"rotateOnce"`
}

type TilesetOption struct {
	TilesetName string       `json:"tilesetName"`
	Tiles       []TileOption `json:"tiles"`
}

func (t *Tile) Rotated(times int) *Tile {
	newSides := t.sides
	for i := 0; i < times; i++ {
		newSides = Sides{newSides.Left, newSides.Top, newSides.Right, newSides.Bottom}
	}

	return &Tile{image: t.image, rotation: (t.rotation + times) % 4, sides: newSides}
}

func (t *Tile) String() string {
	return fmt.Sprintf("Tile{ rotation: %v, sides: %v }", t.rotation, t.sides)
}

func ParseTilesetsFromJSON() []TilesetOption {
	var options []TilesetOption
	data, err := os.ReadFile("tilesets.json")
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(data, &options)
	if err != nil {
		panic(err.Error())
	}

	return options
}

func CreateTileset(name string, options []TilesetOption) []*Tile {
	tiles := []*Tile{}

	for _, tileset := range options {
		if tileset.TilesetName == name {
			for _, tile := range tileset.Tiles {
				img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("tiles/%v/%v", tileset.TilesetName, tile.ImageName))
				if err != nil {
					log.Fatal(err.Error())
				}

				t := &Tile{image: *img, sides: tile.Sides}
				tiles = append(tiles, t)
				if tile.Rotate {
					tiles = append(tiles, t.Rotated(1), t.Rotated(2), t.Rotated(3))
				} else if tile.RotateOnce {
					tiles = append(tiles, t.Rotated(1))
				}
			}
		}
	}
	return tiles
}

func GetTilesetsNames(options []TilesetOption) []string {
	tilesetNames := []string{}
	for _, v := range options {
		tilesetNames = append(tilesetNames, v.TilesetName)
	}
	return tilesetNames
}

func IsTilesetValid(name string, options []TilesetOption) bool {
	for _, v := range options {
		if v.TilesetName == name {
			return true
		}
	}
	return false
}

type Sides struct {
	Top    Side `json:"top"`
	Right  Side `json:"right"`
	Bottom Side `json:"bottom"`
	Left   Side `json:"left"`
}

type Side struct {
	Left   int `json:"left"`
	Middle int `json:"middle"`
	Right  int `json:"right"`
}

const (
	SideEmpty = iota
	SidePipe
)

func (s *Side) ConnectsTo(other *Side) bool {
	res := s.Right == other.Left && s.Middle == other.Middle && s.Left == other.Right
	return res
}
