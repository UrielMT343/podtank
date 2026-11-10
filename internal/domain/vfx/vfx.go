package vfx

import (
	"github.com/UrielMT343/podtank/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type EffectType int

const (
	EffectSmoke EffectType = iota
	EffectExplosion
	EffectMuzzle
)

type VFX struct {
	X, Y       float64
	Type       EffectType
	TicksAlive int
}

type VFXPool struct {
	pool            []*VFX
	active          []*VFX
	smokeFrames     []*ebiten.Image
	explosionFrames []*ebiten.Image
}

func NewVFXPool(size int) *VFXPool {
	vp := &VFXPool{
		pool:   make([]*VFX, 0, size),
		active: make([]*VFX, 0, size),
	}

	for i := 0; i < size; i++ {
		vp.pool = append(vp.pool, &VFX{})
	}

	loadSprite := func(path string) *ebiten.Image {
		img, _, err := ebitenutil.NewImageFromFileSystem(assets.FS, path)
		if err != nil {
			panic(err)
		}
		return img
	}

	vp.smokeFrames = []*ebiten.Image{
		loadSprite("sprites/explosionSmoke1.png"),
		loadSprite("sprites/explosionSmoke2.png"),
		loadSprite("sprites/explosionSmoke3.png"),
		loadSprite("sprites/explosionSmoke4.png"),
		loadSprite("sprites/explosionSmoke5.png"),
	}

	vp.explosionFrames = []*ebiten.Image{
		loadSprite("sprites/explosion2.png"),
		loadSprite("sprites/explosion3.png"),
		loadSprite("sprites/explosion4.png"),
		loadSprite("sprites/explosion5.png"),
	}

	return vp
}

func (vp *VFXPool) Spawn(x, y float64, effectType EffectType) {
	if len(vp.pool) == 0 {
		return
	}

	lastIdx := len(vp.pool) - 1
	v := vp.pool[lastIdx]
	vp.pool = vp.pool[:lastIdx]

	v.X = x
	v.Y = y
	v.Type = effectType
	v.TicksAlive = 0

	vp.active = append(vp.active, v)
}

func (vp *VFXPool) Update() {
	ticksPerFrame := 4

	for i := 0; i < len(vp.active); i++ {
		v := vp.active[i]
		v.TicksAlive++

		var totalFrames int
		switch v.Type {
		case EffectSmoke:
			totalFrames = len(vp.smokeFrames)
		case EffectExplosion:
			totalFrames = len(vp.explosionFrames)
		case EffectMuzzle:
			totalFrames = 1
		}

		if (v.TicksAlive / ticksPerFrame) >= totalFrames {
			vp.pool = append(vp.pool, v)
			lastIdx := len(vp.active) - 1
			vp.active[i] = vp.active[lastIdx]
			vp.active = vp.active[:lastIdx]
			i--
		}
	}
}

func (vp *VFXPool) Draw(screen *ebiten.Image, offsetX, offsetY float64) {
	ticksPerFrame := 4

	for _, v := range vp.active {
		frameIdx := v.TicksAlive / ticksPerFrame
		var img *ebiten.Image

		switch v.Type {
		case EffectSmoke:
			if frameIdx < len(vp.smokeFrames) {
				img = vp.smokeFrames[frameIdx]
			}
		case EffectExplosion:
			if frameIdx < len(vp.explosionFrames) {
				img = vp.explosionFrames[frameIdx]
			}
		}

		if img == nil {
			continue
		}

		bounds := img.Bounds().Size()
		halfW := float64(bounds.X) / 2.0
		halfH := float64(bounds.Y) / 2.0

		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(-halfW, -halfH)
		options.GeoM.Translate(v.X+offsetX, v.Y+offsetY)

		screen.DrawImage(img, options)
	}
}
