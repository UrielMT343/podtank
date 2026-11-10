package ui

import (
	"bytes"

	"github.com/UrielMT343/podtank/assets"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type FontBank struct {
	HUD   *text.GoTextFace
	Title *text.GoTextFace
}

func LoadFonts() *FontBank {
	fontData, err := assets.FS.ReadFile("fonts/PressStart2P-Regular.ttf")
	if err != nil {
		panic(err)
	}

	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		panic(err)
	}

	return &FontBank{
		HUD: &text.GoTextFace{
			Source: fontSource,
			Size:   16,
		},
		Title: &text.GoTextFace{
			Source: fontSource,
			Size:   32,
		},
	}
}
