package tank

type TankDirection int

const (
	Up TankDirection = iota
	Down
	Left
	Right
)

type Tank struct {
	X, Y          float32
	Direction     TankDirection
	HealthPoints  int
	Damage        int
	UpdgradeLevel int
	Speed         int
}

type Movement interface {
	Update()
	Draw()
}

func (t Tank) Update() {
}

func (t Tank) Draw() {
}
