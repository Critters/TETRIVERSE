/*
Contains the logic and rendering for the game
*/

package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var t int = 0
var gameState int = 0 // 0:Normal 1:Level complete (FF)
var currentLevel int = 3
var shake float32 = 0
var shakeX = 0

// The 20x10 board
type boardBlock struct {
	state int //0:Empty 1:Filled 2:Rising
	color color.NRGBA
	fade  float32 // Fades from 1 to 0 to control how much red to show
}

var boardMatrix [200]boardBlock

// This is the image we render the game to, this is then copied to the screen
var boardImage *ebiten.Image = ebiten.NewImage(70, 140)

/*
This is an over optimization for this game, but redrawBoardImage is only set
to true when the boardImage needs to be redrawn, which in this game is very
infrequent. This only gets set to true when the player moves or rotates the shape
or when blocks move up after getting selected, which is twice a second. On a high
end PC this optimization increased the average FPS from 150 to 190 (26% improvement)
*/
var redrawBoardImage bool = true

var upcomingShapes []int
var upcomingShapesImage *ebiten.Image = ebiten.NewImage(32, 128)
var redrawUpcomingShapesImage bool = true
var levelHint string

// The 3x3 matrix of a shape
// The 4 rotations of a shape
// All rotations of all shapes
type shape [9]int
type shapes [4]shape

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
	Hint   string   `json:"hint"`
	Start  [200]int `json:"start"`
	Shapes []int    `json:"shapes"`
}

func gameInit() {
	fmt.Println("gameInit()")
	currentLevel = 7
	boardMatrix = [200]boardBlock{}
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

	LoadLevel(currentLevel)
}

func LoadLevel(level int) {
	fmt.Println("LoadLevel(", level, ")")
	currentLevel = level
	var jsonFile JsonFile
	json.Unmarshal(levelFile, &jsonFile)
	tmpBoard := jsonFile.Levels[currentLevel].Start
	for i := range tmpBoard {
		tmpCol := getColor(0) // Empty
		if tmpBoard[i] == 1 {
			tmpCol = getColor(1) // Filled
		}
		boardMatrix[i].state = tmpBoard[i]
		boardMatrix[i].color = tmpCol
	}
	upcomingShapes = jsonFile.Levels[currentLevel].Shapes
	levelHint = jsonFile.Levels[currentLevel].Hint
	PopShape()
}

// Removes the top shape from the list and makes it the current
func PopShape() {
	if len(upcomingShapes) > 0 {
		currentShape = upcomingShapes[0]
		shapeRotation = 0
		upcomingShapes = upcomingShapes[1:]
		redrawUpcomingShapesImage = true
		stillPossibleResult = stillPossible()
	} else {
		// No more shapes left
		currentShape = -1
		redrawUpcomingShapesImage = true
		gameState = 1
	}
}

func NextLevel() {
	LoadLevel(currentLevel + 1)
	gameState = 0
}

var oldX, oldY int

func gameUpdate() {
	t++

	if shake > 0 {
		shake -= 0.0166
		shakeX = rand.Intn(4) - 2
	} else {
		shakeX = 0
	}

	if gameState == 0 {
		if t%30 == 0 {
			raiseShapes()
		}
	} else if gameState == 1 {
		if t%5 == 0 {
			raisedSomething := raiseShapes()
			if !raisedSomething {
				NextLevel()
			}
		}
	}

	// Fade red to white
	for i := range boardMatrix {
		if boardMatrix[i].fade > 0 {
			colorA := getColor(1)
			colorB := getColor(4)
			f := boardMatrix[i].fade
			boardMatrix[i].color = color.NRGBA{
				uint8(lerp(float32(colorA.R), float32(colorB.R), f)),
				uint8(lerp(float32(colorA.G), float32(colorB.G), f)),
				uint8(lerp(float32(colorA.B), float32(colorB.B), f)),
				255,
			}
			boardMatrix[i].fade -= 0.025
			if boardMatrix[i].fade <= 0 {
				boardMatrix[i].color = getColor(1)
				boardMatrix[i].fade = 0
			}
			redrawBoardImage = true
		}
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
		redrawBoardImage = true
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
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		stillPossibleResult = stillPossible()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		LoadLevel(currentLevel)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		extractShape(shapeX, shapeY, currentShape, shapeRotation)
		redrawBoardImage = true
	}

	// Ensure the shape does not go off the right or bottom edges of the board
	if currentShape > -1 {
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
	}

	// Something changed, check shape again
	if redrawBoardImage {
		checkShape(shapeX, shapeY, currentShape, shapeRotation, false)
	}
}

