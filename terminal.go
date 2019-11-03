package main

import (
	"bufio"
	"fmt"
	"strings"
)

// NewTerminalUI creates a new terminal game session for a human player.
func NewTerminalUI(input *bufio.Reader) *TerminalUI {
	return &TerminalUI{
		input: input,
	}
}

// TerminalUI is a terminal session of battleship.
type TerminalUI struct {
	board Board
	score int

	// Reader for user input
	input *bufio.Reader
}

// GetBoard implements Player.
func (g *TerminalUI) GetBoard() *Board {
	return &g.board
}

// Turn implements Player.
// it asks the player to take a turn and executes it.
func (g *TerminalUI) Turn(remote Link) (won bool, err error) {
	fmt.Println("Current score", g.score)
	fmt.Print(g.board)

	var x, y int
	// Loop until we have a good location
	for {
		fmt.Println("Enter shot location (h for help)")
		str, err := g.input.ReadString('\n')
		if err != nil {
			return false, err
		}
		str = strings.ToLower(strings.TrimSpace(str))

		if str == "h" {
			fmt.Println("Syntax: [location]")
			fmt.Println("location is a-j for vertical position, 1-10 for horizontal position. \n i.e. g6")
			continue
		}

		x, y, err = ParsePosition(str)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if g.board.PlayerHasShot(x, y) {
			fmt.Println("You've already shot that location!")
			continue
		}

		break
	}

	hit, sunk := remote.TakeShot(x, y)
	g.board.PlayerShot(x, y, hit)
	if hit {
		fmt.Println("Hit!")
		if sunk != 0 {
			fmt.Printf("You sunk their %v!\n", shipNames[sunk])
			g.score++
		}
	} else {
		fmt.Println("Miss!")
	}

	fmt.Println("Press enter to finish turn")
	g.input.ReadString('\n')

	return g.score >= 5, nil
}

// SetUp asks the user to place their ships on the board, writing them to it.
func (g *TerminalUI) SetUp() error {
	var ships = []byte{
		shipCarrier,
		shipBattleship,
		shipDestroyer,
		shipSubmarine,
		shipPatrolBoat,
	}

	// repeatedly ask for location and direction of ship placement until a sucessful position is given.
askAgain:
	for i := 0; i < len(ships); {
		fmt.Print(g.board)
		fmt.Printf("Enter %v Location and direction. (h for help)\n", shipNames[ships[i]])
		str, err := g.input.ReadString('\n')
		if err != nil {
			return err
		}
		str = strings.ToLower(strings.TrimSpace(str))

		if str == "h" {
			fmt.Print(`Syntax: [location] [direction]
Possible directions are "up", "down", "left", "right"
location is a-j for vertical position, 1-10 for horizontal position.
i.e h4 down
`)
			continue
		}

		// get our two arguments; position and direction.
		args := strings.Split(str, " ")
		if len(args) != 2 {
			fmt.Println("wrong number of arguments; can only take 2")
			continue
		}

		x, y, err := ParsePosition(args[0])
		if err != nil {
			fmt.Println(err)
			continue
		}

		var direction int
		switch args[1] {
		case "up":
			direction = up
		case "down":
			direction = down
		case "left":
			direction = left
		case "right":
			direction = right
		default:
			fmt.Printf("unknown direction %v\n", direction)
			continue askAgain
		}

		err = g.board.PlaceShip(x, y, direction, ships[i])
		if err != nil {
			fmt.Println(err)
			continue
		}

		// sucess, move to the next ship
		i++
	}

	fmt.Println(g.board)
	fmt.Println("All ships placed")
	return nil
}
