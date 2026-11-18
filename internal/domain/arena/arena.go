package arena

import (
	"bytes"
	"math/rand"
	"time"
)

type Arena struct {
	Height, Width         int
	CellHeight, CellWidth int
	X, Y                  float64
	Collisions            [][]byte
}

func NewArena(height, width, heightGridSize, widthGridSize, screenWidth, screenHeight int) *Arena {
	var arena Arena
	maps := [][]byte{M1, M2, M3}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	ri := r.Intn(len(maps))

	rawMap := maps[ri]
	cleanMap := bytes.TrimSpace(rawMap)
	rows := bytes.Split(cleanMap, []byte("\n"))

	actualRows := len(rows) + 2
	actualCols := len(rows[0]) + 2
	arenaWidth := width * actualCols
	arenaHeight := height * actualRows

	arena.Height = arenaHeight
	arena.Width = arenaWidth
	arena.CellHeight = height
	arena.CellWidth = width

	collisions := make([][]byte, actualRows)

	for rowIdx := range actualRows {
		collisions[rowIdx] = make([]byte, actualCols)

		isTopOrBottom := rowIdx == 0 || rowIdx == actualRows-1

		for colIdx := range actualCols {
			isLeftOrRight := colIdx == 0 || colIdx == actualCols-1

			if isTopOrBottom || isLeftOrRight {
				collisions[rowIdx][colIdx] = '1'
			} else {
				origRow := rowIdx - 1
				origCol := colIdx - 1
				if rows[origRow][origCol] == '1' {
					collisions[rowIdx][colIdx] = '1'
				} else {
					collisions[rowIdx][colIdx] = '0'
				}
			}
		}
	}

	arena.Collisions = collisions

	posX := (screenWidth / 2) - (arenaWidth / 2)
	posY := (screenHeight / 2) - (arenaHeight / 2)
	arena.X = float64(posX)
	arena.Y = float64(posY)

	return &arena
}
