package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
}

var currentScreen = 1

const (
	menuScreen int = iota
	gameScreen
	overScreen
)

const (
	GameWidth  int = 160
	GameHeight int = 152
)

func (g *Game) Update() error {
	uiUpdate()
	switch currentScreen {
	case menuScreen:
		menuUpdate()
	case gameScreen:
		gameUpdate()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, float32(GameWidth), float32(GameHeight), getColor(0), false)
	switch currentScreen {
	case menuScreen:
		menuDraw(screen)
	case gameScreen:
		gameDraw(screen)
	}
	drawDebug(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(GameWidth), int(GameHeight)
}

func main() {
	gameInit()
	uiInit()
	ebiten.SetWindowSize(int(GameWidth)*4, int(GameHeight)*4)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetVsyncEnabled(false)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
