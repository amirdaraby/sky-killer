package enemy

type Enemy struct {
	Alive            bool
	X, Y             int
	MovingX, MovingY int
	ReverseMovement  bool
	Character        string
}
