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
 
// fibParallel demonstrates Leiserson's spawn/sync model.
//
// Each recursive call "spawns" fib(n-1) in a goroutine
// (analogous to cilk_spawn) and computes fib(n-2) in the current goroutine.
// The sync.WaitGroup acts as the "cilk_sync".
//
// Analysis according to Leiserson:
//   T1 (Work) = O(phi^n)     — same total work as sequential
//   T∞ (Span) = O(n)         — critical path is linear in n
//   Parallelism = O(phi^n/n) — grows exponentially!
//
// The cutoff avoids goroutine overhead for small sub-problems.
func fibParallel(n, cutoff int) int {
	if n <= cutoff {
		return fibSequential(n)
	}
 
	var result1 int
	var wg sync.WaitGroup
 
	// SPAWN: equivalent to cilk_spawn
	wg.Add(1)
	go func() {
		defer wg.Done()
		result1 = fibParallel(n-1, cutoff)
	}()
 
	// CONTINUATION: runs in the current goroutine
	result2 := fibParallel(n-2, cutoff)
 
	// SYNC: equivalent to cilk_sync
	wg.Wait()
	return result1 + result2
}
 
func main() {
	n := 42
	cutoff := 20
 
	// Sequential baseline
	start := time.Now()
	seqResult := fibSequential(n)
	seqTime := time.Since(start)
 
	// Parallel execution
	start = time.Now()
	parResult := fibParallel(n, cutoff)
	parTime := time.Since(start)
 
	speedup := float64(seqTime) / float64(parTime)
 
	fmt.Println("════════════════════════════════════════════════════")
	fmt.Println("  PARALLEL FIBONACCI — Spawn/Sync Model")
	fmt.Println("════════════════════════════════════════════════════")
	fmt.Printf("  fib(%d) = %d\n", n, seqResult)
	fmt.Printf("  Cutoff  = %d\n", cutoff)
	fmt.Println()
	fmt.Printf("  Sequential: %v\n", seqTime)
	fmt.Printf("  Parallel:   %v\n", parTime)
	fmt.Printf("  Speedup:    %.2fx\n", speedup)
	fmt.Println()
	fmt.Println("  Leiserson Analysis:")
	fmt.Printf("  T1 (Work)    = %v (sequential time)\n", seqTime)
	fmt.Printf("  T∞ (Span)    = O(%d) (proportional to n)\n", n)
	fmt.Printf("  Parallelism  = O(phi^%d / %d) ≈ millions\n", n, n)
	fmt.Println("  → Work-dominated: adding CPUs HELPS")
	fmt.Println()
	if seqResult != parResult {
		fmt.Println("  ⚠️  ERROR: results diverge!")
	} else {
		fmt.Println("  ✅ Results are identical")
	}
	fmt.Println("════════════════════════════════════════════════════")
}