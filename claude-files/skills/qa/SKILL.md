---
name: qa
version: 2.0.0
description: Quality assurance validation for the ENTIRE completed feature. This skill should be used when the user asks to "verify feature", "QA the implementation", "check requirements coverage", "validate feature completeness", or after ALL sub-features (if any) have passed code review. Checks if the entire implementation matches the original feature request and reports any deviations. Final step before release.
---

# QA - Quality Assurance

Verify the ENTIRE feature implementation matches the original feature request. Report any deviations.

## Workflow Position

```
/feature → /tdd → /breakdown → /engineer → /test-design → /coder → /code-review → /qa
                                                                                    ↑
                                                                                   HERE
                                                                                 (FINAL)
```

## When To Run QA

**Run /qa ONLY when the ENTIRE feature is complete:**

- If NO breakdown: After `/code-review` passes
- If breakdown exists: After ALL sub-features have passed `/code-review`

```
# No breakdown:
/feature → /tdd → /breakdown (skip) → /test-design → /coder → /code-review → /qa ✓

# With breakdown:
/feature → /tdd → /breakdown →
  Sub A: /engineer → /test-design → /coder → /code-review ✓
  Sub B: /engineer → /test-design → /coder → /code-review ✓
  Sub C: /engineer → /test-design → /coder → /code-review ✓
→ /qa ✓ (after ALL sub-features complete)
```

## Process

1. **Read the original feature spec** - `./docs/features/<NNN>-<feature-name>.feature.md`
2. **Review the complete implementation** - All code produced
3. **Compare against requirements** - Every user story, every acceptance criterion
4. **Identify deviations** - What's different from the original request?
5. **Report findings** - Document any gaps or discrepancies

## Validation Checklist

### Requirements Coverage
- [ ] All user stories from feature spec are implemented
- [ ] All acceptance criteria are met
- [ ] All edge cases (as defined in feature spec) are handled
- [ ] Success criteria from feature spec are achievable

### Scope Validation
- [ ] All "In Scope" items from feature spec are delivered
- [ ] No "Out of Scope" items were accidentally included
- [ ] Constraints from feature spec were respected

### Feature Completeness
- [ ] Feature works end-to-end as described
- [ ] All paths (happy, error, edge) function correctly
- [ ] Integration with existing system is complete

## Deviation Report Format

For each deviation found:

| Field | Description |
|-------|-------------|
| **ID** | DEV-001, DEV-002, etc. |
| **Type** | Missing / Different / Extra |
| **Feature Spec Reference** | Which user story or criterion |
| **Expected** | What the feature spec said |
| **Actual** | What was implemented |
| **Severity** | Critical / Major / Minor |
| **Recommendation** | Fix required / Acceptable / Document |

## QA Report Format

```markdown
## QA Report: [Feature Name]

**Feature Spec:** `./docs/features/<NNN>-<feature-name>.feature.md`
**Date:** [YYYY-MM-DD]
**Status:** ✅ Pass / ❌ Fail / ⚠️ Pass with deviations

### Summary
[One-line summary of QA findings]

### User Story Coverage

| Story ID | Description | Status | Notes |
|----------|-------------|--------|-------|
| US-001 | [Description] | ✅ / ❌ / ⚠️ | |
| US-002 | [Description] | ✅ / ❌ / ⚠️ | |

### Acceptance Criteria

| Story | Criterion | Status | Notes |
|-------|-----------|--------|-------|
| US-001 | Given X, When Y, Then Z | ✅ / ❌ | |

### Scope Verification
- **In Scope Items:** [All delivered ✅ / Missing items listed]
- **Out of Scope:** [None included ✅ / Unexpected items listed]

### Deviations Found

| ID | Type | Severity | Description |
|----|------|----------|-------------|
| DEV-001 | Missing | Major | [Description] |

### Verdict
- [ ] ✅ **Approved** - Feature complete, ready for release
- [ ] ❌ **Failed** - Critical deviations must be fixed
- [ ] ⚠️ **Conditional** - Minor deviations documented, can release
```

## Interview Tool

Use to clarify ambiguities in feature spec:

```
The feature spec says "handle errors gracefully". The implementation shows a generic error page. Is this acceptable?
1. Yes, generic error page is fine
2. No, need specific error messages
3. Need to discuss further
4. Other (specify)
```

## Critical Rules

- **Compare to ORIGINAL feature spec** - Not the TDD, not the test-design
- **Check EVERYTHING** - Every user story, every criterion
- **Report ALL deviations** - Even minor ones
- **Be objective** - Don't assume intent, report what's different
- **Run after COMPLETE feature** - All sub-features must be done first

## Deviation Categories

### Missing (Type: Missing)
Something in the feature spec that wasn't implemented.

### Different (Type: Different)
Something implemented differently than specified.

### Extra (Type: Extra)
Something implemented that wasn't in the feature spec.

## Severity Levels

- **Critical** - Feature doesn't work as specified, blocks release
- **Major** - Significant deviation from spec, should fix before release
- **Minor** - Small deviation, can document and release

## Next Step

### If QA Passes:

Prompt the user using AskUserQuestion:

```
QA passed! Feature complete. How would you like to proceed?
1. Ready for merge/release (Recommended)
2. Review documentation before release
3. Other
```

✅ Feature workflow complete!

### If QA Fails:

Prompt the user using AskUserQuestion:

```
QA found deviations. How would you like to proceed?
1. Clear memory and return to /coder to fix issues (Recommended)
2. Continue to fix issues (keep context)
3. Accept deviations and proceed to release
4. Other
```

- If **Option 1**: Inform user to clear context, then invoke /coder
- If **Option 2**: Proceed directly to /coder
- If **Option 3**: Document deviations and proceed to release
- If **Option 4**: Follow user's instructions

After fixes → Re-run /qa

### After Release:

Feature workflow complete. Archive documentation:
```
./docs/features/001-user-authentication.feature.md  ← Keep
./docs/features/001-user-authentication.tdd.md      ← Keep
./docs/features/001-user-authentication.test.md     ← Keep
./docs/features/001-user-authentication.breakdown.md ← Keep (if exists)
./docs/features/001a-*.tdd.md                       ← Keep (if exists)
./docs/features/001a-*.test.md                      ← Keep (if exists)
```
