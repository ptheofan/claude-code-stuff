---
name: engineer
version: 1.0.0
description: Lightweight technical design for sub-features after /breakdown. This skill should be used when the user asks to "design sub-feature", "create sub-feature TDD", "detail implementation approach", or needs to define types, function contracts, and domain errors for a specific sub-feature slice while staying within the parent TDD's architecture decisions. Use after /breakdown, before /test-design.
---

# Sub-Feature Technical Design

Create focused technical design for a sub-feature, staying within architecture decisions already made by the parent TDD.

## Workflow Position

```
/feature → /tdd → /breakdown → /engineer (per sub) → /test-design → /coder → /code-review → /qa
```

## Process

1. **Read parent TDD** - `./docs/features/<NNN>-<feature-name>.tdd.md`
2. **Identify sub-feature scope** - What slice of the parent feature is this?
3. **Interview if needed** - Minimal questions for implementation clarity only
4. **Design execution details** - Types, contracts, errors specific to this sub-feature
5. **Save file** - `./docs/features/<NNN><suffix>-<sub-name>.tdd.md`

## File Naming

Use parent number + letter suffix:
```
./docs/features/001-user-authentication.tdd.md      ← Parent TDD
./docs/features/001a-password-validation.tdd.md     ← Sub-feature
./docs/features/001b-session-management.tdd.md      ← Another sub-feature
```

## What's Already Decided (Parent TDD)

Do NOT re-decide these - reference parent TDD:

- Module boundaries and structure
- Overall architecture and patterns
- Database schema and migrations
- Authentication strategy
- External integrations approach
- Error handling patterns

## What This Skill Defines

Focus on execution details within parent's architecture:

- **Types/Interfaces** specific to this sub-feature slice
- **Function signatures and contracts** - inputs, outputs, error conditions
- **Domain errors** this sub-feature throws
- **Implementation approach notes** - how to execute within parent's design

## Interview Tool

Use sparingly - only for implementation clarity. Most decisions should already be in parent TDD.

Example:
```
For password validation, which approach should we take?
1. Inline validation in service (simple)
2. Separate PasswordValidator class (reusable)
3. Other (specify)
```

## Interview Questions

Only ask about:

### Implementation Details
- How should [specific function] handle [specific case]?
- Which pattern fits this slice better: [A] vs [B]?

### Edge Cases
- What happens when [X] in this specific slice?
- Any special handling needed for [Y]?

## Critical Rules

- **Defer to parent TDD** - Architecture decisions are made, follow them
- **Stay in scope** - Only design this sub-feature, not adjacent ones
- **Reference, don't repeat** - Point to parent TDD sections, don't copy
- **No implementation code** - Types, signatures, contracts only
- **Minimal questions** - Most answers should already exist in parent TDD

## Design Principles

Apply these from CLAUDE.md:

- **Clean Architecture** - Import down, emit events up
- **Strict module boundaries** - Only use exposed APIs
- **Domain Errors** - Specific exception classes with traceable data
- **SOLID, KISS, DRY** - Simple, readable, no duplication

## Output

Save to `./docs/features/<NNN><suffix>-<sub-name>.tdd.md` using template in `assets/TEMPLATE.md`.

## Content Rules

**Include:**
- Sub-feature specific types and interfaces
- Function contracts (signature, purpose, pre/post conditions, throws)
- Domain errors specific to this slice
- Dependencies (what this sub-feature needs from parent architecture)
- Implementation approach notes

**Exclude:**
- Architecture decisions (in parent TDD)
- Database schema (in parent TDD)
- Module structure (in parent TDD)
- Implementation code
- Business logic snippets

## Next Step

After sub-feature TDD is approved → `/test-design`
