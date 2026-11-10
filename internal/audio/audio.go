package audio

import (
	"github.com/UrielMT343/podtank/assets"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

var AudioCtx = audio.NewContext(44100)

type SoundBank struct {
	Shoot     *audio.Player
	HitWall   *audio.Player
	Explosion *audio.Player
}

func LoadSound(path string) *audio.Player {
	file, err := assets.FS.Open(path)
	if err != nil {
		panic(err)
	}

	decoded, err := vorbis.DecodeWithoutResampling(file)
	if err != nil {
		panic(err)
	}

	player, err := AudioCtx.NewPlayer(decoded)

	if err != nil {
		panic(err)
	}

	return player
}

func PlaySFX(p *audio.Player) {
	if p != nil {
		p.Rewind()
		p.Play()
	}
}
