package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	tm "github.com/buger/goterm"
	ts "github.com/kopoli/go-terminal-size"
)

// x -> column
// y -> line (row)

type World struct {
	ScreenX, ScreenY   int
	PlayerX, PlayerY   int
	Map                [][2]int
	MapCharacter       string
	NextStart, NextEnd int
	Bullets            []Bullet
}

type Bullet struct {
	X, Y     int
	GoingToY int
	ShotBy   string
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
		ScreenX:      screenX,
		ScreenY:      screenY,
		PlayerX:      screenX / 2,
		PlayerY:      screenY - 1,
		Map:          make([][2]int, screenY),
		MapCharacter: " ",
		NextStart:    screenX/2 - 20,
		NextEnd:      screenX/2 + 20,
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

	player = tm.MoveTo(player, world.PlayerX, world.PlayerY)

	tm.Print(player)

	tm.Flush()
}

func physics(world *World, gameRunning *bool) {

	for i := 0; i < len(world.Map); i++ {

		if (world.Map[i][0] >= world.PlayerX || world.Map[i][1] <= world.PlayerX) && world.PlayerY == i {
			*gameRunning = false
		}

		for j := 0; j < len(world.Bullets); j++ {

			if (world.Map[i][0] >= world.Bullets[j].X || world.Map[i][1] <= world.Bullets[j].X) && world.Bullets[j].Y == i {
				world.Bullets = append(world.Bullets[:j], world.Bullets[j+1:]...)
			}

		}

	}

	for i := 0; i < len(world.Bullets); i++ {
		if world.Bullets[i].Y == world.Bullets[i].GoingToY {
			world.Bullets = append(world.Bullets[:i], world.Bullets[i+1:]...)
			continue
		}

		if world.Bullets[i].Y >= world.PlayerY && world.Bullets[i].X == world.PlayerX {
			*gameRunning = false
			continue
		}

		if world.Bullets[i].ShotBy == ShotByPlayer {
			world.Bullets[i].Y--
			continue
		}

		if world.Bullets[i].ShotBy == ShotByEnemy {
			world.Bullets[i].Y++
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

}

func listenPlayerMovement(world *World, gameRunning *bool, screeSize ts.Size) {
	keyboard.Listen(func(key keys.Key) (stop bool, err error) {

		if key.Code == keys.Space {
			world.Bullets = append(world.Bullets, Bullet{X: world.PlayerX, Y: world.PlayerY - 1, GoingToY: 0, ShotBy: ShotByPlayer})
		}

		if key.Code == keys.Right && world.PlayerX < screeSize.Width-2 {

			world.PlayerX += 1

		}

		if key.Code == keys.Left && world.PlayerX > 2 {

			world.PlayerX -= 1

		}

		if key.Code == keys.Up && world.PlayerY > 2 && world.PlayerY >= world.ScreenY/2 {

			world.PlayerY -= 1

		}

		if key.Code == keys.Down && world.PlayerY < screeSize.Height-2 {

			world.PlayerY += 1

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
