package tank

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	bulletWidth  = 15
	bulletHeight = 15
)

type Bullet struct {
	x, y, vx, vy     float64
	offsetX, offsetY float64
	Player           Player
}

type BulletPool struct {
	pool      []*Bullet
	active    []*Bullet
	bulletImg *ebiten.Image
}

func NewBulletPool(size int) *BulletPool {
	bp := BulletPool{
		pool:   make([]*Bullet, 0, size),
		active: make([]*Bullet, 0, size),
	}

	for i := 0; i < size; i++ {
		bp.pool = append(bp.pool, &Bullet{})
	}

	bulletImg := ebiten.NewImage(bulletWidth, bulletHeight)
	bulletImg.Fill(color.White)

	bp.bulletImg = bulletImg

	return &bp
}

func (bp *BulletPool) Update(collisions *[][]byte) {
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

		hit := false

		if b.x < 0 || b.y < 0 || b.x >= float64(maxCols*50) || b.y >= float64(maxRows*50) {
			hit = true
		} else {
			col := int(b.x) / 50
			row := int(b.y) / 50

			if grid[row][col] == '1' {
				hit = true
			}
		}

		if hit {
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
		options := &ebiten.DrawImageOptions{}
		options.Filter = ebiten.FilterNearest

		options.GeoM.Translate(-(bulletWidth / 2), -(bulletHeight / 2))

		options.GeoM.Translate(b.x+b.offsetX, b.y+b.offsetY)
		screen.DrawImage(bp.bulletImg, options)
	}
}
