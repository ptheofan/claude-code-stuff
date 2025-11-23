CRITICAL DEVELOPMENT POLICY - NEVER IGNORE OR FORGET:
  
  This policy is MANDATORY for ALL projects and sessions

  # TYPE SAFETY (TypeScript)
  - MINIMIZE 'any' types - only acceptable in extreme cases like catch (error: unknown)
  - NO shortcuts: no @ts-ignore, no casting to unknown unless absolutely necessary
  - EXCEPTION: In tests only, 'as unknown as SomeType' is acceptable for mocking
  - NO disabling TypeScript warnings or eslint rules

  # MANDATORY FEATURE DEVELOPMENT PROCESS
  1. PLAN: Break feature into small testable parts
  2. PRIORITIZE: Always propose next steps that minimize dependencies on unimplemented features
     - Build foundational components first (utilities, data models, core services)
     - Implement dependency-free features before dependent ones
     - Create mock/stub interfaces for missing dependencies when absolutely necessary
  3. ITERATE: Advance in small parts, test each iteration, ensure working condition
  4. DEFINITION OF DONE (ALL must be met):
     ✅ Meets ALL requirements
     ✅ ALL tests pass (no deleting tests to fix builds)
     ✅ 80%+ coverage with SEPARATE unit and integration test suites
     ✅ Unit tests: isolated, fast, mock external dependencies
     ✅ Integration tests: test component interactions, use real dependencies
     ✅ Entire codebase in fully working condition
     ✅ No TypeScript errors/warnings
     ✅ No deprecated functions/APIs used
     ✅ Follows KISS principle (simple, readable solutions)
     ✅ Minimal, justified mocking (only external dependencies in unit tests)
     ✅ Minimal TODOs (only for legitimate missing dependencies with tickets)
     ✅ Proper error handling

  # FORBIDDEN SHORTCUTS
  ❌ Disabling warnings/errors
  ❌ @ts-ignore without justification  
  ❌ Casting to any/unknown for convenience (except 'as unknown as Type' in tests for mocking)
  ❌ Skipping tests
  ❌ Deleting tests to avoid fixing failures (removing redundant tests OK with justification)
  ❌ Over-mocking (mock only external dependencies in unit tests, NOT in integration tests)
  ❌ Mixing unit and integration test patterns
  ❌ Committing broken code
  ❌ Reducing coverage below 80%
  ❌ Using deprecated functions/APIs

  **REMINDER: Forgetting this policy is unacceptable. Quality over speed ALWAYS.**

  ## Unit Tests, Integration Tests and E2E tests strategy
  We use the pyramid layered approach following the **Solitary Unit Testing** pattern.
  ### Unit Tests
   - test the smallest functions
   - test services (mock smaller functions used that are already unit tested)
   - test commands (mock services and smaller functionat that are already unit tested)
   - test queries (mock what is already unit tested)
   - test gql resolvers (intput/output) as internal calls should have been already unit tested with any of the prev mentioned
   - test gql queries (intput/output) as internal calls should have been already unit tested with any of the prev mentioned
   - test REST API controller/actions input/output as internal calls should have been already unit tested with any of the prev mentioned

  ### Integration Tests
   - focus on database ops and other internal integrations (feel free to mock things that are already unit tested)

  ### E2E Tests
  As E2E tests are expensive we are more selective about them.
  - focus on happy critical business paths
  - include selective business critical failure paths
  - exclude paths are covered by unit tests and integration tests
  
  ## Before starting ANY work:
  1. Confirm understanding of current state
  2. Reference this policy
  3. Analyze feature dependencies and propose ground-up build order
  4. Select appropriate agent(s) for the task:
     - typescript-pro for TS development
     - nestjs-core-expert for NestJS fundamentals
     - nestjs-cqrs-expert for CQRS patterns, sagas, and CQRS event-driven workflows (nestjs CQRS)
     - nestjs-eventemitter2-expert for nestjs eventemitter2 architecture (nestjs events bus)
     - nestjs-passport-expert for authentication
     - nestjs-unit-test-expert for unit testing
     - nestjs-database-expert for database integration (TypeORM, Prisma, Sequelize, Mongoose)
     - nestjs-typeorm-expert for TypeORM Data Mapper pattern, mappers, tree entities, and @Transactional()
     - objection-orm-expert for Objection.js ORM with TypeScript, Knex.js, relations, and graph operations
     - nestjs-configuration-expert for configuration management
     - nestjs-validation-expert for input validation and DTOs
     - nestjs-caching-expert for caching strategies
     - nestjs-serialization-expert for response serialization
     - nestjs-task-scheduling-expert for cron jobs and scheduled tasks
     - nestjs-queues-expert for background job processing
     - nestjs-logger-expert for logging
     - nestjs-security-expert for security (Helmet, CORS, CSRF, rate limiting)
     - nestjs-file-upload-expert for file uploads
     - nestjs-streaming-expert for file streaming
     - nestjs-http-module-expert for external API calls
     - nestjs-compression-expert for response compression
     - nestjs-cookies-expert for cookie management
     - nestjs-session-expert for session management
     - nestjs-mvc-expert for server-side rendering
     - nestjs-versioning-expert for API versioning
     - nestjs-sse-expert for Server-Sent Events
     - nx-monorepo-expert for NX monorepo management with NestJS
     - qbittorrent-api-expert for qBittorrent Web API integration
     - apollo-client-setup-expert for Apollo Client configuration and setup
     - apollo-queries-expert for Apollo Client queries (useQuery, useLazyQuery)
     - apollo-mutations-expert for Apollo Client mutations and cache updates
     - apollo-caching-expert for Apollo Client caching strategies
     - apollo-subscriptions-expert for Apollo Client subscriptions and real-time data
     - apollo-testing-expert for testing Apollo Client applications
     - code-reviewer for quality checks
     - test-automator for testing
     - security-engineer for security
     - performance-engineer for optimization
     - multi-agent-coordinator for complex workflows
  5. Plan incremental approach prioritizing dependency-free components
  6. Use TODO comments only for legitimate missing feature dependencies with ticket references

  # SESSION BEHAVIOR REQUIREMENTS
  - ALWAYS check this policy before starting any work
  - REQUIRE explicit confirmation of approach before major changes  
  - ENFORCE 80%+ test coverage standard - no exceptions
  - AGENT SELECTION IS MANDATORY - always specify which agent(s) to use
  - PREFER specialized agents over generic responses
  - ALWAYS validate against quality gates before marking work complete

  # AGENT SELECTION GUIDELINES ===
  When selecting agents for tasks, use these domain-specific guidelines:

  **Code Quality & Architecture:**
  - code-reviewer: Code review and quality assessment
  - architect-reviewer: Architecture and design pattern evaluation  
  - refactoring-specialist: Code refactoring and cleanup

  **Testing & Quality Assurance:**
  - qa-expert: Test strategy and quality validation
  - test-automator: Test implementation and automation
  - accessibility-tester: Accessibility compliance testing
  - performance-engineer: Performance optimization and testing

  **TypeScript & Frontend Development:**
  - typescript-pro: TypeScript development and type safety
  - react-specialist: React component development
  - nextjs-developer: Next.js application development

  **NestJS & Backend Architecture:**
  - nestjs-core-expert: Core NestJS architecture, modules, providers, and dependency injection
  - nestjs-cqrs-expert: CQRS pattern, commands, queries, events, sagas, event sourcing, and workflow orchestration
  - nestjs-eventemitter2-expert: Event-driven architecture and domain events
  - nestjs-passport-expert: Authentication and authorization with Passport.js strategies (JWT, OAuth, sessions, GraphQL auth guards)
  - nestjs-unit-test-expert: Unit testing for NestJS components with Jest and TestingModule
  - nestjs-database-expert: Database integration with TypeORM, Prisma, Sequelize, and Mongoose
  - nestjs-typeorm-expert: TypeORM Data Mapper pattern, custom repositories, mappers, @Transactional() decorator, tree entities, and testing
  - objection-orm-expert: Objection.js ORM with TypeScript focus covering model definition with jsonSchema validation, relation types (HasMany, BelongsToOne, ManyToMany, HasOne, HasOneThrough), query building with Knex.js integration, eager loading (withGraphFetched vs withGraphJoined), graph operations (insertGraph, upsertGraph with #ref/#id), transaction management, lifecycle hooks ($beforeInsert, $afterUpdate, etc.), snake_case to camelCase conversion, raw queries, JSON/JSONB queries, custom query builders, modifiers, and production patterns. Handles complex hierarchical data operations and provides type-safe database access patterns.
  - nestjs-configuration-expert: Configuration management and environment variables
  - nestjs-validation-expert: Input validation, DTOs, and data transformation
  - nestjs-caching-expert: Caching strategies with Redis and in-memory cache
  - nestjs-serialization-expert: Response serialization and data transformation
  - nestjs-task-scheduling-expert: Cron jobs, intervals, and scheduled tasks
  - nestjs-queues-expert: Background job processing with Bull/BullMQ
  - nestjs-logger-expert: Logging with Winston, Pino, and built-in logger
  - nestjs-security-expert: Security best practices (Helmet, CORS, CSRF, rate limiting)
  - nestjs-file-upload-expert: File upload handling with Multer and cloud storage
  - nestjs-streaming-expert: File streaming and large data handling
  - nestjs-http-module-expert: External API integration with HttpModule
  - nestjs-compression-expert: Response compression for performance
  - nestjs-cookies-expert: Cookie management and secure cookie handling
  - nestjs-session-expert: Session management with Redis and other stores
  - nestjs-mvc-expert: Server-side rendering with template engines
  - nestjs-versioning-expert: API versioning strategies
  - nestjs-sse-expert: Server-Sent Events for real-time updates

  **NX Monorepo & Workspace Management:**
  - nx-monorepo-expert: Expert in NX monorepo management for NestJS projects using @nx/nest plugin (latest version). Specializes in workspace initialization (npx create-nx-workspace --preset=nest, nx add @nx/nest), application and library generation (apps, buildable libs, publishable libs with custom import paths), complete NestJS scaffolding (application, library, module, controller, service, resource, pipe, guard, interceptor, middleware, filter, decorator, provider, gateway, resolver, class, interface), dependency graph management (nx graph, project visualization), affected commands for CI/CD optimization (nx affected -t build/test/lint, nx run-many), project removal and migration, tag-based boundaries and @nx/enforce-module-boundaries lint rules (scope and type tags, constraint enforcement), caching strategies (local and remote Nx Replay), monorepo structure best practices (feature-based organization, data-access patterns, ui/util/model library types, 80/20 apps-to-libs ratio), executors configuration (build, serve, test, lint in project.json), project graph visualization (nx graph with task dependencies), VSCode debugging setup, Docker deployment with @nx/node:setup-docker, CLI plugin configuration (transformers for Swagger, etc.), and enterprise-scale monorepo patterns. Always enforces module boundaries and architectural best practices including proper library type classification (feature, data-access, ui, shell, utils, model), scope-based organization, public API enforcement via index.ts, and circular dependency prevention.

  **qBittorrent Integration:**
  - qbittorrent-api-expert: Expert in qBittorrent Web API (v4.1+) integration using @ctrl/qbittorrent TypeScript library (v9.9+). Specializes in both raw API endpoints and normalized library methods including getAllData(), getTorrent(), pauseTorrent(), resumeTorrent(), addTorrent(), removeTorrent(). Covers authentication (cookie-based SID, automatic login), application management (preferences, settings, version info), torrent management (add, delete, pause, resume, files, properties, categories, tags, trackers), transfer info (speeds, limits, quotas), RSS feed management (feeds, rules, article matching), search functionality (plugins, queries), logging (main log, peer log), and sync endpoints. Provides production-ready solutions for qBittorrent automation with type-safe implementations, normalized functions for torrent client interoperability, state management (export/import), error handling, rate limiting, and real-time updates.

  **Apollo GraphQL & Client:**
  - apollo-client-setup-expert: Apollo Client setup, configuration, ApolloProvider, link configuration, cache setup, and authentication
  - apollo-queries-expert: Apollo Client queries including useQuery, useLazyQuery, query options, polling, and refetching
  - apollo-mutations-expert: Apollo Client mutations including useMutation, optimistic responses, and cache updates
  - apollo-caching-expert: Apollo Client caching including InMemoryCache configuration, cache policies, field policies, and cache persistence
  - apollo-subscriptions-expert: Apollo Client subscriptions including useSubscription, WebSocket setup, and real-time updates
  - apollo-testing-expert: Testing Apollo Client applications with MockedProvider, query/mutation testing, and mocking strategies

  **Infrastructure & DevOps:**
  - devops-engineer: Deployment and CI/CD processes
  - cloud-architect: Cloud infrastructure design
  - kubernetes-specialist: Container orchestration

  **Security:**
  - security-engineer: Security implementation and best practices
  - security-auditor: Security assessment and compliance
  - penetration-tester: Security vulnerability testing

  **Debugging & Problem Solving:**
  - debugger: Code debugging and issue resolution
  - error-detective: Error investigation and root cause analysis
  - error-coordinator: Error handling coordination across systems

  **Project Management & Orchestration:**
  - multi-agent-coordinator: Complex workflows requiring multiple agents
  - workflow-orchestrator: Process design and workflow management
  - task-distributor: Task planning and distribution

  **Documentation & Communication:**
  - documentation-engineer: Technical documentation creation
  - technical-writer: Clear technical writing and communication
  - api-documenter: API documentation and specifications

  **Business & Product:**
  - business-analyst: Business requirements and analysis
  - product-manager: Product strategy and prioritization
  - ux-researcher: User experience research and design
  - legal-advisor: Compliance and regulatory requirements

  ## PLANNING STRATEGY REQUIREMENTS
  - ALWAYS prioritize foundational components first
  - MINIMIZE TODOs - build actual functionality instead of leaving notes
  - MANDATORY build order: data_models → utilities → core_services → business_logic → ui_components → integrations
  - USE dependency-aware incremental development approach
  - BUILD in small, testable increments (1-2 hour max per iteration)

  ## PROJECT DEFAULTS
  **TypeScript Standards:**
  - strict: true (non-negotiable)
  - no_any_tolerance: minimal (only for catch blocks and extreme edge cases)
  - warning_tolerance: zero (fix all warnings)

  **Testing Requirements:**
  - Framework preference: jest or vitest
  - Minimum coverage: 80% (enforced)
  - SEPARATE unit and integration test suites:
    - Unit tests: isolated, fast, mock external dependencies
    - Integration tests: test component interactions, use real dependencies
    - NEVER mix unit and integration test patterns

  # QUALITY GATES (ALL MUST PASS)
  ✅ No TypeScript errors
  ✅ All tests passing
  ✅ Coverage >= 80%
  ✅ No disabled warnings
  ✅ No deprecated APIs
  ✅ Minimal TODOs with proper justification
  ✅ Dependencies built before dependents
  ✅ Code reviewed by appropriate agent
  ✅ Working codebase state

  **FAILURE TO MEET QUALITY GATES = WORK IS NOT COMPLETE**