package main

import (
	"fmt"
	"time"
)
 
// fibSequential calculates Fibonacci purely sequentially.
//
// Analysis according to Leiserson:
//   T1 (Work) = O(phi^n) where phi ≈ 1.618
//   T∞ (Span) = O(phi^n) (no parallelism, T∞ = T1)
//   Parallelism = T1/T∞ = 1.0 (no gain possible with more CPUs)
func fibSequential(n int) int {
	if n <= 1 {
		return n
	}
	return fibSequential(n-1) + fibSequential(n-2)
}
 
func main() {
	n := 42
 
	start := time.Now()
	result := fibSequential(n)
	elapsed := time.Since(start)
 
	fmt.Println("════════════════════════════════════════════")
	fmt.Println("  SEQUENTIAL FIBONACCI")
	fmt.Println("  Model: Entire DAG executed on 1 CPU")
	fmt.Println("════════════════════════════════════════════")
	fmt.Printf("  fib(%d)   = %d\n", n, result)
	fmt.Printf("  Time      = %v\n", elapsed)
	fmt.Println()
	fmt.Println("  Leiserson Analysis:")
	fmt.Println("  T1 (Work)    = total time above")
	fmt.Println("  T∞ (Span)    = T1 (no parallelism exploited)")
	fmt.Println("  Parallelism  = 1.0")
	fmt.Println("════════════════════════════════════════════")
}