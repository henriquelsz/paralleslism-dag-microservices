package main

import (
	"fmt"
	"strings"
)

// ServiceNode represents a node in a microservices DAG
type ServiceNode struct {
	Name         string
	LatencyMs    int
	Dependencies []string
}

// DAGAnalyzer computes Leiserson metrics for a services DAG
type DAGAnalyzer struct {
	nodes map[string]*ServiceNode
	order []string // preserves insertion order
}

func NewDAGAnalyzer() *DAGAnalyzer {
	return &DAGAnalyzer{nodes: make(map[string]*ServiceNode)}
}

func (d *DAGAnalyzer) AddNode(name string, latencyMs int, deps ...string) {
	d.nodes[name] = &ServiceNode{
		Name:         name,
		LatencyMs:    latencyMs,
		Dependencies: deps,
	}
	d.order = append(d.order, name)
}

// Work computes T₁: sum of all latencies
func (d *DAGAnalyzer) Work() int {
	total := 0
	for _, node := range d.nodes {
		total += node.LatencyMs
	}
	return total
}

// Span computes T∞: longest path (critical path)
func (d *DAGAnalyzer) Span() int {
	memo := make(map[string]int)
	maxSpan := 0
	for name := range d.nodes {
		s := d.spanOf(name, memo)
		if s > maxSpan {
			maxSpan = s
		}
	}
	return maxSpan
}

func (d *DAGAnalyzer) spanOf(name string, memo map[string]int) int {
	if val, ok := memo[name]; ok {
		return val
	}
	node := d.nodes[name]
	maxDepSpan := 0
	for _, dep := range node.Dependencies {
		depSpan := d.spanOf(dep, memo)
		if depSpan > maxDepSpan {
			maxDepSpan = depSpan
		}
	}
	memo[name] = maxDepSpan + node.LatencyMs
	return memo[name]
}

// CriticalPath returns the nodes on the critical path
func (d *DAGAnalyzer) CriticalPath() []string {
	memo := make(map[string]int)
	for name := range d.nodes {
		d.spanOf(name, memo)
	}

	maxName := ""
	maxVal := 0
	for name, val := range memo {
		if val > maxVal {
			maxVal = val
			maxName = name
		}
	}

	var path []string
	current := maxName
	for current != "" {
		path = append([]string{current}, path...)
		node := d.nodes[current]
		bestDep := ""
		bestVal := 0
		for _, dep := range node.Dependencies {
			if memo[dep] > bestVal {
				bestVal = memo[dep]
				bestDep = dep
			}
		}
		current = bestDep
	}
	return path
}

// Parallelism returns T₁/T∞
func (d *DAGAnalyzer) Parallelism() float64 {
	return float64(d.Work()) / float64(d.Span())
}

// Tp computes estimated time with P processors: T₁/P + T∞
func (d *DAGAnalyzer) Tp(p int) float64 {
	return float64(d.Work())/float64(p) + float64(d.Span())
}

