package main

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	up = iota
	right
	down
	left
)

type Board struct {
	grid   []*Cell
	width  int
	height int
}

func (b *Board) CellAt(x int, y int) *Cell {
	index := x + (b.width)*y
	if x < 0 || y < 0 {
		return nil
	}
	if index >= len(b.grid) || index < 0 {
		return nil
	}
	return b.grid[index]
}

func (b *Board) UpdateEntropies() {
	for i := 0; i < b.width; i++ {
		for j := 0; j < b.height; j++ {
			if b.CellAt(i, j).collapsed {
				continue
			}
			b.calculatePossibleTilesAt(i, j)
		}
	}
	fmt.Println("=== UPDATING DONE ===")
}

func (b *Board) IsFull() bool {
	for i := 0; i < b.width; i++ {
		for j := 0; j < b.height; j++ {
			if !b.CellAt(i, j).collapsed {
				return false
			}
		}
	}
	return true
}

func (b *Board) calculatePossibleTilesAt(cx int, cy int) {
	xDiff := [4]int{0, 1, 0, -1}
	yDiff := [4]int{-1, 0, 1, 0}

	curr := b.CellAt(cx, cy)
	overallPossible := map[int]int{}

	neighborsCount := 0

	for dir := 0; dir < 4; dir++ {
		newX := xDiff[dir] + cx
		newY := yDiff[dir] + cy

		neighbor := b.CellAt(newX, newY)

		if neighbor == nil {
			continue
		}

		if !neighbor.collapsed {
			continue
		}

		fmt.Printf("Cell [%v, %v] has collapsed neighbor at [%v, %v]\n", cx, cy, newX, newY)
		neighborsCount++

		placedTile := neighbor.tile
		newPossibleTiles := []int{}
		for i, thisTile := range curr.possibleTiles {
			fmt.Printf("at [%v, %v]: checking whether tile (%v) %v is possible\n", cx, cy, i, thisTile)
			switch dir {
			case up:
				if thisTile.sides.top.ConnectsTo(&placedTile.sides.bottom) {
					newPossibleTiles = append(newPossibleTiles, i)
					fmt.Printf("Tile (%v) connects to placed tile %v from UP\n", i, *placedTile)
				}
			case right:
				if thisTile.sides.right.ConnectsTo(&placedTile.sides.left) {
					newPossibleTiles = append(newPossibleTiles, i)
					fmt.Printf("Tile (%v) connects to placed tile %v from RIGHT\n", i, *placedTile)
				}
			case down:
				if thisTile.sides.bottom.ConnectsTo(&placedTile.sides.top) {
					newPossibleTiles = append(newPossibleTiles, i)
					fmt.Printf("Tile (%v) connects to placed tile %v from DOWN\n", i, *placedTile)
				}
			case left:
				if thisTile.sides.left.ConnectsTo(&placedTile.sides.right) {
					newPossibleTiles = append(newPossibleTiles, i)
					fmt.Printf("Tile (%v) connects to placed tile %v from LEFT\n", i, *placedTile)
				}
			}
		}
		fmt.Printf("new possible options: %v\n", newPossibleTiles)

		for _, v := range newPossibleTiles {
			overallPossible[v]++
		}
	}

	if neighborsCount == 0 {
		fmt.Printf("at [%v, %v] has no collapsed neighbors, returning...\n", cx, cy)
		return
	}

	fmt.Printf("overall possible options: %v\n", overallPossible)
	fmt.Printf("neighbors: %v\n", neighborsCount)
	result := []*Tile{}
	for k, v := range overallPossible {
		if v == neighborsCount {
			result = append(result, curr.possibleTiles[k])
		}
	}
	curr.possibleTiles = result
	fmt.Printf("new possible tiles at [%v, %v]: %v\n\n", cx, cy, curr.possibleTiles)
}

func (b *Board) GetCellWithLeastEntropy() *Cell {
	selected := make([]*Cell, 0)
	leastEntropy := math.MaxInt

	for i := 0; i < b.width; i++ {
		for j := 0; j < b.height; j++ {
			c := b.CellAt(i, j)
			if c.collapsed {
				continue
			}
			if c.GetEntropy() < leastEntropy {
				leastEntropy = c.GetEntropy()
			}
		}
	}

	for i := 0; i < b.width; i++ {
		for j := 0; j < b.height; j++ {
			c := b.CellAt(i, j)
			if c.collapsed {
				continue
			}

			if c.GetEntropy() == leastEntropy {
				selected = append(selected, c)
			}
		}
	}

	// fmt.Println("getting cell with least entropy")
	// fmt.Printf("possible options: %v\n", selected)

	if len(selected) == 0 {
		return nil
	}

	i := rand.Intn(len(selected))

	return selected[i]
}
