package main

import (
	_ "image/png"

	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func GetRandomPosition(rangeX, rangeY int) mgl64.Vec2 {
	return mgl64.Vec2{
		rand.Float64() * float64(rangeX),
		rand.Float64() * float64(rangeY),
	}
}

func GetRandomDirection() mgl64.Vec2 {
	rotation := rand.Float64() * (math.Pi * 2)
	return mgl64.Vec2{
		math.Cos(rotation),
		math.Sin(rotation),
	}
}

type Actor struct {
	color     color.RGBA
	image     *ebiten.Image
	position  mgl64.Vec2
	direction mgl64.Vec2
	speed     float64
	origin    mgl64.Vec2

	width  int
	height int

	behaviors []Behavior
}

func NewActor(imageFile string, color color.RGBA, width int, height int) *Actor {
	image, _, err := ebitenutil.NewImageFromFile(imageFile)
	if err != nil {
		log.Fatal(err)
	}

	imageWidth, imageHeight := image.Size()
	actor := &Actor{
		image:  image,
		color:  color,
		origin: mgl64.Vec2{float64(imageWidth) / 2, float64(imageHeight) / 2},
		width:  width,
		height: height,
	}
	return actor
}

func (actor *Actor) Update() {
	for _, behavior := range actor.behaviors {
		behavior.Update(actor)
	}

	if actor.direction.Len() > 0 {
		actor.direction = actor.direction.Normalize()
	}

	actor.position = actor.position.Add(actor.direction.Mul(actor.speed))
}

func (actor *Actor) Draw(screen *ebiten.Image) {
	rotation := math.Atan2(actor.direction.Y(), actor.direction.X())

	op := &ebiten.DrawImageOptions{}
	w, h := actor.image.Size()

	widthScale := float64(actor.width) / float64(w)
	heightScale := float64(actor.height) / float64(h)

	sw := float64(w) * widthScale
	sh := float64(h) * heightScale

	op.GeoM.Scale(widthScale, heightScale)

	op.GeoM.Translate(-sw/2.0, -sh/2.0)
	op.GeoM.Rotate(rotation)
	op.GeoM.Translate(sw/2.0, sh/2.0)

	op.GeoM.Translate(actor.position.X(), actor.position.Y())

	// op.ColorM.Apply(actor.color)

	diffR := float64(0xff - actor.color.R)
	diffG := float64(0xff - actor.color.G)
	diffB := float64(0xff - actor.color.B)

	op.ColorM.Translate(-diffR, -diffG, -diffB, 0)

	op.Filter = ebiten.FilterNearest

	screen.DrawImage(actor.image, op)
}
