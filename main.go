package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/maps"
)

var cmds = map[string]func(){
	"simple": simpleTest,
	"evolve": evolve,
	"test":   TestGeneratePositions,
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) == 0 {
		fmt.Println("please provide first arg:", strings.Join(maps.Keys(cmds), ", "))
		os.Exit(1)
	}
	fn, ok := cmds[os.Args[1]]
	if ok {
		fn()
	} else {
		fmt.Println("command must be one of", strings.Join(maps.Keys(cmds), ", "))
		os.Exit(1)
	}
}

func simpleTest() {
	// Generate positions
	fmt.Print("Generating possible board positions... ")
	start := time.Now()
	allPos := getAllPositions()
	elapsed := time.Since(start)
	fmt.Println("Done.")
	fmt.Printf("Generated %d positions in %v\n", len(allPos), elapsed)

	fmt.Print("Generating random strategy... ")
	st := genRandStrat(allPos)
	fmt.Println("Done.")

	// Play it randomly
	fmt.Println("Playing a random game of Tic-Tac-Toe...")
	res := playRandom(st, true)
	fmt.Println("Result:", res)
}
