# Claude Code Configuration

A comprehensive Claude Code setup with strict development policies, specialized agents, and quality enforcement. This repository contains configuration files and deployment scripts to set up a professional development environment with Claude Code.

## 🚀 Quick Start

### Installation
```bash
# Clone this repository
git clone <your-repo-url>
cd claude-code-stuff

# Deploy to your Claude Code directory
./deploy

# For sync mode (removes old files not in this repo)
./deploy --sync

# For custom target directory
./deploy --target /path/to/custom/claude/dir

# Dry run to see what would change
./deploy --dry-run
```

### Deployment Options
- **Basic deploy**: `./deploy` - Merges files to `~/.claude/`
- **Sync mode**: `./deploy --sync` - Removes files not in source (use with caution)
- **Custom target**: `./deploy --target /path/to/dir` - Deploy to custom location
- **Dry run**: `./deploy --dry-run` - Preview changes without applying
- **Help**: `./deploy --help` - Show all options

## 📋 What's Included

### Configuration Files
- **`config.yaml`** - Main Claude Code configuration with development policies
- **`agents/`** - Specialized agent configurations organized by domain
- **`deploy`** - Deployment script with sync, dry-run, and custom target support

### Agent Categories
- **Core Development** (01-core-development/) - Backend, frontend, fullstack, API design
- **Language Specialists** (02-language-specialists/) - TypeScript, Python, React, etc.
- **Infrastructure** (03-infrastructure/) - DevOps, cloud, Kubernetes, security
- **Quality & Security** (04-quality-security/) - Testing, code review, security audit
- **Data & AI** (05-data-ai/) - ML, data science, AI engineering
- **Developer Experience** (06-developer-experience/) - Tooling, documentation, refactoring
- **Specialized Domains** (07-specialized-domains/) - Blockchain, IoT, fintech, gaming
- **Business & Product** (08-business-product/) - PM, UX, legal, marketing
- **Meta Orchestration** (09-meta-orchestration/) - Multi-agent coordination, workflow
- **Research & Analysis** (10-research-analysis/) - Market research, competitive analysis

## 🎯 Development Policy

This configuration enforces strict development standards. Here's what you need to know:

### Core Principles
- **Quality over speed** - Always
- **Type safety first** - Minimal `any` types, no shortcuts
- **Test-driven development** - 80%+ coverage required
- **Incremental development** - Small, testable parts with dependency-aware planning
- **KISS principle** - Simple, readable solutions
- **Ground-up building** - Build dependencies before dependents to minimize TODOs

### Quality Gates
Every feature must pass these checks:
- ✅ All TypeScript errors resolved
- ✅ All tests passing (no deleting tests to fix builds)
- ✅ 80%+ test coverage
- ✅ No disabled warnings or deprecated APIs
- ✅ Minimal TODOs with proper justification and ticket references
- ✅ Dependencies built before dependents (ground-up approach)
- ✅ Code reviewed by appropriate agent
- ✅ Working codebase state

## 🏗️ Dependency-Aware Development Strategy

### Recommended Build Order
To minimize TODOs and build stable foundations:

1. **Data Models** - Types, interfaces, schemas (no dependencies)
2. **Utilities** - Helper functions, constants, pure functions
3. **Core Services** - Business logic, validation, domain services
4. **Business Logic** - Feature-specific implementations
5. **UI Components** - User interface that consumes business logic
6. **Integrations** - External API connections and third-party services

### Planning Questions
Before starting any feature, ask:
- What does this component depend on?
- What dependencies can I build first?
- How can I minimize TODOs in this implementation?
- Which foundational pieces are missing?

