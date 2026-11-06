package main

import (
	"fmt"
	"log"

	"github.com/UrielMT343/podtank/internal/domain/arena"
	"github.com/UrielMT343/podtank/internal/domain/tank"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	cellGridWidth  = 50
	cellGridHeight = 50

	arenaGridHeight = 16
	arenaGridWidth  = 13

	layoutWidth  = 800
	layoutHeight = 900
)

func Start() (battlefield *arena.Arena, tanks map[tank.Player]*tank.Tank, pools map[tank.Player]*tank.BulletPool) {
	b := arena.NewArena(cellGridHeight, cellGridWidth, arenaGridHeight,
		arenaGridWidth, layoutWidth, layoutHeight)

	t1 := tank.NewTank(50, 50, b.X, b.Y, cellGridWidth, cellGridHeight, tank.Down, tank.Player1)
	t2 := tank.NewTank(650, 800, b.X, b.Y, cellGridWidth, cellGridHeight, tank.Up, tank.Player2)

	tanks = map[tank.Player]*tank.Tank{
		tank.Player1: t1,
		tank.Player2: t2,
	}

	pools = map[tank.Player]*tank.BulletPool{
		tank.Player1: tank.NewBulletPool(50, tank.Player1),
		tank.Player2: tank.NewBulletPool(50, tank.Player2),
	}

	return b, tanks, pools
}

func GetLocalInput() tank.CmdState {
	return tank.CmdState{
		Up:    ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp),
		Down:  ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown),
		Left:  ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft),
		Right: ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight),
		Shoot: ebiten.IsKeyPressed(ebiten.KeySpace),
	}
}

type Game struct {
	Battlefield *arena.Arena
	Tanks       map[tank.Player]*tank.Tank
	LocalID     tank.Player
	BulletPools map[tank.Player]*tank.BulletPool
	GameOver    bool
	Winner      string
}

func (g *Game) Update() error {
	if g.GameOver {
		return nil
	}

	localCmd := GetLocalInput()
	if myTank, ok := g.Tanks[g.LocalID]; ok {
		myTank.Update(localCmd, &g.Battlefield.Collisions, g.BulletPools[g.LocalID])
	}

	remoteCmd := tank.CmdState{}
	enemyID := tank.Player2
	if g.LocalID == tank.Player2 {
		enemyID = tank.Player1
	}

	if enemyTank, ok := g.Tanks[enemyID]; ok {
		enemyTank.Update(remoteCmd, &g.Battlefield.Collisions, g.BulletPools[enemyID])
	}

	for _, bp := range g.BulletPools {
		bp.Update(&g.Battlefield.Collisions, g.Tanks)
	}

	for id, t := range g.Tanks {
		if t.IsDead() {
			g.GameOver = true
			if id == tank.Player1 {
				g.Winner = "PLAYER 2"
			} else {
				g.Winner = "PLAYER 1"
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Battlefield.Draw(screen)

	for _, t := range g.Tanks {
		t.Draw(screen)
	}

	for _, bp := range g.BulletPools {
		bp.Draw(screen)
	}

	p1HP := g.Tanks[tank.Player1].HP()
	p2HP := g.Tanks[tank.Player2].HP()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("P1 HP: %d  |  P2 HP: %d", p1HP, p2HP))

	if g.GameOver {
		ebitenutil.DebugPrintAt(screen, "GAME OVER - "+g.Winner+" WINS!", layoutWidth/2-100, layoutHeight/2)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return layoutWidth, layoutHeight
}

func main() {
	_, monitorH := ebiten.Monitor().Size()

	targetHeight := int(float64(monitorH) * 0.70)

	aspectRatio := float64(layoutWidth) / float64(layoutHeight)
	targetWidth := int(float64(targetHeight)*aspectRatio) + 100

	ebiten.SetWindowSize(targetWidth, targetHeight)
	ebiten.SetWindowTitle("Podtank")

	b, tanks, bp := Start()

	if err := ebiten.RunGame(&Game{
		Battlefield: b,
		Tanks:       tanks,
		LocalID:     tank.Player1,
		BulletPools: bp,
	}); err != nil {
		log.Fatal(err)
	}
}
