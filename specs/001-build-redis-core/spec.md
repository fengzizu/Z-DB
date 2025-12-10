# Feature Specification: Build Redis-Core (Non-Relational Engine)

**Feature Branch**: `001-build-redis-core`
**Created**: 2025-12-09
**Status**: Draft
**Input**: User description: "关于非关系性数据库，我们可以模仿https://app.codecrafters.io/courses/redis/overview 这个链接里的内容去构建它，这个链接里是一个收费的教程，我们可以借鉴它有的所有功能"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Basic Server & Protocol (Priority: P1)

As a database client (e.g., redis-cli), I want to connect to the server via TCP and exchange basic commands using the standard RESP protocol, so that I can verify connectivity and protocol compliance.

**Why this priority**: This is the foundation of the entire system. Without a working TCP server and protocol parser, no other features can exist.

**Independent Test**: Can be fully tested by using `redis-cli` or `netcat` to send PING/ECHO and receiving correct responses.

**Acceptance Scenarios**:

1. **Given** the server is running on port 6379, **When** a client connects and sends `PING`, **Then** the server responds with `PONG`.
2. **Given** the server is running, **When** a client sends `ECHO "hello world"`, **Then** the server responds with `"hello world"`.
3. **Given** multiple clients connecting simultaneously, **When** they send commands, **Then** the server handles them concurrently without blocking.

---

### User Story 2 - Key-Value Storage & Expiry (Priority: P1)

As a developer using the database, I want to store and retrieve string values with optional expiration times, so that I can use the system for caching and basic state management.

**Why this priority**: This provides the core utility of a Key-Value store.

**Independent Test**: Verify SET stores data and GET retrieves it. Verify keys disappear after their TTL expires.

**Acceptance Scenarios**:

1. **Given** an empty database, **When** client sends `SET mykey "value"`, **Then** server responds `OK`.
2. **Given** `mykey` exists, **When** client sends `GET mykey`, **Then** server returns `"value"`.
3. **Given** a key set with `PX 100` (100ms expiry), **When** client waits 200ms and sends `GET mykey`, **Then** server returns `(nil)`.

---

### User Story 3 - Persistence via RDB & AOF (Priority: P2)

As a system administrator, I want the database to support both RDB (snapshots) and AOF (append-only file) persistence, so that I can choose between faster recovery (RDB) or maximum data durability (AOF).

**Why this priority**: Data persistence is critical. RDB is standard, but AOF is preferred for minimizing data loss.

**Independent Test**:
1. RDB: Create valid RDB, start server, verify keys.
2. AOF: Perform writes, restart server, verify AOF replay restores keys.

**Acceptance Scenarios**:

1. **Given** an RDB file, **When** server starts, **Then** data is loaded.
2. **Given** AOF enabled, **When** writes occur, **Then** commands are appended to `appendonly.aof`.
3. **Given** a restart with AOF file present, **When** server starts, **Then** it replays commands to restore state.
4. **Given** mixed config, **When** both files exist, **Then** AOF takes precedence (standard Redis behavior).

---

### User Story 4 - Replication (Priority: P2)

As a system administrator, I want to configure a server as a replica of another master, so that I can scale reads and ensure high availability.

**Why this priority**: Replication is a key distributed system concept covered in the curriculum.

**Independent Test**: Configure a master and a replica. Write to master, read from replica.

**Acceptance Scenarios**:

1. **Given** a master and a replica, **When** the replica connects, **Then** it completes the handshake (PING, REPLCONF, PSYNC).
2. **Given** a synced replica, **When** a write command is sent to master, **Then** the replica receives and applies the write propagation.
3. **Given** a replica, **When** `INFO replication` is called, **Then** it reports its role and master connection status.

---

### User Story 5 - Advanced Structures (Streams, Lists, Sets, Hashes, Sorted Sets) (Priority: P3)

As a developer, I want to use advanced data structures like Streams, Hashes, Sets, and Sorted Sets, so that I can model complex data patterns (e.g., event logs, object storage, unique collections, leaderboards).

**Why this priority**: Expands utility but depends on the core KV engine being solid.

