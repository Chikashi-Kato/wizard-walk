package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	MOVEMENT = 5
)

type size struct{
	width int
	height int
}

type location struct{
	x int
	y int
}

type Wizard struct{
	up *ebiten.Image
	left *ebiten.Image
	down *ebiten.Image
	right *ebiten.Image
	loc location
}

type World struct {
	background *ebiten.Image
	worldSize size
	loc location
}

type Game struct{
	world World

	windowSize size
	wizard *Wizard
}

func (g *Game) Update() error {
	for _, p := range inpututil.PressedKeys() {
		switch p {
		case ebiten.KeyArrowUp:
			if g.wizard.loc.y > 0 {
				if g.wizard.loc.y <= g.getWindowCenter().y && g.world.loc.y == 0 ||
					(g.wizard.loc.y > g.getWindowCenter().y && g.world.loc.y >= g.world.worldSize.height - g.windowSize.height) {
					g.wizard.loc.y -= MOVEMENT
				} else if g.world.loc.y > 0 {
					g.world.loc.y -= MOVEMENT
				}
			}
		case ebiten.KeyArrowLeft:
			if g.wizard.loc.x > 0 {
				if g.wizard.loc.x <= g.getWindowCenter().x && g.world.loc.x == 0 ||
					(g.wizard.loc.x > g.getWindowCenter().x && g.world.loc.x >= g.world.worldSize.width - g.windowSize.width){
					g.wizard.loc.x -= MOVEMENT
				} else if g.world.loc.x > 0 {
					g.world.loc.x -= MOVEMENT
				}
			}
		case ebiten.KeyArrowDown:
			if g.wizard.loc.y < g.world.worldSize.height {
				if g.wizard.loc.y < g.getWindowCenter().y && g.world.loc.y == 0 ||
					(g.world.loc.y >= g.world.worldSize.height-g.windowSize.height && g.wizard.loc.y < g.windowSize.height - 75) {
					g.wizard.loc.y += MOVEMENT
				} else if g.world.loc.y < g.world.worldSize.height-g.windowSize.height {
					g.world.loc.y += MOVEMENT
				}
			}
		case ebiten.KeyArrowRight:
			if g.wizard.loc.x < g.getWindowCenter().x ||
				(g.world.loc.x >= g.world.worldSize.width - g.windowSize.width && g.wizard.loc.x < g.windowSize.width - 75) {
				g.wizard.loc.x += MOVEMENT
			} else if g.world.loc.x < g.world.worldSize.width - g.windowSize.width {
				g.world.loc.x += MOVEMENT
			}
		}
	}
	fmt.Println(g.wizard.loc.x, g.wizard.loc.y, g.world.loc.x, g.world.loc.y, g.getWindowCenter().x, g.getWindowCenter().y)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	croppedMap := g.world.background.SubImage(image.Rectangle{
		Min: image.Point{g.world.loc.x, g.world.loc.y},
		Max: image.Point{g.world.loc.x + g.windowSize.width, g.world.loc.y + g.windowSize.height},
	}).(*ebiten.Image)
	screen.DrawImage(croppedMap, nil)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1.5, 1.5)
	op.GeoM.Translate(float64(g.wizard.loc.x), float64(g.wizard.loc.y))
	screen.DrawImage(g.wizard.left, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.windowSize.width, g.windowSize.height
}

func (g Game) getWindowCenter() location {
	return location{g.windowSize.width / 2 - 40,g.windowSize.height / 2 - 40}
}

func main() {
	var err error

	g := &Game{}
	g.windowSize = size{640, 480}

	ebiten.SetWindowSize(g.windowSize.width,g.windowSize.height)
	ebiten.SetWindowTitle("Forgotten Runes Wizard Walk")

	// Load background

	g.world.background, _, err = ebitenutil.NewImageFromFile("./images/world.png");
	if err != nil {
		log.Fatal("background file not found")
	}
	g.world.worldSize.width, g.world.worldSize.height = g.world.background.Size()

	// Load wizard
	wiz, _, err := ebitenutil.NewImageFromFile("./images/wizard/wizard.png")
	if err != nil {
		log.Fatal("wizard file not found")
	}

	g.wizard = &Wizard{
		up: wiz,
		left: wiz,
		down: wiz,
		right: wiz,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}