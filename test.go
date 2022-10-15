package main

import (
	"fmt"
	"time"
)

func TestGeneratePositions() {
	start := time.Now()
	allPos := getAllPositions()
	elapsed := time.Since(start)
	fmt.Printf("Generated %d positions in %v\n", len(allPos), elapsed)
}
