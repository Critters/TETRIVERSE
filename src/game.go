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
var boardImage *ebiten.Image = ebiten.NewImage(70, 140)
var boardShapes []int
var t int = 0

type shape [9]int
type shapes [4]shape

var possibleShapes []shapes
var currentShape int = 1

var shapeX = 4
var shapeY = 16
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
		{0, 1, 0, 0, 1, 1, 0, 1, 0}} // T
	possibleShapes[5] = [4]shape{
		{0, 0, 0, 1, 1, 0, 0, 1, 1},
		{0, 0, 1, 0, 1, 1, 0, 1, 0},
		{0, 0, 0, 1, 1, 0, 0, 1, 1},
		{0, 0, 1, 0, 1, 1, 0, 1, 0}} // Other Z

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
	t++
	if t%30 == 0 {
		raiseShapes()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) && shapeY > 0 {
		shapeY--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) && shapeY < 17 {
		shapeY++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) && shapeX > 0 {
		shapeX--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) && shapeX < 7 {
		shapeX++
	}
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
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		extractShape(shapeX, shapeY, currentShape, shapeRotation)
	}
}

func gameDraw(screen *ebiten.Image) {
	// Frame
	drawOutlinedRect(screen, 63, 0, 71, 141, getColor(3), getColor(0))
	// Board
	boardImage.Clear()
	for i := 0; i < 200; i++ {
		vector.DrawFilledRect(boardImage, float32(i%10)*7, float32((i/10)%20)*7, 6, 6, getColor(boardMatrix[i]), false)
	}
	// Shape
	drawShape(boardImage, shapeX, shapeY, currentShape, shapeRotation, true)
	// Compile
	var dio *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	dio.GeoM.Translate(64, 1)
	screen.DrawImage(boardImage, dio)
}

func drawShape(screen *ebiten.Image, posX int, posY int, shapeID int, shapeRotation int, checkValid bool) bool {
	var shape = possibleShapes[shapeID][shapeRotation]
	var isValid = true
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if shape[x+((y%3)*3)] == 1 {
				var col = getColor(3)
				if boardMatrix[(posX+x)+((posY+y)*10)] == 0 {
					isValid = false
					col = getColor(4)
					// = uint8((0.5 + (math.Sin(t) / 2)) * 255)
					//fmt.Println(col.A)
				}
				vector.DrawFilledRect(screen, float32(posX+x)*7, float32(posY+y)*7, 6, 6, col, false)
			}
		}
	}
	return isValid
}

func checkShape(posX int, posY int, shapeID int, shapeRotation int) bool {
	var shape = possibleShapes[shapeID][shapeRotation]
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if shape[x+((y%3)*3)] == 1 {
				if boardMatrix[(posX+x)+((posY+y)*10)] == 0 {
					return false
				}
			}
		}
	}
	return true
}

func extractShape(posX int, posY int, shapeID int, shapeRotation int) {
	var clear = checkShape(posX, posY, shapeID, shapeRotation)
	if clear {
		var shape = possibleShapes[shapeID][shapeRotation]
		for x := 0; x < 3; x++ {
			for y := 0; y < 3; y++ {
				if shape[x+((y%3)*3)] == 1 {
					boardMatrix[(posX+x)+((posY+y)*10)] = 2
				}
			}
		}
	}
}

func raiseShapes() {
	for i := 0; i < 200; i++ {
		if boardMatrix[i] == 2 {
			boardMatrix[i] = 0
			if i-10 > 0 {
				boardMatrix[i-10] = 2
			}
		}
	}
}
