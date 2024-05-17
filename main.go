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
	ScreenX, ScreenY int
	PlayerX, PlayerY int
	Map              [][2]int
	MapCharacter     string
}

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

	fmt.Printf("%+v", world.Map)
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

	// draw player
	player := tm.Background("P", tm.RED)

	player = tm.MoveTo(player, world.PlayerX, world.PlayerY)

	tm.Print(player)

	tm.Flush()
}

func physics(world *World, gameRunning *bool) {

	for i := 0; i < len(world.Map); i++ {

		if (world.Map[i][0] >= world.PlayerX || world.Map[i][1] <= world.PlayerX) && world.PlayerY == i {
			*gameRunning = false
		}

	}

	//shift the map
	for i := len(world.Map) - 2; i >= 0; i-- {
		world.Map[i+1] = world.Map[i]
	}

	world.Map[0] = [2]int{((world.ScreenX / 2) - randRange(4, 10)), ((world.ScreenX / 2) + randRange(4, 15))}
}

func listenPlayerMovement(world *World, gameRunning *bool, screeSize ts.Size) {
	keyboard.Listen(func(key keys.Key) (stop bool, err error) {

		if key.Code == keys.Right && world.PlayerX < screeSize.Width-2 {

			world.PlayerX += 1

		}

		if key.Code == keys.Left && world.PlayerX > 2 {

			world.PlayerX -= 1

		}

		if key.Code == keys.Up && world.PlayerY > 2 {

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
	return rand.Intn(max-min) + min
}
