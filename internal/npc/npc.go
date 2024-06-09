package npc

type NPC interface {
	X() int
	Y() int
	MoveTo(int, int)
	StepUp()
	StepDown()
	StepLeft()
	StepRight()
	// Shoot()
}

type Enemies struct {
	Units       []NPC
	ShootsFired int
}