### Example: User Registration Feature
```typescript
// ❌ Wrong order - building UI first
function UserRegistrationForm() {
  // TODO: Add validation once validation service exists
  // TODO: Add user creation once user service exists
  // TODO: Add error handling once error service exists
  return <form>...</form>;
}

// ✅ Right order - dependencies first
// 1. Data models
interface User { id: string; email: string; name: string; }
interface RegistrationData { email: string; name: string; password: string; }

// 2. Utilities
function validateEmail(email: string): boolean { /* ... */ }
function hashPassword(password: string): string { /* ... */ }

// 3. Core services
function validateRegistration(data: RegistrationData): ValidationResult { /* ... */ }
function createUser(data: RegistrationData): Promise<User> { /* ... */ }

// 4. UI components (built on stable foundation)
function UserRegistrationForm() {
  // All dependencies exist - no TODOs needed!
  const handleSubmit = async (data: RegistrationData) => {
    const validation = validateRegistration(data);
    if (validation.isValid) {
      return await createUser(data);
    }
    // Handle validation errors
  };
  return <form onSubmit={handleSubmit}>...</form>;
}
```

## 📚 Development Policy Examples

*For team reference - examples of what the config.yaml enforces*

### TypeScript Type Safety

#### ❌ **Bad Examples**
```typescript
// Using 'any' for convenience
function processData(data: any): any {
  return data.someProperty;
}

// Casting for convenience in production code
const user = apiResponse as unknown as User;

// Disabling warnings
// @ts-ignore
const result = someComplexObject.undefinedProperty;
```

#### ✅ **Good Examples**
```typescript
// Proper typing
interface UserData {
  id: number;
  name: string;
}
function processData(data: UserData): string {
  return data.name;
}

// Exception: Extreme cases only
try {
  riskyOperation();
} catch (error: unknown) {
  console.error('Error:', error);
}

// Exception: Test mocking only
const mockUser = { id: 1, name: 'Test' } as unknown as CompleteUserInterface;
```

### KISS Principle

#### ❌ **Too Clever/Complex**
```typescript
// One-liner that's hard to understand
const result = items?.reduce((a,b)=>({...a,[b.id]:b}),[]);

// Over-engineered abstraction
class AbstractFactoryManagerSingleton {
  // 50 lines of complex pattern implementation
}

// Cryptic variable names
const u = users.filter(x => x.a > 5).map(y => y.b);
```

#### ✅ **Keep It Simple**
```typescript
// Clear and readable
const result = {};
for (const item of items || []) {
  result[item.id] = item;
}

// Simple, direct approach
function createUser(data: UserInput): User {
  return new User(data);
}

// Descriptive names
const activeUsers = users
  .filter(user => user.isActive)
  .map(user => user.profile);
```

### Testing Discipline

#### ❌ **Bad Test Practices**
```typescript
// Deleting failing tests to make build pass
// it('should calculate tax correctly', () => {
//   expect(calculateTax(100, 0.1)).toBe(10);
// }); 
// ❌ Commented out because it was failing

// Over-mocking internal logic
const mockCalculator = jest.fn().mockReturnValue(42);
const result = processOrder(mockCalculator); // Mocking business logic

// Skipping tests for convenience
describe.skip('Payment processing', () => {
  // ❌ Skipped because "it's complicated"
});
```

#### ✅ **Good Test Practices**
```typescript
// Fix the failing test properly
it('should calculate tax correctly', () => {
  expect(calculateTax(100, 0.1)).toBe(10); // Fixed the calculation bug
});

// Mock only external dependencies
const mockPaymentAPI = jest.fn();
const result = processPayment(mockPaymentAPI, orderData);

// Remove redundant tests with justification
// Removed duplicate test - same scenario covered in line 15
it('should handle valid orders', () => {
  expect(processOrder(validOrder)).toBeDefined();
});
```

### Deprecated Functions

#### ❌ **Using Deprecated APIs**
```typescript
// Old JavaScript methods
document.write('<script>...</script>');
new Date().getYear(); // Returns 2-digit year
escape('hello world');

// Deprecated React patterns
componentWillMount() {
  this.fetchData();
}

// Deprecated Node.js
new Buffer('data'); // Security vulnerability
fs.exists('/path', callback); // Deprecated
```

#### ✅ **Modern Alternatives**
```typescript
// Modern DOM manipulation
const script = document.createElement('script');
document.head.appendChild(script);
new Date().getFullYear(); // Returns 4-digit year
encodeURIComponent('hello world');

// Modern React patterns
useEffect(() => {
  fetchData();
}, []);

// Modern Node.js
Buffer.from('data');
fs.access('/path', callback); // Or fs.promises.access
```

