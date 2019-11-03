package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// Constants for the byte describing a position on the board.
// The first 3 bits are an integer describing the ship type,
// with the remaining bits used as flags for shots taken by the player and opponent.
// Currently only 3 flags are used; the middle 2 bytes are unused.
const (
	playerShot  = 1 << iota // 00000001
	playerHit               // 00000010
	opponentHit             // 00000100

	shipMask       = (1<<8 - 1) &^ (1<<5 - 1) // 11100000
	shipCarrier    = 1 << 5
	shipBattleship = 2 << 5
	shipDestroyer  = 3 << 5
	shipSubmarine  = 4 << 5
	shipPatrolBoat = 5 << 5
)

var shipNames = map[byte]string{
	shipCarrier:    "Carrier",
	shipBattleship: "Battleship",
	shipDestroyer:  "Destroyer",
	shipSubmarine:  "Submarine",
	shipPatrolBoat: "Patrol Boat",
}

// Constants for defining direction,
// i.e. placement of ships.
const (
	up = 1 + iota // 1 offset, so the nil-case for direction is invalid.
	down
	left
	right
)

// Dimensions of the board.
const boardSize = 10

// IsValid returns true if the given coordinates are a location on the board.
func IsValid(x, y int) bool {
	if x < 0 || x >= boardSize {
		return false
	}
	if y < 0 || y >= boardSize {
		return false
	}
	return true
}

// RandomBoard creates a board with randomly placed ships
func RandomBoard(rng *rand.Rand) Board {
	var b Board
	var ships = []byte{
		shipCarrier,
		shipBattleship,
		shipDestroyer,
		shipSubmarine,
		shipPatrolBoat,
	}

	for i := 0; i < len(ships); {
		x := rng.Intn(boardSize)
		y := rng.Intn(boardSize)
		direction := rng.Intn(4) + 1
		if err := b.PlaceShip(x, y, direction, ships[i]); err == nil {
			// sucessful placement, move on to next ship
			i++
		}
	}

	return b
}

// Board holds information about the players board, with a byte describing the state of the given position.
// This does not hold the entire game's state, but rather just one player's boards.
type Board [boardSize][boardSize]byte

// Clear clears the board
func (b *Board) Clear() {
	for x := 0; x < boardSize; x++ {
		for y := 0; y < boardSize; y++ {
			b[x][y] = 0
		}
	}
}

// OpponentShot executes a shot by the opponent, returning if the shot is a hit and if sunk != 0, the ship that was sunk.
// x,y should be checked for validity beforehand.
func (b *Board) OpponentShot(x, y int) (hit bool, sunk byte) {
	b[x][y] |= opponentHit
	ship := b[x][y] & shipMask
	if ship > 0 {
		// ship bits are non-nil; a hit
		if b.IsSunk(x, y) {
			return true, ship
		}
		return true, 0
	}
	return false, 0
}

// PlayerShot records a shot, either playerShot or playerHit
// x,y should be checked for validity beforehand.
func (b *Board) PlayerShot(x, y int, hit bool) {
	if hit {
		b[x][y] |= playerShot | playerHit
	} else {
		b[x][y] |= playerShot
	}
}

// PlayerHasShot returns true if the player has already shot the given position.
// x,y should be checked for validity beforehand.
func (b *Board) PlayerHasShot(x, y int) bool {
	return b[x][y]&playerShot > 0
}

// PlayerHasHit returns true if the player has hit an enemy ship at the given position.
// x,y should be checked for validity beforehand.
func (b *Board) PlayerHasHit(x, y int) bool {
	return b[x][y]&playerHit > 0
}

// String formats the board as a string.
// It implements fmt.Stringer, so directly passing the board to a print call is a valid way of printing the board.
func (b Board) String() (str string) {
	// this function involves a lot of unoptimised string concatentation,
	// and when used to play the game, isn't the most astetically appealing.
	// Improvements can be made here, or ever a graphical soulution substituted here.

	// Top board; show shots by this player
	str += "  1 2 3 4 5 6 7 8 9 10\n"
	// iterate over y in reverse becuase coordinates start at the bottom left, but we print from top left.
	for y := boardSize - 1; y >= 0; y-- {
		str += string('A'+y) + " "
		for x := 0; x < boardSize; x++ {
			switch {
			case b[x][y]&playerHit > 0:
				str += "X "
			case b[x][y]&playerShot > 0:
				str += "O "
			default:
				str += "  "
			}
		}
		str += "\n"
	}
	str += "\n"

	// Bottom board; show player ships and opponent shots
	str += "  1 2 3 4 5 6 7 8 9 10\n"
	for y := boardSize - 1; y >= 0; y-- {
		str += string('A'+y) + " "
		for x := 0; x < boardSize; x++ {
			switch {
			case b[x][y]&shipMask > 0 && b[x][y]&opponentHit > 0:
				// is a ship, and has been hit
				str += "X "
			case b[x][y]&shipMask == shipCarrier:
				str += "C "
			case b[x][y]&shipMask == shipBattleship:
				str += "B "
			case b[x][y]&shipMask == shipDestroyer:
				str += "D "
			case b[x][y]&shipMask == shipSubmarine:
				str += "S "
			case b[x][y]&shipMask == shipPatrolBoat:
				str += "P "
			default:
				str += "  "
			}
		}
		str += "\n"
	}

	return
}

