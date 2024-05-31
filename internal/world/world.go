package world

import (
	"github.com/amirdaraby/sky-killer/internal/bullet"
	"github.com/amirdaraby/sky-killer/internal/enemy"
	"github.com/amirdaraby/sky-killer/internal/player"
)

// terminal: column -> X , line -> Y
// CurrentMatrix: map is an slice of arrays which have two value, first edge of river and second edge of river (both values are 'X')
type World struct {
	X, Y               int
	CurrentMatrix      [][2]int
	NextStart, NextEnd int
	SandCharacter      string
	RiverCharacter     string
	Player             *player.Player
	Enemies            []enemy.Enemy
	Bullets            []bullet.Bullet
}