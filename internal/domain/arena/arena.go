package arena

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type Arena struct {
	Height, Width int
	Battlefield   *ebiten.Image
	Position      ebiten.GeoM
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
	//i.Fill(color.White)
	//i.Bounds()
	arena.Battlefield = i

	rawMap := ReadMap1()
	cleanMap := strings.TrimSpace(rawMap)
	rows := strings.Split(cleanMap, "\n")
	for rowIdx, row := range rows {
		for colIdx, c := range row {
			posXTile := colIdx * width
			posYTile := rowIdx * height

			geo := ebiten.GeoM{}
			geo.Translate(float64(posXTile), float64(posYTile))

			options := &ebiten.DrawImageOptions{GeoM: geo}
			if c == '1' {
				arena.Battlefield.DrawImage(obstacleTile, options)
			} else {
				arena.Battlefield.DrawImage(floorTile, options)
			}
		}
	}

	posX := (screenWidth / 2) - (arenaWidth / 2)
	posY := (screenHeight / 2) - (arenaHeight / 2)

	geo := ebiten.GeoM{}
	geo.Translate(float64(posX), float64(posY))
	arena.Position = geo

	return arena
}

func (a *Arena) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{
		GeoM: a.Position,
	}

	screen.DrawImage(a.Battlefield, options)
}
