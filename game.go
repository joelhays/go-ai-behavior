package main

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"

	_ "image/jpeg"
	_ "image/png"
)

type Game struct {
	state      *State
	background *ebiten.Image
	static     *ebiten.Image
	dynamic    *ebiten.Image
}

func NewGame(state *State) *Game {
	game := &Game{state: state}

	screenw, screenh := state.screenWidth, state.screenHeight

	game.background = ebiten.NewImage(screenw, screenh)
	game.static = ebiten.NewImage(screenw, screenh)
	game.dynamic = ebiten.NewImage(screenw, screenh)

	return game
}

func (g *Game) Update() error {
	for _, actorGroup := range g.state.actors {
		for _, actor := range actorGroup {
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
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.state.backgroundColor)
	screenw, screenh := screen.Size()

	if g.state.backgroundImage != nil {
		op := &ebiten.DrawImageOptions{}

		w, h := g.state.backgroundImage.Size()
		widthScale := float64(screenw) / float64(w)
		heightScale := float64(screenh) / float64(h)

		diffR := float64(255.0-config.Window.Background.Color.R) / 255.0
		diffG := float64(255.0-config.Window.Background.Color.G) / 255.0
		diffB := float64(255.0-config.Window.Background.Color.B) / 255.0

		op.ColorM.Translate(-diffR, -diffG, -diffB, 0)

		op.GeoM.Scale(widthScale, heightScale)

		game.background.DrawImage(g.state.backgroundImage, op)
	}

	game.static.Clear()
	game.dynamic.Clear()

	for _, actorGroup := range g.state.actors {
		for _, actor := range actorGroup {
			if actor.static {
				actor.Draw(game.static)
			} else {
				actor.Draw(game.dynamic)
			}
		}
	}

	screen.DrawImage(game.background, &ebiten.DrawImageOptions{})
	screen.DrawImage(game.static, &ebiten.DrawImageOptions{})
	screen.DrawImage(game.dynamic, &ebiten.DrawImageOptions{})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
