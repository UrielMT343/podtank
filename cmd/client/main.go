package main

import (
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

func Init() (battlefield *arena.Arena, player1 *tank.Tank) {
	b := arena.NewArena(cellGridHeight, cellGridWidth, arenaGridHeight,
		arenaGridWidth, layoutWidth, layoutHeight)

	t1 := tank.NewTank(b.X, b.Y, cellGridWidth, cellGridHeight, int(tank.None))

	return &b, &t1
}

type Game struct {
	Battlefield *arena.Arena
	Player1     *tank.Tank
}

func (g *Game) Update() error {
	g.Player1.Update(&g.Battlefield.Collisions)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Test")
	screen.Bounds()
	g.Battlefield.Draw(screen)
	g.Player1.Draw(screen)
}

// Takes the original size and converts it to a more logical size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return layoutWidth, layoutHeight
}

func main() {
	_, monitorH := ebiten.Monitor().Size()

	targetHeight := int(float64(monitorH) * 0.70)

	aspectRatio := float64(layoutWidth) / float64(layoutHeight)
	targetWidth := int(float64(targetHeight) * aspectRatio)

	ebiten.SetWindowSize(targetWidth, targetHeight)
	ebiten.SetWindowTitle("Podtank")
	b, t1 := Init()
	if err := ebiten.RunGame(&Game{Battlefield: b, Player1: t1}); err != nil {
		log.Fatal(err)
	}
}
