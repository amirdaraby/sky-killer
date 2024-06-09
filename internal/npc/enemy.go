package npc

type NormalEnemy struct {
	x, y                 int
	MovingToX, MovingToY int
	Reverse              bool
}

// todo do the moving enemy logic here, this give ability to have different enemy interactions
func (e *NormalEnemy) MoveTo(X int, Y int) {
	e.x, e.y = X, Y
}

func NewNormalEnemy(X, Y int) *NormalEnemy {
	return &NormalEnemy{
		x: X,
		y: Y,
	}
}

func (e NormalEnemy) X() int {
	return e.x
}

func (e NormalEnemy) Y() int {
	return e.y
}

func (e *NormalEnemy) StepUp() {
	e.y = e.y - 1
}

func (e *NormalEnemy) StepDown() {
	e.y = e.y + 1
}

func (e *NormalEnemy) StepLeft() {
	e.x = e.x - 1
}

func (e *NormalEnemy) StepRight() {
	e.x = e.x + 1
}
