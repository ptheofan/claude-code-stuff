---
name: feature
description: Feature discovery and specification through user interview. Use when exploring a new feature idea that needs scope definition, user stories, and acceptance criteria.
---

# Feature Discovery

Explore and define features through user interview. Output a complete feature spec.

## Process

1. **Scan `./docs/features/`** - Find highest existing number, increment for new file (3 digits: 001, 002, etc.)
2. **Interview the user** - Use the interview tool to ask questions with suggested answers
3. **Define scope** - What's in, what's explicitly out
4. **Write user stories** - With acceptance criteria
5. **Save file** - `./docs/features/<NNN>-<feature-name>.feature.md`

## Interview Tool

Use Claude Code's built-in interview tool to ask questions. Provide suggested answers with the last option being "Other (specify)".

Example:
```
Who is the primary user?
1. End user / Customer
2. Admin / Back-office
3. Developer / API consumer
4. Other (specify)
```

## File Naming

```bash
# Auto-detect next number (3 digits)
./docs/features/001-user-authentication.feature.md
./docs/features/002-notifications.feature.md
./docs/features/003-payment-integration.feature.md
```

Scan existing files, find highest number, increment by 1.

## Interview Questions

Ask the user about:

### Problem
- What problem are we solving?
- What happens if we don't solve it?
- Who experiences this problem?

### Users
- Who is the primary user?
- What's their goal?
- Are there secondary users?

### Scope
- What MUST be in the MVP?
- What is explicitly OUT of scope?
- Any hard constraints?

### Success
- How do we know this feature succeeded?
- What does "done" look like?

### Edge Cases
- What happens when [X]?
- What if the user does [Y]?

## Critical Rules

- **Minimize assumptions** - If uncertain, ask
- **Interview first** - Don't draft until you understand
- **Explicit scope** - Both IN and OUT must be defined
- **No implementation details** - That's for `/tdd`

## Output

Save to `./docs/features/<NNN>-<feature-name>.feature.md` using template in `references/TEMPLATE.md`.

## Next Step

After feature is approved → `/tdd`
