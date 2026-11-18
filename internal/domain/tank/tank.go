package tank

import (
	"math"
)

type TankDirection int

type Player byte

type CmdState struct {
	Up, Down, Left, Right, Shoot bool
}

const (
	initialHP                  = 5
	initialDmg                 = 14
	initialSpd                 = 75
	cellGridSize               = 50.0
	moveDelay                  = 15
	shotDelay                  = 1000
	None         TankDirection = iota
	Up
	Down
	Left
	Right
	Player1 Player = '1'
	Player2 Player = '2'
)

type Tank struct {
	x, y, offsetX, offsetY                        float64
	direction, lastDirection                      TankDirection
	height, width, healthPoints, damage, speed    int
	Player                                        Player
	shotCooldown                                  int
	upFrames, downFrames, leftFrames, rightFrames int
}

func NewTank(spawnX, spawnY, screenOffsetX, screenOffsetY float64, tankWidth, tankHeight int, initialDirection TankDirection, player Player) *Tank {
	return &Tank{
		healthPoints:  initialHP,
		x:             spawnX,
		y:             spawnY,
		offsetX:       screenOffsetX,
		offsetY:       screenOffsetY,
		height:        tankHeight,
		width:         tankWidth,
		direction:     initialDirection,
		lastDirection: initialDirection,
		damage:        initialDmg,
		speed:         initialSpd,
		Player:        player,
	}
}

func (t *Tank) X() float64                   { return t.x }
func (t *Tank) Y() float64                   { return t.y }
func (t *Tank) OffsetX() float64             { return t.offsetX }
func (t *Tank) OffsetY() float64             { return t.offsetY }
func (t *Tank) Width() int                   { return t.width }
func (t *Tank) Height() int                  { return t.height }
func (t *Tank) Direction() TankDirection     { return t.direction }
func (t *Tank) LastDirection() TankDirection { return t.lastDirection }
func (t *Tank) IsDead() bool                 { return t.healthPoints <= 0 }
func (t *Tank) HP() int                      { return t.healthPoints }

func (t *Tank) checkCollision(collisions *[][]byte) bool {
	grid := *collisions
	maxRows := len(grid)
	if maxRows == 0 {
		return false
	}
	maxCols := len(grid[0])

	leftEdge := t.x
	rightEdge := t.x + float64(t.width) - 0.1
	topEdge := t.y
	bottomEdge := t.y + float64(t.height) - 0.1

	mapPixelWidth := float64(maxCols * cellGridSize)
	mapPixelHeight := float64(maxRows * cellGridSize)

	if leftEdge < 0 || topEdge < 0 || rightEdge >= mapPixelWidth || bottomEdge >= mapPixelHeight {
		return true
	}

	minCol := int(leftEdge) / cellGridSize
	maxCol := int(rightEdge) / cellGridSize
	minRow := int(topEdge) / cellGridSize
	maxRow := int(bottomEdge) / cellGridSize

	for row := minRow; row <= maxRow; row++ {
		for col := minCol; col <= maxCol; col++ {
			if grid[row][col] == '1' {
				return true
			}
		}
	}

	return false
}

func (t *Tank) Shoot(pool *BulletPool) {
	if len(pool.pool) == 0 {
		return
	}

	lastIdx := len(pool.pool) - 1
	b := pool.pool[lastIdx]
	pool.pool = pool.pool[:lastIdx]

	b.x = t.x + float64(t.width)/2.0
	b.y = t.y + float64(t.height)/2.0

	b.offsetX = t.offsetX
	b.offsetY = t.offsetY

	bulletSpeed := 10.0
	dir := t.direction
	if dir == None {
		dir = t.lastDirection
	}

	switch dir {
	case Up:
		b.vx = 0
		b.vy = -bulletSpeed
	case Down:
		b.vx = 0
		b.vy = bulletSpeed
	case Left:
		b.vx = -bulletSpeed
		b.vy = 0
	case Right:
		b.vx = bulletSpeed
		b.vy = 0
	}

	b.dir = dir

	pool.active = append(pool.active, b)
}

func (t *Tank) Update(cmd CmdState, collisions *[][]byte, bp *BulletPool) bool {
	mov := float64(t.speed) / 50.0
	pixelNearest := 5.0
	didShoot := false

	if cmd.Up {
		t.upFrames++
	} else {
		t.upFrames = 0
	}
	if cmd.Down {
		t.downFrames++
	} else {
		t.downFrames = 0
	}
	if cmd.Left {
		t.leftFrames++
	} else {
		t.leftFrames = 0
	}
	if cmd.Right {
		t.rightFrames++
	} else {
		t.rightFrames = 0
	}

	if t.shotCooldown > 0 {
		t.shotCooldown--
	}

	if cmd.Shoot && t.shotCooldown == 0 {
		t.Shoot(bp)
		didShoot = true
		t.shotCooldown = 60
	}

	if t.direction == None {
		if t.upFrames >= 1 {
			t.direction = Up
			t.lastDirection = Up
		} else if t.downFrames >= 1 {
			t.direction = Down
			t.lastDirection = Down
		} else if t.leftFrames >= 1 {
			t.direction = Left
			t.lastDirection = Left
		} else if t.rightFrames >= 1 {
			t.direction = Right
			t.lastDirection = Right
		}
	}

	oldX := t.x
	oldY := t.y

	switch t.direction {
	case Up:
		if t.upFrames > 0 {
			if t.upFrames > moveDelay {
				t.y -= mov

				nearestX := math.Round(t.x/50.0) * 50.0
				if math.Abs(t.x-nearestX) <= pixelNearest {
					t.x = nearestX
				}
			}
		} else {
			t.direction = None
			nearestY := math.Round(t.y/50.0) * 50.0
			if math.Abs(t.y-nearestY) <= pixelNearest {
				t.y = nearestY
			}
		}
	case Down:
		if t.downFrames > 0 {
			if t.downFrames > moveDelay {
				t.y += mov

				nearestX := math.Round(t.x/50.0) * 50.0
				if math.Abs(t.x-nearestX) <= pixelNearest {
					t.x = nearestX
				}
			}
		} else {
			t.direction = None
			nearestY := math.Round(t.y/50.0) * 50.0
			if math.Abs(t.y-nearestY) <= pixelNearest {
				t.y = nearestY
			}
		}
	case Left:
		if t.leftFrames > 0 {
			if t.leftFrames > moveDelay {
				t.x -= mov

				nearestY := math.Round(t.y/50.0) * 50.0
				if math.Abs(t.y-nearestY) <= pixelNearest {
					t.y = nearestY
				}
			}
		} else {
			t.direction = None
			nearestX := math.Round(t.x/50.0) * 50.0
			if math.Abs(t.x-nearestX) <= pixelNearest {
				t.x = nearestX
			}
		}
	case Right:
		if t.rightFrames > 0 {
			if t.rightFrames > moveDelay {
				t.x += mov

				nearestY := math.Round(t.y/50.0) * 50.0
				if math.Abs(t.y-nearestY) <= pixelNearest {
					t.y = nearestY
				}
			}
		} else {
			t.direction = None
			nearestX := math.Round(t.x/50.0) * 50.0
			if math.Abs(t.x-nearestX) <= pixelNearest {
				t.x = nearestX
			}
		}
	}

	if t.checkCollision(collisions) {
		t.x = oldX
		t.y = oldY
	}

	return didShoot
}
