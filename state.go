package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type State struct {
	screenWidth     int
	screenHeight    int
	backgroundImage *ebiten.Image
	backgroundColor color.RGBA
	actors          map[string][]*Actor
	images          map[string]*ebiten.Image
}
