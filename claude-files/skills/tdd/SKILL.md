---
name: tdd
version: 2.0.0
description: Technical Design Document creation for the ENTIRE feature. This skill should be used when the user asks to "create TDD", "write technical design", "design architecture", "plan implementation", "design system", "create ERD", "define data flows", "identify risks", or needs to transform feature requirements into a comprehensive technical blueprint. Use after /feature, before /breakdown.
---

# Technical Design Document

Create a comprehensive technical design document for the ENTIRE feature.

## Workflow Position

```
/feature → /tdd → /breakdown → /engineer → /test-design → /coder → /code-review → /qa
             ↑
            HERE
```

## Process

1. **Read the feature spec** - `./docs/features/<NNN>-<feature-name>.feature.md`
2. **Analyze requirements** - Break down into technical components
3. **Interview the user** - Use interview tool for architectural decisions
4. **Design the complete solution** - All aspects listed below
5. **Save file** - `./docs/features/<NNN>-<feature-name>.tdd.md`

## What This TDD Must Cover

### System Design
- High-level architecture overview
- Component interactions
- Module boundaries and responsibilities
- Integration points with existing system

### Entity Relationship Diagram (ERD)
- Database entities and relationships
- Cardinality (1:1, 1:N, N:M)
- Key attributes per entity
- Indexes needed

### Data Flows
- How data moves through the system
- Input → Processing → Output paths
- Transformation points
- Validation checkpoints

### State Management
- State definitions (if applicable)
- State transitions
- State persistence strategy
- Event triggers

### API Contracts
- Endpoints/methods needed
- Request/response shapes
- Authentication requirements
- Error responses

### Risk Identification
- Technical risks and mitigations
- Dependencies on external systems
- Performance concerns
- Security considerations

### Deployment Strategy
- How will this be deployed?
- Feature flags needed?
- Migration strategy
- Rollback plan

## Interview Tool

Use Claude Code's interview tool (AskUserQuestion) for architectural decisions.

Example:
```
How should we handle authentication for this feature?
1. JWT Bearer token (existing auth system)
2. API Key (for service-to-service)
3. Session-based (for web UI)
4. Other (specify)
```

## File Naming

Match the feature file number:
```
./docs/features/001-user-authentication.feature.md
./docs/features/001-user-authentication.tdd.md    ← Same number
```

## Design Principles

Apply these from CLAUDE.md:

- **Clean Architecture** - Import down, emit events up
- **Strict module boundaries** - Only use exposed APIs
- **Domain Errors** - Specific exception classes with traceable data
- **Centralized integrations** - External calls through dedicated services
- **SOLID, KISS, DRY** - Simple, readable, no duplication

## Critical Rules

- **Minimize assumptions** - If < 90% certain, ask
- **Comprehensive coverage** - All sections must be addressed
- **No implementation code** - Types, interfaces, contracts only
- **Document risks** - Don't hide uncertainties

## Output

Save to `./docs/features/<NNN>-<feature-name>.tdd.md` using template in `assets/TEMPLATE.md`.

## Content Rules

**Include:**
- System architecture diagram (mermaid)
- ERD diagram (mermaid)
- Data flow diagrams
- State diagrams (if applicable)
- API contracts with types
- Risk assessment table
- Deployment considerations

**Exclude:**
- Implementation code
- Business logic snippets
- "How the function works" details

The TDD is the complete blueprint for the ENTIRE feature.

## TDD Completeness Checklist

Before completing, verify:
- [ ] System design documented
- [ ] ERD complete with relationships
- [ ] Data flows mapped
- [ ] States defined (if applicable)
- [ ] API contracts specified
- [ ] Risks identified with mitigations
- [ ] Deployment strategy defined

## Next Step

After TDD is approved, prompt the user using AskUserQuestion:

```
TDD complete. How would you like to proceed?
1. Clear memory and continue with /breakdown (Recommended)
2. Continue with /breakdown
3. Other
```

- If **Option 1**: Inform user to clear context, then invoke /breakdown
- If **Option 2**: Proceed directly to /breakdown
- If **Option 3**: Follow user's instructions
