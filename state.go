package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type State struct {
	backgroundImage *ebiten.Image
	backgroundColor color.RGBA
	actors          map[string][]*Actor
	images          map[string]*ebiten.Image
}
