---
name: coder
description: Implement a feature using TDD approach. Reads TDD and test-design, creates implementation plan, writes tests first then code.
---

# Coder

Implement a feature following TDD methodology. Tests drive the code.

## Process

1. **Read the docs:**
   - `./docs/features/<NNN>-<feature-name>.feature.md` - What we're building
   - `./docs/features/<NNN>-<feature-name>.tdd.md` - How to build it
   - `./docs/features/<NNN>-<feature-name>.test.md` - Test specifications

2. **Create implementation plan** - Use plan mode to break down work

3. **Implement using TDD:**
   - Write test first (from test-design)
   - Run test → verify it fails
   - Write minimal code to pass
   - Refactor if needed
   - Repeat

4. **Follow test pyramid order:**
   - Unit tests for utilities first
   - Then services (stub utilities)
   - Then handlers (stub services)
   - Then controllers (stub handlers)
   - Integration tests
   - E2E tests last

## Critical Rules

### TDD Discipline
- **Tests drive code** - Code matches test specs, not vice versa
- **Red → Green → Refactor** - See test fail before making it pass
- **No code without test** - Every feature has corresponding test

### Code Quality (from CLAUDE.md)
- **Clean Architecture** - Import down, emit events up
- **Strict module boundaries** - Never access another module's internals
- **Domain Errors** - Specific exception classes with traceable data
- **No `any`** - Proper TypeScript types
- **No `@ts-ignore`** - Fix the issue, don't suppress it
- **SOLID** - Single responsibility, Open/closed, Liskov, Interface segregation, Dependency inversion
- **KISS** - Keep it simple
- **DRY** - Don't repeat yourself

### Testing Standards
- **Stub loggers** - Keep test output noise-free
- **Meaningful coverage** - Critical paths, not percentage gaming
- **Controller tests** - Input validation and output shape only
- **Integration tests** - Self-contained, no production data

## Interview Tool

Use the interview tool when uncertain:

```
This service could use Strategy pattern or simple if/else. Which do you prefer?
1. Strategy pattern (more extensible)
2. Simple if/else (KISS, fewer files)
3. Let me explain the context more
```

## Pre-Commit Checklist

Before committing, verify:
- [ ] All tests pass
- [ ] No TS errors/warnings
- [ ] No console.logs or debug code
- [ ] Domain errors used (no generic throws)
- [ ] No circular dependencies
- [ ] Module boundaries respected
- [ ] Migrations included if schema changed

## Next Step

After implementation → `/code-review`
