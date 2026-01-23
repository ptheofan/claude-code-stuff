---
name: breakdown
description: Analyze a TDD and break it into implementable sub-features. Use after /tdd to create a phased implementation plan with clear dependencies.
---

# Feature Breakdown

Analyze a TDD and decompose it into smaller, independently implementable sub-features.

## Workflow Position

```
/feature → /tdd → /breakdown → /engineer (per sub) → /test-design → /coder → /code-review → /qa
```

## Process

1. **Read the TDD** - `./docs/features/<NNN>-<feature-name>.tdd.md`
2. **Identify natural boundaries** - Where can we slice?
3. **Interview the user** - Priorities, constraints, team considerations
4. **Define sub-features** - With dependencies and sequence
5. **Save file** - `./docs/features/<NNN>-<feature-name>.breakdown.md`

## Interview Tool

Use Claude Code's built-in interview tool to clarify priorities and constraints.

Example:
```
Which aspect should we implement first?
1. Core data model and basic CRUD
2. Validation and error handling
3. External integrations
4. Other (specify)
```

## File Naming

Breakdown file uses same number as parent:
```
./docs/features/001-user-authentication.feature.md   ← Feature spec
./docs/features/001-user-authentication.tdd.md       ← Technical design
./docs/features/001-user-authentication.breakdown.md ← This skill's output
```

Sub-features get letter suffixes (created during /engineer):
```
./docs/features/001a-password-validation.tdd.md
./docs/features/001b-session-management.tdd.md
./docs/features/001c-password-reset.tdd.md
```

## Slicing Strategies

### By Layer
- Data model first, then services, then API
- Good for: foundational features

### By User Story
- One story = one sub-feature
- Good for: user-facing features

### By Integration
- Core logic first, external integrations later
- Good for: features with multiple external dependencies

### By Risk
- Risky/uncertain parts first
- Good for: features with unknowns

## Interview Questions

Ask the user about:

### Priorities
- What must ship first?
- What can wait for v1.1?
- Any hard deadlines for specific parts?

### Dependencies
- External systems availability?
- Team dependencies (design, other features)?
- Data migration needs?

### Constraints
- Team size working on this?
- Parallel vs sequential implementation?
- Any parts that MUST be released together?

### Risk
- Which parts are riskiest or most uncertain?
- Should we tackle unknowns first or last?

## Sub-Feature Criteria

Each sub-feature MUST be:

- **Independently testable** - Can write and run tests for just this slice
- **Independently deployable** - Can ship without other sub-features (maybe feature-flagged)
- **Vertically sliced** - Delivers user value, not just "backend for X"
- **Right-sized** - 1-3 days of work for one developer

## Critical Rules

- **No orphan sub-features** - Every sub-feature must connect to parent TDD
- **Clear dependencies** - If B needs A, document it
- **Minimize coupling** - Sub-features should be as independent as possible
- **Preserve architecture** - Don't break module boundaries established in TDD

## Output

Save to `./docs/features/<NNN>-<feature-name>.breakdown.md` using template in `references/TEMPLATE.md`.

## Content Rules

**Include:**
- Sub-feature list with IDs (a, b, c...)
- Scope for each sub-feature
- Dependencies between sub-features
- Suggested implementation order
- Parent TDD reference for each sub-feature

**Exclude:**
- Technical design (that's /engineer)
- Test cases (that's /test-design)
- Implementation code

## Next Step

After breakdown is approved:

1. **Clear context** - Start fresh for each sub-feature
2. **Run `/engineer`** - For the first sub-feature in sequence
3. **Repeat** - `/engineer` → `/test-design` → `/coder` for each sub-feature

```
Sub-feature 001a → /engineer → /test-design → /coder → /code-review → /qa
Sub-feature 001b → /engineer → /test-design → /coder → /code-review → /qa
...
```

Context can be cleared between sub-features since each reads from files.
