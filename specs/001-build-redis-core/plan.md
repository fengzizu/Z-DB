# Implementation Plan: Build Redis-Core (Non-Relational Engine)

**Branch**: `001-build-redis-core` | **Date**: 2025-12-09 | **Spec**: [Feature Spec](spec.md)
**Input**: Feature specification from `/specs/001-build-redis-core/spec.md`

## Summary

Build a Redis-compatible database server in Go from scratch. The system will support a RESP-compliant TCP server, in-memory key-value storage (String, List, Hash, Set, ZSet), RDB persistence, and replication.

## Technical Context

**Language/Version**: Go 1.21+ (utilizing generics if needed, standard library focus)
**Primary Dependencies**: None (Standard Library only as per Constitution)
**Storage**: In-memory `map` + RDB Snapshot + AOF (Append Only File) Log
**Testing**: `testing` package (Unit tests), `redis-cli` / CodeCrafters test suite (Integration tests)
**Target Platform**: Linux/Unix-like systems (Standard sockets)
**Project Type**: Single Backend Service
**Performance Goals**: Support 100+ concurrent clients, minimal latency for O(1) ops.
**Constraints**: No external libraries for core logic. Thread-safe data access.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
| :--- | :--- | :--- |
| **Learning Over Performance** | PASS | Focus on clear implementation of RESP and data structures. |
| **Idiomatic Go** | PASS | Will use `net`, `io`, `sync` and Goroutines. |
| **Zero-Dependency Core** | PASS | Explicitly building parser and storage engine from scratch. |
| **Dual-Paradigm Architecture** | PASS | Architecture will decouple Server/Parser from the KV Engine to allow future SQL Engine. |
| **Test-Driven Confidence** | PASS | Plan includes unit tests for data structures and integration tests for server. |
| **User-Driven Implementation** | PASS | Plan defines structure; User implements logic. |

## Project Structure

### Documentation (this feature)

```text
specs/001-build-redis-core/
├── plan.md              # This file
├── research.md          # Protocol & Data Structure details
├── data-model.md        # Internal memory structures
├── quickstart.md        # Running instructions
├── checklists/          # Requirement validation
└── tasks.md             # Detailed implementation steps
```

### Source Code (repository root)

```text
zdb/                     # Main source directory
├── redis/               # Non-relational (Redis-like) Engine
│   ├── cmd/
│   │   └── server/
│   │       └── main.go  # Entry point
│   ├── config/          # Configuration
│   ├── internal/
│   │   ├── server/      # TCP Listener
│   │   ├── resp/        # RESP Protocol
│   │   ├── core/        # Database Core
│   │   ├── datastruct/  # Data Structures
│   │   ├── persistence/ # RDB/AOF
│   │   └── replication/ # Replication
│   └── tests/           # Integration tests
├── mysql/               # Relational (MySQL-like) Engine (Future Phase)
│   └── ...
└── common/              # Shared utilities (logging, utils) - optional
```

**Structure Decision**: Split top-level directories by paradigm (`redis/` vs `mysql/`) to clearly separate the two distinct learning paths while keeping them in the same monorepo.

## Phase 0: Research & Architecture

### 1. RESP Protocol (Redis Serialization Protocol)
- **Types**: Simple Strings (`+`), Errors (`-`), Integers (`:`), Bulk Strings (`$`), Arrays (`*`).
- **Parsing Strategy**: `bufio.Reader` based state machine. Read byte -> determine type -> read length -> read content.

### 2. Concurrency Model
- **Per-Client**: One Goroutine per active connection.
- **Shared State**: Global `Store` struct protected by `sync.RWMutex`.
- **Granularity**: Initial locking on the whole DB. Future optimization: Key-level locking or Sharding (if needed, but Global Mutex is fine for Phase 1).

### 3. Data Structures
- **KV Store**: `map[string]*RedisObject`.
- **RedisObject**: Struct containing `Type` (String, List, etc.), `Encoding`, `Ptr` (to actual data), and `LRU/Expiry` metadata.
- **Expiration**: Passive (check on access) + Active (background ticker sampling).

### 4. Persistence (RDB & AOF)
- **RDB**:
    - **Format**: Binary dump. Magic Header (`REDIS`), Version, Metadata, DB Selector, Key-Value pairs with Type byte, EOF, Checksum.
    - **Strategy**: `SAVE` (blocking) for simplicity initially.
- **AOF**:
    - **Format**: RESP protocol text stream.
    - **Strategy**: Write commands to `appendonly.aof` immediately after execution.
    - **Replay**: Read AOF line-by-line via RESP parser and re-execute commands on startup.

## Phase 1: Data Model & Contracts

### Core Interfaces

```go
// internal/core/engine.go

type Engine interface {
    Exec(cmd *Command) *Result
    Close() error
}

type Command struct {
    Name string
    Args [][]byte
}

type Result struct {
    Value interface{} // To be serialized to RESP
    Error error
}
```

### Protocol Contract

```go
// internal/resp/decoder.go

type Decoder struct {
    r *bufio.Reader
}

func (d *Decoder) Decode() (Value, error)
```

## Phase 2: Implementation Roadmap

See `tasks.md` for granular breakdown.
1. **Skeleton**: TCP Server + PING/PONG.
2. **Protocol**: Full RESP Parser/Writer.
3. **Engine**: In-memory SET/GET + Expiry.
4. **Data Types**: List -> Hash -> Set -> ZSet.
5. **Persistence**: RDB Load/Save + AOF Log/Replay.
6. **Replication**: Master/Replica handshake.
