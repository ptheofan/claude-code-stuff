---
description: "Validate requirements completeness and quality against development standards"
argument-hint: [requirements-text]
---

Have qa-expert and business-analyst validate these requirements against our development policy:

Completeness criteria:
- All user stories have clear, testable acceptance criteria (80%+ test coverage standard)
- Success metrics are defined and measurable
- Dependencies and assumptions are documented
- Risk factors and mitigation strategies identified
- Build order considers dependency hierarchy

Quality criteria:
- Requirements are testable and verifiable
- Language is clear and unambiguous
- Business value is articulated for each feature
- Technical feasibility is considered
- KISS principle applied (simple, clear requirements)
- No deprecated approaches or technologies

Development alignment:
- Features can be broken into small, testable increments (1-2 hour iterations)
- Each requirement supports incremental development
- Clear error handling requirements specified
- TypeScript-friendly data structures defined

Requirements: $ARGUMENTS
