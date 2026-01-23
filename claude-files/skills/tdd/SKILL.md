---
name: tdd
description: Technical Design Document creation through analysis and user interview. Use after /feature to create implementation blueprint with types, APIs, and architecture decisions.
---

# Technical Design Document

Analyze a feature spec and create a comprehensive technical design through user interview.

## Process

1. **Read the feature spec** - `./docs/features/<NNN>-<feature-name>.feature.md`
2. **Analyze requirements** - Break down into technical components
3. **Interview the user** - Use the interview tool for anything uncertain
4. **Design the solution** - Architecture, APIs, data models, error handling
5. **Save file** - `./docs/features/<NNN>-<feature-name>.tdd.md`

## Interview Tool

Use Claude Code's built-in interview tool to ask questions. Provide suggested answers with the last option being "Other (specify)".

Example:
```
How should we handle authentication for this endpoint?
1. JWT Bearer token (existing auth)
2. API Key
3. Public (no auth)
4. Other (specify)
```

## File Naming

Match the feature file number:
```
./docs/features/001-user-authentication.feature.md
./docs/features/001-user-authentication.tdd.md    ← Same number
```

## Interview Questions

Ask the user about:

### Architecture
- Which modules are affected?
- Any new modules needed?
- How does this integrate with existing code?

### Data
- What data needs to be stored?
- Any migrations required?
- Caching strategy?

### APIs
- What endpoints/methods are needed?
- Input/output shapes?
- Authentication requirements?

### Edge Cases
- What happens when [X] fails?
- How do we handle [Y] scenario?

### Constraints
- Performance requirements?
- Backwards compatibility needs?
- Third-party limitations?

## Critical Rules

- **Minimize assumptions** - If < 90% certain, ask
- **Interview first** - Don't design in the dark
- **Leave no stone unturned** - Thorough analysis
- **It's OK to have questions** - Document them in Open Questions
- **No implementation code** - Types, interfaces, contracts only

## Design Principles

Apply these from CLAUDE.md:

- **Clean Architecture** - Import down, emit events up
- **Strict module boundaries** - Only use exposed APIs
- **Domain Errors** - Specific exception classes with traceable data
- **Centralized integrations** - External calls through dedicated services
- **SOLID, KISS, DRY** - Simple, readable, no duplication

## Output

Save to `./docs/features/<NNN>-<feature-name>.tdd.md` using template in `references/TEMPLATE.md`.

## Content Rules

**Include:**
- Types, interfaces, contracts
- API signatures (inputs/outputs)
- Data models / schemas
- Module boundaries and exposed APIs
- Error types and their data
- Architecture diagrams (mermaid)

**Exclude:**
- Implementation code
- Business logic snippets
- "How the function works" details

The TDD is the blueprint, not the construction.

## Next Step

After TDD is approved → `/test-design`
