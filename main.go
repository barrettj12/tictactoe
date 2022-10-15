package main

import (
	"encoding/json"
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
	// Try load positions from file
	allPos, err := loadPositions()
	if err == nil {
		fmt.Printf("Loaded %d positions from file\n", len(allPos))
	} else {
		fmt.Println("error loading positions:", err)
		// Regenerate positions
		fmt.Print("Generating possible board positions... ")
		start := time.Now()
		allPos := getAllPositions()
		elapsed := time.Since(start)
		fmt.Println("Done.")
		fmt.Printf("Generated %d positions in %v\n", len(allPos), elapsed)

		// Write positions to file
		data, err := json.Marshal(allPos)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile("positions.json", data, os.ModePerm)
		if err != nil {
			panic(err)
		}
		fmt.Println("Positions written to file.")
	}

	// Try load strategy from file
	st, err := loadStrategy()
	if err == nil {
		fmt.Println("Loaded strategy from file")
	} else {
		fmt.Println("error loading strategy:", err)
		fmt.Print("Generating random strategy... ")
		st = genRandStrat(allPos)
		fmt.Println("Done.")

		// Write strategy to file
		data, err := json.Marshal(st)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile("strategy.json", data, os.ModePerm)
		if err != nil {
			panic(err)
		}
		fmt.Println("Strategy written to file.")
	}

	// Play it randomly
	fmt.Println("Playing a random game of Tic-Tac-Toe...")
	res := playRandom(st, true)
	fmt.Println("Result:", res)
}

func loadPositions() (allPos []Position, err error) {
	data, err := os.ReadFile("positions.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &allPos)
	return
}

func loadStrategy() (st Strategy, err error) {
	data, err := os.ReadFile("strategy.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &st)
	return
}
