<!--
Sync Impact Report:
- Version change: 1.0.0 -> 1.1.0
- List of modified principles: Added User-Driven Implementation, Added Collaboration Protocol
- Added sections: Collaboration Protocol
- Removed sections: None
- Templates requiring updates: Ensure tasks in .specify/templates/tasks-template.md allow for user implementation.
- Follow-up: None.
-->

# Project Constitution: Z-DB

> **Mission**: To build a comprehensive database learning project in Go, exploring both non-relational (Redis-like) and relational (MySQL-like) paradigms to understand database internals.

## 1. Metadata & Governance

| Attribute | Value |
| :--- | :--- |
| **Project Name** | Z-DB |
| **Constitution Version** | 1.1.0 |
| **Ratification Date** | 2025-12-09 |
| **Last Amended Date** | 2025-12-09 |

### Amendment Policy
1. **Proposal**: Changes to this constitution must be submitted via Pull Request.
2. **Versioning**:
   - **MAJOR**: Fundamental changes to the project mission or technology stack.
   - **MINOR**: Addition of new principles or significant process changes.
   - **PATCH**: Clarifications, formatting, or minor adjustments.
3. **Compliance**: All project specifications and code must verify alignment with the current constitution.

## 2. Core Principles

### 2.1. Learning Over Performance
**Rule**: Prioritize code clarity, readability, and conceptual understanding over raw performance or complex micro-optimizations.
**Rationale**: The primary objective is education. While performance is important in databases, it should not compromise the ability to understand the underlying mechanisms.

### 2.2. Idiomatic Go (Go Way)
**Rule**: Code MUST adhere to standard Go conventions ("Effective Go"). Use standard formatting (`gofmt`), error handling patterns, and project layout conventions.
**Rationale**: Correct usage of the language tools and patterns is essential for a maintainable and readable codebase in Go.

### 2.3. Zero-Dependency Core
**Rule**: Avoid external libraries for core database functionalities (storage, indexing, networking) unless absolutely necessary.
**Rationale**: Implementing these from scratch (using only the Go standard library) provides the deepest learning experience about how they work.

### 2.4. Dual-Paradigm Architecture
**Rule**: The system architecture SHOULD support modularity to allow switching or implementing different engines (KV Store vs. Relational Table).
**Rationale**: To fulfill the goal of supporting both Redis-like and MySQL-like features, the core components (network handling, parsing) should be decoupled from the execution engine.

### 2.5. Test-Driven Confidence
**Rule**: Every feature MUST include unit tests. Integration tests are REQUIRED for network and persistence layers.
**Rationale**: Database systems require rigorous correctness. Testing is the only way to verify that custom implementations of complex algorithms (like B-Trees or WAL) function correctly.

### 2.6. User-Driven Implementation
**Rule**: The primary code implementation MUST be written by the user. The AI Assistant acts as a guide, providing analysis, specifications, and teaching support, but avoids directly generating full solution code for core logic unless explicitly requested for debugging or educational illustration.
**Rationale**: Deep learning and mastery of database internals requires the user to personally grapple with the implementation challenges.

## 3. Implementation Guidelines

### 3.1. Phase 1: Non-Relational (Redis-like)
- Implement an in-memory Key-Value store.
- Support basic commands (SET, GET, DEL).
- Implement a custom protocol or RESP (Redis Serialization Protocol).

### 3.2. Phase 2: Relational (MySQL-like)
- Implement a SQL parser (subset).
- Implement a persistent storage engine (e.g., B+ Tree).
- Implement basic table operations (INSERT, SELECT, UPDATE).

## 4. Documentation Standards
**Rule**: Code must be documented with "Why" and "How" comments for complex logic.
**Rationale**: As a learning project, the documentation serves as the study notes for the author and future readers.

## 5. Collaboration Protocol

### 5.1. AI as PM & Analyst
**Role**: The AI assumes the role of Product Manager and Requirements Analyst.
**Responsibility**: Proactively define feature requirements, break down specifications into manageable tasks, and ensure the project roadmap is followed.

### 5.2. AI as Teacher
**Role**: The AI assumes the role of a Technical Mentor/Teacher.
**Responsibility**: When the user encounters technical difficulties, provide conceptual explanations, pseudo-code strategies, and debugging assistance. Encourage the user to solve the problem to maximize learning.
