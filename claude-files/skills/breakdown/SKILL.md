---
name: breakdown
version: 2.0.0
description: Analyze feature and TDD to determine IF breakdown is needed. This skill should be used when the user asks to "break down feature", "split into sub-features", "check if breakdown needed", "analyze feature size", or after completing /tdd to determine next steps. Decides whether to proceed with sub-features (/engineer) or skip directly to /test-design. Use after /tdd.
---

# Feature Breakdown

Analyze feature and TDD to determine IF breakdown into sub-features is needed.

## Workflow Position

```
/feature → /tdd → /breakdown → /engineer → /test-design → /coder → /code-review → /qa
                       ↑
                      HERE
                       │
                       ├─→ IF breakdown needed: → /engineer (per sub)
                       │
                       └─→ IF small enough: → /test-design (skip /engineer)
```

## Process

1. **Read the feature spec** - `./docs/features/<NNN>-<feature-name>.feature.md`
2. **Read the TDD** - `./docs/features/<NNN>-<feature-name>.tdd.md`
3. **Analyze complexity** - Apply decision criteria below
4. **Decide: breakdown needed?**
   - **YES** → Create breakdown file, proceed to `/engineer`
   - **NO** → Skip to `/test-design` directly
5. **If breakdown needed** - Save file `./docs/features/<NNN>-<feature-name>.breakdown.md`

## Decision Criteria: Is Breakdown Needed?

### Breakdown IS Needed When:
- Feature touches 3+ distinct modules
- Implementation would take > 3 days
- Multiple independent user stories exist
- Clear natural boundaries are visible
- Risk of merge conflicts if done as single unit
- Different team members could work on parts in parallel

### Breakdown NOT Needed When:
- Feature is contained within 1-2 modules
- Implementation would take ≤ 3 days
- Code changes are tightly coupled (can't be separated)
- Breaking down would create artificial boundaries
- All changes must be deployed together anyway

### Key Question

> "Can parts of this feature be implemented, tested, and potentially deployed independently without breaking the codebase?"

If YES → Breakdown needed
If NO → Skip to /test-design

## Breakdown Principles

When breakdown IS needed:

### Feature Perspective
- Each sub-feature delivers identifiable user value
- Sub-features can be understood in isolation
- Clear acceptance criteria per sub-feature

### Code Perspective
- Parts that should be written in a single ticket stay together
- Don't break natural code boundaries
- Each sub-feature can have its own PR
- No circular dependencies between sub-features

### Autonomy Criteria

Each sub-feature MUST be:
- **Independently testable** - Tests can run for just this slice
- **Independently deployable** - Can ship without others (maybe feature-flagged)
- **Coherent unit** - Makes sense as a standalone piece of work
- **Right-sized** - 1-3 days of work for one developer

## Interview Tool

Use to clarify breakdown decisions:

```
This feature has 3 distinct parts. How should we approach?
1. Break into 3 sub-features (parallel work possible)
2. Keep as single feature (tightly coupled)
3. Break into 2 sub-features (specify grouping)
4. Other (specify)
```

## File Naming

If breakdown needed:
```
./docs/features/001-user-authentication.feature.md   ← Feature spec
./docs/features/001-user-authentication.tdd.md       ← Technical design
./docs/features/001-user-authentication.breakdown.md ← This skill's output
```

Sub-features get letter suffixes (created during /engineer):
```
./docs/features/001a-password-validation.tdd.md
./docs/features/001b-session-management.tdd.md
```

## Output Format (When Breakdown Needed)

Save to `./docs/features/<NNN>-<feature-name>.breakdown.md` using template in `assets/TEMPLATE.md`.

**Include:**
- Decision rationale (why breakdown was needed)
- Sub-feature list with IDs (a, b, c...)
- Scope for each sub-feature
- Dependencies between sub-features
- Suggested implementation order

## Critical Rules

- **Read BOTH files** - Feature spec AND TDD before deciding
- **Don't force breakdown** - If it's small enough, skip to /test-design
- **Don't avoid breakdown** - If it's too big, split it properly
- **Code coherence** - Don't break apart code that belongs together
- **Document the decision** - Whether YES or NO, explain why

## Next Steps

### If Breakdown Needed:

After saving breakdown file, prompt the user using AskUserQuestion:

```
Breakdown complete. How would you like to proceed?
1. Clear memory and continue with /engineer for first sub-feature (Recommended)
2. Continue with /engineer
3. Other
```

- If **Option 1**: Inform user to clear context, then invoke /engineer
- If **Option 2**: Proceed directly to /engineer
- If **Option 3**: Follow user's instructions

Then repeat: `/engineer` → `/test-design` → `/coder` → `/code-review` for each sub-feature.
After ALL sub-features complete → `/qa`

### If No Breakdown Needed:

Prompt the user using AskUserQuestion:

```
No breakdown needed. How would you like to proceed?
1. Clear memory and continue with /test-design (Recommended)
2. Continue with /test-design
3. Other
```

- If **Option 1**: Inform user to clear context, then invoke /test-design
- If **Option 2**: Proceed directly to /test-design
- If **Option 3**: Follow user's instructions

Continue: `/test-design` → `/coder` → `/code-review` → `/qa`
