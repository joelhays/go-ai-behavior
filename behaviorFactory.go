package main

import "github.com/go-gl/mathgl/mgl64"

type BehaviorFactory struct {
}

func (f *BehaviorFactory) Create(behaviorType string, config interface{}) []Behavior {

	data := config.(map[interface{}]interface{})

	var behaviors []Behavior

	switch behaviorType {
	case "constant":
		weight := data["weight"].(float64)
		directionData := data["direction"].([]interface{})
		direction := mgl64.Vec2{
			directionData[0].(float64),
			directionData[1].(float64),
		}
		behaviors = append(behaviors, NewConstantBehavior(weight, direction))
	case "keyboard":
		weight := data["weight"].(float64)
		behaviors = append(behaviors, NewKeyboardBehavior(weight))
	case "wander":
		weight := data["weight"].(float64)
		changeInterval := data["changeInterval"].(int)
		behaviors = append(behaviors, NewWanderBehavior(weight, changeInterval))
	case "seek":
		weight := data["weight"].(float64)
		targetName := data["target"].(string)
		for _, actor := range game.actors {
			if actor.name == targetName {
				behaviors = append(behaviors, NewSeekBehavior(weight, actor))
			}
		}
	case "avoid":
		weight := data["weight"].(float64)
		radius := data["radius"].(float64)
		targetName := data["target"].(string)
		for _, actor := range game.actors {
			if actor.name == targetName {
				behaviors = append(behaviors, NewAvoidBehavior(weight, actor, radius))
			}
		}
	}

	return behaviors
}
