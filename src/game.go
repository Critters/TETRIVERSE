/*
Contains the logic and rendering for the game
*/

package main

import (
	_ "embed"
	"encoding/json"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// The 20x10 board
var boardMatrix [200]int

// This is the image we render the game to, this is then copied to the screen
var boardImage *ebiten.Image = ebiten.NewImage(70, 140)

/*
This is an over optimization for this game, but redrawBoardImage is only set
to true when the boardImage needs to be redrawn, which in this game is very
infrequent. This only gets set to true when the player moves or rotates the shape
or when blocks move up after getting selected, whcih is twice a second. On a high
end PC this optimization increased the average FPS from 150 to 190 (26% improvement)
*/
var redrawBoardImage bool = true

var boardShapes []int
var t int = 0

// The 3x3 matrix of a shape
type shape [9]int

// The 4 rotations of a shape
type shapes [4]shape

// All rotations of all shapes
var possibleShapes []shapes

// The shape currently selected, it's position on the board, and its rotation
var currentShape int = 4
var shapeX = 7
var shapeY = 17
var shapeRotation = 2

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
		{1, 1, 0, 1, 0, 0, 1, 0, 0}} // L
	possibleShapes[2] = [4]shape{
		{0, 0, 0, 1, 1, 1, 1, 0, 0},
		{1, 1, 0, 0, 1, 0, 0, 1, 0},
		{0, 0, 1, 1, 1, 1, 0, 0, 0},
		{1, 0, 0, 1, 0, 0, 1, 1, 0}} // J
	possibleShapes[3] = [4]shape{
		{0, 0, 0, 0, 1, 1, 1, 1, 0},
		{1, 0, 0, 1, 1, 0, 0, 1, 0},
		{0, 0, 0, 0, 1, 1, 1, 1, 0},
		{1, 0, 0, 1, 1, 0, 0, 1, 0}} // Z
	possibleShapes[4] = [4]shape{
		{0, 0, 0, 1, 1, 1, 0, 1, 0},
		{0, 1, 0, 1, 1, 0, 0, 1, 0},
		{0, 1, 0, 1, 1, 1, 0, 0, 0},
		{1, 0, 0, 1, 1, 0, 1, 0, 0}} // T
	possibleShapes[5] = [4]shape{
		{0, 0, 0, 1, 1, 0, 0, 1, 1},
		{0, 1, 0, 1, 1, 0, 1, 0, 0},
		{0, 0, 0, 1, 1, 0, 0, 1, 1},
		{0, 1, 0, 1, 1, 0, 1, 0, 0}} // Other Z

	loadLevel(0)
}

func loadLevel(level int) {
	var jsonFile JsonFile
	json.Unmarshal(levelFile, &jsonFile)
	boardMatrix = jsonFile.Levels[level].Start
	boardShapes = jsonFile.Levels[level].Shapes
}

var oldX, oldY int

