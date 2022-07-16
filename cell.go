package main

import (
	"fmt"
	"math/rand"
)

type Cell struct {
	tile          *Tile
	possibleTiles []*Tile
	collapsed     bool
}

func (c *Cell) GetEntropy() int {
	return len(c.possibleTiles)
}

func (c *Cell) Collapse() {
	if c.collapsed {
		panic("cell is already collapsed")
	}

	if len(c.possibleTiles) == 0 {
		panic("cell cannot be collapsed, has no options")
	}

	c.tile = c.possibleTiles[rand.Intn(len(c.possibleTiles))]
	c.collapsed = true
}

func (c *Cell) String() string {
	return fmt.Sprintf("Cell{ Tile{ %v } }", c.tile)
}
