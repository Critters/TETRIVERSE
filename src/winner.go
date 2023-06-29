package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var flasherCount int
var flasherFrequency int
var colorCount int

func winnerInit() {
	soundPlay(sfx_winner)
	flasherCount = 0
	colorCount = 0
}

func winnerUpdate() {
	flasherFrequency++
	if flasherFrequency%60 == 0 {
		colorCount++
		currentPallette++
		currentPallette = currentPallette % 5
	}
	if flasherFrequency%30 == 0 {
		flasherCount++
	}
	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	if clicked || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		switchScreen(0)
	}
}

func winnerDraw(screen *ebiten.Image) {
	text.Draw(screen, "TETRI", fontEarlyGameBoyLarge, 3, 32, getColor(3))
	text.Draw(screen, "VERSE", fontEarlyGameBoyLarge, 82, 32, getColor(4))
	text.Draw(screen, "MASTER!", fontEarlyGameBoyLarge, 32, 46, getColor(1+(flasherCount%4)))

	text.Draw(screen, "Now try ENDLESS!", fontTkachevica, 42, 80, getColor(3))
	text.Draw(screen, "[ESC] Back", fontTkachevica, 57, 90, getColor(4))
}
