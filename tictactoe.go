package main

import (
	"fmt"
	"math/rand"
)

type Square byte

const (
	Empty Square = ' '
	SqX   Square = 'X'
	SqO   Square = 'O'
)

// Board positions
type Position [9]Square

var StartPos Position = [9]Square{
	Empty, Empty, Empty,
	Empty, Empty, Empty,
	Empty, Empty, Empty,
}

// Position needs to implement encoding.TextMarshaler
func (p Position) MarshalText() (text []byte, err error) {
	str := []byte{}
	for i := 0; i < 9; i++ {
		str = append(str, byte(p[i]))
	}
	return str, nil
}

func (p *Position) UnmarshalText(text []byte) error {
	if len(text) != 9 {
		return fmt.Errorf("wrong length %d != 9", len(text))
	}

	for i := 0; i < 9; i++ {
		sq := Square(text[i])
		if sq == Empty || sq == SqX || sq == SqO {
			p[i] = sq
		} else {
			return fmt.Errorf("invalid square value %v", sq)
		}
	}
	return nil
}

// Pretty print the current board position
func printPos(pos Position) {
	fmt.Printf(`
 %c │ %c │ %c 
───┼───┼───
 %c │ %c │ %c 
───┼───┼───
 %c │ %c │ %c 
`[1:], pos[0], pos[1], pos[2], pos[3], pos[4], pos[5], pos[6], pos[7], pos[8])
}

// Count how many Xs and Os on a board
// i.e. which turn it is
func countTurn(pos Position) int {
	count := 0
	for _, sq := range pos {
		if sq == SqX || sq == SqO {
			count++
		}
	}
	return count
}

// Get all indices of blank squares on board
func getBlanks(pos Position) []int {
	blanks := []int{}
	for i := 0; i < 9; i++ {
		if pos[i] == Empty {
			blanks = append(blanks, i)
		}
	}
	return blanks
}

// Get all possible (playable) board positions
func getAllPositions() []Position {
	allPos := []Position{StartPos}
	seen := 0

	for turn := 0; turn < 8; turn++ {
		start := seen
		seen = len(allPos)

		// Determine what mark (X or O) to add
		var mark Square
		switch turn {
		case 0, 2, 4, 6, 8:
			// Add an X
			mark = SqX
		case 1, 3, 5, 7:
			// Add an O
			mark = SqO
		default:
			panic(fmt.Sprintf("turn %d not valid", turn))
		}

		for _, pos := range allPos[start:seen] {
			// Generate possible next positions
			blanks := getBlanks(pos)
			for _, i := range blanks {
				newPos := pos
				newPos[i] = mark
				allPos = append(allPos, newPos)
			}

		}
	}

	return allPos
}

type Result string

const (
	Draw Result = "draw"
	XWon Result = "X won"
	OWon Result = "O won"
)

// Return result for a given board
func result(pos Position) Result {
	// all lines on board
	lines := [][]int{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
		{0, 3, 6},
		{1, 4, 7},
		{2, 5, 8},
		{0, 4, 8},
		{2, 4, 6},
	}

	for _, l := range lines {
		if pos[l[0]] == pos[l[1]] && pos[l[0]] == pos[l[2]] {
			if pos[l[0]] == SqX {
				return XWon
			}
			if pos[l[0]] == SqO {
				return OWon
			}
		}
	}

	return Draw
}

type Strategy map[Position]int

// Play the given strategy s against a random strategy
// s == X plays first
func playRandom(s Strategy, print bool) Result {
	pos := StartPos

	for {
		turn := countTurn(pos)

		// Check if there's a winner
		res := result(pos)
		if res == XWon || res == OWon || turn == 9 {
			return res
		}

		switch turn {
		case 0, 2, 4, 6, 8:
			// X's (s's) turn
			// Get next index to play
			ind := s[pos]
			if pos[ind] != Empty {
				panic(fmt.Sprintf("index %v in pos %v is already filled", ind, pos))
			}
			pos[ind] = SqX
		case 1, 3, 5, 7:
			// O's (random) turn
			blanks := getBlanks(pos)
			ind := blanks[rand.Intn(len(blanks))]
			pos[ind] = SqO
		default:
			panic(fmt.Sprintf("turn %d not valid", turn))
		}

		if print {
			fmt.Printf("\nTurn %d:\n", turn+1)
			printPos(pos)
		}
	}
}

// generate random strategy
func genRandStrat(allPositions []Position) Strategy {
	st := make(Strategy, len(allPositions))

	for _, pos := range allPositions {
		blanks := getBlanks(pos)
		st[pos] = blanks[rand.Intn(len(blanks))]
	}

	return st
}
