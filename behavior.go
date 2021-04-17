package main

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

type Behavior interface {
	Update(actor *Actor)
}

type ConstantBehavior struct {
	weight    float64
	direction mgl64.Vec2
}

func NewConstantBehavior(weight float64, direction mgl64.Vec2) Behavior {
	return &ConstantBehavior{
		weight:    weight,
		direction: direction.Normalize(),
	}
}

func (b *ConstantBehavior) Update(actor *Actor) {
	actor.direction = actor.direction.Add(b.direction.Mul(b.weight))
}

type KeyboardBehavior struct {
	weight float64
}

func NewKeyboardBehavior(weight float64) Behavior {
	return &KeyboardBehavior{
		weight: weight,
	}
}

func (b *KeyboardBehavior) Update(actor *Actor) {
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		actor.direction = actor.direction.Add(mgl64.Vec2{0, 1}.Mul(b.weight))
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		actor.direction = actor.direction.Add(mgl64.Vec2{0, -1}.Mul(b.weight))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		actor.direction = actor.direction.Add(mgl64.Vec2{-1, 0}.Mul(b.weight))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		actor.direction = actor.direction.Add(mgl64.Vec2{1, 0}.Mul(b.weight))
	}
}

type WanderBehavior struct {
	weight         float64
	changeInterval int
	tick           int
	direction      mgl64.Vec2
}

func NewWanderBehavior(weight float64, changeInterval int) Behavior {
	b := &WanderBehavior{
		weight:         weight,
		changeInterval: changeInterval,
	}
	return b
}

func (b *WanderBehavior) Update(actor *Actor) {
	if b.tick == 0 {
		b.direction = GetRandomDirection()
	}
	b.tick++
	b.tick %= b.changeInterval

	actor.direction = actor.direction.Add(b.direction.Mul(b.weight))
}

type SeekBehavior struct {
	weight float64
	target *Actor
}

//todo add radius for seek
func NewSeekBehavior(weight float64, target *Actor) Behavior {
	b := &SeekBehavior{
		weight: weight,
		target: target,
	}
	return b
}

func (b *SeekBehavior) Update(actor *Actor) {
	targetDirection := b.target.position.Sub(actor.position).Normalize()
	actor.direction = actor.direction.Add(targetDirection.Mul(b.weight))
}

type AvoidBehavior struct {
	weight float64
	target *Actor
	radius float64
}

func NewAvoidBehavior(weight float64, target *Actor, radius float64) Behavior {
	b := &AvoidBehavior{
		weight: weight,
		target: target,
		radius: radius,
	}
	return b
}

func (b *AvoidBehavior) Update(actor *Actor) {
	targetDirection := actor.position.Sub(b.target.position)
	if targetDirection.Len() > b.radius {
		return
	}
	targetDirection = targetDirection.Normalize()
	actor.direction = actor.direction.Add(targetDirection.Mul(b.weight))
}

type RotateBehavior struct {
	weight float64
	speed  float64
}

func NewRotateBehavior(weight float64, speed float64) Behavior {
	b := &RotateBehavior{
		weight: weight,
		speed:  speed,
	}
	return b
}

func (b *RotateBehavior) Update(actor *Actor) {
	actor.rotation += b.speed * 1.0 / 60.0
}
