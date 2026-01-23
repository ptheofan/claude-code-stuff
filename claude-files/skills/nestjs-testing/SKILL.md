---
name: nestjs-testing
description: NestJS testing patterns using Test Pyramid with Solitary Unit Testing. Use when writing unit tests, integration tests, or E2E tests for NestJS applications.
---

# NestJS Testing Patterns

Use **context7** for Jest/NestJS testing API docs. This skill defines OUR testing methodology.

## Test Pyramid (Solitary Unit Testing)

Test from bottom up. Stub already-tested dependencies.

```
        /\
       /E2E\        ← Critical paths only
      /------\
     /Integration\   ← DB + internal integrations
    /--------------\
   /   Unit Tests   \ ← Foundation: smallest → largest
  /------------------\
```

## Unit Test Layers

### 1. Smallest Functions First
```typescript
// utils/email.util.spec.ts
describe('validateEmail', () => {
  it('should return true for valid email', () => {
    expect(validateEmail('test@example.com')).toBe(true);
  });
  
  it('should return false for invalid email', () => {
    expect(validateEmail('invalid')).toBe(false);
  });
});
```

### 2. Services (stub utilities)
```typescript
// users.service.spec.ts
describe('UsersService', () => {
  let service: UsersService;
  let repository: jest.Mocked<UserRepository>;

  beforeEach(async () => {
    const module = await Test.createTestingModule({
      providers: [
        UsersService,
        {
          provide: UserRepository,
          useValue: {
            findById: jest.fn(),
            create: jest.fn(),
          },
        },
      ],
    }).compile();

    service = module.get(UsersService);
    repository = module.get(UserRepository);
  });

  // validateEmail already tested → stub if needed
  // Repository methods → mock (external dependency)
});
```

### 3. Command/Query Handlers (stub services)
```typescript
describe('CreateUserHandler', () => {
  let handler: CreateUserHandler;
  let usersService: jest.Mocked<UsersService>;

  beforeEach(async () => {
    const module = await Test.createTestingModule({
      providers: [
        CreateUserHandler,
        {
          provide: UsersService,
          useValue: {
            create: jest.fn(), // Already unit tested
          },
        },
      ],
    }).compile();
  });

  it('should create user and return id', async () => {
    usersService.create.mockResolvedValue({ id: '123' });
    
    const result = await handler.execute(
      new CreateUserCommand('test@example.com', 'Test'),
    );
    
    expect(result).toBe('123');
  });
});
```

### 4. Controllers/Resolvers (input/output only)
```typescript
describe('UsersController', () => {
  // Focus on:
  // - Input validation (DTOs)
  // - Output shape
  // - HTTP status codes
  // - NOT business logic (already tested in services/handlers)

  it('should return 201 on successful creation', async () => {
    commandBus.execute.mockResolvedValue('123');

    const result = await controller.create(validDto);

    expect(result).toEqual({ id: '123' });
    expect(commandBus.execute).toHaveBeenCalledWith(
      expect.any(CreateUserCommand),
    );
  });

  it('should validate input dto', async () => {
    const invalidDto = { email: 'invalid' }; // missing name
    
    await expect(controller.create(invalidDto))
      .rejects.toThrow(BadRequestException);
  });
});
```

## Integration Tests

Focus on DB operations and internal integrations. Self-contained, no production data.

```typescript
describe('UserRepository (Integration)', () => {
  let repository: UserRepository;
  let dataSource: DataSource;

  beforeAll(async () => {
    // Use test database
    const module = await Test.createTestingModule({
      imports: [
        TypeOrmModule.forRoot(testDbConfig),
        TypeOrmModule.forFeature([UserEntity]),
      ],
      providers: [UserRepository],
    }).compile();

    repository = module.get(UserRepository);
    dataSource = module.get(DataSource);
  });

  beforeEach(async () => {
    // Clean slate for each test
    await dataSource.synchronize(true);
  });

  afterAll(async () => {
    await dataSource.destroy();
  });

  it('should persist and retrieve user', async () => {
    const user = await repository.create({
      email: 'test@example.com',
      name: 'Test',
    });

    const found = await repository.findById(user.id);

    expect(found).toMatchObject({
      email: 'test@example.com',
      name: 'Test',
    });
  });
});
```

## E2E Tests

Critical business paths only. Skip what's covered by unit/integration.

```typescript
describe('User Registration (E2E)', () => {
  // Happy path
  it('should register user and return JWT', async () => {
    const response = await request(app.getHttpServer())
      .post('/auth/register')
      .send({ email: 'new@example.com', password: 'Password123!' })
      .expect(201);

    expect(response.body).toHaveProperty('access_token');
  });

  // Critical failure path
  it('should reject duplicate email', async () => {
    await createUser({ email: 'existing@example.com' });

    await request(app.getHttpServer())
      .post('/auth/register')
      .send({ email: 'existing@example.com', password: 'Password123!' })
      .expect(409);
  });
});
```

## Logger Stubbing

Keep test output noise-free. Use logger as test evidence.

```typescript
describe('PaymentService', () => {
  let logger: jest.Mocked<Logger>;

  beforeEach(() => {
    logger = {
      log: jest.fn(),
      error: jest.fn(),
      warn: jest.fn(),
    } as any;

    // Inject mocked logger
  });

  it('should log error on payment failure', async () => {
    paymentGateway.charge.mockRejectedValue(new Error('Declined'));

    await expect(service.processPayment(data)).rejects.toThrow();

    expect(logger.error).toHaveBeenCalledWith(
      expect.stringContaining('Payment failed'),
      expect.objectContaining({ orderId: data.orderId }),
    );
  });
});
```

## Test Coverage

Prioritize meaningful coverage over percentage:
- ✅ Critical paths thoroughly covered
- ✅ Edge cases and error scenarios
- ✅ Domain exceptions tested
- ❌ NOT: trivial getters/setters just to hit numbers
