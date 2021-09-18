package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type gameState string

const (
	LogoSpeed = 0.005
	LogoWait  = 60 // frame
	WalkSpeed = 5

	InputDelay = 15

	LogoState          = gameState("logo")
	WizardIdInputState = gameState("wizard")
	PlayState          = gameState("play")
)

var alagardFont font.Face

type size struct {
	width  int
	height int
}

type location struct {
	x int
	y int
}

type Logo struct {
	logo     *ebiten.Image
	scale    float64
	duration int // frame
}

type Game struct {
	state      gameState
	logo       Logo
	windowSize size

	world    *World
	wizard   *Wizard
	wizardId string

	// Wait function
	isNop           bool
	nopFrameCounter uint
	nopFrameLength  uint
}

func (g *Game) Update() error {
	if g.checkNop() {
		return nil
	}

	switch g.state {
	case WizardIdInputState:
		for _, p := range inpututil.PressedKeys() {
			fmt.Println(p.String())

			switch p {
			case ebiten.KeyNumpad0,
				ebiten.KeyDigit0,
				ebiten.KeyNumpad1,
				ebiten.KeyDigit1,
				ebiten.KeyNumpad2,
				ebiten.KeyDigit2,
				ebiten.KeyNumpad3,
				ebiten.KeyDigit3,
				ebiten.KeyNumpad4,
				ebiten.KeyDigit4,
				ebiten.KeyNumpad5,
				ebiten.KeyDigit5,
				ebiten.KeyNumpad6,
				ebiten.KeyDigit6,
				ebiten.KeyNumpad7,
				ebiten.KeyDigit7,
				ebiten.KeyNumpad8,
				ebiten.KeyDigit8,
				ebiten.KeyNumpad9,
				ebiten.KeyDigit9:
				if len(g.wizardId) < 4 {
					_, num := trimLastChar(p.String())
					g.wizardId += num
					g.wait(InputDelay)
				}
			case ebiten.KeyBackspace:
				if len(g.wizardId) > 0 {
					g.wizardId, _ = trimLastChar(g.wizardId)
					g.wait(InputDelay)
				}
			case ebiten.KeyEnter:
				if len(g.wizardId) > 0 {
					// Initialize wizard
					wizId, err := strconv.Atoi(g.wizardId)
					if err != nil {
						log.Fatal("Invalid wizard Id")
					}

					if g.wizard, err = getWizard(uint(wizId)); err != nil {
						log.Fatal("initializing wizard failed")
					}

					g.state = PlayState
				}
			}
		}
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
	case WizardIdInputState:
		text.Draw(screen, "Type Wizard ID, and then 'Enter'", alagardFont, g.windowSize.width/5, g.windowSize.height/3, color.White)
		if len(g.wizardId) > 0 {
			text.Draw(screen, g.wizardId, alagardFont, g.windowSize.width/2-len(g.wizardId)*5, g.windowSize.height/2, color.White)
		}
	case LogoState:
		op := &ebiten.DrawImageOptions{}
		x, y := g.logo.logo.Size()
		x = g.windowSize.width/2 - x/2
		y = g.windowSize.height/2 - y/2
		op.GeoM.Translate(float64(x), float64(y))

		if g.logo.scale < 1 {
			g.logo.scale += LogoSpeed
			op.ColorM.Scale(g.logo.scale, g.logo.scale, g.logo.scale, 1)
		} else if g.logo.duration < LogoWait {
			g.logo.duration += 1
		} else {
			g.state = WizardIdInputState
		}

		screen.DrawImage(g.logo.logo, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.windowSize.width, g.windowSize.height
}

func (g *Game) getWindowCenter() location {
	return location{g.windowSize.width/2 - 40, g.windowSize.height/2 - 40}
}

func (g *Game) wait(frames uint) {
	g.nopFrameCounter = 0
	g.nopFrameLength = frames
	g.isNop = true
}

func (g *Game) checkNop() bool {
	if g.isNop {
		g.nopFrameCounter += 1
		if g.nopFrameCounter > g.nopFrameLength {
			g.isNop = false
		}
	}

	return g.isNop
}

//go:embed alagard.ttf
var fontData []byte

func main() {
	var err error

	g := &Game{}
	g.state = LogoState
	g.windowSize = size{640, 480}

	ebiten.SetWindowSize(g.windowSize.width, g.windowSize.height)
	ebiten.SetWindowTitle("Forgotten Runes Wizard Walk")

	// Load background
	if g.world, err = getWorld(); err != nil {
		log.Fatal("initializing world failed")
	}

	// Load logo
	logoImage, err := loadEbitenImageFromUrl("https://www.forgottenrunes.com/static/img/forgotten-runes-logo.png")
	if err != nil {
		log.Fatal("logo file not found")
	}

	g.logo.logo = logoImage

	// font
	tt, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	alagardFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