### TODO Comments and Exceptions

#### ❌ **Invalid TODO Usage**
```typescript
// Lazy TODO without context
// TODO: fix this
const result = hackyImplementation(); // ❌ No ticket, no context

// TODO without timeline or dependency justification
// TODO: make this better
const messyFunction = () => { /* ... */ }; // ❌ Vague, no plan

// TODO to avoid doing work
// TODO: add tests later
function importantFeature() { /* no tests written */ } // ❌ Shortcut
```

#### ✅ **Valid TODO Usage**
```typescript
// Valid TODO for missing feature dependency
// TODO: Replace with UserService.validateEmail() once auth module is complete (TICKET-456)
const isValidEmail = email.includes('@'); // Temporary validation

// Valid TODO for dependent feature implementation
// TODO: Implement proper caching once Redis service is deployed (INFRA-789)
const getCachedData = (key: string) => {
  return localStorage.getItem(key); // Temporary browser storage
};

// Valid TODO for external dependency
// TODO: Use PaymentProcessor.charge() once payment service API is finalized (PAY-123)
const processPayment = async (amount: number) => {
  return { success: true, id: 'temp-' + Date.now() }; // Mock implementation
};
```

### Agent Selection Examples

#### ❌ **Wrong Agent for Task**
```typescript
// Using general agent for TypeScript-specific work
"Fix this TypeScript compilation error" 
// → Should use typescript-pro agent

// Using single agent for complex multi-domain task
"Build a full authentication system with security audit"
// → Should use multiple agents: backend-developer + security-engineer + test-automator
```

#### ✅ **Right Agent for Task**
```typescript
// TypeScript issues → typescript-pro
"Have typescript-pro fix these type compilation errors"

// Code quality → code-reviewer
"Have code-reviewer analyze this component for best practices"

// Complex workflow → multi-agent-coordinator
"Have multi-agent-coordinator plan the authentication system using backend-developer, security-engineer, and test-automator"

// Performance issues → performance-engineer
"Have performance-engineer optimize this database query"
```

### Definition of Done Examples

#### ❌ **Not Ready for Delivery**
```typescript
// ❌ Has TypeScript errors
function processUser(user) { // Missing type annotation
  return user.invalidProperty; // Property doesn't exist
}

// ❌ No tests
// No test file exists

// ❌ Uses deprecated function
const year = new Date().getYear(); // Deprecated

// ❌ Over-complex
const result = data?.reduce((acc, item) => 
  item.type === 'valid' ? {...acc, [item.id]: item.data?.value || 0} : acc, {});
```

#### ✅ **Ready for Delivery**
```typescript
// ✅ Proper TypeScript
interface User {
  id: string;
  name: string;
}

function processUser(user: User): string {
  return user.name;
}

// ✅ Has comprehensive tests
describe('processUser', () => {
  it('should return user name', () => {
    const user = { id: '1', name: 'John' };
    expect(processUser(user)).toBe('John');
  });
});

// ✅ Uses modern APIs
const year = new Date().getFullYear();

// ✅ Simple and readable
const validItems = data?.filter(item => item.type === 'valid') || [];
const result = {};
for (const item of validItems) {
  result[item.id] = item.data?.value || 0;
}
```

### Emergency Exception Examples

#### ❌ **Invalid Exceptions**
```typescript
// "It's faster this way"
const data = apiResponse as any; // ❌ Convenience, not emergency

// "We'll fix it later"
// @ts-ignore - TODO: fix types
return someFunction(); // ❌ Not documented properly

// "The test is flaky"
// describe.skip('Integration tests', () => { // ❌ Avoiding real issues
```

#### ✅ **Valid Exceptions**
```typescript
// Production incident requiring immediate fix
// EMERGENCY: Critical bug fix for payment system
// Technical debt ticket: TECH-123 created to properly type this
const urgentFix = legacyPaymentData as any;

// Well-documented third-party integration issue
// @ts-expect-error - Library @types/old-lib missing return type definition
// See issue: https://github.com/old-lib/types/issues/456
const result = oldLibraryFunction();

// Temporary migration step with clear timeline
// MIGRATION: Remove this cast after API v2 migration (Sprint 23)
const userData = legacyApiResponse as unknown as NewUserFormat;
```

