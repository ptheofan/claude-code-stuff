---
name: code-review
version: 1.0.0
description: Review code changes using git diff against quality standards. This skill should be used when the user asks to "review code", "check my changes", "validate implementation", "review PR", "code review", or needs to validate that implementation meets architecture, code quality, testing, and error handling standards before merging. Use after /coder, before /qa.
---

# Code Review

Review code changes against project standards and best practices.

## Process

1. **Get the diff** - Run `git diff` or `git diff --staged` to see changes
2. **Review against standards** - Check each item in the checklist
3. **Report findings** - Summarize issues and suggestions
4. **Request changes or approve**

## Review Checklist

### Architecture
- [ ] **Clean Architecture** - Dependencies flow downward only
- [ ] **Module boundaries** - No direct imports of another module's internals
- [ ] **Centralized integrations** - External calls through dedicated services
- [ ] **No circular dependencies**

### TypeScript Quality
- [ ] **No `any`** - Proper types used
- [ ] **No `@ts-ignore`** - Issues fixed, not suppressed
- [ ] **No suppressed warnings**
- [ ] **Types/interfaces in separate files**

### Code Quality
- [ ] **SOLID** - Single responsibility, proper abstractions
- [ ] **KISS** - Simple, readable solutions
- [ ] **DRY** - No code duplication
- [ ] **Meaningful names** - Clear variable/function names
- [ ] **Minimal comments** - Only for complex/non-obvious logic
- [ ] **No console.logs** - Debug code removed

### Error Handling
- [ ] **Domain Errors** - Specific exception classes used
- [ ] **Traceable data** - Errors include IDs/context for debugging
- [ ] **Errors bubble** - Handled at appropriate boundaries

### Testing
- [ ] **Tests exist** - New code has corresponding tests
- [ ] **Test pyramid followed** - Unit → Integration → E2E
- [ ] **Meaningful coverage** - Critical paths covered
- [ ] **Controller tests** - Focus on input/output, not business logic
- [ ] **Logger stubbed** - Test output is noise-free

### Database
- [ ] **N+1 awareness** - Queries optimized
- [ ] **Migrations included** - Schema changes have migrations
- [ ] **Pagination** - List endpoints paginated

### Performance & Reliability
- [ ] **Idempotency** - APIs safely retryable where applicable

## Review Output Format

```markdown
## Code Review: [Feature/PR Name]

### Summary
[One-line summary]

### ✅ Good
- [What's done well]

### ❌ Issues (Must Fix)
- [ ] [Critical issues that block merge]

### ⚠️ Suggestions (Optional)
- [ ] [Improvements that would be nice]

### Verdict
- [ ] ✅ Approved
- [ ] 🔄 Request changes
```

## Commands

```bash
# Review staged changes
git diff --staged

# Review all uncommitted changes
git diff

# Review specific file
git diff path/to/file.ts

# Review against main branch
git diff main...HEAD
```

## Next Step

After code review passes → `/qa`
