## Deployment

1. use the `./deploy create-config` to initialize the configuration file.
2. Read the configuration file carefully and tweak it as you wish
3. use the `./deploy --sync` to deploy the contents of the repo to your ~/.claude folder. (will ask you before doing anything)
4. by default it will also create a backup which you can rollback to

## Skills Overview

### Workflow

```
/feature → /tdd → /test-design → /coder → /code-review → /qa
    ↓         ↓          ↓
  .feature.md  .tdd.md   .test.md
```

### File Naming

All feature documents saved to `./docs/features/` with 3-digit auto-increment:

```
./docs/features/
├── 001-user-authentication.feature.md
├── 001-user-authentication.tdd.md
├── 001-user-authentication.test.md
├── 002-notifications.feature.md
├── 002-notifications.tdd.md
└── ...
```

### Workflow Skills

| Command | Purpose | Output |
|---------|---------|--------|
| `/feature` | Discover & define feature via interview | `<NNN>-name.feature.md` |
| `/tdd` | Technical design via interview | `<NNN>-name.tdd.md` |
| `/test-design` | Test cases (unit/integration/E2E) | `<NNN>-name.test.md` |
| `/coder` | Implement using TDD approach | Code |
| `/code-review` | Review changes via git diff | Pass/Fail |
| `/qa` | Validate against feature spec | Pass/Fail |

### Tech-Specific Skills (Auto-Trigger)

| Skill | Purpose |
|-------|---------|
| `nestjs-core` | NestJS patterns: TypeORM Data Mapper, CQRS, auth, module structure |
| `nestjs-testing` | NestJS test pyramid with solitary unit testing |
| `apollo-graphql` | Apollo Client: queries, mutations, caching, subscriptions |
| `react-router` | React Router v7: routing, loaders, actions, forms |
| `flutter-riverpod` | Riverpod 3.x + flutter_hooks patterns |

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
