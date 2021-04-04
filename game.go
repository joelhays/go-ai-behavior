package main

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	actors []*Actor
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) AddActor(actor *Actor) {
	g.actors = append(g.actors, actor)
}

func (g *Game) Update() error {
	for _, actor := range g.actors {
		actor.Update()

		if actor.position.X() > float64(config.Window.Width) {
			actor.position = mgl64.Vec2{
				0,
				actor.position.Y(),
			}
		} else if actor.position.X() < 0.0 {
			actor.position = mgl64.Vec2{
				float64(config.Window.Width),
				actor.position.Y(),
			}
		}

		if actor.position.Y() > float64(config.Window.Height) {
			actor.position = mgl64.Vec2{
				actor.position.X(),
				0,
			}
		} else if actor.position.Y() < 0.0 {
			actor.position = mgl64.Vec2{
				actor.position.X(),
				float64(config.Window.Height),
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x60, 0x60, 0x6f, 0xff})

	for _, actor := range g.actors {
		actor.Draw(screen)
	}

	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
