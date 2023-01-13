package main

import (
	"log"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func main() {
	img, _, err := ebitenutil.NewImageFromFile("gopher.png")
	if err != nil {
		log.Fatalf("unable to load image due: %v", err)
	}
	ebiten.SetWindowSize(240, 240)
	ebiten.SetWindowDecorated(false)
	ebiten.SetScreenTransparent(true)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&Game{charaImg: img}); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	charaImg *ebiten.Image
}

var initCursorX, initCursorY int

func (g *Game) Update() error {
	isPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	if isPressed {
		cursorX, cursorY := ebiten.CursorPosition()
		isJustPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
		if isJustPressed {
			initCursorX = cursorX
			initCursorY = cursorY
		}

		windowX, windowY := ebiten.WindowPosition()

		// log.Printf("cursorX: %v, cursorY: %v", cursorX, cursorY)
		// log.Printf("windowX: %v, windowY: %v", windowX, windowY)

		actualCursorX := windowX + cursorX
		actualCursorY := windowY + cursorY
		// log.Printf("actualCursorX: %v, actualCursorY: %v", actualCursorX, actualCursorY)

		ebiten.SetWindowPosition(actualCursorX-initCursorX, actualCursorY-initCursorY)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// screen.Fill(color.RGBA{0xff, 0, 0, 0xff})
	screen.DrawImage(g.charaImg, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 240, 240
}
