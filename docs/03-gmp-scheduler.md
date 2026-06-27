# GMP Scheduler: Goroutine Runtime

## Overview

Go uses an M:N scheduling model: M goroutines mapped to N OS threads.
The scheduler is cooperative + preemptive (async preemption since Go 1.14).

## GMP Components

```
┌─────────────────────────────────────────────────────────────┐
│                        Go Runtime                           │
│                                                             │
│  ┌─────┐  ┌─────┐  ┌─────┐  ┌─────┐  ┌─────┐            │
│  │  G  │  │  G  │  │  G  │  │  G  │  │  G  │  Goroutines │
│  └──┬──┘  └──┬──┘  └──┬──┘  └──┬──┘  └──┬──┘            │
│     │        │        │        │        │                  │
│  ┌──┴────────┴──┐  ┌──┴────────┴──┐                       │
│  │   P (proc)   │  │   P (proc)   │   Logical Processors  │
│  │ local queue  │  │ local queue  │   (GOMAXPROCS)        │
│  └──────┬───────┘  └──────┬───────┘                       │
│         │                  │                               │
│  ┌──────┴───────┐  ┌──────┴───────┐                       │
│  │   M (thread) │  │   M (thread) │   OS Threads          │
│  └──────────────┘  └──────────────┘                       │
│                                                             │
│  ┌──────────────────────────────────────┐                  │
│  │         Global Run Queue             │                  │
│  │  (overflow from local queues)        │                  │
│  └──────────────────────────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

| Component | Role | Count |
|-----------|------|-------|
| **G** (Goroutine) | Unit of work, ~2KB stack | Millions possible |
| **M** (Machine) | OS thread, executes G | Limited (~10,000) |
| **P** (Processor) | Scheduling context + local queue | GOMAXPROCS (default: CPU cores) |

## Scheduling Rules

1. **G needs P to run**: A goroutine must be assigned to a P's local queue
2. **P needs M to execute**: P binds to an M (OS thread) for execution
3. **Local queue capacity**: Each P holds up to 256 Gs
4. **Overflow → global queue**: Excess Gs go to the global run queue

## Work Stealing

When a P's local queue is empty:
1. Check global run queue (grab batch of G/n)
2. Check network poller for ready Gs
3. **Steal half** from another P's local queue

This ensures all cores stay busy without central bottleneck.

## Handoff (Syscall Handling)

```
Goroutine makes syscall (e.g., file I/O):
  1. M blocks on syscall (can't run other Gs)
  2. P detaches from M ("handoff")
  3. P finds/creates new M to keep running Gs
  4. When syscall returns, G re-enters a P's queue
```

Network I/O uses **netpoller** (epoll/kqueue) instead non-blocking, no handoff.

## Preemption

| Version | Mechanism |
|---------|-----------|
| < Go 1.14 | Cooperative: only at function calls |
| ≥ Go 1.14 | Async: OS signals interrupt tight loops |

Async preemption solves: `for {}` no longer starves other goroutines.
The runtime sends SIGURG to preempt long-running Gs at safe points.

## Key Parameters

| Setting | Default | Purpose |
|---------|---------|---------|
| `GOMAXPROCS` | CPU cores | Number of Ps (parallelism level) |
| `GOMAXPROCS=1` | | Serializes all goroutines (debugging) |
| `runtime.LockOSThread()` | | Pin G to M (for C interop, UI libs) |
