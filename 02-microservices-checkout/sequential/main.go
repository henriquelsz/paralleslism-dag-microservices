package main
 
import (
	"context"
	"fmt"
	"time"
)
 
// callService simulates a microservice call with fixed latency
func callService(ctx context.Context, name string, latency time.Duration) (string, error) {
	select {
	case <-time.After(latency):
		return fmt.Sprintf("[%s: OK in %v]", name, latency), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
 
// checkoutSequential executes ALL calls in series.
//
// Resulting DAG (linear chain):
//   GW → Auth → Config → Order → Inventory → Price → Shipping → Payment → Notification
//
// Leiserson Analysis:
//   T1  = 5+30+20+40+35+25+45+60+15 = 275ms
//   T∞  = 275ms (equal to T1 — everything is sequential)
//   Parallelism = T1/T∞ = 1.0 (NONE!)
func checkoutSequential(ctx context.Context) ([]string, time.Duration, error) {
	start := time.Now()
	var results []string
 
	services := []struct {
		name    string
		latency time.Duration
	}{
		{"Gateway",      5 * time.Millisecond},
		{"Auth",         30 * time.Millisecond},
		{"Config",       20 * time.Millisecond},
		{"Order",        40 * time.Millisecond},
		{"Inventory",    35 * time.Millisecond},
		{"Price",        25 * time.Millisecond},
		{"Shipping",     45 * time.Millisecond},
		{"Payment",      60 * time.Millisecond},
		{"Notification", 15 * time.Millisecond},
	}
 
	for _, svc := range services {
		r, err := callService(ctx, svc.name, svc.latency)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, r)
	}
	return results, time.Since(start), nil
}
 
func main() {
	ctx := context.Background()
	results, elapsed, _ := checkoutSequential(ctx)
 
	fmt.Println("══════════════════════════════════════════════════════")
	fmt.Println("  SEQUENTIAL CHECKOUT — Anti-Pattern")
	fmt.Println("  All services called in a chain")
	fmt.Println("══════════════════════════════════════════════════════")
	for _, r := range results {
		fmt.Println("  ", r)
	}
	fmt.Println()
	fmt.Printf("  Total time:    %v\n", elapsed)
	fmt.Println()
	fmt.Println("  Leiserson Analysis:")
	fmt.Println("  T1 (Work)      = 275ms")
	fmt.Println("  T∞ (Span)      = 275ms (everything sequential)")
	fmt.Println("  Parallelism    = 1.0")
	fmt.Println()
	fmt.Println("  ⚠️  Adding CPUs/instances does NOT help at all!")
	fmt.Println("  Span equals Work — zero opportunity for")
	fmt.Println("  parallelism in this structure.")
	fmt.Println("══════════════════════════════════════════════════════")
}
