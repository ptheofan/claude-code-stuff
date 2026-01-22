---
name: tdd
description: Create technical design documents (TDD) or system design documents. Use when user asks to design a feature, write a TDD, create a technical spec, or plan system architecture.
---

# Technical Design Document Skill

Generate comprehensive technical design documents as markdown files.

## Process

1. **Clarify scope** - Ask questions if requirements are unclear (< 90% certain)
2. **Use template** - Follow structure in `references/TEMPLATE.md`
3. **Apply project principles** - Reference CLAUDE.md for architecture standards
4. **Output as markdown** - Save to appropriate location (e.g., `docs/design/` or as specified)

## Key Principles to Embed

- Clean Architecture (import down, emit up)
- Strict module boundaries
- Domain Errors with traceable data
- Centralized external integrations
- Test strategy (unit → integration → E2E)

## Content Rules

**Include:**
- Types, interfaces, contracts
- API signatures (inputs/outputs)
- Data models / schemas
- Module boundaries and exposed APIs
- Error types and their data

**Exclude:**
- Implementation code
- Business logic snippets
- "How the function works" details

The TDD is the blueprint, not the construction.

## Quality Checks

Before finalizing, verify:
- [ ] Problem is clearly stated
- [ ] Alternatives have pros/cons
- [ ] Aligns with Clean Architecture
- [ ] Module boundaries respected
- [ ] Error handling uses Domain Errors
- [ ] Test strategy defined
- [ ] Open questions listed (if any)

## Template

See `references/TEMPLATE.md` for the full structure.
