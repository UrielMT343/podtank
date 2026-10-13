package arena

import (
	"bytes"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Arena struct {
	Height, Width int
	Battlefield   *ebiten.Image
	X, Y          float64
	Collisions    [][]byte
}

func NewArena(height, width, heightGridSize, widthGridSize, screenWidth, screenHeight int) Arena {
	arenaWidth := width * widthGridSize
	arenaHeight := height * heightGridSize

	var arena Arena
	arena.Height = arenaHeight
	arena.Width = arenaWidth

	obstacleTile := ebiten.NewImage(width, height)
	obstacleTile.Fill(color.RGBA{255, 0, 0, 1.0})

	floorTile := ebiten.NewImage(width, height)
	floorTile.Fill(color.RGBA{0, 255, 0, 1.0})

	i := ebiten.NewImage(arenaWidth, arenaHeight)
	arena.Battlefield = i

	maps := [][]byte{M1, M2, M3}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	ri := r.Intn(len(maps))

	var collisions [][]byte

	rawMap := maps[ri]
	cleanMap := bytes.TrimSpace(rawMap)
	rows := bytes.Split(cleanMap, []byte("\n"))

	collisions = make([][]byte, len(rows))

	for rowIdx, rowBytes := range rows {
		collisions[rowIdx] = make([]byte, 0, len(rowBytes))
		posYTile := rowIdx * height

		for colIdx, c := range rowBytes {
			posXTile := colIdx * width

			geo := ebiten.GeoM{}
			geo.Translate(float64(posXTile), float64(posYTile))
			options := &ebiten.DrawImageOptions{GeoM: geo}

			if c == '1' {
				arena.Battlefield.DrawImage(obstacleTile, options)
				collisions[rowIdx] = append(collisions[rowIdx], '1')
			} else {
				arena.Battlefield.DrawImage(floorTile, options)
				collisions[rowIdx] = append(collisions[rowIdx], '0')
			}
		}
	}

	arena.Collisions = collisions

	posX := (screenWidth / 2) - (arenaWidth / 2)
	posY := (screenHeight / 2) - (arenaHeight / 2)

	arena.X = float64(posX)
	arena.Y = float64(posY)

	return arena
}

func (a *Arena) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(a.X, a.Y)

	screen.DrawImage(a.Battlefield, options)
}
