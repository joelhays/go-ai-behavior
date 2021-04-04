package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/viper"
)

var (
	game   *Game
	config appConfig
)

func init() {
	rand.Seed(time.Now().Unix())

	loadConfig()

	fmt.Printf("%+v\n", config)

	ebiten.SetWindowTitle(config.Window.Title)
	ebiten.SetWindowSize(config.Window.Width, config.Window.Height)

	game = NewGame()

	for _, actorCfg := range config.Actors {
		for i := 0; i < actorCfg.Count; i++ {
			actor := NewActor(actorCfg.Image, actorCfg.Color, actorCfg.Width, actorCfg.Height)
			actor.speed = actorCfg.Speed
			actor.direction = GetRandomDirection()
			actor.position = GetRandomPosition(config.Window.Width, config.Window.Height)

			game.AddActor(actor)
		}
	}

	setActorBehavior()

	// config.Actors = append(config.Actors,
	// 	struct {
	// 		Name      string
	// 		Count     int
	// 		Image     string
	// 		Width     int
	// 		Height    int
	// 		Speed     float64
	// 		Color     color.RGBA
	// 		Behaviors []struct {
	// 			Type string
	// 			Data interface{}
	// 		}
	// 	}{
	// 		Name:   "whatever",
	// 		Count:  0,
	// 		Image:  "arrow.png",
	// 		Width:  1,
	// 		Height: 1,
	// 		Speed:  1,
	// 		Color:  color.RGBA{1, 2, 3, 1},
	// 		Behaviors: []struct {
	// 			Type string
	// 			Data interface{}
	// 		}{
	// 			{
	// 				Type: "BehaviorConstant",
	// 				Data: struct {
	// 					Weight    float64
	// 					Direction mgl64.Vec2
	// 				}{
	// 					Weight:    0.5,
	// 					Direction: GetRandomDirection(),
	// 				},
	// 			},
	// 			{
	// 				Type: "BehaviorKeyboard",
	// 				Data: struct {
	// 					Weight float64
	// 				}{
	// 					Weight: 0.1,
	// 				},
	// 			},
	// 		},
	// 	})
	// fmt.Printf("%+v\n", config)

	// viper.Set("Actors", config.Actors)
	// viper.WriteConfig()
}

func setActorBehavior() {
	for _, actorCfg := range config.Actors {
		for i := 0; i < actorCfg.Count; i++ {
			game.actors[i].speed = actorCfg.Speed

			var behaviors []Behavior
			for _, behaviorCfg := range actorCfg.Behaviors {
				data := behaviorCfg.Data.(map[interface{}]interface{})

				switch behaviorCfg.Type {
				case "BehaviorConstant":
					weight := data["weight"].(float64)
					directionData := data["direction"].([]interface{})
					direction := mgl64.Vec2{
						directionData[0].(float64),
						directionData[1].(float64),
					}
					behaviors = append(behaviors, NewBehaviorConstant(weight, direction))
				case "BehaviorKeyboard":
					weight := data["weight"].(float64)
					behaviors = append(behaviors, NewBehaviorKeyboard(weight))
				case "BehaviorWander":
					weight := data["weight"].(float64)
					changeInterval := data["changeInterval"].(int)
					behaviors = append(behaviors, NewBehaviorWander(weight, changeInterval))
				}
			}
			game.actors[i].behaviors = behaviors
		}
	}
}

func viperConfigChanged(e fsnotify.Event) {
	newConfig := appConfig{}
	err := viper.Unmarshal(&newConfig)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	config = newConfig

	setActorBehavior()
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(viperConfigChanged)
}

func main() {
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
