package arena

import (
	_ "embed"
)

//go:embed maps/map1.txt
var M1 []byte

//go:embed maps/map2.txt
var M2 []byte

//go:embed maps/map3.txt
var M3 []byte
