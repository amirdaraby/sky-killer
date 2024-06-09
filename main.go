package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/amirdaraby/sky-killer/internal/npc"
	tm "github.com/buger/goterm"
	ts "github.com/kopoli/go-terminal-size"
)

// x -> column
// y -> line (row)

type World struct {
	Enemies            npc.Enemies
	ScreenX, ScreenY   int
	Player             *Player
	Map                [][2]int
	MapCharacter       string
	NextStart, NextEnd int
	Bullets            []Bullet
	MaxEnemyInScreen   int
}

type Bullet struct {
	X, Y     int
	GoingToY int
	ShotBy   string
}

type Player struct {
	X, Y int
}

const (
	ShotByPlayer = "shot_by_player"
	ShotByEnemy  = "shot_by_enemy"
)

func main() {
	// initialize

	screenSize, err := ts.GetSize()

	if err != nil {
		panic(err)
	}

	var screenX int = screenSize.Width
	var screenY int = screenSize.Height

	world := World{
		ScreenX: screenX,
		ScreenY: screenY,
		Player: &Player{
			X: screenX / 2,
			Y: screenY - 1,
		},
		Map:              make([][2]int, screenY),
		MapCharacter:     " ",
		NextStart:        screenX/2 - 20,
		NextEnd:          screenX/2 + 20,
		MaxEnemyInScreen: 5,
	}

	for i := range world.Map {
		world.Map[i] = [2]int{(screenX / 2) - 10, (screenX / 2) + 10}
	}

	cursor.Hide()

	// game is running right now, this becomes false when escape button pressed
	gameRunning := true

	go listenPlayerMovement(&world, &gameRunning, screenSize)

	for gameRunning {
		time.Sleep(time.Millisecond * 100)
		physics(&world, &gameRunning)
		draw(&world)
	}

	tm.Clear()

	tm.Print("Thanks for playing <3")

	tm.Flush()

	cursor.Show()
}

func draw(world *World) {

	tm.Clear()

	for i := 0; i < world.ScreenY; i++ {
		// draw river
		tm.Print(tm.MoveTo(tm.Background(strings.Repeat(" ", world.ScreenX-(world.Map[i][0]+(world.ScreenX-world.Map[i][1]))), tm.BLUE), world.Map[i][0], i))

		// draw river edge
		tm.Print(tm.MoveTo(tm.Background(strings.Repeat(world.MapCharacter, world.Map[i][0]), tm.GREEN), 0, i))
		tm.Print(tm.MoveTo(tm.Background(strings.Repeat(world.MapCharacter, world.ScreenX-world.Map[i][1]), tm.GREEN), world.Map[i][1], i))

	}

	for i := 0; i < len(world.Bullets); i++ {
		tm.Print(tm.MoveTo(tm.Background("|", tm.CYAN), world.Bullets[i].X, world.Bullets[i].Y))
	}

	// draw player
	player := tm.Background(" ", tm.RED)

	player = tm.MoveTo(player, world.Player.X, world.Player.Y)

	tm.Print(player)

	// draw enemies
	for i := 0; i < len(world.Enemies.Units); i++ {
		enemy := tm.Background(" ", tm.BLACK)

		enemy = tm.MoveTo(enemy, world.Enemies.Units[i].X(), world.Enemies.Units[i].Y())

		tm.Print(enemy)
	}

	tm.Flush()
}