// PrintAnalysis prints the full Leiserson analysis
func (d *DAGAnalyzer) PrintAnalysis(title string) {
	work := d.Work()
	span := d.Span()
	parallelism := d.Parallelism()
	critPath := d.CriticalPath()

	fmt.Println()
	fmt.Printf("  %s\n", title)
	fmt.Println()
	fmt.Println("  ┌────────────────────────────────────────────────┐")
	fmt.Println("  │         DAG ANALYSIS — Leiserson Model         │")
	fmt.Println("  ├────────────────────────────────────────────────┤")
	fmt.Printf("  │  T₁ (Total work)      = %-4d ms               │\n", work)
	fmt.Printf("  │  T∞ (Span / Critical) = %-4d ms               │\n", span)
	fmt.Printf("  │  Parallelism (T₁/T∞)  = %-5.2f                │\n", parallelism)
	fmt.Println("  ├────────────────────────────────────────────────┤")
	fmt.Printf("  │  Critical Path:                                │\n")
	fmt.Printf("  │  %s\n", strings.Join(critPath, " → "))
	fmt.Println("  ├────────────────────────────────────────────────┤")
	fmt.Println("  │  Tp by number of processors:                   │")

	for _, p := range []int{1, 2, 4, 8, 16, 100, 1000} {
		tp := d.Tp(p)
		barLen := int(tp / 8)
		if barLen > 40 {
			barLen = 40
		}
		bar := strings.Repeat("█", barLen)
		fmt.Printf("  │  P=%-5d → Tp = %6.1f ms  %s\n", p, tp, bar)
	}

	fmt.Println("  ├────────────────────────────────────────────────┤")

	if parallelism < 2 {
		fmt.Println("  │  ⚠️  WARNING: Parallelism < 2                 │")
		fmt.Println("  │  → SPAN-dominated                              │")
		fmt.Println("  │  → More CPUs/instances DO NOT help            │")
		fmt.Println("  │  → Reduce sequential dependencies!            │")
	} else if parallelism < 8 {
		fmt.Printf("  │  ⚡ Moderate parallelism                       │\n")
		fmt.Printf("  │  → Up to ~%d CPUs bring real benefits          \n", int(parallelism))
		fmt.Println("  │  → Beyond that, returns diminish quickly      │")
	} else {
		fmt.Println("  │  ✅ High parallelism!                          │")
		fmt.Println("  │  → Horizontal scaling brings real gains       │")
	}

	fmt.Println("  └────────────────────────────────────────────────┘")
}

func main() {
	fmt.Println("══════════════════════════════════════════════════════")
	fmt.Println("  DAG ANALYZER — Leiserson Metrics")
	fmt.Println("  Tp = T₁/P + T∞")
	fmt.Println("══════════════════════════════════════════════════════")

	// ── Scenario 1: E-Commerce Checkout (real DAG) ──
	dag1 := NewDAGAnalyzer()
	dag1.AddNode("Gateway", 5)
	dag1.AddNode("Auth", 30, "Gateway")
	dag1.AddNode("Config", 20, "Gateway")
	dag1.AddNode("Order", 40, "Auth", "Config")
	dag1.AddNode("Inventory", 35, "Order")
	dag1.AddNode("Price", 25, "Order")
	dag1.AddNode("Shipping", 45, "Order")
	dag1.AddNode("Payment", 60, "Inventory", "Price", "Shipping")
	dag1.AddNode("Notification", 15, "Payment")

	dag1.PrintAnalysis("SCENARIO 1: E-Commerce Checkout (DAG with parallelism)")

	// ── Scenario 2: Same services, fully sequential ──
	dag2 := NewDAGAnalyzer()
	dag2.AddNode("Gateway", 5)
	dag2.AddNode("Auth", 30, "Gateway")
	dag2.AddNode("Config", 20, "Auth")
	dag2.AddNode("Order", 40, "Config")
	dag2.AddNode("Inventory", 35, "Order")
	dag2.AddNode("Price", 25, "Inventory")
	dag2.AddNode("Shipping", 45, "Price")
	dag2.AddNode("Payment", 60, "Shipping")
	dag2.AddNode("Notification", 15, "Payment")

	dag2.PrintAnalysis("SCENARIO 2: Same services in a sequential chain")

	// ── Scenario 3: Highly parallel DAG ──
	dag3 := NewDAGAnalyzer()
	dag3.AddNode("Entrance", 5)
	dag3.AddNode("SvcA", 30, "Entrance")
	dag3.AddNode("SvcB", 25, "Entrance")
	dag3.AddNode("SvcC", 35, "Entrance")
	dag3.AddNode("SvcD", 20, "Entrance")
	dag3.AddNode("SvcE", 40, "Entrance")
	dag3.AddNode("SvcF", 15, "Entrance")
	dag3.AddNode("SvcG", 28, "Entrance")
	dag3.AddNode("SvcH", 33, "Entrance")
	dag3.AddNode("Exit", 5, "SvcA", "SvcB", "SvcC", "SvcD", "SvcE", "SvcF", "SvcG", "SvcH")

	dag3.PrintAnalysis("SCENARIO 3: Max fan-out (8 parallel services)")

	fmt.Println()
	fmt.Println("══════════════════════════════════════════════════════")
}