---
name: coder
version: 2.0.0
description: Implement a feature using TDD methodology. This skill should be used when the user asks to "implement feature", "start coding", "write the code", "build the feature", or is ready to implement after test-design is complete. ALWAYS starts in planning mode. Reads TDD and test-design, then implements code that satisfies the tests. Use after /test-design, before /code-review.
---

# Coder

> **⚠️ ALWAYS START IN PLANNING MODE.** Enter plan mode before writing any code.

Implement a feature using Test-Driven Development. Code is written to satisfy the tests.

## Workflow Position

```
/feature → /tdd → /breakdown → /engineer → /test-design → /coder → /code-review → /qa
                                                             ↑
                                                            HERE
```

## Process

1. **ENTER PLANNING MODE** - Always plan before coding
2. **Read the documents:**
   - Feature spec: `./docs/features/<NNN>-<feature-name>.feature.md`
   - TDD (or mini-TDD): `./docs/features/<NNN>-<feature-name>.tdd.md`
   - Test design: `./docs/features/<NNN>-<feature-name>.test.md`
3. **Create implementation plan** - Break down the work in plan mode
4. **Get plan approved** - Exit plan mode only after approval
5. **Implement using TDD** - Red → Green → Refactor

## Planning Mode (MANDATORY)

**Always enter planning mode first.** The plan should include:

- Order of implementation (which tests first)
- Files to create/modify
- Dependencies between components
- Potential risks or blockers

Only proceed to implementation after plan is approved.

## TDD Implementation Cycle

For each test case from test-design:

```
1. RED    → Write the test (from test-design spec)
2. RED    → Run test → Verify it FAILS
3. GREEN  → Write MINIMAL code to make test pass
4. GREEN  → Run test → Verify it PASSES
5. REFACTOR → Clean up code (if needed)
6. REPEAT → Next test case
```

## Implementation Order

Follow the test pyramid from bottom up:

1. **Unit tests first**
   - Utilities/helpers
   - Then services (stub utilities)
   - Then handlers (stub services)
   - Then controllers (stub handlers)

2. **Integration tests**
   - Database operations
   - Module interactions

3. **E2E tests last**
   - Critical paths only

## Code Quality Principles

Apply these from CLAUDE.md:

### Architecture
- **Clean Architecture** - Import down, emit events up
- **Strict module boundaries** - Never access another module's internals
- **Centralized integrations** - External calls through dedicated services

### Code Standards
- **SOLID** - Single responsibility, Open/closed, Liskov, Interface segregation, Dependency inversion
- **KISS** - Keep it simple, stupid
- **DRY** - Don't repeat yourself

### TypeScript
- **No `any`** - Proper types always
- **No `@ts-ignore`** - Fix issues, don't suppress
- **Types in separate files** - Keep them organized

### Error Handling
- **Domain Errors** - Specific exception classes
- **Traceable data** - Include IDs/context for debugging
- **Throw early, catch late** - Handle at boundaries

## Interview Tool

Use when implementation decisions arise:

```
This service could use Strategy pattern or simple if/else. Which do you prefer?
1. Strategy pattern (more extensible)
2. Simple if/else (KISS, fewer files)
3. Let me explain the context more
```

## Pre-Commit Checklist

Before completing, verify:
- [ ] All tests pass
- [ ] No TS errors/warnings
- [ ] No console.logs or debug code
- [ ] Domain errors used (no generic throws)
- [ ] No circular dependencies
- [ ] Module boundaries respected
- [ ] Migrations included if schema changed
- [ ] Code matches the TDD design

## Critical Rules

- **ALWAYS plan first** - Enter planning mode before coding
- **Tests drive code** - Code exists to make tests pass
- **Red before green** - See test fail before making it pass
- **Minimal code** - Write just enough to pass the test
- **Follow the TDD** - Implementation matches the design

## Output

Working code that:
- Passes all test cases from test-design
- Follows the architecture from TDD
- Adheres to CLAUDE.md principles

## Next Step

After implementation complete, prompt the user using AskUserQuestion:

```
Implementation complete. How would you like to proceed?
1. Clear memory and continue with /code-review (Recommended)
2. Continue with /code-review
3. Other
```

- If **Option 1**: Inform user to clear context, then invoke /code-review
- If **Option 2**: Proceed directly to /code-review
- If **Option 3**: Follow user's instructions
