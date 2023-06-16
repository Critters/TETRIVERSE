package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var highlighted int = 0

// W & S change which menu item is highlighted
func uiMenuUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		highlighted--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		highlighted++
	}
	// Wraps around when pressing W on the top menu item
	if highlighted < 0 {
		highlighted = 2
	}
	// Wraps around when pressing S on the bottom menu item
	highlighted = highlighted % 3
}

func uiMenuDraw(screen *ebiten.Image) {
	var col color.RGBA = getColor(1)
	var dotX int = 60
	var dotY int = 64
	if highlighted == 0 {
		col = getColor(3)
	}
	text.Draw(screen, "PLAY", fontEarlyGameBoy, 60, 64, col)
	col = getColor(1)
	if highlighted == 1 {
		dotX = 48
		dotY = 80
		col = getColor(3)
	}
	text.Draw(screen, "OPTIONS", fontEarlyGameBoy, 48, 80, col)
	col = getColor(1)
	if highlighted == 2 {
		dotX = 48
		dotY = 96
		col = getColor(3)
	}
	text.Draw(screen, "CREDITS", fontEarlyGameBoy, 48, 96, col)

	vector.DrawFilledCircle(screen, float32(dotX-6), float32(dotY-3), 4, getColor(1), false)

}
