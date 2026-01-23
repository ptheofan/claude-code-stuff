---
name: qa
version: 1.0.0
description: Quality assurance validation that implementation matches feature requirements and scope. This skill should be used when the user asks to "verify feature", "QA the implementation", "check requirements coverage", "validate acceptance criteria", "ensure feature completeness", or needs to confirm that all user stories and acceptance criteria from the feature spec are satisfied. Final step before release, use after /code-review.
---

# QA - Quality Assurance

Verify the implementation matches the feature requirements and scope is fulfilled.

## Process

1. **Read the feature spec** - `./docs/features/<NNN>-<feature-name>.feature.md`
2. **Review implementation** - Compare against requirements
3. **Validate acceptance criteria** - Each criterion must pass
4. **Check scope boundaries** - Nothing missing, nothing extra
5. **Report findings**

## Validation Checklist

### Requirements Coverage
- [ ] All user stories implemented
- [ ] All acceptance criteria met
- [ ] Edge cases handled (as defined in feature spec)

### Scope Validation
- [ ] All "In Scope" items delivered
- [ ] No "Out of Scope" items accidentally included
- [ ] Constraints respected

### User Story Verification

For each user story in the feature spec:

| Story ID | Description | Status | Notes |
|----------|-------------|--------|-------|
| US-001 | | ✅ / ❌ / ⚠️ | |

### Acceptance Criteria Verification

For each acceptance criterion:

| Story | Criterion | Status | Notes |
|-------|-----------|--------|-------|
| US-001 | Given X, When Y, Then Z | ✅ / ❌ | |

## Interview Tool

Use the interview tool to clarify edge cases:

```
The feature spec mentions "handle invalid input gracefully". What should happen?
1. Show validation error message
2. Silently ignore and use default
3. Redirect to help page
4. Other (specify)
```

## QA Report Format

```markdown
## QA Report: [Feature Name]

**Feature Spec:** `./docs/features/<NNN>-<feature-name>.feature.md`  
**Date:** [YYYY-MM-DD]  
**Status:** Pass / Fail / Partial

### Requirements Coverage

| User Story | Status | Notes |
|------------|--------|-------|
| US-001 | ✅ | |
| US-002 | ❌ | Missing X |

### Acceptance Criteria

| ID | Criterion | Status |
|----|-----------|--------|
| AC-001 | | ✅ |

### Scope Verification
- **In Scope:** All delivered ✅
- **Out of Scope:** None included ✅

### Issues Found
1. [Issue description]

### Verdict
- [ ] ✅ Approved - Ready for release
- [ ] ❌ Failed - Issues must be fixed
- [ ] ⚠️ Partial - Minor issues, can release with follow-up
```

## Next Step

After QA passes → Ready for merge/release 🎉
