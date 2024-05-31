package player

// TODO
// const (
// 	playerCharacter int = ' '
// 	playerMoveForward = 'w'
// 	playerMoveLeft = 'a'
// 	playerMoveRight = 'd'
// 	playerMoveBack = 's'
// 	playerShoot = 'v'
// )

type Player struct {
	Alive     bool
	X, Y      int
	Character string
}

