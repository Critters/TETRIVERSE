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
	gamePuzzleScreen
	gameEndlessScreen
	creditsScreen
	winnerScreen
)

func switchScreen(screen int) {
	fmt.Printf("switchScreen(%d)\n", screen)
	switch screen {
	case 1:
		soundPlay(sfx_start)
		gameInit(0)
	case 2:
		soundPlay(sfx_start)
		gameInit(1)
	case 3:
		creditsInit()
	case 4:
		winnerInit()
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
	case gamePuzzleScreen, gameEndlessScreen:
		gameUpdate()
	case creditsScreen:
		creditsUpdate()
	case winnerScreen:
		winnerUpdate()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, float32(GameWidth), float32(GameHeight), getColor(0), false)
	switch currentScreen {
	case menuScreen:
		menuDraw(screen)
	case gamePuzzleScreen, gameEndlessScreen:
		gameDraw(screen)
	case creditsScreen:
		creditsDraw(screen)
	case winnerScreen:
		winnerDraw(screen)
	}
	//drawDebug(screen)
	drawFooter(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(GameWidth), int(GameHeight)
}

func main() {
	soundInit()
	uiInit()
	switchScreen(0)
	ebiten.SetWindowSize(int(GameWidth)*5, int(GameHeight)*5)
	ebiten.SetWindowTitle("TETRIVERSE")
	ebiten.SetVsyncEnabled(false)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
