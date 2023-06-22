package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
}

var currentScreen int

const (
	menuScreen int = iota
	gameScreen
	overScreen
)

func switchScreen(screen int) {
	fmt.Println("switchScreen()")
	switch screen {
	case 1:
		gameInit()
	}
	currentScreen = screen
}

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
	uiInit()
	switchScreen(1)
	ebiten.SetWindowSize(int(GameWidth)*5, int(GameHeight)*5)
	ebiten.SetWindowTitle("REVERSTRIS")
	ebiten.SetVsyncEnabled(false)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
