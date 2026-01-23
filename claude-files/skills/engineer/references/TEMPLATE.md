# [Sub-Feature Name] - Technical Design

**Author:** [Name]
**Date:** [YYYY-MM-DD]
**Status:** Draft | In Review | Approved
**Parent TDD:** `./docs/features/<NNN>-<feature-name>.tdd.md`
**Sub-Feature ID:** [NNN][suffix] (e.g., 001a)

---

## 1. Overview

### Purpose
Brief description of what this sub-feature accomplishes.

### Scope
| In Scope | Not In Scope |
|----------|--------------|
|          |              |

---

## 2. Parent TDD References

| Aspect | Parent TDD Section | Notes |
|--------|-------------------|-------|
| Module | §2 Architecture | |
| Schema | §3 Data Model | |
| Errors | §5 Error Handling | |
| API | §4 API Design | |

---

## 3. Types & Interfaces

Sub-feature specific types only. Reference parent TDD for shared types.

```typescript
// Example: Password validation specific
interface PasswordValidationResult {
  valid: boolean;
  errors: PasswordValidationError[];
}

interface PasswordValidationError {
  code: string;
  message: string;
}
```

---

## 4. Function Contracts

### [Function/Method Name]

| Aspect | Details |
|--------|---------|
| **Purpose** | What this function does |
| **Location** | Module/class where it belongs |
| **Signature** | `functionName(param: Type): ReturnType` |

**Preconditions:**
-

**Postconditions:**
-

**Throws:**
- `ErrorClass` - When condition

---

## 5. Domain Errors

| Error Class | When Thrown | Data Included |
|-------------|-------------|---------------|
|             |             |               |

---

## 6. Implementation Notes

Approach guidance (not code) for implementing within parent's architecture.

### Key Decisions
-

### Constraints
-

---

## 7. Dependencies

### Internal (from parent architecture)
| Dependency | Purpose |
|------------|---------|
|            |         |

### External (if any)
| Dependency | Purpose | Already in parent TDD? |
|------------|---------|----------------------|
|            |         | Yes/No               |

---

## 8. Open Questions

| # | Question | Answer |
|---|----------|--------|
| 1 |          |        |

---

## Sign-off

- [ ] Aligns with parent TDD architecture
- [ ] Types and contracts complete
- [ ] Domain errors defined
- [ ] Ready for Test Design (`/test-design`)
