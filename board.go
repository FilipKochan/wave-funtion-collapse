package main

import (
	"math"
	"math/rand"
	"sort"
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
	if x < 0 || y < 0 || x >= b.width || y >= b.height {
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
	// fmt.Println("=== UPDATING DONE ===")
}

func (b *Board) IsFull() bool {
	for i := 0; i < b.width; i++ {
		for j := 0; j < b.height; j++ {
			c := b.CellAt(i, j)
			if !c.collapsed && c.GetEntropy() > 0 {
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

		neighborsCount++

		placedTile := neighbor.tile
		newPossibleTiles := []int{}
		for i, thisTile := range curr.possibleTiles {
			// fmt.Printf("at [%v, %v]: checking whether tile (%v) %v is possible\n", cx, cy, i, thisTile)
			switch dir {
			case up:
				if thisTile.sides.Top.ConnectsTo(&placedTile.sides.Bottom) {
					newPossibleTiles = append(newPossibleTiles, i)
					// fmt.Printf("Tile (%v) connects to placed tile %v from UP\n", i, *placedTile)
				}
			case right:
				if thisTile.sides.Right.ConnectsTo(&placedTile.sides.Left) {
					newPossibleTiles = append(newPossibleTiles, i)
					// fmt.Printf("Tile (%v) connects to placed tile %v from RIGHT\n", i, *placedTile)
				}
			case down:
				if thisTile.sides.Bottom.ConnectsTo(&placedTile.sides.Top) {
					newPossibleTiles = append(newPossibleTiles, i)
					// fmt.Printf("Tile (%v) connects to placed tile %v from DOWN\n", i, *placedTile)
				}
			case left:
				if thisTile.sides.Left.ConnectsTo(&placedTile.sides.Right) {
					newPossibleTiles = append(newPossibleTiles, i)
					// fmt.Printf("Tile (%v) connects to placed tile %v from LEFT\n", i, *placedTile)
				}
			}
		}
		// fmt.Printf("new possible options: %v\n", newPossibleTiles)

		for _, v := range newPossibleTiles {
			overallPossible[v]++
		}
	}

	if neighborsCount == 0 {
		// fmt.Printf("at [%v, %v] has no collapsed neighbors, returning...\n", cx, cy)
		return
	}

	result := []*Tile{}
	overallPossibleList := []int{}

	for k, v := range overallPossible {
		if v == neighborsCount {
			overallPossibleList = append(overallPossibleList, k)
		}
	}

	sort.Slice(overallPossibleList, func(i, j int) bool {
		return overallPossibleList[i] < overallPossibleList[j]
	})

	for _, v := range overallPossibleList {
		result = append(result, curr.possibleTiles[v])
	}

	curr.possibleTiles = result
}

func (b *Board) GetCellWithLeastEntropy(randGenerator *rand.Rand) *Cell {
	selected := []*Cell{}
	leastEntropy := math.MaxInt

	for i := 0; i < b.width; i++ {
		for j := 0; j < b.height; j++ {
			c := b.CellAt(i, j)
			if c.collapsed {
				continue
			}
			e := c.GetEntropy()
			if e < leastEntropy && e > 0 {
				leastEntropy = e
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

	if len(selected) == 0 {
		return nil
	}

	i := randGenerator.Intn(len(selected))
	result := selected[i]
	return result
}
