package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func init() {
	flag.BoolVar(&hideAI, "no-show-ai", false, "no-show-ai hides the AI's board during gameplay")
}

func main() {
	flag.Parse()

	input := bufio.NewReader(os.Stdin)
	player1, player2, err := gameSetup(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p1Link := NewLocalLink(player1)
	p2Link := NewLocalLink(player2)

	for {
		fmt.Println("Player 1 Turn")
		won, err := player1.Turn(p2Link)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if won {
			fmt.Println("Player 1 Won!")
			return
		}

		fmt.Println("Player 2 Turn")
		won, err = player2.Turn(p1Link)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if won {
			fmt.Println("Player 2 Won!")
			return
		}
	}
}

// gameSetup sets up the game as per user preference,
// returning an AI, TerminalUI, or some combination of the two.
func gameSetup(input *bufio.Reader) (p1 Player, p2 Player, err error) {
	fmt.Println("Player 1:")
	p1, err = askAndCreatePlayer(input)
	if err != nil {
		return
	}
	fmt.Println("Player 2:")
	p2, err = askAndCreatePlayer(input)
	return
}

// askAndCreatePlayer asks the user what kind of player they want to create, and creates it
func askAndCreatePlayer(input *bufio.Reader) (Player, error) {
	var err error
	var str string

	for {
		fmt.Println("Enter \"ai\" or \"player\"")
		str, err = input.ReadString('\n')
		if err != nil {
			return nil, err
		}
		str = strings.ToLower(strings.TrimSpace(str))

		if str == "ai" {
			return NewAI(), nil
		}
		if str == "player" {
			tui := NewTerminalUI(input)
			tui.SetUp()
			return tui, nil
		}
	}
}
