package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var highlighted int = 0
var uiVisualImage *ebiten.Image
var uiVisualMatrix [18]int

func uiMenuInit() {
	uiVisualImage = ebiten.NewImage(128, 64)
	uiVisualMatrix = [18]int{
		1, 0, 1, 0, 0, 0,
		1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 0}
}

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
	// The three options: Play, Options, Credits
	var col color.RGBA = getColor(1)
	var dotX int = 60
	var dotY int = 64
	if highlighted == 0 {
		col = getColor(3)
		uiVisualMatrix = [18]int{
			2, 0, 1, 0, 0, 0,
			2, 1, 1, 1, 1, 0,
			2, 2, 1, 1, 1, 0}
	}
	text.Draw(screen, "PLAY", fontEarlyGameBoy, 60, 64, col)
	col = getColor(1)
	if highlighted == 1 {
		dotX = 48
		dotY = 80
		col = getColor(3)
		uiVisualMatrix = [18]int{
			1, 0, 2, 0, 0, 0,
			1, 2, 2, 1, 1, 0,
			1, 1, 2, 1, 1, 0}
	}
	text.Draw(screen, "OPTIONS", fontEarlyGameBoy, 48, 80, col)
	col = getColor(1)
	if highlighted == 2 {
		dotX = 48
		dotY = 96
		col = getColor(3)
		uiVisualMatrix = [18]int{
			1, 0, 1, 0, 0, 0,
			1, 1, 1, 2, 2, 0,
			1, 1, 1, 2, 2, 0}
	}
	text.Draw(screen, "CREDITS", fontEarlyGameBoy, 48, 96, col)

	// The dot
	vector.DrawFilledCircle(screen, float32(dotX-6), float32(dotY-300), 4, getColor(1), false)

	// The visualization
	var visX int = 48
	var visY int = 112
	var dio *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	dio.GeoM.Translate(float64(visX), float64(visY))

	for i := 0; i < 18; i++ {
		drawBlock(uiVisualImage, float32((i%6)*9), float32(((i/6)%6)*9), uiVisualMatrix[i])
	}

	screen.DrawImage(uiVisualImage, dio)
}
