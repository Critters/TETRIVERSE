/*
Contains the logic and rendering for the main menu
*/

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func creditsInit() {

}

func creditsUpdate() {
	// Mouse
	//var cursorX, cursorY = ebiten.CursorPosition()
	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	//if cursorX > 0 && cursorX < 16 {
	//	if cursorY > 0 && cursorY < 16 {
	if clicked {
		switchScreen(0)
	}
	//	}
	//}
	// TODO: Add keyboard support
}

func creditsDraw(screen *ebiten.Image) {
	text.Draw(screen, "Created by:", fontTkachevica, 6, 16, getColor(1))
	text.Draw(screen, "David Scott", fontTkachevica, 64, 16, getColor(2))

	text.Draw(screen, "Fonts:", fontTkachevica, 6, 30, getColor(1))
	text.Draw(screen, "Jimmy Campbell", fontTkachevica, 64, 30, getColor(2))
	text.Draw(screen, "Tkachevica", fontTkachevica, 64, 38, getColor(2))

	text.Draw(screen, "Tech:", fontTkachevica, 6, 50, getColor(1))
	text.Draw(screen, "Golang / Ebitengine", fontTkachevica, 64, 50, getColor(2))
}
