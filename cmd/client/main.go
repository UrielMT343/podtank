package main

import (
	"log"

	"github.com/UrielMT343/podtank/internal/domain/arena"
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

func Init() (battlefield *arena.Arena, errInit error) {
	b := arena.NewArena(cellGridHeight, cellGridWidth, arenaGridHeight,
		arenaGridWidth, layoutWidth, layoutHeight)

	return &b, nil
}

type Game struct {
	Battlefield *arena.Arena
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Test")
	screen.Bounds()
	g.Battlefield.Draw(screen)
}

// Takes the originla size and converts it to a more logical size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return layoutWidth, layoutHeight
}

func main() {
	ebiten.SetWindowSize(layoutWidth, layoutHeight)
	ebiten.SetWindowTitle("Podtank")
	b, _ := Init()
	if err := ebiten.RunGame(&Game{Battlefield: b}); err != nil {
		log.Fatal(err)
	}
}
