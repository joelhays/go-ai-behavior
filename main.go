package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/spf13/viper"
)

var (
	state  *State
	game   *Game
	config appConfig
)

func init() {
	rand.Seed(time.Now().Unix())

	loadConfig()

	ebiten.SetWindowTitle(config.Window.Title)
	ebiten.SetWindowSize(config.Window.Width, config.Window.Height)

	configureState()

	game = NewGame(state)
}

func loadImage(name string) *ebiten.Image {
	if len(name) == 0 {
		return nil
	}

	if img, ok := state.images[name]; ok {
		return img
	}

	image, _, err := ebitenutil.NewImageFromFile("images/" + name)
	if err != nil {
		log.Printf("%v", err)
		return nil
	}
	state.images[name] = image
	return image
}

func configureState() {
	if state == nil {
		state = &State{}
		state.images = make(map[string]*ebiten.Image)
		state.actors = make(map[string][]*Actor)
	}

	state.backgroundImage = loadImage(config.Window.Background.Image)
	state.backgroundColor = config.Window.Background.Color

	for _, actorCfg := range config.Actors {
		if len(state.actors[actorCfg.Name]) > actorCfg.Count {
			state.actors[actorCfg.Name] = state.actors[actorCfg.Name][:actorCfg.Count]
		} else if len(state.actors[actorCfg.Name]) < actorCfg.Count {
			numToAdd := actorCfg.Count - len(state.actors[actorCfg.Name])
			for i := 0; i < numToAdd; i++ {
				actor := &Actor{}

				if actorCfg.Position.Type != "static" {
					actor.position = GetRandomPosition(config.Window.Width, config.Window.Height)
				}
				actor.direction = GetRandomDirection()
				state.actors[actorCfg.Name] = append(state.actors[actorCfg.Name], actor)
			}
		}

		for i := 0; i < actorCfg.Count; i++ {
			actor := state.actors[actorCfg.Name][i]

			actor.name = actorCfg.Name
			actor.image = loadImage(actorCfg.Image)

			actor.width = actorCfg.Width
			actor.height = actorCfg.Height
			actor.name = actorCfg.Name
			actor.color = actorCfg.Color
			actor.speed = actorCfg.Speed

			actor.static = false
			if actorCfg.Position.Type == "static" {
				actor.static = true
				directionData := actorCfg.Position.Data.([]interface{})
				actor.position = mgl64.Vec2{
					directionData[0].(float64),
					directionData[1].(float64),
				}
			}
		}
	}

	for _, actorCfg := range config.Actors {
		for i := 0; i < actorCfg.Count; i++ {
			actor := state.actors[actorCfg.Name][i]

			var behaviors []Behavior
			factory := BehaviorFactory{}

			for _, behaviorCfg := range actorCfg.Behaviors {
				newBehaviors := factory.Create(behaviorCfg.Type, behaviorCfg.Data)
				behaviors = append(behaviors, newBehaviors...)
			}

			actor.behaviors = behaviors
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

	configureState()

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

// func debugSaveConfig() {
// 	config.Actors = append(config.Actors,
// 		struct {
// 			Name      string
// 			Count     int
// 			Image     string
// 			Width     int
// 			Height    int
// 			Speed     float64
// 			Color     color.RGBA
// 			Behaviors []struct {
// 				Type string
// 				Data interface{}
// 			}
// 		}{
// 			Name:   "whatever",
// 			Count:  0,
// 			Image:  "arrow.png",
// 			Width:  1,
// 			Height: 1,
// 			Speed:  1,
// 			Color:  color.RGBA{1, 2, 3, 1},
// 			Behaviors: []struct {
// 				Type string
// 				Data interface{}
// 			}{
// 				{
// 					Type: "BehaviorConstant",
// 					Data: struct {
// 						Weight    float64
// 						Direction mgl64.Vec2
// 					}{
// 						Weight:    0.5,
// 						Direction: GetRandomDirection(),
// 					},
// 				},
// 				{
// 					Type: "BehaviorKeyboard",
// 					Data: struct {
// 						Weight float64
// 					}{
// 						Weight: 0.1,
// 					},
// 				},
// 			},
// 		})
// 	fmt.Printf("%+v\n", config)

// 	viper.Set("Actors", config.Actors)
// 	viper.WriteConfig()
// }
