package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func creditsInit() {

}

func creditsUpdate() {
	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	if clicked || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		switchScreen(0)
	}
}

func creditsDraw(screen *ebiten.Image) {
	text.Draw(screen, "Created by:", fontTkachevica, 6, 16, getColor(1))
	text.Draw(screen, "David Scott", fontTkachevica, 64, 16, getColor(2))

	text.Draw(screen, "Fonts:", fontTkachevica, 6, 30, getColor(1))
	text.Draw(screen, "Jimmy Campbell", fontTkachevica, 64, 30, getColor(2))
	text.Draw(screen, "Tkachevica", fontTkachevica, 64, 38, getColor(2))

	text.Draw(screen, "Tech:", fontTkachevica, 6, 50, getColor(1))
	text.Draw(screen, "Golang / Ebitengine", fontTkachevica, 64, 50, getColor(2))

	text.Draw(screen, "[ESC] Back:", fontTkachevica, 6, 70, getColor(4))
}
