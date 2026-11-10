package tank

import (
	"image/color"
	"math"

	"github.com/UrielMT343/podtank/assets"
	"github.com/UrielMT343/podtank/internal/audio"
	"github.com/UrielMT343/podtank/internal/domain/vfx"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Bullet struct {
	x, y, vx, vy     float64
	offsetX, offsetY float64
	dir              TankDirection
}

type BulletPool struct {
	pool                     []*Bullet
	active                   []*Bullet
	bulletImg, trailImg      *ebiten.Image
	Owner                    Player
	scaleX, scaleY           float64
	trailScaleX, trailScaleY float64
}

func NewBulletPool(size int, owner Player) *BulletPool {
	bp := BulletPool{
		pool:   make([]*Bullet, 0, size),
		active: make([]*Bullet, 0, size),
		Owner:  owner,
	}

	for i := 0; i < size; i++ {
		bp.pool = append(bp.pool, &Bullet{})
	}

	var imagePath string
	if owner == Player1 {
		imagePath = "sprites/bullet_blue.png"
	} else {
		imagePath = "sprites/bullet_red.png"
	}

	sprite, img, err := ebitenutil.NewImageFromFileSystem(assets.FS, imagePath)
	if err == nil {
		bp.bulletImg = sprite
		bounds := img.Bounds().Size()
		bp.scaleX = 15.0 / float64(bounds.X)
		bp.scaleY = 15.0 / float64(bounds.Y)
	} else {
		fallback := ebiten.NewImage(15, 15)
		fallback.Fill(color.White)
		bp.bulletImg = fallback
		bp.scaleX = 1.0
		bp.scaleY = 1.0
	}

	trailSprite, trailImgObj, err := ebitenutil.NewImageFromFileSystem(assets.FS, "sprites/shotThin.png")
	if err == nil {
		bp.trailImg = trailSprite
		tbounds := trailImgObj.Bounds().Size()
		bp.trailScaleX = 15.0 / float64(tbounds.X)
		bp.trailScaleY = bp.trailScaleX
	}

	return &bp
}

func (bp *BulletPool) Update(collisions *[][]byte, tanks map[Player]*Tank, vp *vfx.VFXPool, sfx *audio.SoundBank) {
	grid := *collisions
	maxRows := len(grid)
	maxCols := 0
	if maxRows > 0 {
		maxCols = len(grid[0])
	}

	for i := 0; i < len(bp.active); i++ {
		b := bp.active[i]

		b.x += b.vx
		b.y += b.vy

		hitWall := false
		hitTank := false

		if b.x < 0 || b.y < 0 || b.x >= float64(maxCols*50) || b.y >= float64(maxRows*50) {
			hitWall = true
		} else {
			col := int(b.x) / 50
			row := int(b.y) / 50
			if grid[row][col] == '1' {
				hitWall = true
			}
		}

		if !hitWall {
			for playerID, targetTank := range tanks {
				if bp.Owner == playerID {
					continue
				}

				if b.x >= targetTank.x && b.x <= targetTank.x+float64(targetTank.width) &&
					b.y >= targetTank.y && b.y <= targetTank.y+float64(targetTank.height) {
					targetTank.healthPoints--
					hitTank = true
					break
				}
			}
		}

		if hitWall {
			if vp != nil {
				vp.Spawn(b.x, b.y, vfx.EffectSmoke)
			}
			if sfx != nil && sfx.HitWall != nil {
				audio.PlaySFX(sfx.HitWall)
			}
		} else if hitTank {
			if vp != nil {
				vp.Spawn(b.x, b.y, vfx.EffectExplosion)
			}
			if sfx != nil && sfx.Explosion != nil {
				audio.PlaySFX(sfx.Explosion)
			}
		}

		if hitWall || hitTank {
			bp.pool = append(bp.pool, b)
			lastIdx := len(bp.active) - 1
			bp.active[i] = bp.active[lastIdx]
			bp.active = bp.active[:lastIdx]
			i--
		}
	}
}

func (bp *BulletPool) Draw(screen *ebiten.Image) {
	for _, b := range bp.active {
		if bp.trailImg != nil {
			tbounds := bp.trailImg.Bounds().Size()
			tW := float64(tbounds.X) * bp.trailScaleX
			tH := float64(tbounds.Y) * bp.trailScaleY

			trailOpt := &ebiten.DrawImageOptions{}
			trailOpt.Filter = ebiten.FilterNearest

			trailOpt.GeoM.Scale(bp.trailScaleX, bp.trailScaleY)
			trailOpt.GeoM.Translate(-tW/2.0, -tH/2.0)
			trailOpt.GeoM.Translate(0, 10.0)

			switch b.dir {
			case Up:
			case Down:
				trailOpt.GeoM.Rotate(math.Pi)
			case Left:
				trailOpt.GeoM.Rotate(-math.Pi / 2)
			case Right:
				trailOpt.GeoM.Rotate(math.Pi / 2)
			}

			trailOpt.GeoM.Translate(b.x+b.offsetX, b.y+b.offsetY)
			screen.DrawImage(bp.trailImg, trailOpt)
		}

		options := &ebiten.DrawImageOptions{}
		options.Filter = ebiten.FilterNearest

		options.GeoM.Scale(bp.scaleX, bp.scaleY)
		options.GeoM.Translate(-7.5, -7.5)

		switch b.dir {
		case Up:
		case Down:
			options.GeoM.Rotate(math.Pi)
		case Left:
			options.GeoM.Rotate(-math.Pi / 2)
		case Right:
			options.GeoM.Rotate(math.Pi / 2)
		}

		options.GeoM.Translate(7.5, 7.5)
		options.GeoM.Translate(b.x+b.offsetX-7.5, b.y+b.offsetY-7.5)

		screen.DrawImage(bp.bulletImg, options)
	}
}
