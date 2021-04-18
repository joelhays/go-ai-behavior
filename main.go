package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
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

	f, err := os.Create("./go-ai-behavior.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)

	runtime.LockOSThread()

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
		state.screenWidth = config.Window.Width
		state.screenHeight = config.Window.Height
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
	defer pprof.StopCPUProfile()
	defer runtime.UnlockOSThread()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
