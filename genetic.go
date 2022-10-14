package main

import (
	"fmt"
	"math/rand"
	"sort"
)

// Parameters for genetic algorithm
const (
	// How large each generation is
	GENERATION_SIZE = SURVIVORS*SURVIVORS + SURVIVORS

	// How many of each generation will survive to reproduce
	SURVIVORS = 5

	// How many random games are played to determine the test score
	NUM_GAMES = 100

	// What percentage of genes are mutated
	MUTATION_RATE = 0.02

	// How many generations to run the algorithm for
	NUM_GENERATIONS = 500
)

var testCounter = 0

// Returns a score evaluating how good this strategy is
func test(st Strategy) int {
	recoverTest := func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("testCounter=%d, st=%v", testCounter, st))
		}
	}
	defer recoverTest()

	testCounter++
	wins := 0
	for i := 0; i < NUM_GAMES; i++ {
		if playRandom(st, false) == XWon {
			wins++
		}
	}
	return wins
}

// Given a "generation" of strategies, determine a next generation
func nextGen(gen []Strategy) []Strategy {
	// Whittle down generation
	survivors := choose(gen)
	nextGenr := make([]Strategy, 0, GENERATION_SIZE)
	// All survivors to next generation
	for _, st := range survivors {
		nextGenr = append(nextGenr, st)
	}

	// Create children of each pair
	for _, s1 := range survivors {
		for _, s2 := range survivors {
			// Create random mix of s1 and s2
			child := make(Strategy, len(s1))
			for pos := range s1 {
				if rand.Float64() < MUTATION_RATE {
					// Mutate this gene - pick random
					blanks := getBlanks(pos)
					child[pos] = blanks[rand.Intn(len(blanks))]
				} else if rand.Intn(2) == 0 {
					child[pos] = s1[pos]
				} else {
					child[pos] = s2[pos]
				}
			}

			nextGenr = append(nextGenr, child)
		}
	}

	return nextGenr
}

// Choose survivors from generation
func choose(gen []Strategy) []Strategy {
	scores := make([]int, len(gen))

	for i, st := range gen {
		scores[i] = test(st)
	}
	sort.Slice(gen, func(i, j int) bool { return scores[i] > scores[j] })

	// Print top scores
	sort.Slice(scores, func(i, j int) bool { return scores[i] > scores[j] })
	fmt.Println(scores[:SURVIVORS])

	return gen[:SURVIVORS]
}

// Run the algorithm
func evolve() {
	allPos := getAllPositions()

	// Initialise first gen randomly
	gen := make([]Strategy, 0, GENERATION_SIZE)
	fmt.Println("Creating first gen...")
	for k := 0; k < GENERATION_SIZE; k++ {
		st := genRandStrat(allPos)
		gen = append(gen, st)
		fmt.Printf("\r   \r%d", k)
	}
	fmt.Println("\r   \rDone :)")

	// Run the algorithm
	for n := 0; n < NUM_GENERATIONS; n++ {
		gen = nextGen(gen)
	}
}
