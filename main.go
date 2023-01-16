package main

import (
	"log"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 72
	screenHeight = 64
)

func main() {
	// initialize game
	game, err := NewGame(GameConfig{
		ActionIdleSource: &ActionSource{
			ImagePaths: []string{
				"assets/idle1.png",
				"assets/idle2.png",
				"assets/idle3.png",
				"assets/idle4.png",
			},
		},
		ActionSleepSource: &ActionSource{
			ImagePaths: []string{
				"assets/zzz1.png",
				"assets/zzz2.png",
				"assets/zzz3.png",
				"assets/zzz4.png",
			},
		},
		ExitButtonImagePath: "assets/close.png",
		ScreenDimension: Dimension{
			Width:  screenWidth,
			Height: screenHeight,
		},
	})
	if err != nil {
		log.Fatalf("unable to initialize game due: %v", err)
	}
	// run game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
