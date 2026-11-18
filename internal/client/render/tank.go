package render

import (
	"math"

	"github.com/UrielMT343/podtank/assets"
	"github.com/UrielMT343/podtank/internal/domain/tank"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TankRenderer struct {
	sprites map[tank.Player]*ebiten.Image
}

func NewTankRenderer() *TankRenderer {
	loadSprite := func(path string) *ebiten.Image {
		img, _, err := ebitenutil.NewImageFromFileSystem(assets.FS, path)
		if err != nil {
			panic(err)
		}
		return img
	}

	return &TankRenderer{
		sprites: map[tank.Player]*ebiten.Image{
			tank.Player1: loadSprite("sprites/tank_blue.png"),
			tank.Player2: loadSprite("sprites/tank_red.png"),
		},
	}
}

func (tr *TankRenderer) Draw(screen *ebiten.Image, t *tank.Tank) {
	if t.IsDead() {
		return
	}

	sprite, ok := tr.sprites[t.Player]
	if !ok {
		return
	}

	s := sprite.Bounds().Size()
	scaleX := float64(t.Width()) / float64(s.X)
	scaleY := float64(t.Height()) / float64(s.Y)

	options := &ebiten.DrawImageOptions{}
	options.Filter = ebiten.FilterNearest
	options.GeoM.Scale(scaleX, scaleY)

	dir := t.Direction()
	if dir == tank.None {
		dir = t.LastDirection()
	}

	switch dir {
	case tank.Up:
		options.GeoM.Translate(-25, -25)
		options.GeoM.Rotate(math.Pi)
		options.GeoM.Translate(25, 25)
	case tank.Down:
		options.GeoM.Translate(-25, -25)
		options.GeoM.Rotate(0.0)
		options.GeoM.Translate(25, 25)
	case tank.Left:
		options.GeoM.Translate(-25, -25)
		options.GeoM.Rotate(math.Pi / 2)
		options.GeoM.Translate(25, 25)
	case tank.Right:
		options.GeoM.Translate(-25, -25)
		options.GeoM.Rotate((3 * math.Pi) / 2)
		options.GeoM.Translate(25, 25)
	}

	options.GeoM.Translate(t.X()+t.OffsetX(), t.Y()+t.OffsetY())

	screen.DrawImage(sprite, options)
}
