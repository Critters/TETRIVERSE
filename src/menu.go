/*
Contains the logic and rendering for the main menu
*/

package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var highlighted int = 0
var lastHighlighted int = -1

var uiVisualImage *ebiten.Image
var uiVisualMatrix [200]int

func menuInit() {
	uiVisualImage = ebiten.NewImage(128, 64)
	uiVisualMatrix = [200]int{
		1, 0, 1, 0, 0, 0,
		1, 1, 1, 1, 1, 0,
		1, 1, 1, 1, 1, 1}
	//logicInit(uiVisualMatrix)
}

// Squares in the background thing
var menuBackground []float32 = make([]float32, 3040)

// W & S change which menu item is highlighted
func menuUpdate() {
	var cursorX, cursorY = ebiten.CursorPosition()

	// Menu Background
	v := int(cursorX/9) + ((cursorY / 9) * 160)
	if v < 0 {
		v = 0
	} else if v >= len(menuBackground) {
		v = len(menuBackground) - 1
	}
	menuBackground[v] = 0.25
	for i := range menuBackground {
		if menuBackground[i] > 0 {
			menuBackground[i] -= 0.01
		}
	}

	// Keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		highlighted--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		highlighted++
	}
	// Mouse
	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	if cursorX > 54 && cursorX < 110 {
		if cursorY > 52 && cursorY < 62 {
			highlighted = 0
			if clicked {
				switchScreen(1)
				return
			}
		}
		if cursorY > 61 && cursorY < 74 {
			highlighted = 1
			if clicked {
				switchScreen(2)
				return
			}
		}
		if cursorY > 73 && cursorY < 86 {
			highlighted = 2
			if clicked {
				switchScreen(3)
				return
			}
		}
	}
	// Wraps around when pressing W on the top menu item
	if highlighted < 0 {
		highlighted = 2
	}
	// Wraps around when pressing S on the bottom menu item
	highlighted = highlighted % 3

	if highlighted != lastHighlighted {
		lastHighlighted = highlighted
		soundPlay(sfx_menu)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switchScreen(highlighted + 1)
	}
}

func menuDraw(screen *ebiten.Image) {
	// Menu Background
	for i := range menuBackground {
		if menuBackground[i] > 0 {
			colA := getColor(0)
			colB := getColor(1)
			colC := color.NRGBA{
				uint8(lerp(float32(colA.R), float32(colB.R), menuBackground[i])),
				uint8(lerp(float32(colA.G), float32(colB.G), menuBackground[i])),
				uint8(lerp(float32(colA.B), float32(colB.B), menuBackground[i])), 255}
			vector.DrawFilledRect(screen, float32((i%160)*9), float32(((i/160)%160)*9), 7, 7, colC, false)
		}
	}

	text.Draw(screen, "TETRI", fontEarlyGameBoyLarge, 3, 32, getColor(3))
	text.Draw(screen, "VERSE", fontEarlyGameBoyLarge, 82, 32, getColor(4))

	// The four options: Puzzle, Arcade, Options, Credits
	var col color.NRGBA = getColor(1)
	var dotX int = 60
	var dotY int = 60
	if highlighted == 0 {
		col = getColor(3)
		uiVisualMatrix = [200]int{
			0, 0, 1, 0, 0, 0,
			3, 3, 1, 1, 1, 0,
			3, 3, 1, 1, 1, 1}
	}
	text.Draw(screen, "PUZZLE", fontEarlyGameBoy, 60, 60, col)

	col = getColor(1)
	if highlighted == 1 {
		dotX = 56
		dotY = 72
		col = getColor(3)
		uiVisualMatrix = [200]int{
			0, 0, 3, 0, 0, 0,
			1, 1, 3, 1, 1, 0,
			1, 1, 3, 3, 1, 1}
	}
	text.Draw(screen, "ENDLESS", fontEarlyGameBoy, 56, 72, col)

	/*
		col = getColor(1)
		if highlighted == 2 {
			dotX = 68
			dotY = 84
			col = getColor(3)
			uiVisualMatrix = [200]int{
				0, 0, 3, 0, 0, 0,
				1, 3, 3, 3, 1, 0,
				1, 1, 1, 1, 1, 1}
		}
		text.Draw(screen, "HELP", fontEarlyGameBoy, 68, 84, col)
	*/

	col = getColor(1)
	if highlighted == 2 {
		dotX = 56
		dotY = 84 //96
		col = getColor(3)
		uiVisualMatrix = [200]int{
			0, 0, 1, 0, 0, 0,
			1, 1, 1, 3, 3, 0,
			1, 1, 1, 1, 3, 3}
	}
	text.Draw(screen, "CREDITS", fontEarlyGameBoy, 56, 84, col)

	// The dot
	vector.DrawFilledCircle(screen, float32(dotX-8), float32(dotY-3000), 4, getColor(3), false)

	// The visualization
	var visX int = 54
	var visY int = 108
	var dio *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	dio.GeoM.Translate(float64(visX), float64(visY))

	for i := 0; i < 18; i++ {
		if uiVisualMatrix[i] > 0 {
			drawBlock(uiVisualImage, float32((i%6)*9), float32(((i/6)%6)*9), uiVisualMatrix[i])
		}
	}

	screen.DrawImage(uiVisualImage, dio)
}