**Independent Test**: Use specific commands (`XADD`, `LPUSH`, `HSET`, `SADD`, `ZADD`) and verify structure integrity.

**Acceptance Scenarios**:

1. **Given** a stream key, **When** `XADD` is called with an ID, **Then** the entry is appended and the ID returned.
2. **Given** a hash key, **When** `HSET` field value is called, **Then** `HGET` returns the correct value.
3. **Given** a set key, **When** `SADD` adds members, **Then** `SMEMBERS` returns all unique members.
4. **Given** a sorted set, **When** `ZRANGE` is called, **Then** members are returned in correct score order.

---

### User Story 6 - Transactions & Pub/Sub (Priority: P3)

As a developer, I want to execute atomic transaction blocks and use Publish/Subscribe messaging, so that I can build robust and reactive applications.

**Why this priority**: Adds atomicity and messaging capabilities.

**Independent Test**: Send MULTI, commands, EXEC. Subscribe to a channel and publish to it.

**Acceptance Scenarios**:

1. **Given** a transaction started with `MULTI`, **When** commands are queued and `EXEC` is called, **Then** all commands execute atomically.
2. **Given** a subscriber on channel "news", **When** a publisher sends message to "news", **Then** the subscriber receives the message.

### Edge Cases

- **Malformed Protocol**: System MUST gracefully close connection if client sends invalid RESP.
- **Concurrent Writes**: System MUST ensure thread safety when multiple clients write to the same key (e.g., INCR).
- **Network Partitions**: Replica MUST retry connection if Master becomes unreachable.
- **Large Payloads**: System SHOULD protect against OOM by limiting max command size (e.g., 512MB default).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST implement a TCP server listening on port 6379 (configurable).
- **FR-002**: System MUST parse and serialize messages using the Redis Serialization Protocol (RESP).
- **FR-003**: System MUST support concurrent client connections using Go routines (not single-threaded event loop unless explicitly chosen for learning, but Go conventions suggest concurrency).
- **FR-004**: System MUST implement an in-memory Key-Value store supporting String types.
- **FR-005**: System MUST support Passive and Active expiration of keys.
- **FR-006**: System MUST parse RDB file format (version 6-9 compatible) AND support AOF (Append Only File) persistence with fsync policies (no-appendfsync-on-rewrite not required for v1).
- **FR-007**: System MUST implement Replication logic (Master-Replica handshake, Command Propagation, ACK).
- **FR-008**: System MUST implement List operations (`LPUSH`, `RPUSH`, `LPOP`, `LRANGE`).
- **FR-009**: System MUST implement Stream operations (`XADD`, `XRANGE`, `XREAD` with blocking).
- **FR-010**: System MUST implement Sorted Set operations (`ZADD`, `ZRANGE`, `ZRANK`).
- **FR-011**: System MUST implement Transaction commands (`MULTI`, `EXEC`, `DISCARD`, `WATCH`).
- **FR-012**: System MUST implement Pub/Sub commands (`SUBSCRIBE`, `PUBLISH`).
- **FR-013**: System MUST implement Hash operations (`HSET`, `HGET`, `HGETALL`).
- **FR-014**: System MUST implement Set operations (`SADD`, `SMEMBERS`, `SISMEMBER`).

### Key Entities

- **ClientConnection**: Represents an active TCP connection and its buffering state.
- **Command**: Represents a parsed Redis command with arguments.
- **Store**: The global thread-safe data structure holding the keys and values.
- **Value**: The polymorphic data holder (String, List, Hash, Set, Stream, SortedSet).
- **ReplicaState**: Tracks the synchronization offset and capabilities of a connected replica.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Fully compatible with standard `redis-cli` for all implemented commands.
- **SC-002**: Successfully passes a simulated suite of "CodeCrafters" stage checks (functionally equivalent).
- **SC-003**: Can handle at least 100 concurrent idle connections and 10 active concurrent request streams without panic or deadlock.
- **SC-004**: Persistence system correctly handles both RDB snapshots and AOF log replay, ensuring data consistency after restart.
- **SC-005**: Replication lag is observable via `INFO` command and propagation occurs within sub-second latency for local instances.
