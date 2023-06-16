package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	GameWindowOffsetX int = 32
	GameWindowOffsetY int = 16
	GameWindowWidth   int = GameWidth - 32
	GameWindowHeight  int = GameHeight - 16
)

//go:embed assets/EarlyGameBoy.ttf
var fontEarlyGameBoy_embed []byte
var fontEarlyGameBoy font.Face

type palette []color.RGBA

func getColor(index int) color.RGBA {
	return colorPallettes[currentPallette][index]
}

var (
	currentPallette int = 3
	colorPallettes  []palette
)

func uiInit() {
	colorPallettes = make([]palette, 4)
	colorPallettes[0] = []color.RGBA{
		{51, 44, 80, 255},
		{70, 135, 143, 255},
		{148, 227, 68, 255},
		{226, 243, 228, 255},
	}
	colorPallettes[1] = []color.RGBA{
		{53, 51, 63, 255},
		{218, 52, 103, 255},
		{255, 164, 154, 255},
		{241, 224, 205, 255},
	}
	colorPallettes[2] = []color.RGBA{
		{124, 63, 8, 255},
		{235, 107, 111, 255},
		{249, 168, 117, 255},
		{255, 246, 211, 255},
	}
	colorPallettes[3] = []color.RGBA{
		{15, 15, 27, 255},
		{86, 90, 117, 255},
		{198, 183, 190, 255},
		{250, 251, 246, 255},
	}
	tt, err := opentype.ParseReaderAt(bytes.NewReader(fontEarlyGameBoy_embed))
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	fontEarlyGameBoy, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    8,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func uiUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		currentPallette = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		currentPallette = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		currentPallette = 2
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		currentPallette = 3
	}
}

func drawDebug(screen *ebiten.Image) {
	text.Draw(screen, fmt.Sprint(math.Floor(ebiten.ActualFPS())), fontEarlyGameBoy, 1, 8, getColor(1))
	vector.DrawFilledRect(screen, 46, 1, 4, 4, getColor(0), false)
	vector.DrawFilledRect(screen, 50, 1, 4, 4, getColor(1), false)
	vector.DrawFilledRect(screen, 54, 1, 4, 4, getColor(2), false)
	vector.DrawFilledRect(screen, 58, 1, 4, 4, getColor(3), false)
}

func drawBackground(screen *ebiten.Image) {
	// Game Frame
	DrawOutlinedRect(screen, float32(GameWindowOffsetX), float32(GameWindowOffsetY), float32(GameWindowWidth), float32(GameWindowHeight), getColor(3), getColor(0))
}

func DrawOutlinedRect(screen *ebiten.Image, x float32, y float32, width float32, height float32, border color.RGBA, fill color.RGBA) {
	vector.DrawFilledRect(screen, x, y, width, height, border, false)
	vector.DrawFilledRect(screen, x+1, y+1, width-2, height-2, fill, false)
}
