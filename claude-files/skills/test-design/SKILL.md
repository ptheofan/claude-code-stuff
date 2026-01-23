---
name: test-design
description: Design test cases for unit, integration, and E2E testing based on TDD. Use after /tdd to prepare test specifications before implementation.
---

# Test Design

Prepare comprehensive test cases that will drive implementation through TDD.

## Process

1. **Read the TDD** - `./docs/features/<NNN>-<feature-name>.tdd.md`
2. **Read the feature spec** - `./docs/features/<NNN>-<feature-name>.feature.md`
3. **Design test cases** - Unit, integration, E2E
4. **Interview the user** - Use interview tool for edge cases and priorities
5. **Save file** - `./docs/features/<NNN>-<feature-name>.test.md`

## File Naming

Match the feature file number:
```
./docs/features/001-user-authentication.feature.md
./docs/features/001-user-authentication.tdd.md
./docs/features/001-user-authentication.test.md   ← Same number
```

## Test Pyramid (Solitary Unit Testing)

### Unit Tests
Start with smallest functions, stub already-tested dependencies:

1. **Utilities/Helpers** - Pure functions, no dependencies
2. **Services** - Mock repositories, external calls
3. **Command/Query Handlers** - Mock services
4. **Controllers/Resolvers** - Input validation, output shape only

### Integration Tests
- Database operations
- Internal module integrations
- Self-contained (no production data)

### E2E Tests (Selective)
- Critical happy paths only
- Business-critical failure paths
- Skip what's covered by unit/integration

## Interview Questions

Use the interview tool to clarify:

```
Which paths are business-critical for E2E testing?
1. User registration + login flow
2. Payment processing
3. Core feature X
4. Other (specify)
```

```
What edge cases should we prioritize?
1. Network failures
2. Invalid input handling
3. Concurrent access
4. Other (specify)
```

## Test Case Format

For each test case, specify:
- **Description** - What are we testing?
- **Type** - Unit / Integration / E2E
- **Given** - Preconditions
- **When** - Action
- **Then** - Expected outcome
- **Priority** - High / Medium / Low

## Critical Rules

- **Tests drive code** - Write test specs first, code matches tests (not vice versa)
- **Meaningful coverage** - Prioritize critical paths over percentage
- **Stub loggers** - Keep test output noise-free
- **Domain errors tested** - Verify correct exceptions thrown

## Output

Save to `./docs/features/<NNN>-<feature-name>.test.md` using template in `references/TEMPLATE.md`.

## Next Step

After test design is approved → `/coder`
