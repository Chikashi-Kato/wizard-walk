package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	background *ebiten.Image
	worldSize  size
	loc        location
}

func getWorld() (*World, error) {
	w := World{}

	// Get the data
	worldMap, err := loadEbitenImageFromUrl("https://www.forgottenrunes.com/static/img/map/map.png")
	if err != nil {
		return nil, fmt.Errorf("failed to create ebiten image: %w", err)
	}

	w.background = worldMap
	w.worldSize.width, w.worldSize.height = w.background.Size()

	return &w, nil
}
