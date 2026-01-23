---
name: engineer
version: 2.0.0
description: Create a mini-TDD for a specific sub-feature. This skill should be used when the user asks to "design sub-feature", "create sub-feature TDD", "engineer sub-feature", or needs to create a focused technical design for one sub-feature slice. Uses the same principles and quality as the main /tdd skill. Use after /breakdown, before /test-design. Only used when breakdown was needed.
---

# Sub-Feature Technical Design (Mini-TDD)

Create a focused technical design document for a specific sub-feature, using the same principles and quality standards as the main /tdd skill.

## Workflow Position

```
/feature → /tdd → /breakdown → /engineer → /test-design → /coder → /code-review → /qa
                                   ↑
                                  HERE
                                   │
                            (per sub-feature)
```

**Note:** This skill is ONLY used when /breakdown determined that sub-features are needed. If no breakdown was needed, skip directly from /breakdown to /test-design.

## Process

1. **Read the parent TDD** - `./docs/features/<NNN>-<feature-name>.tdd.md`
2. **Read the breakdown** - `./docs/features/<NNN>-<feature-name>.breakdown.md`
3. **Identify this sub-feature's scope** - What slice are we designing?
4. **Design with full TDD rigor** - Apply same standards as main TDD
5. **Save file** - `./docs/features/<NNN><suffix>-<sub-name>.tdd.md`

## Same Quality as Main TDD

This is NOT a lightweight design. Apply the SAME principles and thoroughness as /tdd:

### What This Mini-TDD Must Cover

#### System Design (for this slice)
- Components affected by this sub-feature
- How this slice integrates with the parent architecture
- Module boundaries for this slice

#### ERD (for this slice)
- Entities this sub-feature touches
- New entities introduced (if any)
- Relationships relevant to this slice

#### Data Flows (for this slice)
- Input → Processing → Output for this sub-feature
- How data enters and exits this slice

#### State Management (if applicable)
- States this sub-feature manages
- Transitions triggered by this slice

#### API Contracts (for this slice)
- Endpoints/methods this sub-feature implements
- Request/response shapes
- Error responses specific to this slice

#### Risk Identification
- Risks specific to this sub-feature
- Dependencies on other sub-features
- Integration risks

## File Naming

Use parent number + letter suffix:
```
./docs/features/001-user-authentication.tdd.md      ← Parent TDD
./docs/features/001a-password-validation.tdd.md     ← Sub-feature A
./docs/features/001b-session-management.tdd.md      ← Sub-feature B
```

## What's Already Decided (Reference Parent TDD)

These are established - reference but don't repeat:
- Overall architecture decisions
- Database schema structure
- Authentication strategy
- External integrations approach
- Module boundaries

## What This Skill Defines (New for Sub-Feature)

Focus on execution details for THIS slice:
- **Types/Interfaces** specific to this sub-feature
- **Function contracts** - inputs, outputs, pre/post conditions
- **Domain errors** this sub-feature throws
- **Implementation approach** within parent's architecture
- **Test boundaries** - what's tested at this slice level

## Interview Tool

Use when implementation decisions need clarification:

```
For password validation in this sub-feature, which approach?
1. Inline validation in service (simple)
2. Separate PasswordValidator class (reusable)
3. Validation decorator pattern (AOP)
4. Other (specify)
```

## Design Principles

Apply these from CLAUDE.md (same as main TDD):

- **Clean Architecture** - Import down, emit events up
- **Strict module boundaries** - Only use exposed APIs
- **Domain Errors** - Specific exception classes with traceable data
- **SOLID, KISS, DRY** - Simple, readable, no duplication

## Critical Rules

- **Same rigor as /tdd** - This is a full design, just scoped smaller
- **Reference parent, don't repeat** - Point to parent TDD sections
- **Stay in scope** - Only design THIS sub-feature
- **No implementation code** - Types, contracts, diagrams only

## Output

Save to `./docs/features/<NNN><suffix>-<sub-name>.tdd.md` using template in `assets/TEMPLATE.md`.

## Mini-TDD Completeness Checklist

Before completing, verify:
- [ ] System design for this slice documented
- [ ] ERD additions/changes specified
- [ ] Data flows for this slice mapped
- [ ] API contracts for this slice defined
- [ ] Risks for this slice identified
- [ ] References to parent TDD included
- [ ] Clear scope boundaries defined

## Next Step

After sub-feature TDD is approved, prompt the user using AskUserQuestion:

```
Mini-TDD complete. How would you like to proceed?
1. Clear memory and continue with /test-design (Recommended)
2. Continue with /test-design
3. Other
```

- If **Option 1**: Inform user to clear context, then invoke /test-design
- If **Option 2**: Proceed directly to /test-design
- If **Option 3**: Follow user's instructions

Then repeat for next sub-feature: `/engineer` → `/test-design` → `/coder` → `/code-review`

After ALL sub-features complete → `/qa`
