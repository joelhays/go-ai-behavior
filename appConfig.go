package main

import "image/color"

type appConfig struct {
	Window struct {
		Title      string
		Width      int
		Height     int
		Background struct {
			Image string
			Color color.RGBA
		}
	}
	Actors []struct {
		Name      string
		Count     int
		Image     string
		Width     int
		Height    int
		Speed     float64
		Color     color.RGBA
		Behaviors []struct {
			Type string
			Data interface{}
		}
	}
}
