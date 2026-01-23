---
name: nestjs-core
description: NestJS development patterns including TypeORM (Data Mapper), CQRS, authentication, validation, and module architecture. Use when building NestJS features, services, modules, or APIs.
---

# NestJS Core Patterns

Use **context7** for NestJS/TypeORM API documentation. This skill defines OUR conventions.

## Architecture Rules

- **Clean Architecture**: Import down, emit events up
- **Strict module boundaries**: Never import another module's entities/repositories directly
- **Expose via module API**: Services, DTOs, interfaces only
- **Domain Errors**: Throw specific exceptions, not generic errors

## TypeORM Patterns

### Data Mapper Only (NO Active Record)
```typescript
// ✅ Repository pattern
@Injectable()
export class UserRepository {
  constructor(
    @InjectRepository(UserEntity)
    private readonly repository: Repository<UserEntity>,
  ) {}
}

// ❌ NEVER Active Record
class User extends BaseEntity { } // FORBIDDEN
```

### Transaction Management
```typescript
import { Transactional } from 'typeorm-transactional';

@Injectable()
export class UserService {
  @Transactional()
  async createWithProfile(data: CreateUserDto): Promise<User> {
    const user = await this.userRepository.create(data);
    await this.profileRepository.create({ userId: user.id });
    return user;
  }
}
```

### Entity ↔ DTO Mapping
```typescript
// Separate mapper class - entities never leak outside module
@Injectable()
export class UserMapper {
  toDto(entity: UserEntity): UserDto { ... }
  toEntity(dto: CreateUserDto): Partial<UserEntity> { ... }
}
```

## CQRS Patterns

### Commands (Write)
```typescript
// commands/create-user.command.ts
export class CreateUserCommand {
  constructor(
    public readonly email: string,
    public readonly name: string,
  ) {}
}

// commands/handlers/create-user.handler.ts
@CommandHandler(CreateUserCommand)
export class CreateUserHandler implements ICommandHandler<CreateUserCommand> {
  async execute(command: CreateUserCommand): Promise<string> {
    // Returns ID or result
  }
}
```

### Queries (Read)
```typescript
@QueryHandler(GetUserQuery)
export class GetUserHandler implements IQueryHandler<GetUserQuery> {
  async execute(query: GetUserQuery): Promise<UserDto> {
    // Returns DTO, never entity
  }
}
```

### Domain Events
```typescript
// Emit after state change
export class UserCreatedEvent {
  constructor(public readonly userId: string) {}
}

// In aggregate/service
this.eventBus.publish(new UserCreatedEvent(user.id));
```

## Authentication Patterns

### JWT + Guards
```typescript
@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
  async validate(payload: JwtPayload): Promise<User> {
    const user = await this.usersService.findById(payload.sub);
    if (!user) throw new UnauthorizedException();
    return user; // Attached to request.user
  }
}
```

### GraphQL Auth Guard
```typescript
@Injectable()
export class GqlAuthGuard extends AuthGuard('jwt') {
  getRequest(context: ExecutionContext) {
    return GqlExecutionContext.create(context).getContext().req;
  }
}
```

## Validation

Always use class-validator + class-transformer:
```typescript
@Controller('users')
export class UsersController {
  @Post()
  @UsePipes(new ValidationPipe({ transform: true, whitelist: true }))
  create(@Body() dto: CreateUserDto) { }
}
```

## Module Structure

```
users/
├── users.module.ts           # Public API
├── users.service.ts          # Orchestration
├── users.repository.ts       # Data access
├── users.mapper.ts           # Entity ↔ DTO
├── commands/
│   ├── create-user.command.ts
│   └── handlers/
├── queries/
├── events/
├── dto/
│   ├── create-user.dto.ts
│   └── user.dto.ts
├── entities/
│   └── user.entity.ts        # NEVER export from module
├── interfaces/
└── exceptions/
    └── user-not-found.exception.ts
```

## Domain Exceptions

```typescript
export class UserNotFoundException extends NotFoundException {
  constructor(userId: string) {
    super({
      code: 'USER_NOT_FOUND',
      message: `User ${userId} not found`,
      userId, // Traceable data
    });
  }
}
```
