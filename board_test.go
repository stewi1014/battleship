package main

import (
	"testing"
)

// It's mostly possible to test this program through use.
// Many different invalid inputs have been tested, and the game played through with player vs player, player vs ai, and ai vs ai.

func TestPlaceShip(t *testing.T) {
	type args struct {
		x, y      int
		direction int
		shipType  byte
	}
	testCases := []struct {
		desc       string
		startBoard Board
		wantBoard  Board
		wantError  bool // if wantError is true, the modified board is checked to be unmodified; error cases must not modify the board.
		args       args
	}{
		{
			desc: "Destroyer placed upwards",
			args: args{
				x:         4,
				y:         2,
				direction: up,
				shipType:  shipDestroyer,
			},
			wantBoard: func() Board {
				var b Board
				b[4][2] = shipDestroyer
				b[4][3] = shipDestroyer
				b[4][4] = shipDestroyer
				return b
			}(),
		},
		{
			desc: "Carrier placed downwards",
			args: args{
				x:         8,
				y:         5,
				direction: down,
				shipType:  shipCarrier,
			},
			wantBoard: func() Board {
				var b Board
				b[8][5] = shipCarrier
				b[8][4] = shipCarrier
				b[8][3] = shipCarrier
				b[8][2] = shipCarrier
				b[8][1] = shipCarrier
				return b
			}(),
		},
		{
			desc: "Battleship placed leftwards",
			args: args{
				x:         5,
				y:         5,
				direction: left,
				shipType:  shipBattleship,
			},
			wantBoard: func() Board {
				var b Board
				b[5][5] = shipBattleship
				b[4][5] = shipBattleship
				b[3][5] = shipBattleship
				b[2][5] = shipBattleship
				return b
			}(),
		},
		{
			desc: "Submarine placed rightwards",
			args: args{
				x:         0,
				y:         0,
				direction: right,
				shipType:  shipSubmarine,
			},
			wantBoard: func() Board {
				var b Board
				b[0][0] = shipSubmarine
				b[1][0] = shipSubmarine
				b[2][0] = shipSubmarine
				return b
			}(),
		},
		{
			desc: "PatrolBoat placed downwards at edge of board",
			args: args{
				x:         boardSize - 1,
				y:         boardSize - 1,
				direction: down,
				shipType:  shipPatrolBoat,
			},
			wantBoard: func() Board {
				var b Board
				b[boardSize-1][boardSize-1] = shipPatrolBoat
				b[boardSize-1][boardSize-2] = shipPatrolBoat
				return b
			}(),
		},
		{
			desc: "top edge error",
			args: args{
				x:         5,
				y:         boardSize - 1,
				direction: up,
				shipType:  shipPatrolBoat,
			},
			wantError: true,
		},
		{
			desc: "bottom edge error",
			args: args{
				x:         5,
				y:         0,
				direction: down,
				shipType:  shipPatrolBoat,
			},
			wantError: true,
		},
		{
			desc: "left edge error",
			args: args{
				x:         0,
				y:         5,
				direction: left,
				shipType:  shipPatrolBoat,
			},
			wantError: true,
		},
		{
			desc: "right edge error",
			args: args{
				x:         boardSize - 1,
				y:         5,
				direction: right,
				shipType:  shipPatrolBoat,
			},
			wantError: true,
		},
		{
			desc: "existing ship error",
			args: args{
				x:         7,
				y:         3,
				direction: right,
				shipType:  shipPatrolBoat,
			},
			startBoard: func() Board {
				var b Board
				b[8][5] = shipCarrier
				b[8][4] = shipCarrier
				b[8][3] = shipCarrier
				b[8][2] = shipCarrier
				b[8][1] = shipCarrier
				return b
			}(),
			wantError: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			board := tC.startBoard

			err := board.PlaceShip(tC.args.x, tC.args.y, tC.args.direction, tC.args.shipType)
			if tC.wantError {
				if err == nil {
					t.Fatalf("PlaceShip(%v, %v, %v, %v): wanted error, got \n%v", tC.args.x, tC.args.y, tC.args.direction, tC.args.shipType, board)
				}
				if board != tC.startBoard {
					t.Fatalf("PlaceShip returned error (good), but it modified the board!\n%v", board)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if board != tC.wantBoard {
					t.Fatalf("PlaceShip() didn't return what was wanted \nwant:\n%v\ngot:\n%v", tC.wantBoard, board)
				}
			}
		})
	}
}

func TestParsePosition(t *testing.T) {
	testCases := []struct {
		location string
		wantx    int
		wanty    int
		wanterr  bool // if checking for returned error, wantx and wanty are ignored
	}{
		{
			location: "a1",
			wantx:    0,
			wanty:    0,
		},
		{
			location: "A1",
			wantx:    0,
			wanty:    0,
		},
		{
			location: "A10",
			wantx:    9,
			wanty:    0,
		},
		{
			location: "a0",
			wanterr:  true,
		},
		{
			location: "x1",
			wanterr:  true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.location, func(t *testing.T) {
			x, y, err := ParsePosition(tC.location)
			if tC.wanterr {
				if err == nil {
					t.Fatalf("wanted error but didn't get one, got x:%v, y:%v", x, y)
				}
			} else {
				if err != nil {
					t.Fatalf("got error %v", err)
				}
				if x != tC.wantx || y != tC.wanty {
					t.Fatalf("wanted x:%v, y:%v, got x:%v, y:%v", tC.wantx, tC.wanty, x, y)
				}
			}
		})
	}
}
