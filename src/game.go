package main

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var boardMatrix [200]int
var boardImage *ebiten.Image = ebiten.NewImage(10*7, 20*7)

//go:embed assets/levels.json
var levelFile []byte

type JsonFile struct {
	Levels []Level `json:"levels"`
}
type Level struct {
	Start  []int `json:"start"`
	Shapes []int `json:"shapes"`
}

func gameInit() {
	boardMatrix = [200]int{}
	loadLevel(0)
}

func loadLevel(level int) {
	var jsonFile JsonFile
	json.Unmarshal(levelFile, &jsonFile)
	fmt.Print(jsonFile)
	fmt.Print(jsonFile.Levels[level].Shapes)
}

func gameUpdate() {

}

func gameDraw(screen *ebiten.Image) {
	drawOutlinedRect(screen, 0, 0, 10*7, 20*7, getColor(3), getColor(0))
}
