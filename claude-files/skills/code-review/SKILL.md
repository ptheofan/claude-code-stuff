---
name: code-review
version: 2.0.0
description: Review code produced by /coder against quality standards. This skill should be used when the user asks to "review code", "check my changes", "validate implementation", "review PR", or after /coder completes. Checks if code solves the problem, runs linter/tooling, and spawns a coder agent to fix any issues found. Use after /coder, before /qa.
---

# Code Review

Review code produced by /coder. Identify issues and spawn an agent to fix them if needed.

## Workflow Position

```
/feature → /tdd → /breakdown → /engineer → /test-design → /coder → /code-review → /qa
                                                                        ↑
                                                                       HERE
```

## Process

1. **Read the TDD** - Understand what was supposed to be built
2. **Read the test-design** - Understand expected behaviors
3. **Get the code diff** - `git diff` or `git diff --staged`
4. **Run linter and tooling** - Execute codebase tooling (lint, typecheck, test)
5. **Review against checklist** - Check each item below
6. **Report findings** - Document issues and suggestions
7. **If issues found** - Spawn coder agent to fix them
8. **If approved** - Proceed to next step prompt

## Always Run These Commands

```bash
# Run linter
npm run lint
# or: yarn lint, pnpm lint, etc.

# Run type check
npm run typecheck
# or: tsc --noEmit

# Run tests
npm run test

# Check for any other project-specific tooling
# Read package.json scripts section
```

## Review Checklist

### Does It Solve The Problem?
- [ ] Code implements what the TDD designed
- [ ] All test cases from test-design are covered
- [ ] Feature requirements are met
- [ ] Nothing is missing from the implementation

### Did The Coder Forget Anything?
- [ ] All endpoints/methods implemented
- [ ] Error handling complete
- [ ] Edge cases covered
- [ ] Migrations included (if schema changed)

### Architecture
- [ ] **Clean Architecture** - Dependencies flow downward only
- [ ] **Module boundaries** - No direct imports of another module's internals
- [ ] **Centralized integrations** - External calls through dedicated services
- [ ] **No circular dependencies**

### Code Quality
- [ ] **SOLID** - Single responsibility, proper abstractions
- [ ] **KISS** - Simple, readable solutions
- [ ] **DRY** - No code duplication
- [ ] **Meaningful names** - Clear variable/function names
- [ ] **No console.logs** - Debug code removed

### TypeScript
- [ ] **No `any`** - Proper types used
- [ ] **No `@ts-ignore`** - Issues fixed, not suppressed
- [ ] **No suppressed warnings**
- [ ] **Types/interfaces in separate files**

### Error Handling
- [ ] **Domain Errors** - Specific exception classes used
- [ ] **Traceable data** - Errors include IDs/context
- [ ] **Errors bubble** - Handled at appropriate boundaries

### Testing
- [ ] **Tests exist** - New code has corresponding tests
- [ ] **Tests pass** - All tests green
- [ ] **Meaningful coverage** - Critical paths covered

## Review Output Format

```markdown
## Code Review: [Feature/Sub-feature Name]

### Summary
[One-line summary of review findings]

### Tooling Results
- Linter: ✅ Pass / ❌ X issues
- TypeCheck: ✅ Pass / ❌ X errors
- Tests: ✅ Pass / ❌ X failing

### ✅ Good
- [What's done well]

### ❌ Issues (Must Fix)
- [ ] [Critical issues that block approval]

### ⚠️ Suggestions (Optional)
- [ ] [Improvements that would be nice]

### Verdict
- [ ] ✅ Approved - Proceed to /qa
- [ ] 🔄 Issues found - Spawning coder agent to fix
```

## If Issues Found: Spawn Coder Agent

When issues are identified that must be fixed:

1. Document the specific issues clearly
2. **Spawn a Task agent** with the coder skill to fix them:

```
Task: Fix code review issues
- Issue 1: [description]
- Issue 2: [description]
- ...

Fix these issues following the existing TDD and test-design.
```

3. After agent completes, re-review the fixes
4. Repeat until all issues resolved

## Critical Rules

- **Always run tooling** - Linter, typecheck, tests before manual review
- **Check completeness** - Did coder forget anything?
- **Spawn agent for fixes** - Don't just report issues, get them fixed
- **Re-review after fixes** - Verify issues are actually resolved

## Git Commands

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

After code review passes (all issues fixed), prompt the user using AskUserQuestion:

### If more sub-features remain:

```
Code review passed. More sub-features remain. How would you like to proceed?
1. Clear memory and continue with /engineer for next sub-feature (Recommended)
2. Continue with /engineer
3. Other
```

- If **Option 1**: Inform user to clear context, then invoke /engineer
- If **Option 2**: Proceed directly to /engineer
- If **Option 3**: Follow user's instructions

### If this was the last sub-feature (or no breakdown):

```
Code review passed. All sub-features complete. How would you like to proceed?
1. Clear memory and continue with /qa (Recommended)
2. Continue with /qa
3. Other
```

- If **Option 1**: Inform user to clear context, then invoke /qa
- If **Option 2**: Proceed directly to /qa
- If **Option 3**: Follow user's instructions
