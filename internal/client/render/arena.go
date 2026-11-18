package render

import (
	"github.com/UrielMT343/podtank/assets"
	"github.com/UrielMT343/podtank/internal/domain/arena"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ArenaRenderer struct {
	Battlefield *ebiten.Image
}

func NewArenaRenderer(a *arena.Arena) *ArenaRenderer {
	battlefield := ebiten.NewImage(a.Width, a.Height)

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

	cellW := float64(a.CellWidth)
	cellH := float64(a.CellHeight)

	scaleBoundX := cellW / 16.0
	scaleBoundY := cellH / 16.0

	obsBounds := obstacleSprite.Bounds().Size()
	scaleObsX := cellW / float64(obsBounds.X)
	scaleObsY := cellH / float64(obsBounds.Y)

	sandBounds := sandSprite.Bounds().Size()
	scaleSandX := cellW / float64(sandBounds.X)
	scaleSandY := cellH / float64(sandBounds.Y)

	actualRows := len(a.Collisions)
	actualCols := len(a.Collisions[0])

	for rowIdx := range actualRows {
		posYTile := float64(rowIdx * a.CellHeight)

		for colIdx := range actualCols {
			posXTile := float64(colIdx * a.CellWidth)
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

				battlefield.DrawImage(tileImg, options)
			} else if a.Collisions[rowIdx][colIdx] == '1' {
				options.GeoM.Scale(scaleObsX, scaleObsY)
				options.GeoM.Translate(posXTile, posYTile)
				battlefield.DrawImage(obstacleSprite, options)
			} else {
				options.GeoM.Scale(scaleSandX, scaleSandY)
				options.GeoM.Translate(posXTile, posYTile)
				battlefield.DrawImage(sandSprite, options)
			}
		}
	}

	return &ArenaRenderer{
		Battlefield: battlefield,
	}
}

func (ar *ArenaRenderer) Draw(screen *ebiten.Image, a *arena.Arena) {
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(a.X, a.Y)

	screen.DrawImage(ar.Battlefield, options)
}
