# Learning Plan

This repo doubles as a learning project for a Go backend developer roadmap. The source plan
lives at `C:\Users\User\Desktop\New Go\План обучения Нурсултан.xlsx` (sheet "ПЛАН") — this file
is a snapshot of it for agents working in this repo, organized by the plan's own study order
(column "Порядок изучения", groups 1→11).

When helping with tasks in this repo, use this to know what topic space is in scope and roughly
where a given task sits in the sequence — see also `## Collaboration Mode (Learning Phase)` in
`AGENTS.md` for how to work on these tasks (review, don't implement by default).

## Group 1 — Base types
- Constant
- Map
- Array, slice
- Pointers
- String, rune, byte
- Int, uint, float

## Group 2 — OOP
- Polymorphism
- Inheritance, composition, embedding
- Encapsulation
- Struct
- Interface
- OOP, SOLID, DRY

## Group 3 — Foundations
- Git
- Data structures, Big-O, basic algorithms
- Lambda, variadic
- Recursion, closure
- First-class function

## Group 4 — Application plumbing
- Context
- Errors handling
- Main/init, defer

## Group 5 — Tools
- Godoc, gofmt, go vet
- Pprof
- Linters (GolangCI, golint)
- Go mod
- Go test/benchmark

## Group 6 — Concurrency
- Garbage collector
- Package sync: atomic, mutex, waitgroup, once, pool, map
- Channel, select
- Goroutine, scheduler
- Race condition

## Group 7 — Testing
- Testify
- Mocking
- Package httptest

## Group 8 — Database
- SQL
- Indexes
- Transactions, ACID
- Cache
- Locking (optimistic/pessimistic, table/row/advisory locks)
- Sharding, partitioning, replication

## Group 9 — Network
- gRPC
- WebSocket
- API (REST vs SOAP, methods)
- OSI model
- OS processes and threads
- Virtualization, containerization

## Group 10 — Software Engineering Practices
- Design patterns
- Microservices vs Monolith
- Code standards
- Code review, code smells, refactoring

## Group 11 — Delivery practices
- CI/CD
- Tracing, logging, metrics
- Agile, Scrum, Kanban

## Next up (not yet in this repo)
- Redis — likely a cache layer alongside/instead of the existing `sturdyc` in-process cache in
  `business/domain/productbus`, or session/token storage
- Kafka — an event-driven extension, e.g. publishing domain events on writes and having a
  consumer (or the `metrics` service) react to them

## Status

Progress tracking (what's already covered vs. not) is kept in the assistant's own session
memory, not in this file — this file is topic scope, not a live progress tracker. Ask the
assistant if you want a refresher on what's been covered.
