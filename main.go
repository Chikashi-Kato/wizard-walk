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

type gameState string

const (
	LogoSpeed = 0.005
	LogoWait = 60 // frame
	WalkSpeed = 5

	LogoState = gameState("logo")
	PlayState = gameState("play")
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

type Logo struct {
	logo     *ebiten.Image
	scale    float64
	duration int // frame
}

type Game struct{
	state gameState
	world World
	logo Logo

	windowSize size
	wizard *Wizard
}

func (g *Game) Update() error {
	switch g.state {
	case PlayState:
		for _, p := range inpututil.PressedKeys() {
			switch p {
			case ebiten.KeyArrowUp:
				if g.wizard.loc.y > 0 {
					if g.wizard.loc.y <= g.getWindowCenter().y && g.world.loc.y == 0 ||
						(g.wizard.loc.y > g.getWindowCenter().y && g.world.loc.y >= g.world.worldSize.height-g.windowSize.height) {
						g.wizard.loc.y -= WalkSpeed
					} else if g.world.loc.y > 0 {
						g.world.loc.y -= WalkSpeed
					}
				}
			case ebiten.KeyArrowLeft:
				if g.wizard.loc.x > 0 {
					if g.wizard.loc.x <= g.getWindowCenter().x && g.world.loc.x == 0 ||
						(g.wizard.loc.x > g.getWindowCenter().x && g.world.loc.x >= g.world.worldSize.width-g.windowSize.width) {
						g.wizard.loc.x -= WalkSpeed
					} else if g.world.loc.x > 0 {
						g.world.loc.x -= WalkSpeed
					}
				}
			case ebiten.KeyArrowDown:
				if g.wizard.loc.y < g.world.worldSize.height {
					if g.wizard.loc.y < g.getWindowCenter().y && g.world.loc.y == 0 ||
						(g.world.loc.y >= g.world.worldSize.height-g.windowSize.height && g.wizard.loc.y < g.windowSize.height-75) {
						g.wizard.loc.y += WalkSpeed
					} else if g.world.loc.y < g.world.worldSize.height-g.windowSize.height {
						g.world.loc.y += WalkSpeed
					}
				}
			case ebiten.KeyArrowRight:
				if g.wizard.loc.x < g.getWindowCenter().x ||
					(g.world.loc.x >= g.world.worldSize.width-g.windowSize.width && g.wizard.loc.x < g.windowSize.width-75) {
					g.wizard.loc.x += WalkSpeed
				} else if g.world.loc.x < g.world.worldSize.width-g.windowSize.width {
					g.world.loc.x += WalkSpeed
				}
			}
		}
		fmt.Println(g.wizard.loc.x, g.wizard.loc.y, g.world.loc.x, g.world.loc.y, g.getWindowCenter().x, g.getWindowCenter().y)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case PlayState:
		croppedMap := g.world.background.SubImage(image.Rectangle{
			Min: image.Point{g.world.loc.x, g.world.loc.y},
			Max: image.Point{g.world.loc.x + g.windowSize.width, g.world.loc.y + g.windowSize.height},
		}).(*ebiten.Image)
		screen.DrawImage(croppedMap, nil)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1.5, 1.5)
		op.GeoM.Translate(float64(g.wizard.loc.x), float64(g.wizard.loc.y))
		screen.DrawImage(g.wizard.left, op)
	case LogoState:
		op := &ebiten.DrawImageOptions{}
		x, y := g.logo.logo.Size()
		x = g.windowSize.width / 2 - x / 2
		y = g.windowSize.height / 2 - y / 2
		op.GeoM.Translate(float64(x), float64(y))

		if g.logo.scale < 1 {
			g.logo.scale += LogoSpeed
			op.ColorM.Scale(g.logo.scale, g.logo.scale, g.logo.scale, 1)
		} else if g.logo.duration < LogoWait {
			g.logo.duration += 1
		} else {
			g.state = PlayState
		}

		screen.DrawImage(g.logo.logo, op)
	}
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
	g.state = LogoState
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

	// Load logo
	g.logo.logo, _, err = ebitenutil.NewImageFromFile("./images/forgotten-runes-logo.png")

	if err != nil {
		log.Fatal("logo file not found")
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}