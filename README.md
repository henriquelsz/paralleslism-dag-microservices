# Parallel Computation and Microservices — Leiserson's DAG Theory
 
Repository with practical Go examples from the article
**"Parallel Computation and Microservices: An Analysis Through the Lens of Charles E. Leiserson's DAG Theory"**.
 
📖 [Read the full article on Medium](https://medium.com/@henrique.lsouza2001/computa%C3%A7%C3%A3o-paralela-e-microservi%C3%A7os-uma-an%C3%A1lise-sob-a-%C3%B3tica-da-teoria-de-dags-de-charles-e-ef041fbc7805?postPublishedType=repub)
 
---
 
## Concepts Demonstrated
 
| Example | Leiserson Concept | What It Demonstrates |
|---------|------------------|----------------------|
| ` + "`01-fibonacci/sequential`" + ` | T₁ (Work) | Sequential baseline — all work on 1 CPU |
| ` + "`01-fibonacci/parallel`" + ` | Spawn/Sync, Cutoff | Parallelism with goroutines, analogous to ` + "`cilk_spawn`" + ` |
| ` + "`01-fibonacci/benchmark`" + ` | Overhead vs. Useful Work | How to find the optimal cutoff empirically |
| ` + "`02-microservices-checkout/sequential`" + ` | T₁ = T∞ (Maximum Span) | Synchronous chain — parallelism = 1.0 |
| ` + "`02-microservices-checkout/parallel`" + ` | Span Reduction | Optimized DAG — independent calls in parallel |
| ` + "`03-dag-analyzer`" + ` | Tp = T₁/P + T∞ | Tool that calculates Work, Span, Parallelism, and Critical Path |
 
---
 
## Prerequisites
 
- Go 1.21+
 
---
 
## How to Run
 
` + "```" + `bash
# Clone the repository
git clone https://github.com/seu-usuario/paralleslism-dag-microservices.git
cd paralleslism-dag-microservices
 
# Sequential Fibonacci
go run 01-fibonacci/sequential/main.go
 
# Parallel Fibonacci
go run 01-fibonacci/parallel/main.go
 
# Cutoff benchmark (find the optimal value)
go run 01-fibonacci/benchmark/main.go
 
# Sequential checkout (microservices chain)
go run 02-microservices-checkout/sequential/main.go
 
# Parallel checkout (optimized DAG)
go run 02-microservices-checkout/parallel/main.go
 
# DAG analyzer (calculates T₁, T∞, Parallelism, Critical Path)
go run 03-dag-analyzer/main.go
