package arena

import (
	_ "embed"
)

//go:embed map1.txt
var m1 string

func ReadMap1() string {
	return m1
}
