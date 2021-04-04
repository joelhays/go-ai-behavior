package main

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

type Behavior interface {
	GetWeight() float64
	Update(actor *Actor)
}

type BehaviorConstant struct {
	weight    float64
	direction mgl64.Vec2
}

func NewBehaviorConstant(weight float64, direction mgl64.Vec2) Behavior {
	return &BehaviorConstant{
		weight:    weight,
		direction: direction.Normalize(),
	}
}

func (b *BehaviorConstant) GetWeight() float64 {
	return b.weight
}

func (b *BehaviorConstant) Update(actor *Actor) {
	actor.direction = actor.direction.Add(b.direction.Mul(b.weight))
}

type BehaviorKeyboard struct {
	weight float64
}

func NewBehaviorKeyboard(weight float64) Behavior {
	return &BehaviorKeyboard{
		weight: weight,
	}
}

func (b *BehaviorKeyboard) GetWeight() float64 {
	return b.weight
}

func (b *BehaviorKeyboard) Update(actor *Actor) {
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

type BehaviorWander struct {
	weight         float64
	changeInterval int
	tick           int
	direction      mgl64.Vec2
}

func NewBehaviorWander(weight float64, changeInterval int) Behavior {
	b := &BehaviorWander{
		weight:         weight,
		changeInterval: changeInterval,
	}
	return b
}

func (b *BehaviorWander) GetWeight() float64 {
	return b.weight
}

func (b *BehaviorWander) Update(actor *Actor) {
	if b.tick == 0 {
		b.direction = GetRandomDirection()
	}
	b.tick++
	b.tick %= b.changeInterval

	actor.direction = actor.direction.Add(b.direction.Mul(b.weight))
}

type BehaviorSeek struct {
	weight float64
	target *Actor
}

//todo add radius for seek
func NewBehaviorSeek(weight float64, target *Actor) Behavior {
	b := &BehaviorSeek{
		weight: weight,
		target: target,
	}
	return b
}

func (b *BehaviorSeek) GetWeight() float64 {
	return b.weight
}

func (b *BehaviorSeek) Update(actor *Actor) {
	targetDirection := b.target.position.Sub(actor.position).Normalize()
	actor.direction = actor.direction.Add(targetDirection.Mul(b.weight))
}

type BehaviorAvoid struct {
	weight float64
	target *Actor
	radius float64
}

func NewBehaviorAvoid(weight float64, target *Actor, radius float64) Behavior {
	b := &BehaviorAvoid{
		weight: weight,
		target: target,
		radius: radius,
	}
	return b
}

func (b *BehaviorAvoid) GetWeight() float64 {
	return b.weight
}

func (b *BehaviorAvoid) Update(actor *Actor) {
	targetDirection := actor.position.Sub(b.target.position)
	if targetDirection.Len() > b.radius {
		return
	}
	targetDirection = targetDirection.Normalize()
	actor.direction = actor.direction.Add(targetDirection.Mul(b.weight))
}
