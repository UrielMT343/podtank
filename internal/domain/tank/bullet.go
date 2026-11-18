package tank

type ImpactType int

const (
	ImpactWall ImpactType = iota
	ImpactTank
)

type ImpactEvent struct {
	Type ImpactType
	X, Y float64
}

type Bullet struct {
	x, y, vx, vy     float64
	offsetX, offsetY float64
	dir              TankDirection
}

func (b *Bullet) X() float64         { return b.x }
func (b *Bullet) Y() float64         { return b.y }
func (b *Bullet) OffsetX() float64   { return b.offsetX }
func (b *Bullet) OffsetY() float64   { return b.offsetY }
func (b *Bullet) Dir() TankDirection { return b.dir }

type BulletPool struct {
	pool   []*Bullet
	active []*Bullet
	Owner  Player
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

	return &bp
}

func (bp *BulletPool) Active() []*Bullet {
	return bp.active
}

func (bp *BulletPool) Update(collisions *[][]byte, tanks map[Player]*Tank) []ImpactEvent {
	var impacts []ImpactEvent

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
			impacts = append(impacts, ImpactEvent{Type: ImpactWall, X: b.x, Y: b.y})
		} else if hitTank {
			impacts = append(impacts, ImpactEvent{Type: ImpactTank, X: b.x, Y: b.y})
		}

		if hitWall || hitTank {
			bp.pool = append(bp.pool, b)
			lastIdx := len(bp.active) - 1
			bp.active[i] = bp.active[lastIdx]
			bp.active = bp.active[:lastIdx]
			i--
		}
	}

	return impacts
}
