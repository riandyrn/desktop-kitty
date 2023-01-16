package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"gopkg.in/validator.v2"
)

type GameConfig struct {
	ActionSourceIdle         *ActionSource `validate:"nonnil"`
	ActionSourceSleep        *ActionSource `validate:"nonnil"`
	ActionSourceWalkingLeft  *ActionSource `validate:"nonnil"`
	ActionSourceWalkingRight *ActionSource `validate:"nonnil"`
	ExitButtonImagePath      string        `validate:"min=1"`
	ScreenDimension          Dimension     `validate:"nonzero"`
}

func (c GameConfig) Validate() error {
	return validator.Validate(c)
}

func NewGame(cfg GameConfig) (*Game, error) {
	// validate config
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %+v", cfg)
	}
	// load exit button image
	exitButtonImage, _, err := ebitenutil.NewImageFromFile(cfg.ExitButtonImagePath)
	if err != nil {
		return nil, fmt.Errorf("unable to load exit button image due: %w", err)
	}
	// load actions
	actionIdle, err := cfg.ActionSourceIdle.ToAction(ActionTypeIdle)
	if err != nil {
		return nil, fmt.Errorf("unable to load idle action due: %w", err)
	}
	actionSleep, err := cfg.ActionSourceSleep.ToAction(ActionTypeSleep)
	if err != nil {
		return nil, fmt.Errorf("unable to load sleep action due: %w", err)
	}
	actionWalkingLeft, err := cfg.ActionSourceWalkingLeft.ToAction(ActionTypeWalkingLeft)
	if err != nil {
		return nil, fmt.Errorf("unable to load walking left action due: %w", err)
	}
	actionWalkingRight, err := cfg.ActionSourceWalkingRight.ToAction(ActionTypeWalkingRight)
	if err != nil {
		return nil, fmt.Errorf("umable to load walking right action due: %w", err)
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
		actionIdle:         actionIdle,
		actionSleep:        actionSleep,
		actionWalkingLeft:  actionWalkingLeft,
		actionWalkingRight: actionWalkingRight,
		exitButtonImage:    exitButtonImage,
		windowPos:          windowPos,
		screenDimension:    cfg.ScreenDimension,
		currentAction:      actionWalkingRight,
	}

	return g, nil
}

type Game struct {
	actionIdle         *Action
	actionSleep        *Action
	actionWalkingLeft  *Action
	actionWalkingRight *Action
	exitButtonImage    *ebiten.Image
	displayImgTick     int
	windowPos          Point
	screenDimension    Dimension
	currentAction      *Action

	lastLeftClickPos Point
	displayImage     *ebiten.Image
	lastImgIdx       int
}

func (g *Game) Update() error {
	// get current cursor position
	cursorX, cursorY := ebiten.CursorPosition()
	cursorPos := Point{X: cursorX, Y: cursorY}
	// increment display image tick
	g.updateDisplayImage()
	// wake up kitty if necessary, we put this here because we want
	// to quickly catch the click event, if we put this in Draw() it
	// will be more slower
	g.handleWakeUpKittyIfNecessary()
	// check whether user click the exit button
	g.handleExitIfNecessary(cursorPos)
	// update window position
	g.updateWindowPosOnLeftClick(cursorPos)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// set window position according to calculation
	ebiten.SetWindowPosition(g.windowPos.X, g.windowPos.Y)
	// draw character image
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(0, 16)
	screen.DrawImage(g.displayImage, opt)
	// draw exit button, we want to position it on top right
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(g.screenDimension.Width-g.exitButtonImage.Bounds().Dx()), 0)
	screen.DrawImage(g.exitButtonImage, opt)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) updateDisplayImage() {
	// increment image tick
	g.displayImgTick++

	// update display image
	imgIdx := (g.displayImgTick / 30) % len(g.currentAction.Images)
	animLoopCount := (g.displayImgTick / 30) / len(g.currentAction.Images)
	g.displayImage = &g.currentAction.Images[imgIdx]

	// move window position if necessary
	if g.lastImgIdx != imgIdx {
		switch g.currentAction.Type {
		case ActionTypeWalkingLeft:
			g.windowPos.X -= 4
		case ActionTypeWalkingRight:
			g.windowPos.X += 4
		}
	}

	// if animation loop has finished, determine next action
	if imgIdx == 0 && animLoopCount > 0 {
		switch g.currentAction.Type {
		case ActionTypeIdle:
			if animLoopCount > 5 {
				// determine next action: sleep, walking left, walking right
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				nextActionTypes := []ActionType{
					ActionTypeSleep,
					ActionTypeWalkingLeft,
					ActionTypeWalkingRight,
				}
				nextActionType := nextActionTypes[r.Intn(len(nextActionTypes))]
				g.updateCurrentAction(nextActionType)
			}
		case ActionTypeSleep:
			if animLoopCount > 15 {
				g.updateCurrentAction(ActionTypeIdle)
			}
		case ActionTypeWalkingLeft, ActionTypeWalkingRight:
			if animLoopCount > 2 {
				g.updateCurrentAction(ActionTypeIdle)
			}
		}
	}
	g.lastImgIdx = imgIdx
}

func (g *Game) updateCurrentAction(actType ActionType) {
	act := g.actionIdle
	switch actType {
	case ActionTypeSleep:
		act = g.actionSleep
	case ActionTypeWalkingLeft:
		act = g.actionWalkingLeft
	case ActionTypeWalkingRight:
		act = g.actionWalkingRight
	}
	g.currentAction = act
	g.displayImgTick = 0
}

func (g *Game) handleWakeUpKittyIfNecessary() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.updateCurrentAction(ActionTypeIdle)
	}
}

func (g *Game) handleExitIfNecessary(cursorPos Point) {
	// check if the cursor position is above exit button
	isAboveButton := g.isCursorAboveExitButton(cursorPos)
	if isAboveButton && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// when the cursor is above exit button & user click it, this means user
		// want to exit the program, so we do it.
		os.Exit(0)
	}
}

func (g *Game) isCursorAboveExitButton(cursorPos Point) bool {
	btnDimension := Dimension{
		Width:  g.exitButtonImage.Bounds().Dx(),
		Height: g.exitButtonImage.Bounds().Dy(),
	}
	return cursorPos.X >= (g.screenDimension.Width-btnDimension.Width) &&
		cursorPos.X <= g.screenDimension.Width &&
		cursorPos.Y >= 0 && cursorPos.Y <= btnDimension.Height
}

func (g *Game) updateWindowPosOnLeftClick(cursorPos Point) {
	// if left mouse button is clicked, this means user is currently trying to drag
	// the cat, this means we need to make the window position follow this with some
	// adjustments.
	isLeftClick := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if isLeftClick {
		// get current position of the mouse cursor, in ebitengine we could only get
		// cursor position relative to the game window, so later we need some adjustment
		// to this cursor position since what we need is actual mouse cursor position
		isJustPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
		if isJustPressed {
			// record the position of left click when it is just happening, this is useful
			// for our window position final adjustment
			g.lastLeftClickPos = cursorPos
		}

		// get actual position of cursor
		curWindowX, curWindowY := ebiten.WindowPosition()
		actualCursorPos := Point{
			X: curWindowX + cursorPos.X,
			Y: curWindowY + cursorPos.Y,
		}

		// update window position to follow actual cursor position with some adjustment
		g.windowPos.X = actualCursorPos.X - g.lastLeftClickPos.X
		g.windowPos.Y = actualCursorPos.Y - g.lastLeftClickPos.Y
	}
}
