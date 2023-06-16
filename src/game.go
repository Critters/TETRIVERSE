package main

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var boardMatrix [200]int
var boardImage *ebiten.Image = ebiten.NewImage(10*7, 20*7)
var boardShapes []int

type shape [9]int
type shapes [4]shape

var possibleShapes []shapes
var currentShape int = 1

var shapeX = 0
var shapeY = 0
var shapeRotation = 0

//go:embed assets/levels.json
var levelFile []byte

type JsonFile struct {
	Levels []Level `json:"levels"`
}
type Level struct {
	Start  [200]int `json:"start"`
	Shapes []int    `json:"shapes"`
}

func gameInit() {
	boardMatrix = [200]int{}
	possibleShapes = make([]shapes, 6)
	possibleShapes[0] = [4]shape{
		{1, 1, 0, 1, 1, 0, 0, 0, 0},
		{1, 1, 0, 1, 1, 0, 0, 0, 0},
		{1, 1, 0, 1, 1, 0, 0, 0, 0},
		{1, 1, 0, 1, 1, 0, 0, 0, 0}} // Square
	possibleShapes[1] = [4]shape{
		{0, 0, 0, 1, 1, 1, 0, 0, 1},
		{0, 1, 0, 0, 1, 0, 1, 1, 0},
		{1, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 0, 1, 0, 0, 1, 0}} // L
	possibleShapes[2] = [4]shape{
		{0, 0, 0, 1, 1, 1, 1, 0, 0},
		{1, 1, 0, 0, 1, 0, 0, 1, 0},
		{0, 0, 1, 1, 1, 1, 0, 0, 0},
		{0, 1, 0, 0, 1, 0, 0, 1, 1}} // J
	possibleShapes[3] = [4]shape{
		{0, 0, 0, 0, 1, 1, 1, 1, 0},
		{0, 1, 0, 0, 1, 1, 0, 0, 1},
		{0, 0, 0, 0, 1, 1, 1, 1, 0},
		{0, 1, 0, 0, 1, 1, 0, 0, 1}} // Z
	possibleShapes[4] = [4]shape{
		{0, 0, 0, 1, 1, 1, 0, 1, 0},
		{0, 1, 0, 1, 1, 0, 0, 1, 0},
		{0, 1, 0, 1, 1, 1, 0, 0, 0},
		{0, 1, 0, 0, 1, 1, 0, 1, 0}}
	possibleShapes[5] = [4]shape{
		{0, 0, 0, 1, 1, 0, 0, 1, 1},
		{0, 0, 1, 0, 1, 1, 0, 1, 0},
		{0, 0, 0, 1, 1, 0, 0, 1, 1},
		{0, 0, 1, 0, 1, 1, 0, 1, 0}}

	loadLevel(0)
}

func loadLevel(level int) {
	var jsonFile JsonFile
	json.Unmarshal(levelFile, &jsonFile)
	boardMatrix = jsonFile.Levels[level].Start
	boardShapes = jsonFile.Levels[level].Shapes
	//fmt.Print(boardMatrix, boardShapes)
}

func gameUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		shapeRotation++
		shapeRotation = shapeRotation % 4
		fmt.Print(shapeRotation)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		shapeRotation--
		if shapeRotation < 0 {
			shapeRotation = 3
		}
		fmt.Print(shapeRotation)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		currentShape++
		currentShape = currentShape % 5
	}
}

func gameDraw(screen *ebiten.Image) {
	drawOutlinedRect(screen, 63, 0, (10*7)+2, (20*7)+2, getColor(3), getColor(0))
	for i := 0; i < 200; i++ {
		vector.DrawFilledRect(boardImage, float32(i%10)*7, float32((i/10)%20)*7, 6, 6, getColor(boardMatrix[i]), false)
	}
	var shape = possibleShapes[currentShape][shapeRotation]
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if shape[x+((y%3)*3)] == 1 {
				vector.DrawFilledRect(boardImage, float32(x)*7, float32(y)*7, 6, 6, getColor(3), false)
			}
		}
	}
	var dio *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	dio.GeoM.Translate(64, 1)
	screen.DrawImage(boardImage, dio)
}
