---
name: test-design
version: 2.0.0
description: Design behavioral test cases that will drive TDD implementation. This skill should be used when the user asks to "design test cases", "plan tests", "create test specifications", "write test scenarios", or needs to create test cases before coding. Tests describe expected behavior - code is written to satisfy tests, NOT tests to satisfy code. Use after /tdd or /engineer, before /coder.
---

# Test Design

Generate meaningful behavioral test cases that will drive Test-Driven Development.

## Workflow Position

```
/feature → /tdd → /breakdown → /engineer → /test-design → /coder → /code-review → /qa
                                                ↑
                                               HERE
```

## Core Philosophy

> **Tests drive development. Write code to satisfy tests, NOT tests to satisfy code.**

Code can have bugs. Tests define the expected behavior. The coder will implement code that makes these tests pass.

## Process

1. **Determine which TDD to read:**
   - If sub-feature exists: Read `./docs/features/<NNN><suffix>-<sub-name>.tdd.md`
   - If no breakdown: Read `./docs/features/<NNN>-<feature-name>.tdd.md`
2. **Read the feature spec** - `./docs/features/<NNN>-<feature-name>.feature.md`
3. **Design behavioral test cases** - What should happen when X occurs?
4. **Interview for edge cases** - Use interview tool for priorities
5. **Save file** - `./docs/features/<NNN>-<feature-name>.test.md` (or `<NNN><suffix>-<sub-name>.test.md`)

## Behavioral Test Case Format

Test cases describe BEHAVIOR, not implementation:

### Format
```
WHEN [action/trigger occurs]
THEN [expected outcome should happen]
GIVEN [preconditions, if any]
```

### Examples

**UI Behavior:**
```
WHEN user clicks "Submit" button with valid form data
THEN form should be submitted and success message displayed
GIVEN user is logged in and on the registration page
```

**API Behavior:**
```
WHEN POST /api/users is called with valid payload
THEN 201 response with user object should be returned
GIVEN authentication token is valid
```

**Business Logic:**
```
WHEN password validation is triggered
THEN error should be returned if password < 8 characters
GIVEN password input is provided
```

**Event/Trigger Behavior:**
```
WHEN order status changes to "shipped"
THEN email notification should be sent to customer
GIVEN customer has email notifications enabled
```

## Test Categories

### Unit Tests
Test individual functions/methods in isolation:
- Pure functions, utilities, helpers
- Service methods (mock dependencies)
- Validators, transformers
- Domain logic

### Integration Tests
Test components working together:
- Database operations
- Module-to-module interactions
- Service-to-repository interactions

### E2E Tests (Selective)
Test complete user flows:
- Critical happy paths only
- Business-critical failure scenarios
- Skip what's covered by unit/integration

## Test Case Template

For each test case:

| Field | Description |
|-------|-------------|
| **ID** | TC-001, TC-002, etc. |
| **Type** | Unit / Integration / E2E |
| **When** | The action or trigger |
| **Then** | Expected outcome |
| **Given** | Preconditions (if any) |
| **Priority** | High / Medium / Low |

## Interview Tool

Use to clarify test priorities:

```
Which scenarios are most critical to test?
1. Happy path - successful user registration
2. Validation failures - invalid input handling
3. Edge cases - concurrent access, timeouts
4. Other (specify)
```

```
For error scenarios, what's the expected behavior?
1. Return specific error code and message
2. Log error and return generic message
3. Throw exception to be caught upstream
4. Other (specify)
```

## File Naming

Match the TDD file:
```
# If no breakdown:
./docs/features/001-user-authentication.tdd.md
./docs/features/001-user-authentication.test.md   ← Same number

# If sub-feature:
./docs/features/001a-password-validation.tdd.md
./docs/features/001a-password-validation.test.md  ← Same suffix
```

## Critical Rules

- **Behavior, not implementation** - Describe WHAT should happen, not HOW
- **Tests come first** - These drive the code, not the other way around
- **Meaningful coverage** - Prioritize critical paths over percentage
- **Domain errors tested** - Verify correct exceptions are thrown
- **Edge cases covered** - Boundaries, nulls, empties, limits

## Output

Save to `./docs/features/<NNN>-<feature-name>.test.md` using template in `assets/TEMPLATE.md`.

**Include:**
- Test cases grouped by type (Unit, Integration, E2E)
- Clear WHEN/THEN/GIVEN for each case
- Priority for each test
- Coverage of happy paths, error paths, edge cases

## Test Design Checklist

Before completing, verify:
- [ ] Happy paths covered
- [ ] Error scenarios defined
- [ ] Edge cases identified
- [ ] Domain errors have test cases
- [ ] Test priorities assigned
- [ ] Behavioral format used (WHEN/THEN/GIVEN)

## Next Step

After test design is approved, prompt the user using AskUserQuestion:

```
Test design complete. How would you like to proceed?
1. Clear memory and continue with /coder (Recommended)
2. Continue with /coder
3. Other
```

- If **Option 1**: Inform user to clear context, then invoke /coder
- If **Option 2**: Proceed directly to /coder
- If **Option 3**: Follow user's instructions

The coder will implement code to make these tests pass.
