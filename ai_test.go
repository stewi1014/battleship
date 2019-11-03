package main

import (
	"fmt"
	"os"
	"testing"
)

// Testing of the AI is limited to checking for runtime errors.
// Unless the AI actually crashes the program or enters an infinite loop,
// and as long as the gameplay is enjoyable (tested), errors within the AI simply give the player an advantage.
// Apart from that, there's no "correct" way to play the game either.

func TestMain(m *testing.M) {
	hideAI = true // force AI to hide its board during tests.
	os.Exit(m.Run())
}

// With this passing, we can say with a reasonably high degree of accuracy that the AI is stable.
// For a game, this should suffice.
func TestAI(t *testing.T) {
	tests := 10000
	maxTurns := boardSize * boardSize

newGame:
	for i := 0; i < tests; i++ {
		ai1, ai2 := NewAI(), NewAI() // boards are randomly generated
		l1, l2 := NewLocalLink(ai1), NewLocalLink(ai2)

		for turns := 0; ; turns++ {
			won, err := ai1.Turn(l2)
			if err != nil {
				t.Fatal(err)
			}
			if won {
				continue newGame
			}

			won, err = ai2.Turn(l1)
			if err != nil {
				t.Fatal(err)
			}
			if won {
				continue newGame
			}

			if turns > maxTurns {
				t.Fatalf("Maximum number of turns(%v) reached", maxTurns)
				fmt.Println("ai1 board: \n", ai1.board)
				fmt.Println("ai2 board: \n", ai2.board)
			}
		}
	}
}
