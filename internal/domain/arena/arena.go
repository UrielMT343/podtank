package arena

import (
	"bytes"
	"math/rand"
	"time"

	"github.com/UrielMT343/podtank/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Arena struct {
	Height, Width int
	Battlefield   *ebiten.Image
	X, Y          float64
	Collisions    [][]byte
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
	arena.Battlefield = ebiten.NewImage(arenaWidth, arenaHeight)

	loadSprite := func(path string) *ebiten.Image {
		img, _, err := ebitenutil.NewImageFromFileSystem(assets.FS, path)
		if err != nil {
			panic(err)
		}
		return img
	}

	tble := loadSprite("sprites/tile_bound_left_end.png")
	tbrg := loadSprite("sprites/tile_bound_right_end.png")
	tbh := loadSprite("sprites/tile_bound_horizontal.png")
	tbue := loadSprite("sprites/tile_bound_up_end.png")
	tbde := loadSprite("sprites/tile_bound_down_end.png")
	tbv := loadSprite("sprites/tile_bound_vertical.png")

	obstacleSprite := loadSprite("sprites/tile_obstacle.png")
	sandSprite := loadSprite("sprites/tile_sand.png")

	scaleBoundX := float64(width) / 16.0
	scaleBoundY := float64(height) / 16.0

	obsBounds := obstacleSprite.Bounds().Size()
	scaleObsX := float64(width) / float64(obsBounds.X)
	scaleObsY := float64(height) / float64(obsBounds.Y)

	sandBounds := sandSprite.Bounds().Size()
	scaleSandX := float64(width) / float64(sandBounds.X)
	scaleSandY := float64(height) / float64(sandBounds.Y)

	var collisions [][]byte
	collisions = make([][]byte, actualRows)

	for rowIdx := 0; rowIdx < actualRows; rowIdx++ {
		collisions[rowIdx] = make([]byte, actualCols)
		posYTile := float64(rowIdx * height)

		for colIdx := 0; colIdx < actualCols; colIdx++ {
			posXTile := float64(colIdx * width)
			options := &ebiten.DrawImageOptions{}

			isTopOrBottom := rowIdx == 0 || rowIdx == actualRows-1
			isLeftOrRight := colIdx == 0 || colIdx == actualCols-1

			if isTopOrBottom || isLeftOrRight {
				options.GeoM.Scale(scaleBoundX, scaleBoundY)
				options.GeoM.Translate(posXTile, posYTile)

				var tileImg *ebiten.Image
				if isTopOrBottom {
					if colIdx == 0 {
						tileImg = tble
					} else if colIdx == actualCols-1 {
						tileImg = tbrg
					} else {
						tileImg = tbh
					}
				} else {
					if rowIdx == 1 {
						tileImg = tbue
					} else if rowIdx == actualRows-2 {
						tileImg = tbde
					} else {
						tileImg = tbv
					}
				}

				arena.Battlefield.DrawImage(tileImg, options)
				collisions[rowIdx][colIdx] = '1'

			} else {
				origRow := rowIdx - 1
				origCol := colIdx - 1
				c := rows[origRow][origCol]

				if c == '1' {
					options.GeoM.Scale(scaleObsX, scaleObsY)
					options.GeoM.Translate(posXTile, posYTile)
					arena.Battlefield.DrawImage(obstacleSprite, options)
					collisions[rowIdx][colIdx] = '1'
				} else {
					options.GeoM.Scale(scaleSandX, scaleSandY)
					options.GeoM.Translate(posXTile, posYTile)
					arena.Battlefield.DrawImage(sandSprite, options)
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

func (a *Arena) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(a.X, a.Y)

	screen.DrawImage(a.Battlefield, options)
}
