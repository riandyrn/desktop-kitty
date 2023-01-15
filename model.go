package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Action struct {
	Name   string
	Images []ebiten.Image
}

type ActionSource struct {
	Name       string
	ImagePaths []string
	Priority   int
}

func (src ActionSource) ToActions() ([]Action, error) {
	var acts []Action
	for i := 0; i < src.Priority; i++ {
		images := make([]ebiten.Image, 0, len(src.ImagePaths))
		for _, imgPath := range src.ImagePaths {
			img, _, err := ebitenutil.NewImageFromFile(imgPath)
			if err != nil {
				return nil, fmt.Errorf("unable to load image due: %v", err)
			}
			images = append(images, *img)
		}
		acts = append(acts, Action{
			Name:   src.Name,
			Images: images,
		})
	}
	return acts, nil
}

type Point struct {
	X int
	Y int
}

type Dimension struct {
	Width  int
	Height int
}
