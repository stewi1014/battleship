package main

import (
	"fmt"
	"math/rand"
	"time"
)

// NewAI returns a new AI with randomly placed ships, and using the current Unix time as a rng seed.
func NewAI() *AI {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &AI{
		rng:   rng,
		board: RandomBoard(rng),
	}
}

// AI is a simple battleship-playing AI.
type AI struct {
	board Board
	score int

	rng *rand.Rand
}

// showAI will show the AIs board during gameplay
var hideAI = false

// GetBoard implements Player.
func (a *AI) GetBoard() *Board {
	return &a.board
}

// Turn implements Player.
func (a *AI) Turn(remote Link) (won bool, err error) {
	// print board after we return if hideAI is false
	defer func() {
		if !hideAI {
			fmt.Println("AI board")
			fmt.Print(a.board)
		}
	}()

	// try to hit a previously hit ship.
	for x := 0; x < boardSize; x++ {
		for y := 0; y < boardSize; y++ {
			if shootx, shooty, ok := a.findShot(x, y); ok {
				a.shoot(shootx, shooty, remote)
				return a.score >= 5, nil
			}
		}
	}

	// no luck, take a random shot
	shootx, shooty := a.getRandomShot()
	a.shoot(shootx, shooty, remote)
	return a.score >= 5, nil
}

// findShot checks if a point on the board is on a previous hit, and if so,
// tries to hit the ship again. If sucessful, ok is true.
func (a *AI) findShot(x, y int) (shootx int, shooty int, ok bool) {
	if !a.board.PlayerHasHit(x, y) {
		// this point hasn't been hit
		return 0, 0, false
	}

	// we've hit this point before. Check surrounding points to see if we can figure out if the ship is placed vertically or horizontally.
	var verticalShip, horizontalShip bool
	// left and right
	if x-1 >= 0 && a.board.PlayerHasHit(x-1, y) ||
		x+1 < boardSize && a.board.PlayerHasHit(x+1, y) {
		// we've hit a point left or right of the position, it might be a horizontally placed ship.
		shootx, ok = a.findHorizontalShot(x, y)
		if ok {
			return shootx, y, true
		}

		horizontalShip = true
	}
	// up and down
	if y-1 >= 0 && a.board.PlayerHasHit(x, y-1) ||
		y+1 < boardSize && a.board.PlayerHasHit(x, y+1) {
		// we've hit a point above or below the position, it might be a vertically placed ship.
		shooty, ok = a.findVerticalShot(x, y)
		if ok {
			return x, shooty, true
		}

		verticalShip = true
	}

	if horizontalShip && verticalShip {
		// we found hits in a horizontal and vertical line.
		// it's possible two vertical, or two horizontal ships are next to eachother.
		// if so, we need to hit the diagonals.
		shootx, shooty, ok := a.findDiagonalShot(x, y)
		if ok {
			return shootx, shooty, true
		}
	}

	if horizontalShip || verticalShip {
		// we've already handled cases for points found in a line, so if a line is what was found,
		// we know there's nothing to do at this point.
		return 0, 0, false
	}

	// the point is hit, but no line was found.
	// try hitting adjacent points.
	return a.findAdjacentShot(x, y)
}

// findHorizontalShot tries to take a shot at a possibly horizontally placed ship at the given position.
// it returns true if a shot is found.
func (a *AI) findHorizontalShot(x, y int) (int, bool) {
	// check right
	for ix := x; ix < boardSize; ix++ {
		if a.board.PlayerHasHit(ix, y) {
			// we've hit this point before
			continue
		}
		if a.board.PlayerHasShot(ix, y) {
			// we've missed at this position.
			break
		}
		// found a point to shoot
		return ix, true
	}

	// check left; same thing in reverse
	for ix := x; ix >= 0; ix-- {
		if a.board.PlayerHasHit(ix, y) {
			continue
		}
		if a.board.PlayerHasShot(ix, y) {
			break
		}
		return ix, true
	}

	return 0, false
}

// findVerticalShot tries to find a shot at a possibly vertically placed ship at the given position.
// it returns true if a shot is found.
func (a *AI) findVerticalShot(x, y int) (int, bool) {
	for iy := y; iy < boardSize; iy++ {
		if a.board.PlayerHasHit(x, iy) {
			continue
		}
		if a.board.PlayerHasShot(x, iy) {
			break
		}
		return iy, true
	}

	for iy := y; iy >= 0; iy-- {
		if a.board.PlayerHasHit(x, iy) {
			continue
		}
		if a.board.PlayerHasShot(x, iy) {
			break
		}
		return iy, true
	}

	return 0, false
}

// findAdjacentShot searches for a shot in the positions next to the given position.
// it returns true if a shot was found.
func (a *AI) findAdjacentShot(x, y int) (int, int, bool) {
	//up
	if y+1 < boardSize && !a.board.PlayerHasShot(x, y+1) {
		return x, y + 1, true
	}
	// down
	if y-1 >= 0 && !a.board.PlayerHasShot(x, y-1) {
		return x, y - 1, true
	}
	// left
	if x-1 >= 0 && !a.board.PlayerHasShot(x-1, y) {
		return x - 1, y, true
	}
	// right
	if x+1 < boardSize && !a.board.PlayerHasShot(x+1, y) {
		return x + 1, y, true
	}

	// all adjacent points are already hit
	return 0, 0, false
}

// findDiagonalShot searches for a shot in the positions immediately diagonal to the given position.
// it returns true if a shot was found.
func (a *AI) findDiagonalShot(x, y int) (int, int, bool) {
	// bottom left
	if x-1 >= 0 && y-1 >= 0 && !a.board.PlayerHasShot(x-1, y-1) {
		return x - 1, y - 1, true
	}
	// bottom right
	if x+1 < boardSize && y-1 >= 0 && !a.board.PlayerHasShot(x+1, y-1) {
		return x + 1, y - 1, true
	}
	// top right
	if x+1 < boardSize && y+1 < boardSize && !a.board.PlayerHasShot(x+1, y+1) {
		return x + 1, y + 1, true
	}
	// top left
	if x-1 >= 0 && y+1 < boardSize && !a.board.PlayerHasShot(x-1, y+1) {
		return x - 1, y + 1, true
	}

	// all diagonal points are already hit
	return 0, 0, false
}

func (a *AI) getRandomShot() (x int, y int) {
	for {
		x = a.rng.Intn(boardSize)
		y = a.rng.Intn(boardSize)
		if !a.board.PlayerHasShot(x, y) {
			break
		}
	}
	return
}

func (a *AI) shoot(x, y int, remote Link) {
	if a.board.PlayerHasShot(x, y) {
		panic("ai tried to hit point it already shot!")
	}

	hit, sunk := remote.TakeShot(x, y)
	a.board.PlayerShot(x, y, hit)
	if sunk != 0 {
		a.score++
	}
}