func gameUpdate() {
	t++
	if t%30 == 0 {
		raiseShapes()
	}

	// Mouse
	var cursorX, cursorY = ebiten.CursorPosition()
	var newShapeX = ((cursorX - 64) / 7) - 1
	var newShapeY = (cursorY / 7) - 1
	if newShapeX < 0 {
		newShapeX = 0
	}
	if newShapeY < 1 {
		newShapeY = 1
	}
	_, dx := ebiten.Wheel()
	if dx != 0 {
		shapeRotation += int(-dx)
		if shapeRotation < 0 {
			shapeRotation = 3
		}
		shapeRotation = shapeRotation % 4
		redrawBoardImage = true
	}
	if newShapeX != oldX || newShapeY != oldY {
		oldX = newShapeX
		oldY = newShapeY
		shapeX = newShapeX
		shapeY = newShapeY
		redrawBoardImage = true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		extractShape(shapeX, shapeY, currentShape, shapeRotation)
	}

	// Keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeyW) && shapeY > 0 {
		shapeY--
		redrawBoardImage = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		shapeY++
		redrawBoardImage = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) && shapeX > 0 {
		shapeX--
		redrawBoardImage = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		shapeX++
		redrawBoardImage = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		shapeRotation++
		shapeRotation = shapeRotation % 4
		redrawBoardImage = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		shapeRotation--
		if shapeRotation < 0 {
			shapeRotation = 3
		}
		redrawBoardImage = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		currentShape++
		currentShape = currentShape % 5
		redrawBoardImage = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		extractShape(shapeX, shapeY, currentShape, shapeRotation)
		redrawBoardImage = true
	}

	// Ensure the shape does not go off the right or bottom edges of the board
	var shape = possibleShapes[currentShape][shapeRotation]
	// Right
	if shapeX >= 8 {
		if shape[2] == 0 && shape[5] == 0 && shape[8] == 0 {
			shapeX = 8
		} else {
			shapeX = 7
		}
	}
	// Bottom
	if shapeY >= 18 {
		if shape[6] == 0 && shape[7] == 0 && shape[8] == 0 {
			shapeY = 18
		} else {
			shapeY = 17
		}
	}

	// Something changed, check shape again
	if redrawBoardImage {
		checkShape(shapeX, shapeY, currentShape, shapeRotation)
	}
}

func gameDraw(screen *ebiten.Image) {
	// Frame
	drawOutlinedRect(screen, 63, 0, 71, 141, getColor(3), getColor(0))
	if redrawBoardImage {
		vector.DrawFilledRect(screen, 0, 0, 6, 6, getColor(4), false)
		// Board
		boardImage.Clear()
		for i := 0; i < 200; i++ {
			vector.DrawFilledRect(boardImage, float32(i%10)*7, float32((i/10)%20)*7, 6, 6, getColor(boardMatrix[i]), false)
		}
		// Shape
		drawShape(boardImage, shapeX, shapeY, currentShape, shapeRotation, true)
		redrawBoardImage = false
	}
	// Compile
	var dio *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	dio.GeoM.Translate(64, 1)
	screen.DrawImage(boardImage, dio)

}

func drawShape(screen *ebiten.Image, posX int, posY int, shapeID int, shapeRotation int, checkValid bool) {
	var shape = possibleShapes[shapeID][shapeRotation]
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if shape[x+((y%3)*3)] == 1 {
				// White
				var col = getColor(3)
				if (posX+x)+((posY+y)*10) < 199 && boardMatrix[(posX+x)+((posY+y)*10)] == 0 {
					// Red
					col = getColor(4)
				}
				vector.DrawFilledRect(screen, float32(posX+x)*7, float32(posY+y)*7, 6, 6, col, false)
			}
		}
	}
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

// Checks if the shape can be extracted
func checkShape(posX int, posY int, shapeID int, shapeRotation int) bool {
	// Reset board to white
	for i := 0; i < 200; i++ {
		if boardMatrix[i] == 4 {
			boardMatrix[i] = 1
		}
	}

	var shape = possibleShapes[shapeID][shapeRotation]
	var extractable = true
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if shape[x+((y%3)*3)] == 1 {
				// A part of the shape is not on a block
				if (posX+x)+((posY+y)*10) < 199 && boardMatrix[(posX+x)+((posY+y)*10)] == 0 {
					extractable = false
				}
				// Any block above the shape would prevent it from been able to be extracted
				// Top row (special case)
				if y == 0 {
					if boardMatrix[(posX+x)+((posY+y-1)*10)] == 1 {
						boardMatrix[(posX+x)+((posY+y-1)*10)] = 4
						extractable = false
					}
				}
				// Second and third row has to exclude blocks covered by the first row
				if y == 1 || y == 2 {
					if boardMatrix[(posX+x)+((posY+y-1)*10)] == 1 && shape[x+(((y%3)-1)*3)] == 0 {
						boardMatrix[(posX+x)+((posY+y-1)*10)] = 4
						extractable = false
					}
				}
			}
		}
	}
	return extractable
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
	redrawBoardImage = true
}
