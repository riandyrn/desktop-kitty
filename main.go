package main

import (
	"fmt"
	"log"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"gopkg.in/validator.v2"
)

const (
	screenWidth  = 72
	screenHeight = 64
)

func main() {
	// initialize game
	game, err := NewGame(GameConfig{
		ImagePaths: []string{
			"assets/idle1.png",
			"assets/idle2.png",
			"assets/idle3.png",
			"assets/idle4.png",
		},
		ScreenDimension: Dimension{Width: screenWidth, Height: screenHeight},
	})
	if err != nil {
		log.Fatalf("unable to initialize game due: %v", err)
	}
	// run game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func NewGame(cfg GameConfig) (*Game, error) {
	// validate config
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %+v", cfg)
	}
	// load images
	images := make([]ebiten.Image, 0, len(cfg.ImagePaths))
	for _, imgPath := range cfg.ImagePaths {
		img, _, err := ebitenutil.NewImageFromFile(imgPath)
		if err != nil {
			return nil, fmt.Errorf("unable to load image due: %v", err)
		}
		images = append(images, *img)
	}
	// adjust window properties
	ebiten.SetWindowSize(cfg.ScreenDimension.Width, cfg.ScreenDimension.Height)
	ebiten.SetWindowDecorated(false)
	ebiten.SetScreenTransparent(true)
	ebiten.SetWindowFloating(true)
	// determine window start position, we start on center of the main display
	maxScreenWidth, maxScreenHeight := ebiten.ScreenSizeInFullscreen()
	windowPos := Point{
		X: (maxScreenWidth / 2) - cfg.ScreenDimension.Width,
		Y: (maxScreenHeight / 2) - cfg.ScreenDimension.Height,
	}

	// initialize game
	g := &Game{
		images:    images,
		windowPos: windowPos,
	}

	return g, nil
}

type GameConfig struct {
	ImagePaths      []string  `validate:"min=1"`
	ScreenDimension Dimension `validate:"nonzero"`
}

func (c GameConfig) Validate() error {
	return validator.Validate(c)
}

type Game struct {
	images           []ebiten.Image
	currImgIdx       int
	windowPos        Point
	lastLeftClickPos Point
}

func (g *Game) Update() error {
	g.currImgIdx = (g.currImgIdx + 1) % 200
	g.updateWindowPosOnLeftClick()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw character image
	imgIdx := (g.currImgIdx / 10) % len(g.images)
	screen.DrawImage(&g.images[imgIdx], nil)
	// set window position according to calculation
	ebiten.SetWindowPosition(g.windowPos.X, g.windowPos.Y)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) updateWindowPosOnLeftClick() {
	// if left mouse button is clicked, this means user is currently trying to drag
	// the cat, this means we need to make the window position follow this with some
	// adjustments.
	isLeftClick := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if isLeftClick {
		// get current position of the mouse cursor, in ebitengine we could only get
		// cursor position relative to the game window, so later we need some adjustment
		// to this cursor position since what we need is actual mouse cursor position
		cursorX, cursorY := ebiten.CursorPosition()
		isJustPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
		if isJustPressed {
			// record the position of left click when it is just happening, this is useful
			// for our window position final adjustment
			g.lastLeftClickPos.X = cursorX
			g.lastLeftClickPos.Y = cursorY
		}

		// get actual position of cursor
		curWindowX, curWindowY := ebiten.WindowPosition()
		actualCursorPos := Point{
			X: curWindowX + cursorX,
			Y: curWindowY + cursorY,
		}

		// update window position to follow actual cursor position with some adjustment
		g.windowPos.X = actualCursorPos.X - g.lastLeftClickPos.X
		g.windowPos.Y = actualCursorPos.Y - g.lastLeftClickPos.Y
	}
}

type Point struct {
	X int
	Y int
}

type Dimension struct {
	Width  int
	Height int
}
