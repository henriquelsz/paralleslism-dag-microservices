package main

import (
	"fmt"
	"sync"
	"time"
)
 
func fibSequential(n int) int {
	if n <= 1 {
		return n
	}
	return fibSequential(n-1) + fibSequential(n-2)
}
 
func fibParallel(n, cutoff int) int {
	if n <= cutoff {
		return fibSequential(n)
	}
	var r1 int
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		r1 = fibParallel(n-1, cutoff)
	}()
	r2 := fibParallel(n-2, cutoff)
	wg.Wait()
	return r1 + r2
}
 
// measureGoroutineOverhead measures the average cost of spawning a goroutine
func measureGoroutineOverhead() time.Duration {
	const iterations = 100_000
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
		}()
	}
	wg.Wait()
	return time.Since(start) / iterations
}
 
func main() {
	n := 42
 
	fmt.Println("══════════════════════════════════════════════════════════")
	fmt.Println("  BENCHMARK: Finding the Optimal Cutoff")
	fmt.Println("  Theory: T_subproblem must be >= 10x T_overhead")
	fmt.Println("══════════════════════════════════════════════════════════")
	fmt.Println()
 
	// Step 1: Measure goroutine overhead
	overhead := measureGoroutineOverhead()
	fmt.Printf("  Average overhead per goroutine: %v\n\n", overhead)
 
	// Step 2: Measure sequential time for different values of n
	fmt.Println("  ┌───────┬──────────────┬──────────────┬─────────────────────┐")
	fmt.Println("  │   n   │    Time      │ vs Overhead  │ Worth parallelizing?│")
	fmt.Println("  ├───────┼──────────────┼──────────────┼─────────────────────┤")
 
	for _, testN := range []int{5, 10, 12, 15, 18, 20, 22, 25, 28, 30} {
		const runs = 1000
		start := time.Now()
		for i := 0; i < runs; i++ {
			fibSequential(testN)
		}
		avg := time.Since(start) / runs
		ratio := float64(avg) / float64(overhead)
		worth := "❌ NO  (< 10x)"
		if ratio >= 10 {
			worth = "✅ YES (>= 10x)"
		}
		fmt.Printf("  │ %5d │ %12v │ %10.1fx  │ %-19s │\n", testN, avg, ratio, worth)
	}
	fmt.Println("  └───────┴──────────────┴──────────────┴─────────────────────┘")
 
	// Step 3: Cutoff benchmark
	fmt.Println()
	fmt.Printf("  Searching for optimal cutoff for fib(%d)...\n\n", n)
 
	// Baseline
	start := time.Now()
	fibSequential(n)
	seqTime := time.Since(start)
 
	fmt.Println("  ┌──────────┬──────────────┬──────────────┬───────────────────┐")
	fmt.Println("  │  Cutoff  │    Time      │   Speedup    │ Status            │")
	fmt.Println("  ├──────────┼──────────────┼──────────────┼───────────────────┤")
	fmt.Printf("  │ seq      │ %12v │    1.00x     │ baseline          │\n", seqTime)
 
	bestCutoff := 0
	bestTime := seqTime
 
	for cutoff := 5; cutoff <= 35; cutoff += 5 {
		start := time.Now()
		fibParallel(n, cutoff)
		parTime := time.Since(start)
 
		speedup := float64(seqTime) / float64(parTime)
		status := "          "
		if parTime < bestTime {
			bestTime = parTime
			bestCutoff = cutoff
			status = "🏆 best!  "
		} else if parTime < seqTime {
			status = "✅ better "
		} else {
			status = "❌ worse  "
		}
 
		fmt.Printf("  │ %8d │ %12v │   %6.2fx     │ %s│\n",
			cutoff, parTime, speedup, status)
	}
 
	fmt.Println("  └──────────┴──────────────┴──────────────┴───────────────────┘")
	fmt.Println()
	fmt.Printf("  🎯 Optimal cutoff: %d (time: %v, speedup: %.2fx)\n",
		bestCutoff, bestTime, float64(seqTime)/float64(bestTime))
	fmt.Println()
	fmt.Println("  Leiserson's rule: only spawn when the sub-problem's work")
	fmt.Println("  is >= 10x the spawn overhead.")
	fmt.Println("══════════════════════════════════════════════════════════")
}