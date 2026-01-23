# CLAUDE.md

## Core Principles

- **Clean Architecture**: Import downward, emit events upward. No circular dependencies.
- **Strict module boundaries**: NEVER access internals of another module. Only use exposed public APIs. Never import another module's entities directly.
- **Simple, reusable code**: Prefer clarity over cleverness.
- **Centralize external integrations**: All calls to Service X go through one class/function.
- **Minimal comments**: Only for complex/non-obvious logic.
- **No deprecated code**: Ever.
- **No TS overrides**: No `@ts-ignore`, `any` casts, or suppressed warnings (tests excepted when necessary).
- **Types/Interfaces**: Keep in separate files.

## Testing

- **TDD**: Design tests first, then implement.
- **Test Pyramid (Isolation)**: Unit and integration tests both follow this. Start with smallest functions, stub already-tested dependencies when testing higher layers.
- **Controller tests**: Focus on input validation and output shape, not business logic.
- **E2E tests**: Cover critical happy paths and selective business-critical failures. Skip what's already covered by unit/integration.
- **Test coverage**: Prioritize meaningful coverage over percentage targets. Cover critical paths and edge cases thoroughly.
- **Stub loggers**: Keep test output noise-free. Use logger calls as test evidence when needed.

## Error Handling

- **Throw early, catch late**
- **Domain Errors**: Create specific exception classes (e.g., `UserNotFoundError`, `InsufficientBalanceError`) - no generic throws.
- **Traceable exceptions**: Include relevant data (IDs, context) for debugging.
- **Let errors bubble**: Handle at appropriate boundaries (controllers, UI), not deep in business logic.

## Decision Making

- **< 90% certain? Ask.** Use interview tool.
- **New library/package?** Present pros/cons first, get approval.
- **Minimize assumptions**: When unsure, ask. We solve together.

## Performance & Reliability

- **N+1 awareness**: Always consider query patterns. Use eager loading or batching.
- **Pagination**: Default on all list endpoints.
- **Idempotency**: Design APIs and jobs to be safely retryable.

## Pre-Commit Checklist

- [ ] Tests pass
- [ ] No TS errors/warnings
- [ ] No console.logs or debug code
- [ ] Domain errors used (no generic throws)
- [ ] No circular dependencies
- [ ] Module boundaries respected
- [ ] Migrations included if schema changed

## Stack-Specific

### TypeScript / NestJS / NextJS
- Strict mode, no implicit any
- Use dependency injection
- Migrations for all schema changes
- **NestJS**: Use class-validator + class-transformer pipes for validation/transformation

### React / React Native
- Functional components + hooks
- Colocate related code

### Flutter
- **Riverpod** for state management

## Standards

- Structured logging
- Environment variables for config/secrets
- Conventional commits
