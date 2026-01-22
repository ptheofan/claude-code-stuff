---
name: code-review
description: Code review against project standards. Use when reviewing PRs, checking code quality, or validating implementations against architecture rules.
---

# Code Review Checklist

Review code against CLAUDE.md standards and project conventions.

## Architecture Review

### Clean Architecture
- [ ] Dependencies flow downward only
- [ ] Events/callbacks emit upward
- [ ] No circular dependencies

### Module Boundaries
- [ ] No direct imports of another module's entities
- [ ] No direct imports of another module's repositories
- [ ] Only exposed APIs used (services, DTOs, interfaces)

### Centralized Integrations
- [ ] External API calls go through dedicated service/class
- [ ] No scattered HTTP calls throughout codebase

## Code Quality

### TypeScript
- [ ] No `any` types (except catch blocks when necessary)
- [ ] No `@ts-ignore` or `@ts-expect-error` without justification
- [ ] No suppressed warnings
- [ ] Strict mode compliant
- [ ] Types/interfaces in separate files

### Simplicity (KISS)
- [ ] Code is readable without mental gymnastics
- [ ] No over-engineering or premature abstraction
- [ ] Clear variable and function names
- [ ] No clever one-liners that sacrifice readability

### Comments
- [ ] Only present for complex/non-obvious logic
- [ ] No commented-out code
- [ ] No redundant comments explaining obvious code

## Error Handling

### Domain Errors
- [ ] Specific exception classes used (not generic Error/throw)
- [ ] Exceptions include traceable data (IDs, context)
- [ ] Errors bubble to appropriate boundaries

### Example Check
```typescript
// ❌ Bad
throw new Error('User not found');

// ✅ Good
throw new UserNotFoundException(userId);
```

## Testing

### Test Pyramid
- [ ] Unit tests for smallest functions first
- [ ] Higher-level tests stub already-tested dependencies
- [ ] Integration tests are self-contained (no production data)

### Controller/Resolver Tests
- [ ] Focus on input validation
- [ ] Focus on output shape
- [ ] NOT testing business logic (that's service layer)

### Test Quality
- [ ] Meaningful coverage (not percentage gaming)
- [ ] Edge cases covered
- [ ] Error scenarios tested
- [ ] Logger stubbed (noise-free output)

## Performance & Reliability

### Database
- [ ] N+1 queries identified and fixed
- [ ] Eager loading used appropriately
- [ ] Pagination on list endpoints

### APIs
- [ ] Idempotent where applicable
- [ ] Proper error responses

## Pre-Merge Checklist

Before approving:

- [ ] Tests pass
- [ ] No TS errors/warnings
- [ ] No console.logs or debug code
- [ ] Domain errors used
- [ ] No circular dependencies
- [ ] Module boundaries respected
- [ ] Migrations included if schema changed
- [ ] No deprecated code/APIs

## Review Response Template

```markdown
## Summary
[One-line summary of the change]

## ✅ Good
- [What's done well]

## ⚠️ Concerns
- [Issues that need addressing]

## 💡 Suggestions
- [Optional improvements]

## Verdict
[ ] ✅ Approve
[ ] 🔄 Request changes
[ ] 💬 Comment only
```
