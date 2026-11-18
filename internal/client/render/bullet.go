package render

import (
	"image/color"
	"math"

	"github.com/UrielMT343/podtank/assets"
	"github.com/UrielMT343/podtank/internal/domain/tank"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type BulletRenderer struct {
	bulletSprites            map[tank.Player]*ebiten.Image
	trailImg                 *ebiten.Image
	scaleX, scaleY           float64
	trailScaleX, trailScaleY float64
}

func NewBulletRenderer() *BulletRenderer {
	loadBullet := func(path string) (*ebiten.Image, float64, float64) {
		sprite, img, err := ebitenutil.NewImageFromFileSystem(assets.FS, path)
		if err == nil {
			bounds := img.Bounds().Size()
			return sprite, 15.0 / float64(bounds.X), 15.0 / float64(bounds.Y)
		}
		fallback := ebiten.NewImage(15, 15)
		fallback.Fill(color.White)
		return fallback, 1.0, 1.0
	}

	b1, sx, sy := loadBullet("sprites/bullet_blue.png")
	b2, _, _ := loadBullet("sprites/bullet_red.png")

	br := &BulletRenderer{
		bulletSprites: map[tank.Player]*ebiten.Image{
			tank.Player1: b1,
			tank.Player2: b2,
		},
		scaleX: sx,
		scaleY: sy,
	}

	trailSprite, trailImgObj, err := ebitenutil.NewImageFromFileSystem(assets.FS, "sprites/shotThin.png")
	if err == nil {
		br.trailImg = trailSprite
		tbounds := trailImgObj.Bounds().Size()
		br.trailScaleX = 15.0 / float64(tbounds.X)
		br.trailScaleY = br.trailScaleX
	}

	return br
}

func (br *BulletRenderer) Draw(screen *ebiten.Image, bp *tank.BulletPool) {
	bulletImg := br.bulletSprites[bp.Owner]

	for _, b := range bp.Active() {
		if br.trailImg != nil {
			tbounds := br.trailImg.Bounds().Size()
			tW := float64(tbounds.X) * br.trailScaleX
			tH := float64(tbounds.Y) * br.trailScaleY

			trailOpt := &ebiten.DrawImageOptions{}
			trailOpt.Filter = ebiten.FilterNearest

			trailOpt.GeoM.Scale(br.trailScaleX, br.trailScaleY)
			trailOpt.GeoM.Translate(-tW/2.0, -tH/2.0)
			trailOpt.GeoM.Translate(0, 10.0)

			switch b.Dir() {
			case tank.Up:
			case tank.Down:
				trailOpt.GeoM.Rotate(math.Pi)
			case tank.Left:
				trailOpt.GeoM.Rotate(-math.Pi / 2)
			case tank.Right:
				trailOpt.GeoM.Rotate(math.Pi / 2)
			}

			trailOpt.GeoM.Translate(b.X()+b.OffsetX(), b.Y()+b.OffsetY())
			screen.DrawImage(br.trailImg, trailOpt)
		}

		options := &ebiten.DrawImageOptions{}
		options.Filter = ebiten.FilterNearest

		options.GeoM.Scale(br.scaleX, br.scaleY)
		options.GeoM.Translate(-7.5, -7.5)

		switch b.Dir() {
		case tank.Up:
		case tank.Down:
			options.GeoM.Rotate(math.Pi)
		case tank.Left:
			options.GeoM.Rotate(-math.Pi / 2)
		case tank.Right:
			options.GeoM.Rotate(math.Pi / 2)
		}

		options.GeoM.Translate(7.5, 7.5)
		options.GeoM.Translate(b.X()+b.OffsetX()-7.5, b.Y()+b.OffsetY()-7.5)

		screen.DrawImage(bulletImg, options)
	}
}
