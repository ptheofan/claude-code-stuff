# [Feature/System Name] - Technical Design Document

**Author:** [Name]  
**Date:** [YYYY-MM-DD]  
**Status:** Draft | In Review | Approved

---

## 1. Problem Statement

What problem are we solving? Why now?

## 2. Goals & Non-Goals

### Goals
- 

### Non-Goals
- 

## 3. Proposed Solution

High-level overview of the approach.

### 3.1 Architecture

How does this fit into Clean Architecture? Which layers are affected?

```
[Diagram or description of component relationships]
```

### 3.2 Module Boundaries

Which modules are involved? What APIs are exposed between them?

| Module | Exposes | Consumes |
|--------|---------|----------|
|        |         |          |

### 3.3 Data Model

Key entities, types, and relationships. Schema definitions only.

```typescript
// Example: entity type
interface User {
  id: string;
  email: string;
  createdAt: Date;
}
```

### 3.4 API Design

Endpoints/methods, inputs, outputs. Signatures only — no implementation.

```typescript
// Example: interface or type definition
interface CreateUserInput {
  email: string;
  name: string;
}

interface CreateUserOutput {
  id: string;
  createdAt: Date;
}
```

### 3.5 Error Handling

Domain errors to introduce:

| Error Class | When Thrown | Data Included |
|-------------|-------------|---------------|
|             |             |               |

## 4. Alternatives Considered

| Option | Pros | Cons | Verdict |
|--------|------|------|---------|
|        |      |      |         |

## 5. Testing Strategy

### Unit Tests
- 

### Integration Tests
- 

### E2E Tests (if applicable)
- 

## 6. Migration / Rollout Plan

How do we ship this safely?

- [ ] Feature flag?
- [ ] Database migrations?
- [ ] Backward compatibility?

## 7. Open Questions

- 

## 8. References

- Related docs, tickets, prior art
