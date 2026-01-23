## Deployment

1. use the `./deploy create-config` to initialize the configuration file.
2. Read the configuration file carefully and tweak it as you wish
3. use the `./deploy --sync` to deploy the contents of the repo to your ~/.claude folder. (will ask you before doing anything)
4. by default it will also create a backup which you can rollback to

## Skills Overview

### Workflow

```
/feature → /tdd → /breakdown → /engineer (per sub) → /test-design → /coder → /code-review → /qa
    ↓         ↓         ↓            ↓
  .feature  .tdd    .breakdown    .tdd (sub)
```

For large features, `/breakdown` splits the TDD into sub-features:

```
001-user-auth.tdd.md
        ↓
   /breakdown
        ↓
001a-password-validation → /engineer → /test-design → /coder → ...
001b-session-management  → /engineer → /test-design → /coder → ...
001c-password-reset      → /engineer → /test-design → /coder → ...
```

### File Naming

All feature documents saved to `./docs/features/` with 3-digit auto-increment:

```
./docs/features/
├── 001-user-authentication.feature.md    # Feature spec
├── 001-user-authentication.tdd.md        # Technical design
├── 001-user-authentication.breakdown.md  # Sub-feature breakdown
├── 001a-password-validation.tdd.md       # Sub-feature TDD
├── 001a-password-validation.test.md      # Sub-feature tests
├── 001b-session-management.tdd.md
├── 001b-session-management.test.md
├── 002-notifications.feature.md
└── ...
```

### Workflow Skills

| Command | Purpose | Output |
|---------|---------|--------|
| `/feature` | Discover & define feature via interview | `<NNN>-name.feature.md` |
| `/tdd` | Technical design via interview | `<NNN>-name.tdd.md` |
| `/breakdown` | Split TDD into sub-features | `<NNN>-name.breakdown.md` |
| `/engineer` | Lightweight TDD for sub-features | `<NNN><x>-name.tdd.md` |
| `/test-design` | Test cases (unit/integration/E2E) | `<NNN>-name.test.md` |
| `/coder` | Implement using TDD approach | Code |
| `/code-review` | Review changes via git diff | Pass/Fail |
| `/qa` | Validate against feature spec | Pass/Fail |

### Context Clearing

| Transition | Clear Context? | Why |
|------------|----------------|-----|
| `/feature` → `/tdd` | ✅ Yes | Feature spec saved to file |
| `/tdd` → `/breakdown` | ❌ No | Breakdown needs TDD context fresh |
| `/breakdown` → `/engineer` | ✅ Yes | Starting fresh, reads from files |
| `/engineer` → `/test-design` | ✅ Yes | Test design reads from files |
| `/test-design` → `/coder` | ✅ Yes | Coder reads from files |
| `/coder` → `/code-review` | ❌ No | Reviewer needs to see what was coded |
| `/code-review` → `/qa` | ❌ No | QA benefits from review context |

**Rule of thumb:** Clear when transitioning to a skill that reads everything from files.

### Tech-Specific Skills (Auto-Trigger)

| Skill | Purpose |
|-------|---------|
| `nestjs-core` | NestJS patterns: TypeORM Data Mapper, CQRS, auth, module structure |
| `nestjs-testing` | NestJS test pyramid with solitary unit testing |
| `apollo-graphql` | Apollo Client: queries, mutations, caching, subscriptions |
| `react-router` | React Router v7: routing, loaders, actions, forms |
| `riverpod` | Riverpod 3.x + flutter_hooks patterns |

### Usage

#### Start a New Feature
```
/feature
```
Claude will interview you and create the feature spec.

#### Continue to Technical Design
```
/tdd
```
Claude reads the feature spec and designs the solution.

#### Break Down Large Features
```
/breakdown
```
Claude analyzes the TDD and proposes sub-features with dependencies.

#### Design Sub-Feature (Lightweight TDD)
```
/engineer
```
Claude creates focused technical design for a sub-feature, referencing parent TDD.

#### Plan Tests
```
/test-design
```
Claude creates test specifications from the TDD.

#### Implement
```
/coder
```
Claude implements using TDD (tests first, then code).

#### Review
```
/code-review
```
Claude reviews git diff against quality standards.

#### Validate
```
/qa
```
Claude validates implementation matches feature requirements.
