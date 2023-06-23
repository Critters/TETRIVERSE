/*
Contains some common UI elements used by the menu + game
*/

package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"math"
	"strings"

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

//go:embed assets/Tkachevica-4pxRegular.ttf
var fontTkachevica_embed []byte
var fontTkachevica font.Face

type palette []color.NRGBA

var currentPallette int = 1
var colorPallettes []palette

func getColor(index int) color.NRGBA {
	col := colorPallettes[currentPallette][index]
	return col
}

func uiInit() {
	colorPallettes = make([]palette, 5)
	colorPallettes[0] = []color.NRGBA{
		{51, 44, 80, 255},
		{70, 135, 143, 255},
		{148, 227, 68, 255},
		{226, 243, 228, 255},
		{218, 52, 103, 255},
	}
	colorPallettes[1] = []color.NRGBA{
		{15, 15, 27, 255},
		{86, 90, 117, 255},
		{198, 183, 190, 255},
		{250, 251, 246, 255},
		{218, 52, 103, 255},
	}
	colorPallettes[2] = []color.NRGBA{
		{0, 0, 0, 255},
		{108, 115, 82, 255},
		{142, 153, 111, 255},
		{196, 208, 160, 255},
		{218, 52, 103, 255},
	}
	colorPallettes[3] = []color.NRGBA{
		{33, 11, 27, 255},
		{77, 34, 44, 255},
		{157, 101, 76, 255},
		{207, 171, 81, 255},
		{218, 52, 103, 255},
	}
	colorPallettes[4] = []color.NRGBA{
		{22, 11, 27, 255},
		{56, 40, 67, 255},
		{124, 109, 128, 255},
		{199, 198, 198, 255},
		{218, 52, 103, 255},
	}

	const dpi = 72

	tt, err := opentype.ParseReaderAt(bytes.NewReader(fontEarlyGameBoy_embed))
	if err != nil {
		log.Fatal(err)
	}
	fontEarlyGameBoy, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    8,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	tt, err = opentype.ParseReaderAt(bytes.NewReader(fontTkachevica_embed))
	if err != nil {
		log.Fatal(err)
	}
	fontTkachevica, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    8,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	menuInit()
}

func uiUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		currentPallette = 0
		redrawBoardImage = true
		redrawUpcomingShapesImage = true
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		currentPallette = 1
		redrawBoardImage = true
		redrawUpcomingShapesImage = true
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		currentPallette = 2
		redrawBoardImage = true
		redrawUpcomingShapesImage = true
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		currentPallette = 3
		redrawBoardImage = true
		redrawUpcomingShapesImage = true
	}
	if inpututil.IsKeyJustPressed(ebiten.Key5) {
		currentPallette = 4
		redrawBoardImage = true
		redrawUpcomingShapesImage = true
	}
}

func drawDebug(screen *ebiten.Image) {
	text.Draw(screen, fmt.Sprint(math.Floor(ebiten.ActualFPS()), "fps"), fontTkachevica, 1, 8, getColor(1))
	var cursorX, cursorY = ebiten.CursorPosition()
	text.Draw(screen, fmt.Sprint(cursorX, ", ", cursorY), fontTkachevica, 1, 16, getColor(1))

	vector.DrawFilledRect(screen, 0, 145, 160, 8, color.NRGBA{0, 0, 0, 255}, false)
	vector.DrawFilledRect(screen, 1, 146, 4, 5, getColor(0), false)
	vector.DrawFilledRect(screen, 5, 146, 4, 5, getColor(1), false)
	vector.DrawFilledRect(screen, 9, 146, 4, 5, getColor(2), false)
	vector.DrawFilledRect(screen, 13, 146, 4, 5, getColor(3), false)
	vector.DrawFilledRect(screen, 17, 146, 4, 5, getColor(4), false)
	text.Draw(screen, "WASD, Q-E, SPACE, [R]eset", fontTkachevica, 22, 151, color.NRGBA{128, 128, 128, 255})
}

/*
func drawBackground(screen *ebiten.Image) {
	// Game Frame
	drawOutlinedRect(screen, float32(GameWindowOffsetX), float32(GameWindowOffsetY), float32(GameWindowWidth), float32(GameWindowHeight), getColor(3), getColor(0))
}
*/

func drawOutlinedRect(screen *ebiten.Image, x float32, y float32, width float32, height float32, border color.NRGBA, fill color.NRGBA) {
	vector.DrawFilledRect(screen, x, y, width, height, border, false)
	vector.DrawFilledRect(screen, x+1, y+1, width-2, height-2, fill, false)
}

func drawBlock(screen *ebiten.Image, x float32, y float32, state int) {
	// 0: inactive 2: highlighted 3: blocking
	vector.DrawFilledRect(screen, x, y, 7, 7, getColor(state), false)
}

func lerp(a float32, b float32, f float32) float32 {
	return a*(1.0-f) + (b * f)
}

func drawCenteredText(screen *ebiten.Image, txt string, centerX int, topY int, col color.NRGBA) {
	words := strings.Split(txt, ",")
	//maxWidth, totalHeight := 0, 0
	// Find the height of all the words, and the max width
	//for i := range words {
	//	rect := text.BoundString(fontEarlyGameBoy, words[i])
	//	totalHeight += rect.Dy()
	//	if rect.Dx() > maxWidth {
	//		maxWidth = rect.Dx()
	//	}
	//}
	// Draw the words
	for i := range words {
		rect := text.BoundString(fontEarlyGameBoy, words[i])
		text.Draw(screen, words[i], fontEarlyGameBoy, centerX-int(rect.Dx()/2), topY+rect.Dy()+(i*(rect.Dy()+2)), col)
	}

}
