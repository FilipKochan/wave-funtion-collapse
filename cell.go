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

func (c *Cell) Collapse(randGenerator *rand.Rand) {
	if c.collapsed {
		panic("cell is already collapsed")
	}

	if len(c.possibleTiles) == 0 {
		panic("cell cannot be collapsed, has no options")
	}

	i := randGenerator.Intn(len(c.possibleTiles))
	resultTile := c.possibleTiles[i]
	c.tile = resultTile
	c.collapsed = true
}

func (c *Cell) String() string {
	return fmt.Sprintf("Cell{ Tile{ %v } }", c.tile)
}
