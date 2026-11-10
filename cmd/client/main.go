package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/UrielMT343/podtank/internal/audio"
	"github.com/UrielMT343/podtank/internal/domain/arena"
	"github.com/UrielMT343/podtank/internal/domain/tank"
	"github.com/UrielMT343/podtank/internal/domain/ui"
	"github.com/UrielMT343/podtank/internal/domain/vfx"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	cellGridWidth  = 50
	cellGridHeight = 50

	arenaGridHeight = 16
	arenaGridWidth  = 13

	layoutWidth  = 1000
	layoutHeight = 900
)

func Start() (*arena.Arena, map[tank.Player]*tank.Tank, map[tank.Player]*tank.BulletPool, *vfx.VFXPool, *audio.SoundBank, *ui.FontBank) {
	b := arena.NewArena(cellGridHeight, cellGridWidth, arenaGridHeight,
		arenaGridWidth, layoutWidth, layoutHeight)

	t1 := tank.NewTank(50, 50, b.X, b.Y, cellGridWidth, cellGridHeight, tank.Down, tank.Player1)

	t2 := tank.NewTank(650, 800, b.X, b.Y, cellGridWidth, cellGridHeight, tank.Up, tank.Player2)

	tanks := map[tank.Player]*tank.Tank{
		tank.Player1: t1,
		tank.Player2: t2,
	}

	pools := map[tank.Player]*tank.BulletPool{
		tank.Player1: tank.NewBulletPool(50, tank.Player1),
		tank.Player2: tank.NewBulletPool(50, tank.Player2),
	}

	vfxPool := vfx.NewVFXPool(50)

	sounds := &audio.SoundBank{
		Shoot:     audio.LoadSound("sfx/shoot.ogg"),
		HitWall:   audio.LoadSound("sfx/hit_wall.ogg"),
		Explosion: audio.LoadSound("sfx/explosion.ogg"),
	}

	fonts := ui.LoadFonts()

	return b, tanks, pools, vfxPool, sounds, fonts
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
	VFXPool     *vfx.VFXPool
	Sounds      *audio.SoundBank
	Fonts       *ui.FontBank
	GameOver    bool
	Winner      string
}

func (g *Game) Update() error {
	if g.GameOver {
		return nil
	}

	localCmd := GetLocalInput()
	if myTank, ok := g.Tanks[g.LocalID]; ok {
		myTank.Update(localCmd, &g.Battlefield.Collisions, g.BulletPools[g.LocalID], g.Sounds)
	}

	remoteCmd := tank.CmdState{}
	enemyID := tank.Player2
	if g.LocalID == tank.Player2 {
		enemyID = tank.Player1
	}

	if enemyTank, ok := g.Tanks[enemyID]; ok {
		enemyTank.Update(remoteCmd, &g.Battlefield.Collisions, g.BulletPools[enemyID], g.Sounds)
	}

	for _, bp := range g.BulletPools {
		bp.Update(&g.Battlefield.Collisions, g.Tanks, g.VFXPool, g.Sounds)
	}

	g.VFXPool.Update()

	for id, t := range g.Tanks {
		if t.IsDead() && !g.GameOver {
			g.GameOver = true
			if g.Sounds != nil && g.Sounds.Explosion != nil {
				audio.PlaySFX(g.Sounds.Explosion)
			}
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

	g.VFXPool.Draw(screen, g.Battlefield.X, g.Battlefield.Y)

	p1HP := g.Tanks[tank.Player1].HP()
	p2HP := g.Tanks[tank.Player2].HP()

	p1Str := fmt.Sprintf("P1 HP:%d", p1HP)
	p2Str := fmt.Sprintf("P2 HP:%d", p2HP)

	p1Opts := &text.DrawOptions{}
	p1Opts.GeoM.Translate(0, 50)
	p1Opts.ColorScale.SetWithColor(color.RGBA{100, 150, 255, 255})
	text.Draw(screen, p1Str, g.Fonts.HUD, p1Opts)

	wP2, _ := text.Measure(p2Str, g.Fonts.HUD, g.Fonts.HUD.Size)
	p2Opts := &text.DrawOptions{}
	p2Opts.GeoM.Translate(float64(layoutWidth)-wP2-0, 50)
	p2Opts.ColorScale.SetWithColor(color.RGBA{255, 100, 100, 255})
	text.Draw(screen, p2Str, g.Fonts.HUD, p2Opts)

	if g.GameOver {
		gameOverStr := "GAME OVER"
		winnerStr := g.Winner + " WINS!"

		wTitle, _ := text.Measure(gameOverStr, g.Fonts.Title, g.Fonts.Title.Size)
		wWinner, _ := text.Measure(winnerStr, g.Fonts.HUD, g.Fonts.HUD.Size)

		xTitle := (float64(layoutWidth) - wTitle) / 2
		xWinner := (float64(layoutWidth) - wWinner) / 2

		titleOpts := &text.DrawOptions{}
		titleOpts.GeoM.Translate(xTitle, float64(layoutHeight)/2-20)
		titleOpts.ColorScale.SetWithColor(color.RGBA{255, 50, 50, 255})
		text.Draw(screen, gameOverStr, g.Fonts.Title, titleOpts)

		winnerOpts := &text.DrawOptions{}
		winnerOpts.GeoM.Translate(xWinner, float64(layoutHeight)/2+30)
		winnerOpts.ColorScale.SetWithColor(color.White)
		text.Draw(screen, winnerStr, g.Fonts.HUD, winnerOpts)
	}
}

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

	b, tanks, pools, vfxPool, sounds, fonts := Start()

	if err := ebiten.RunGame(&Game{
		Battlefield: b,
		Tanks:       tanks,
		LocalID:     tank.Player1,
		BulletPools: pools,
		VFXPool:     vfxPool,
		Sounds:      sounds,
		Fonts:       fonts,
	}); err != nil {
		log.Fatal(err)
	}
}