func physics(world *World, gameRunning *bool) {

	for i := 0; i < len(world.Map); i++ {

		if (world.Map[i][0] >= world.Player.X || world.Map[i][1] <= world.Player.X) && world.Player.Y == i {
			*gameRunning = false
		}

		for j := 0; j < len(world.Bullets); j++ {

			if (world.Map[i][0] >= world.Bullets[j].X || world.Map[i][1] <= world.Bullets[j].X) && world.Bullets[j].Y == i {
				world.Bullets = append(world.Bullets[:j], world.Bullets[j+1:]...)
			}

		}

	}

	for i := 0; i < len(world.Bullets); i++ {

		for j := 0; j < len(world.Enemies.Units); j++ {
			if (world.Bullets[i].Y == world.Enemies.Units[j].Y() || world.Bullets[i].Y == world.Enemies.Units[j].Y()+1) && world.Bullets[i].X == world.Enemies.Units[j].X() {
				world.Enemies.Units = append(world.Enemies.Units[:j], world.Enemies.Units[j+1:]...)
				world.Bullets = append(world.Bullets[:i], world.Bullets[i+1:]...)
				continue
			}
		}

		if world.Bullets[i].Y == world.Bullets[i].GoingToY || world.Bullets[i].Y >= world.ScreenY || world.Bullets[i].Y <= 1 {
			world.Bullets = append(world.Bullets[:i], world.Bullets[i+1:]...)
			continue
		}

		if (world.Bullets[i].Y == world.Player.Y || world.Bullets[i].Y == world.Player.Y+1) && world.Bullets[i].X == world.Player.X {
			*gameRunning = false
			break
		}

		if world.Bullets[i].ShotBy == ShotByPlayer {
			world.Bullets[i].Y -= 2
			continue
		} else {
			world.Bullets[i].Y += 2
		}
	}

	// shift the map
	for i := len(world.Map) - 2; i >= 0; i-- {
		world.Map[i+1] = world.Map[i]
	}

	// randomize map
	if world.NextEnd < world.Map[0][1] {
		world.Map[0][1] -= 1
	}

	if world.NextEnd > world.Map[0][1] {
		world.Map[0][1] += 1
	}

	if world.NextStart < world.Map[0][0] {
		world.Map[0][0] -= 1
	}

	if world.NextStart > world.Map[0][0] {
		world.Map[0][0] += 1
	}

	if world.NextStart == world.Map[0][0] && world.NextEnd == world.Map[0][1] {

		if randRange(0, 4) == 1 {

			world.NextStart = randRange(world.ScreenX/2-(world.ScreenX/6), randRange(world.ScreenX/2-(world.ScreenX/6)+1, world.ScreenX-10))
			world.NextEnd = randRange(world.NextStart, world.ScreenX-10)

			if world.NextEnd-world.NextStart <= 15 {
				world.NextStart -= 15
			}

		}

	}

	// spawn enemies and move them ?
	for i := 0; i < len(world.Enemies.Units); i++ {

		if world.Enemies.Units[i].Y() >= world.ScreenY-1 {
			world.Enemies.Units = append(world.Enemies.Units[i:], world.Enemies.Units[i+1:]...)
			continue
		}

		if world.Enemies.ShootsFired <= 30 {
			if randRange(0, 20) == 5 {
				world.Bullets = append(world.Bullets, Bullet{X: world.Enemies.Units[i].X(), Y: world.Enemies.Units[i].Y() + 4, GoingToY: world.ScreenY, ShotBy: ShotByEnemy})
				world.Enemies.ShootsFired++
			}
		}

		world.Enemies.Units[i].StepDown()
	}

	if randRange(0, 10) == 2 {

		if len(world.Enemies.Units) <= world.MaxEnemyInScreen {
			NewEnemy := npc.NewNormalEnemy(randRange(world.Map[0][0], world.Map[0][1]), 1)

			world.Enemies.Units = append(world.Enemies.Units, NewEnemy)
		}
	}

}

func listenPlayerMovement(world *World, gameRunning *bool, screeSize ts.Size) {
	keyboard.Listen(func(key keys.Key) (stop bool, err error) {

		if key.Code == keys.Space {
			world.Bullets = append(world.Bullets, Bullet{X: world.Player.X, Y: world.Player.Y - 1, GoingToY: 0, ShotBy: ShotByPlayer})
		}

		if key.Code == keys.Right && world.Player.X < screeSize.Width-2 {

			world.Player.X += 1

		}

		if key.Code == keys.Left && world.Player.X > 2 {

			world.Player.X -= 1

		}

		if key.Code == keys.Up && world.Player.Y > 2 && world.Player.Y >= world.ScreenY/2 {

			world.Player.Y -= 1

		}

		if key.Code == keys.Down && world.Player.Y < screeSize.Height-2 {

			world.Player.Y += 1

		}

		if key.Code == keys.Escape {
			*gameRunning = false
		}

		return false, nil // Return false to continue listening
	})
}

func randRange(min, max int) int {

	defer recoverIntn(min, max)

	return rand.Intn(max-min) + min
}

func recoverIntn(min, max int) {
	if r := recover(); r != nil {
		panic(fmt.Sprintf("min: %d \nmax: %d", min, max))
	}
}
