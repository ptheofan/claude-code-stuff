---
name: feature
version: 2.0.0
description: Feature discovery and specification through user interview. This skill should be used when the user asks to "define a feature", "create feature spec", "write user stories", "explore requirements", "scope a feature", "start new feature", "define acceptance criteria", or mentions needing to understand what to build. First step in the workflow chain.
---

# Feature Discovery

Explore and define features through user interview until mature enough for development.

## Workflow Position

```
/feature → /tdd → /breakdown → /engineer → /test-design → /coder → /code-review → /qa
   ↑
  HERE
```

## Definition of "Mature Enough"

A feature is mature enough for development when:
- ✅ All questions answered - no ambiguity remains
- ✅ All paths known - happy paths, error paths, edge cases documented
- ✅ Scope well defined - explicit IN and OUT boundaries
- ✅ Success criteria clear - measurable outcomes defined
- ✅ Users identified - who benefits and how

## Process

1. **Scan `./docs/features/`** - Find highest existing number, increment for new file (3 digits: 001, 002, etc.)
2. **Interview the user** - Deep exploration using interview tool
3. **Iterate until mature** - Keep asking until all unknowns are resolved
4. **Define scope** - What's explicitly IN, what's explicitly OUT
5. **Write user stories** - With acceptance criteria
6. **Save file** - `./docs/features/<NNN>-<feature-name>.feature.md`

## Interview Tool

Use Claude Code's built-in interview tool (AskUserQuestion) to ask questions. Provide suggested answers with the last option being "Other (specify)".

Example:
```
Who is the primary user?
1. End user / Customer
2. Admin / Back-office
3. Developer / API consumer
4. Other (specify)
```

## Interview Strategy

### Phase 1: Problem Understanding
- What problem are we solving?
- What happens if we don't solve it?
- Who experiences this problem most acutely?
- Is there an existing workaround?

### Phase 2: User Understanding
- Who is the primary user?
- What's their goal when using this feature?
- Are there secondary users?
- What's the user's technical level?

### Phase 3: Scope Definition
- What MUST be in the MVP?
- What is explicitly OUT of scope?
- Any hard constraints (time, tech, regulations)?
- Dependencies on other features/systems?

### Phase 4: Path Mapping
- What's the happy path?
- What errors can occur?
- What edge cases exist?
- What happens at boundaries?

### Phase 5: Success Definition
- How do we know this feature succeeded?
- What metrics matter?
- What does "done" look like?
- How will users know it's working?

## Critical Rules

- **Minimize assumptions** - If < 90% certain, ask
- **Interview until mature** - Don't stop early
- **Explicit scope** - Both IN and OUT must be defined
- **No implementation details** - That's for `/tdd`
- **Document uncertainties** - Open questions section for unresolved items

## File Naming

```bash
# Auto-detect next number (3 digits)
./docs/features/001-user-authentication.feature.md
./docs/features/002-notifications.feature.md
./docs/features/003-payment-integration.feature.md
```

Scan existing files, find highest number, increment by 1.

## Output

Save to `./docs/features/<NNN>-<feature-name>.feature.md` using template in `assets/TEMPLATE.md`.

## Maturity Checklist

Before completing, verify:
- [ ] All user questions answered
- [ ] All paths documented (happy, error, edge)
- [ ] Scope explicitly defined (IN and OUT)
- [ ] Success criteria measurable
- [ ] No remaining ambiguity
- [ ] User stories have acceptance criteria

## Next Step

After feature is mature and approved, prompt the user using AskUserQuestion:

```
Feature complete. How would you like to proceed?
1. Clear memory and continue with /tdd (Recommended)
2. Continue with /tdd
3. Other
```

- If **Option 1**: Inform user to clear context, then invoke /tdd
- If **Option 2**: Proceed directly to /tdd
- If **Option 3**: Follow user's instructions
