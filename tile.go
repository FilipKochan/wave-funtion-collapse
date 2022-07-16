package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	image    ebiten.Image
	sides    Sides
	rotation int
}

func (t *Tile) Rotated(times int) *Tile {
	newSides := t.sides
	for i := 0; i < times; i++ {
		newSides = Sides{newSides.left, newSides.top, newSides.right, newSides.bottom}
	}

	return &Tile{image: t.image, rotation: (t.rotation + times) % 4, sides: newSides}
}

func (t *Tile) String() string {
	return fmt.Sprintf("Tile{ rotation: %v, sides: %v, image: %v }", t.rotation, t.sides, t.image)
}

type Sides struct {
	top    Side
	right  Side
	bottom Side
	left   Side
}

type Side struct {
	left   int
	middle int
	right  int
}

const (
	SideEmpty = iota
	SidePipe
)

func (s *Side) ConnectsTo(other *Side) bool {
	res := s.right == other.left && s.middle == other.middle && s.left == other.right
	fmt.Printf("checking connectivity between cells: %v and %v | connects: %v\n", *s, *other, res)
	return res
}