func gameDraw(screen *ebiten.Image) {
	// Frames
	drawOutlinedRect(screen, 63+float32(shakeX), 2, 71, 141, getColor(3), getColor(0))
	if stillPossibleResult {
		drawCenteredText(screen, levelHint, 63+35+shakeX, 16, getColor(1))
	} else {
		drawCenteredText(screen, "Press R,to reset", 63+35+shakeX, 16, getColor(4))
	}

	drawOutlinedRect(screen, 23, 2, 34, 32, getColor(3), getColor(0))
	drawOutlinedRect(screen, 23, 33, 34, 110, getColor(3), getColor(0))
	//text.Draw(screen, "NEXT", fontEarlyGameBoy, 8, 31, getColor(1))

	if redrawBoardImage {
		vector.DrawFilledRect(screen, 0, 0, 6, 6, getColor(4), false)
		// Board
		boardImage.Clear()
		for i := 0; i < 200; i++ {
			if boardMatrix[i].state > 0 {
				vector.DrawFilledRect(boardImage, float32(i%10)*7, float32((i/10)%20)*7, 6, 6, boardMatrix[i].color, false)
			}
		}
		// Shape
		drawShape(boardImage, shapeX, shapeY, 1, currentShape, shapeRotation, true, 0, 0)
		redrawBoardImage = false

	}

	if redrawUpcomingShapesImage {
		upcomingShapesImage.Clear()
		for i := range upcomingShapes {
			if i == 0 {
				var offsetX, offsetY float32 = 0, 0
				switch upcomingShapes[i] {
				case 0:
					offsetX = 2
					offsetY = 2
				case 3, 5, 4, 1, 2:
					offsetY = -6

				}
				drawShape(upcomingShapesImage, 1, 1+(i*4), 1, upcomingShapes[i], shapeRotation, false, offsetX, offsetY)
			} else {
				drawShape(upcomingShapesImage, 2, 3+(i*4), 0.66, upcomingShapes[i], shapeRotation, false, 0, 0)
			}
		}
		redrawUpcomingShapesImage = false
	}

	// Compile
	var dio *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	dio.GeoM.Translate(24, 2)
	screen.DrawImage(upcomingShapesImage, dio)

	dio = &ebiten.DrawImageOptions{}
	dio.GeoM.Translate(64+float64(shakeX), 3)
	screen.DrawImage(boardImage, dio)
}

func drawShape(screen *ebiten.Image, posX int, posY int, scale float32, shapeID int, shapeRotation int, checkValid bool, offsetX float32, offsetY float32) {
	if shapeID == -1 {
		return
	}
	var shape = possibleShapes[shapeID][shapeRotation]
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if shape[x+((y%3)*3)] == 1 {
				// White
				var col = getColor(3)
				if checkValid && (posX+x)+((posY+y)*10) < 199 && boardMatrix[(posX+x)+((posY+y)*10)].state != 1 {
					// Red
					col = getColor(4)
				}
				vector.DrawFilledRect(screen, offsetX+float32(posX+x)*(7*scale), offsetY+float32(posY+y)*(7*scale), (6 * scale), (6 * scale), col, false)
			}
		}
	}
}

func extractShape(posX int, posY int, shapeID int, shapeRotation int) {
	var clear = checkShape(posX, posY, shapeID, shapeRotation, true)
	if clear {
		var shape = possibleShapes[shapeID][shapeRotation]
		for x := 0; x < 3; x++ {
			for y := 0; y < 3; y++ {
				if shape[x+((y%3)*3)] == 1 {
					boardMatrix[(posX+x)+((posY+y)*10)].state = 2 // Rising
				}
			}
		}
		PopShape()
	} else {
		shake = 0.33
	}
}

// Checks if the shape can be extracted
func checkShape(posX int, posY int, shapeID int, shapeRotation int, highlightBlocking bool) bool {
	if shapeID == -1 {
		return false
	}
	/*
		// Reset board to white
		for i := 0; i < 200; i++ {
			if boardMatrix[i].state == 4 {
				boardMatrix[i].state = 1
			}
		}
	*/

	var shape = possibleShapes[shapeID][shapeRotation]
	var extractable = true
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if shape[x+((y%3)*3)] == 1 {
				// A part of the shape is not on a block
				if (posX+x)+((posY+y)*10) < 200 && boardMatrix[(posX+x)+((posY+y)*10)].state != 1 {
					extractable = false
				}
				// Any block above the shape would prevent it from been able to be extracted
				// Top row (special case)
				pos := (posX + x) + ((posY + y - 1) * 10)
				if pos < 0 {
					extractable = false
				}
				if pos > 199 {
					extractable = false
				}
				if y == 0 {
					if boardMatrix[pos].state == 1 {
						if highlightBlocking {
							boardMatrix[pos].fade = 1
						}
						extractable = false
					}
				}
				// Second and third row has to exclude blocks covered by the first row
				if y == 1 || y == 2 {
					if boardMatrix[pos].state == 1 && shape[x+(((y%3)-1)*3)] == 0 {
						if highlightBlocking {
							boardMatrix[pos].fade = 1
						}
						extractable = false
					}
				}
			}
		}
	}
	return extractable
}

func raiseShapes() (raisedSomething bool) {
	raisedSomething = false
	for i := 0; i < 200; i++ {
		if boardMatrix[i].state == 2 {
			raisedSomething = true
			boardMatrix[i].state = 0
			boardMatrix[i].color = getColor(0)
			if i-10 > 0 {
				boardMatrix[i-10].state = 2
				boardMatrix[i-10].color = getColor(2)
			}
		}
	}
	redrawBoardImage = true
	return raisedSomething
}

var stillPossibleResult bool

func stillPossible() bool {
	for x := 0; x < 9; x++ {
		for y := 1; y < 19; y++ {
			for r := 0; r < 4; r++ {
				y2 := y
				if y >= 18 {
					var shape = possibleShapes[currentShape][r]
					if shape[6] == 0 && shape[7] == 0 && shape[8] == 0 {
						y2 = 18
					} else {
						y2 = 17
					}
				}
				if checkShape(x, y2, currentShape, r, false) {
					fmt.Println("Found spot ", currentShape, x, y2, r)
					drawShape(boardImage, x, y2-3, 1, currentShape, r, false, 0, 0)
					return true
				}
			}
		}
	}
	fmt.Println("Still possible? NO")
	return false
}
