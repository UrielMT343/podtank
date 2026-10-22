package tank

import (
	"math"

	"github.com/UrielMT343/podtank/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TankDirection int

type Player byte

const (
	initialHP                  = 5
	initialDmg                 = 14
	initialSpd                 = 50
	cellGridSize               = 50.0
	moveDelay                  = 15
	None         TankDirection = iota
	Up
	Down
	Left
	Right
	Player1 Player = '1'
	Player2 Player = '2'
)

type Tank struct {
	X, Y             float64
	OffsetX, OffsetY float64
	Height, Width    int
	Direction        TankDirection
	lastDirection    TankDirection
	HealthPoints     int
	Damage           int
	Speed            int
	Image            *ebiten.Image
	Player           Player
	scaleX, scaleY   float64
}

func NewTank(startX, startY float64, tankWidth, tankHeight int, initialDirection TankDirection, player Player) Tank {
	var tank Tank
	tank.HealthPoints = initialHP

	tank.X = 0
	tank.Y = 0

	tank.OffsetX = startX
	tank.OffsetY = startY

	tank.Height = tankHeight
	tank.Width = tankWidth

	tank.Direction = initialDirection
	tank.lastDirection = initialDirection
	tank.Damage = initialDmg
	tank.Speed = initialSpd

	tank.Player = player

	if player == Player1 {
		spriteP1, img, _ := ebitenutil.NewImageFromFileSystem(assets.FS, "sprites/tank_blue.png")
		s := img.Bounds().Size()

		scaleX := float64(tank.Width) / float64(s.X)
		scaleY := float64(tank.Height) / float64(s.Y)

		tank.scaleX = scaleX
		tank.scaleY = scaleY

		tank.Image = spriteP1
	}

	return tank
}

func (t *Tank) checkCollision(collisions *[][]byte) bool {
	grid := *collisions
	maxRows := len(grid)
	if maxRows == 0 {
		return false
	}
	maxCols := len(grid[0])

	leftEdge := t.X
	rightEdge := t.X + float64(t.Width) - 0.1
	topEdge := t.Y
	bottomEdge := t.Y + float64(t.Height) - 0.1

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

func (t *Tank) Update(collisions *[][]byte) {
	mov := float64(t.Speed) / 50.0
	pixelNearest := 5.0

	upFrames := inpututil.KeyPressDuration(ebiten.KeyArrowUp)
	downFrames := inpututil.KeyPressDuration(ebiten.KeyArrowDown)
	leftFrames := inpututil.KeyPressDuration(ebiten.KeyArrowLeft)
	rightFrames := inpututil.KeyPressDuration(ebiten.KeyArrowRight)

	if t.Direction == None {
		if upFrames >= 1 {
			t.Direction = Up
			t.lastDirection = Up
		} else if downFrames >= 1 {
			t.Direction = Down
			t.lastDirection = Down
		} else if leftFrames >= 1 {
			t.Direction = Left
			t.lastDirection = Left
		} else if rightFrames >= 1 {
			t.Direction = Right
			t.lastDirection = Right
		}
	}

	oldX := t.X
	oldY := t.Y

	switch t.Direction {
	case Up:
		if upFrames > 0 {
			if upFrames > moveDelay {
				t.Y -= mov

				nearestX := math.Round(t.X/50.0) * 50.0
				if math.Abs(t.X-nearestX) <= pixelNearest {
					t.X = nearestX
				}
			}
		} else {
			t.Direction = None
			nearestY := math.Round(t.Y/50.0) * 50.0
			if math.Abs(t.Y-nearestY) <= pixelNearest {
				t.Y = nearestY
			}
		}
	case Down:
		if downFrames > 0 {
			if downFrames > moveDelay {
				t.Y += mov

				nearestX := math.Round(t.X/50.0) * 50.0
				if math.Abs(t.X-nearestX) <= pixelNearest {
					t.X = nearestX
				}
			}
		} else {
			t.Direction = None
			nearestY := math.Round(t.Y/50.0) * 50.0
			if math.Abs(t.Y-nearestY) <= pixelNearest {
				t.Y = nearestY
			}
		}
	case Left:
		if leftFrames > 0 {
			if leftFrames > moveDelay {
				t.X -= mov

				nearestY := math.Round(t.Y/50.0) * 50.0
				if math.Abs(t.Y-nearestY) <= pixelNearest {
					t.Y = nearestY
				}
			}
		} else {
			t.Direction = None
			nearestX := math.Round(t.X/50.0) * 50.0
			if math.Abs(t.X-nearestX) <= pixelNearest {
				t.X = nearestX
			}
		}
	case Right:
		if rightFrames > 0 {
			if rightFrames > moveDelay {
				t.X += mov

				nearestY := math.Round(t.Y/50.0) * 50.0
				if math.Abs(t.Y-nearestY) <= pixelNearest {
					t.Y = nearestY
				}
			}
		} else {
			t.Direction = None
			nearestX := math.Round(t.X/50.0) * 50.0
			if math.Abs(t.X-nearestX) <= pixelNearest {
				t.X = nearestX
			}
		}
	}

	if t.checkCollision(collisions) {
		t.X = oldX
		t.Y = oldY
	}
}

func (t *Tank) Draw(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	options.Filter = ebiten.FilterNearest

	options.GeoM.Scale(t.scaleX, t.scaleY)

	dir := t.Direction
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

	options.GeoM.Translate(t.X+t.OffsetX, t.Y+t.OffsetY)

	screen.DrawImage(t.Image, options)
}
