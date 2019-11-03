package main

// Player is an interface to a single player's game session.
// It allows any kind of player (human game interfaces, AIs) to be generalised into a single interface.
// Typically a caller would iterate over Turn() to play the game.
type Player interface {
	// GetBoard returns the Player's board.
	GetBoard() *Board

	// Turn takes a turn in the game, returning true if the player won during the turn.
	Turn(Link) (won bool, err error)
}

// Link is an interface, used by a Player, for querying information about the other Player.
// The idea is to provide a single point of comminication between Players, allowing new implementations,
// with new features (i.e. networking), to be made and dropped into existing code.
type Link interface {
	// TakeShot is called when a player takes a shot at the other player, returning if it was a hit or not, and if sunk != 0, she ship that was sunk.
	// Players should keep track of their own score.
	TakeShot(x, y int) (hit bool, sunk byte)
}

// NewLocalLink returns a Link for communicating with the Player p.
func NewLocalLink(p Player) Link {
	return &localLink{
		p: p,
	}
}

type localLink struct {
	p Player
}

// TakeShot implements Link
func (ll *localLink) TakeShot(x, y int) (bool, byte) {
	return ll.p.GetBoard().OpponentShot(x, y)
}