### Coverage and Testing Mix

#### ❌ **Poor Test Coverage**
```typescript
// Only testing happy path
it('should process payment', () => {
  expect(processPayment(validData)).toBe(true);
});
// Missing: error cases, edge cases, validation failures
```

#### ✅ **Comprehensive Test Coverage**
```typescript
describe('processPayment', () => {
  it('should process valid payment', () => {
    expect(processPayment(validData)).toBe(true);
  });
  
  it('should reject invalid card', () => {
    expect(() => processPayment(invalidCard)).toThrow('Invalid card');
  });
  
  it('should handle network timeout', async () => {
    mockAPI.mockRejectedValue(new TimeoutError());
    await expect(processPayment(validData)).rejects.toThrow('Timeout');
  });
  
  it('should validate required fields', () => {
    expect(() => processPayment({})).toThrow('Missing required fields');
  });
});
```

## 🤖 Agent Selection Guidelines

### By Task Type
- **Code Quality**: `code-reviewer`, `architect-reviewer`, `refactoring-specialist`
- **Testing**: `qa-expert`, `test-automator`, `accessibility-tester`, `performance-engineer`
- **TypeScript**: `typescript-pro`, `react-specialist`, `nextjs-developer`
- **Infrastructure**: `devops-engineer`, `cloud-architect`, `kubernetes-specialist`
- **Security**: `security-engineer`, `security-auditor`, `penetration-tester`
- **Debugging**: `debugger`, `error-detective`, `error-coordinator`
- **Orchestration**: `multi-agent-coordinator`, `workflow-orchestrator`, `task-distributor`
- **Documentation**: `documentation-engineer`, `technical-writer`, `api-documenter`

### Best Practices
- Use specialized agents for domain-specific tasks
- Coordinate multiple agents for complex workflows
- Always select appropriate agent before starting work
- Reference agent guidelines in `config.yaml`

## 🔧 Configuration Details

### Session Behavior
- `always_check_policy: true` - Enforces policy compliance
- `require_explicit_confirmation: true` - Requires confirmation for major changes
- `enforce_testing_standards: true` - Enforces testing requirements
- `agent_selection_required: true` - Requires agent selection before work
- `prefer_specialized_agents: true` - Prefers domain-specific agents
- `minimum_coverage: 80` - Enforces 80% test coverage

### Project Defaults
- **TypeScript**: Strict mode, minimal `any` tolerance, zero warning tolerance
- **Testing**: Required, 80% minimum coverage, Jest/Vitest preferred
- **Development**: Dependency-aware incremental approach with validation

## 📖 Remember

- **Quality over speed** - always
- **Fix the problem, don't hide it**
- **Simple is better than clever**
- **Tests are documentation** - keep them meaningful
- **Use the right tool for the job** - specialized agents exist for a reason
- **Build dependencies before dependents** - minimize TODOs

## 🤝 Contributing

This configuration is designed for professional development teams. When contributing:

1. Follow the development policy strictly
2. Use appropriate agents for your tasks
3. Maintain test coverage above 80%
4. Keep code simple and readable
5. Document any exceptions properly

## 📄 License

### Third-Party Components
**Important**: Parts of this repository are copied from other repositories and should be used according to their original licenses:

- **Agent configurations** from [awesome-claude-code-subagents](https://github.com/VoltAgent/awesome-claude-code-subagents) - Please check their LICENSE file for usage terms
- **Any other copied content** - Always verify and comply with the original repository's license

When using this configuration, you are responsible for ensuring compliance with all applicable licenses from the source repositories.

## 🔗 Resources

- [Claude Code Documentation](https://claude.ai/code)
- [Awesome Claude Code Subagents](https://github.com/VoltAgent/awesome-claude-code-subagents)
- [TypeScript Best Practices](https://typescript-eslint.io/rules/)
- [Testing Best Practices](https://jestjs.io/docs/getting-started)