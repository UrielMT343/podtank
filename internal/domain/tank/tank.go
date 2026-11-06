package tank

import (
	"math"

	"github.com/UrielMT343/podtank/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TankDirection int

type Player byte

type CmdState struct {
	Up, Down, Left, Right, Shoot bool
}

const (
	initialHP                  = 5
	initialDmg                 = 14
	initialSpd                 = 50
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
	x, y, offsetX, offsetY, scaleX, scaleY        float64
	direction, lastDirection                      TankDirection
	height, width, healthPoints, damage, speed    int
	image                                         *ebiten.Image
	Player                                        Player
	shotCooldown                                  int
	upFrames, downFrames, leftFrames, rightFrames int
}

func NewTank(spawnX, spawnY, screenOffsetX, screenOffsetY float64, tankWidth, tankHeight int, initialDirection TankDirection, player Player) *Tank {
	var tank Tank
	tank.healthPoints = initialHP

	tank.x = spawnX
	tank.y = spawnY

	tank.offsetX = screenOffsetX
	tank.offsetY = screenOffsetY

	tank.height = tankHeight
	tank.width = tankWidth

	tank.direction = initialDirection
	tank.lastDirection = initialDirection
	tank.damage = initialDmg
	tank.speed = initialSpd

	tank.Player = player

	var imagePath string
	if player == Player1 {
		imagePath = "sprites/tank_blue.png"
	} else {
		imagePath = "sprites/tank_red.png"
	}

	sprite, img, err := ebitenutil.NewImageFromFileSystem(assets.FS, imagePath)
	if err != nil {
		panic(err)
	}

	s := img.Bounds().Size()

	tank.scaleX = float64(tank.width) / float64(s.X)
	tank.scaleY = float64(tank.height) / float64(s.Y)
	tank.image = sprite

	return &tank
}

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

	bulleDir := t.direction
	if bulleDir == None {
		bulleDir = t.lastDirection
	}

	b.dir = bulleDir

	pool.active = append(pool.active, b)
}

func (t *Tank) IsDead() bool {
	return t.healthPoints <= 0
}

func (t *Tank) HP() int {
	return t.healthPoints
}

func (t *Tank) Update(cmd CmdState, collisions *[][]byte, bp *BulletPool) {
	mov := float64(t.speed) / 50.0
	pixelNearest := 5.0

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
}

func (t *Tank) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.Filter = ebiten.FilterNearest

	options.GeoM.Scale(t.scaleX, t.scaleY)

	dir := t.direction
	if dir == None {
		dir = t.lastDirection
	}

	switch dir {
	case Up:
		options.GeoM.Translate(-25, -25)
		options.GeoM.Rotate(math.Pi)
		options.GeoM.Translate(25, 25)
	case Down:
		options.GeoM.Translate(-25, -25)
		options.GeoM.Rotate(0.0)
		options.GeoM.Translate(25, 25)
	case Left:
		options.GeoM.Translate(-25, -25)
		options.GeoM.Rotate(math.Pi / 2)
		options.GeoM.Translate(25, 25)
	case Right:
		options.GeoM.Translate(-25, -25)
		options.GeoM.Rotate((3 * math.Pi) / 2)
		options.GeoM.Translate(25, 25)
	}

	options.GeoM.Translate(t.x+t.offsetX, t.y+t.offsetY)

	screen.DrawImage(t.image, options)
}
