# Task List: Build Redis-Core (Non-Relational Engine)

**Branch**: `001-build-redis-core` | **Spec**: [Feature Spec](spec.md) | **Plan**: [Implementation Plan](plan.md)

## Phase 1: Foundation (Protocol & Server)

- [ ] **1.1. Project Initialization** <!-- id: 1.1 -->
  - **Context**: Set up the Go project structure and entry point.
  - **Inputs**: `plan.md` (Project Structure).
  - **Output**: `redis/` directory with `go.mod`, `cmd/server/main.go`, `internal/server`, `internal/resp`.
  - **Step-by-step**:
    1. Initialize `go.mod` (e.g., `github.com/yourname/zdb`).
    2. Create directory `redis/` and subdirectories defined in Plan.
    3. Create `redis/cmd/server/main.go` to verify build.

- [ ] **1.2. RESP Value Types Definition** <!-- id: 1.2 -->
  - **Context**: Define the data structures for Redis Serialization Protocol (RESP).
  - **Inputs**: `spec.md` (FR-002), RESP Protocol research.
  - **Output**: `redis/internal/resp/value.go`
  - **Step-by-step**:
    1. Create `Value` struct with `Type` (byte) and content fields.
    2. Define constants for RESP types.

- [ ] **1.3. RESP Reader Implementation** <!-- id: 1.3 -->
  - **Context**: Implement the stream parser for RESP messages.
  - **Inputs**: `plan.md` (Research section).
  - **Output**: `redis/internal/resp/reader.go`
  - **Step-by-step**:
    1. Create `NewReader(io.Reader)`.
    2. Implement `ReadLine()`, `ReadInteger()`, `ReadArray()`, `ReadBulk()`.
    3. Unit Test with raw byte strings.

- [ ] **1.4. RESP Writer Implementation** <!-- id: 1.4 -->
  - **Context**: Implement the serializer to send responses back to clients.
  - **Inputs**: `spec.md` (FR-002).
  - **Output**: `redis/internal/resp/writer.go`
  - **Step-by-step**:
    1. Create `NewWriter(io.Writer)`.
    2. Implement `Write(v Value)`.
    3. Unit Test.

- [ ] **1.5. TCP Server & Event Loop** <!-- id: 1.5 -->
  - **Context**: Bind port 6379 and accept connections.
  - **Inputs**: `spec.md` (FR-001, FR-003).
  - **Output**: `redis/internal/server/server.go`
  - **Step-by-step**:
    1. Implement `net.Listen("tcp", ":6379")`.
    2. Loop `listener.Accept()` and spawn `go handleConnection(conn)`.
    3. In `handleConnection`: Create RESP Reader/Writer, loop read/write.
    4. Integration Test with `nc`.

## Phase 2: Core Engine & KV Store

- [ ] **2.1. In-Memory Store & SET/GET** <!-- id: 2.1 -->
  - **Context**: Implement the actual key-value storage.
  - **Inputs**: `spec.md` (FR-004), `plan.md` (Store struct).
  - **Output**: `redis/internal/core/store.go`, `redis/internal/core/eval.go`
  - **Step-by-step**:
    1. Define `Store` struct with map and mutex.
    2. Implement `put` and `get` handlers.
    3. Connect connection handler to `EvalCommand`.

- [ ] **2.2. Command Routing & PING/ECHO** <!-- id: 2.2 -->
  - **Context**: Standardize command execution dispatch.
  - **Inputs**: `spec.md` (User Story 1).
  - **Output**: `redis/internal/core/handler.go`
  - **Step-by-step**:
    1. Create command map.
    2. Register `PING`, `ECHO`, `SET`, `GET`.
    3. Update Server loop.

- [ ] **2.3. Key Expiration (TTL)** <!-- id: 2.3 -->
  - **Context**: Support volatile keys.
  - **Inputs**: `spec.md` (FR-005).
  - **Output**: `redis/internal/core/store.go` updates.
  - **Step-by-step**:
    1. Update `RedisObject` with `ExpiresAt`.
    2. Update `GET` to check expiration.
    3. Implement `PX` arg in `SET`.

## Phase 3: Advanced Data Structures

- [ ] **3.1. List Type** <!-- id: 3.1 -->
  - **Context**: Linked list implementation.
  - **Inputs**: `spec.md` (FR-008).
  - **Output**: `redis/internal/datastruct/list/`, `redis/internal/core/commands_list.go`
  - **Step-by-step**:
    1. Implement Doubly Linked List.
    2. Implement List commands.

- [ ] **3.2. Hash Type** <!-- id: 3.2 -->
  - **Context**: Dictionary implementation.
  - **Inputs**: `spec.md` (FR-013).
  - **Output**: `redis/internal/datastruct/dict/`, `redis/internal/core/commands_hash.go`
  - **Step-by-step**:
    1. Define Hash structure.
    2. Implement Hash commands.

- [ ] **3.3. Set Type** <!-- id: 3.3 -->
  - **Context**: Unique collection.
  - **Inputs**: `spec.md` (FR-014).
  - **Output**: `redis/internal/datastruct/set/`, `redis/internal/core/commands_set.go`
  - **Step-by-step**:
    1. Define Set structure.
    2. Implement Set commands.

- [ ] **3.4. Sorted Set Type** <!-- id: 3.4 -->
  - **Context**: Scored set using Skiplist.
  - **Inputs**: `spec.md` (FR-010).
  - **Output**: `redis/internal/datastruct/sortedset/skiplist.go`
  - **Step-by-step**:
    1. Implement Skiplist.
    2. Implement Sorted Set commands.

## Phase 4: Persistence

- [ ] **4.1. AOF Implementation** <!-- id: 4.1 -->
  - **Context**: Log every write command.
  - **Inputs**: `spec.md` (FR-006, User Story 3).
  - **Output**: `redis/internal/persistence/aof.go`
  - **Step-by-step**:
    1. Create `Aof` struct.
    2. Write commands to AOF file.
    3. Implement `LoadAof`.

- [ ] **4.2. RDB Loader (Basic)** <!-- id: 4.2 -->
  - **Context**: Read snapshot files.
  - **Inputs**: `spec.md` (FR-006).
  - **Output**: `redis/internal/persistence/rdb.go`
  - **Step-by-step**:
    1. Implement binary parser.
    2. Load keys into `Store`.

## Phase 5: Replication (Optional/Final)

- [ ] **5.1. Master-Replica Handshake** <!-- id: 5.1 -->
  - **Context**: Replica connecting to Master.
  - **Inputs**: `spec.md` (FR-007).
  - **Output**: `redis/internal/replication/`
  - **Step-by-step**:
    1. Implement `REPLCONF`, `PSYNC`.
    2. Implement `slaveof`.