// PlaceShip places a ship on the board. If err is non-nil, the board will have not been modified.
func (b *Board) PlaceShip(x, y, direction int, shipType byte) error {
	// mx and my make a vector, added to x,y length times to traverse positions on the ship.
	var length, mx, my int

	switch shipType {
	case shipCarrier:
		length = 5
	case shipBattleship:
		length = 4
	case shipDestroyer:
		length = 3
	case shipSubmarine:
		length = 3
	case shipPatrolBoat:
		length = 2
	default:
		return errors.New("invalid ship type")
	}

	switch direction {
	case left:
		mx = -1
	case right:
		mx = 1
	case down:
		my = -1
	case up:
		my = 1
	default:
		return errors.New("invalid direction")
	}

	// check we're not overwriting another ship, and that all positions are valid.
	// we must be certain everything is fine before we modify the board, else we leave half-written ships on it.
	ix, iy := x, y // can't modify x and y yet, we need it later.
	for i := 0; i < length; i++ {
		if !IsValid(ix, iy) {
			return errors.New("ship is off the board")
		}
		if b[ix][iy]&shipMask > 0 {
			return fmt.Errorf("there is already a ship at %v", FormatPosition(ix, iy))
		}
		ix += mx
		iy += my
	}

	// Everything is good to go, place the ship on the board.
	for i := 0; i < length; i++ {
		b[x][y] = shipType
		x += mx
		y += my
	}

	return nil
}

// IsSunk returns true if the ship at the location x, y has been sunk.
// If x,y is not on a ship, it returns false.
// x,y should be checked for validity beforehand.
func (b *Board) IsSunk(x, y int) bool {
	shipType := b[x][y] & shipMask
	if shipType == 0 { // Not a ship
		return false
	}

	// Figure out if the ship is placed horizontally or verically, and traverse its positions.
	var vertical bool
	// Check if the position above and below are of the same ship type
	// if they are, we know the ship is placed vertically.
	if y+1 < boardSize && b[x][y+1]&shipMask == shipType ||
		y-1 >= 0 && b[x][y-1]&shipMask == shipType {
		vertical = true
	}

	// check positions in positive direction, including x,y
	ix, iy := x, y
	for i := 0; ; i++ {
		if !IsValid(ix, iy) {
			// we're off the board
			break
		}

		// travel in the positive direction
		if b[ix][iy]&shipMask != shipType {
			// we're no longer on the ship.
			break
		}

		if b[ix][iy]&opponentHit == 0 {
			// ship is not hit at this location.
			return false
		}

		if vertical {
			iy++
		} else {
			ix++
		}
	}

	// traverse in negative direction, skipping x,y
	ix, iy = x, y
	for i := 0; ; i++ {
		// We've already checked the starting position in the last loop,
		// so do the increment here and check from the next position onwards.
		if vertical {
			iy--
		} else {
			ix--
		}

		if !IsValid(ix, iy) {
			// we're off the board
			break
		}

		// travel in the negative direction
		if b[ix][iy]&shipMask != shipType {
			// we're no longer on the ship.
			break
		}

		if b[ix][iy]&opponentHit == 0 {
			// ship is not hit at this location.
			return false
		}
	}

	// No positions on the ship were found to have not been hit
	return true
}

// ParsePosition takes a simple two character string and converts it into x, y coordinates on the board.
// i.e. "b6" = 2, 6
func ParsePosition(location string) (x int, y int, err error) {
	if location == "" {
		return 0, 0, errors.New("no location specified")
	}

	location = strings.TrimSpace(strings.ToLower(location))
	y = int(location[0] - 'a')          // first byte of string, offset by the ascii code for 'a'
	x, err = strconv.Atoi(location[1:]) // string->int conversion for remaining string bytes.
	if err != nil {
		return 0, 0, fmt.Errorf("invalid location %v", location)
	}
	x-- // our printed position format starts at 1, but coordinates start at 0.
	if !IsValid(x, y) {
		return 0, 0, fmt.Errorf("invalid location %v", location)
	}
	return
}

// FormatPosition takes an x and y coordinate, and returns the string representation of it.
// i.e. 2, 6 = "b6"
func FormatPosition(x, y int) string {
	return string('a'+rune(y)) + strconv.Itoa(x+1)
}
