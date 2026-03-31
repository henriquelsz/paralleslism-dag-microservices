package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ServiceResult struct {
	Name   string
	Result string
	Err    error
}

func callService(ctx context.Context, name string, latency time.Duration) (string, error) {
	select {
	case <-time.After(latency):
		return fmt.Sprintf("[%s: OK in %v]", name, latency), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// spawnService starts a service call in a goroutine (Leiserson SPAWN)
func spawnService(
	ctx context.Context,
	wg *sync.WaitGroup,
	ch chan<- ServiceResult,
	name string,
	latency time.Duration,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		result, err := callService(ctx, name, latency)
		ch <- ServiceResult{Name: name, Result: result, Err: err}
	}()
}

// collectResults waits for the spawns and collects results (Leiserson SYNC)
func collectResults(wg *sync.WaitGroup, ch chan ServiceResult) []ServiceResult {
	go func() {
		wg.Wait()
		close(ch)
	}()

	var results []ServiceResult
	for r := range ch {
		results = append(results, r)
	}
	return results
}

// checkoutParallel explores the DAG's REAL parallelism.
//
// Optimized DAG:
//
//   Gateway(5ms)
//       ├──→ Auth(30ms) ──┐        ← parallel SPAWN
//       └──→ Config(20ms)─┘        ← parallel SPAWN
//               │
//          Order(40ms)              ← SYNC (waits for Auth and Config)
//        ┌──────┼──────┐
//   Inventory  Price   Shipping         ← SPAWN of 3 goroutines
//    (35ms)  (25ms)  (45ms)
//        └──────┼──────┘
//         Payment(60ms)             ← SYNC (waits for all 3)
//               │
//       Notification(15ms)
//
// T₁ = 275ms (same total work)
// T∞ = 5 + max(30,20) + 40 + max(35,25,45) + 60 + 15 = 195ms
// Parallelism = 275/195 = 1.41
func checkoutParallel(ctx context.Context) ([]string, time.Duration, error) {
	start := time.Now()
	var allResults []string

	// ── Phase 1: Gateway ──
	gw, err := callService(ctx, "Gateway", 5*time.Millisecond)
	if err != nil {
		return nil, 0, err
	}
	allResults = append(allResults, gw)

	// ── Phase 2: Auth and Config in PARALLEL (spawn) ──
	// Phase span: max(30, 20) = 30ms (instead of 50ms sequential)
	phase2Ch := make(chan ServiceResult, 2)
	var wg2 sync.WaitGroup

	spawnService(ctx, &wg2, phase2Ch, "Auth", 30*time.Millisecond)
	spawnService(ctx, &wg2, phase2Ch, "Config", 20*time.Millisecond)

	for _, r := range collectResults(&wg2, phase2Ch) {
		if r.Err != nil {
			return nil, 0, r.Err
		}
		allResults = append(allResults, r.Result)
	}

	// ── Phase 3: Order (depends on Auth + Config) ──
	pedido, err := callService(ctx, "Pedido", 40*time.Millisecond)
	if err != nil {
		return nil, 0, err
	}
	allResults = append(allResults, pedido)

	// ── Phase 4: Inventory, Price and Shipping in PARALLEL (spawn 3) ──
	// Phase span: max(35, 25, 45) = 45ms (instead of 105ms sequential)
	phase4Ch := make(chan ServiceResult, 3)
	var wg4 sync.WaitGroup

	spawnService(ctx, &wg4, phase4Ch, "Estoque", 35*time.Millisecond)
	spawnService(ctx, &wg4, phase4Ch, "Preço", 25*time.Millisecond)
	spawnService(ctx, &wg4, phase4Ch, "Frete", 45*time.Millisecond)

	for _, r := range collectResults(&wg4, phase4Ch) {
		if r.Err != nil {
			return nil, 0, r.Err
		}
		allResults = append(allResults, r.Result)
	}

	// ── Phase 5: Payment (depends on Inventory + Price + Shipping) ──
	pgto, err := callService(ctx, "Pagamento", 60*time.Millisecond)
	if err != nil {
		return nil, 0, err
	}
	allResults = append(allResults, pgto)

	// ── Phase 6: Notification ──
	notif, err := callService(ctx, "Notificação", 15*time.Millisecond)
	if err != nil {
		return nil, 0, err
	}
	allResults = append(allResults, notif)

	return allResults, time.Since(start), nil
}

func main() {
	ctx := context.Background()
	results, elapsed, _ := checkoutParallel(ctx)

	fmt.Println("══════════════════════════════════════════════════════")
	fmt.Println("  PARALLEL CHECKOUT — Optimized DAG")
	fmt.Println("  Independent calls executed in parallel")
	fmt.Println("══════════════════════════════════════════════════════")
	for _, r := range results {
		fmt.Println("  ", r)
	}
	fmt.Println()
	fmt.Printf("  Total time:      %v\n", elapsed)
	fmt.Println()
	fmt.Println("  Leiserson analysis:")
	fmt.Println("  T₁ (Work)       = 275ms")
	fmt.Println("  T∞ (Span)       = 195ms")
	fmt.Println("  Parallelism     = 1.41")
	fmt.Println()
	fmt.Println("  DAG phases:")
	fmt.Println("  Gateway:           5ms  (sequential)")
	fmt.Println("  Auth ∥ Config:    30ms  (max(30,20) — saves 20ms)")
	fmt.Println("  Order:           40ms  (sequential)")
	fmt.Println("  Invent ∥ Price ∥ Ship: 45ms  (max(35,25,45) — saves 60ms)")
	fmt.Println("  Payment:        60ms  (sequential)")
	fmt.Println("  Notification:      15ms  (sequential)")
	fmt.Println("  ─────────────────────────")
	fmt.Println("  Expected T∞:     195ms")
	fmt.Println()
	fmt.Println("  Tp with P=∞:    195ms  (= T∞, absolute lower bound)")
	fmt.Println("  Tp with P=100:  197ms  (almost the same as P=∞)")
	fmt.Println("  Tp with P=3:    287ms  (little improvement)")
	fmt.Println()
	fmt.Println("  ⚠️  Span dominates! To improve beyond 195ms,")
	fmt.Println("  you must REDUCE DEPENDENCIES, not add CPUs.")
	fmt.Println("══════════════════════════════════════════════════════")
}