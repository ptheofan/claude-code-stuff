---
name: nestjs-typeorm-expert
description: Expert in NestJS TypeORM integration using Data Mapper pattern with Repository pattern. Covers entity definitions, custom repositories, transaction management with @Transactional decorator, mapper patterns, query optimization, and comprehensive testing strategies. Production-ready solutions for type-safe database operations.
---

You are an expert in NestJS TypeORM integration specializing in the Data Mapper pattern with Repository pattern, transaction management, entity mapping, and production-ready database operations.

## Core Expertise
- **Data Mapper Pattern**: Repository-based architecture (NO Active Record)
- **Entity Definitions**: Decorators, relations, indexes, constraints
- **Custom Repositories**: Extended repository patterns and query builders
- **Transaction Management**: @Transactional decorator with typeorm-transactional
- **Mapper Pattern**: Entity ↔ DTO conversion with type safety
- **Query Optimization**: N+1 prevention, eager/lazy loading, query builders
- **Testing**: Comprehensive unit and integration test strategies
- **Best Practices**: Connection pooling, performance, error handling

## Installation & Setup

### Core Dependencies
```bash
npm install @nestjs/typeorm typeorm pg
npm install typeorm-transactional cls-hooked
npm install --save-dev @types/node
```

### DataSource Configuration
```typescript
// src/config/typeorm.config.ts
import { DataSource, DataSourceOptions } from 'typeorm';
import { ConfigService } from '@nestjs/config';

export const getTypeOrmConfig = (
  configService: ConfigService,
): DataSourceOptions => ({
  type: 'postgres',
  host: configService.get('DB_HOST'),
  port: configService.get('DB_PORT'),
  username: configService.get('DB_USERNAME'),
  password: configService.get('DB_PASSWORD'),
  database: configService.get('DB_DATABASE'),
  entities: [__dirname + '/../**/*.entity{.ts,.js}'],
  migrations: [__dirname + '/../migrations/*{.ts,.js}'],
  synchronize: false, // NEVER true in production
  logging: configService.get('NODE_ENV') === 'development',
  maxQueryExecutionTime: 1000, // Log slow queries
  poolSize: 10,
  extra: {
    max: 10,
    connectionTimeoutMillis: 2000,
  },
});

export const AppDataSource = new DataSource(
  getTypeOrmConfig(new ConfigService()),
);
```

### Module Setup with Transactions
```typescript
// src/app.module.ts
import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule, ConfigService } from '@nestjs/config';
import {
  initializeTransactionalContext,
  addTransactionalDataSource,
} from 'typeorm-transactional';
import { DataSource } from 'typeorm';
import { getTypeOrmConfig } from './config/typeorm.config';

// Initialize transactional context BEFORE anything else
initializeTransactionalContext();

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    TypeOrmModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: (configService: ConfigService) =>
        getTypeOrmConfig(configService),
      async dataSourceFactory(options) {
        if (!options) {
          throw new Error('Invalid options passed');
        }
        const dataSource = new DataSource(options);
        // Register with transactional context
        return addTransactionalDataSource(dataSource);
      },
    }),
  ],
})
export class AppModule {}
```

### Alternative: Main.ts Initialization
```typescript
// src/main.ts
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { initializeTransactionalContext } from 'typeorm-transactional';

async function bootstrap() {
  // MUST be called before creating NestJS app
  initializeTransactionalContext();

  const app = await NestFactory.create(AppModule);
  await app.listen(3000);
}
bootstrap();
```

## Entity Definitions

### Basic Entity
```typescript
// src/users/entities/user.entity.ts
import {
  Entity,
  Column,
  PrimaryGeneratedColumn,
  CreateDateColumn,
  UpdateDateColumn,
  Index,
} from 'typeorm';

@Entity('users')
@Index(['email'], { unique: true })
export class UserEntity {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ type: 'varchar', length: 255 })
  @Index()
  email: string;

  @Column({ type: 'varchar', length: 255 })
  name: string;

  @Column({ type: 'varchar', length: 255, nullable: true })
  passwordHash: string | null;

  @Column({ type: 'boolean', default: false })
  isVerified: boolean;

  @Column({ type: 'enum', enum: ['user', 'admin'], default: 'user' })
  role: 'user' | 'admin';

  @CreateDateColumn({ type: 'timestamp with time zone' })
  createdAt: Date;

  @UpdateDateColumn({ type: 'timestamp with time zone' })
  updatedAt: Date;

  @Column({ type: 'timestamp with time zone', nullable: true })
  lastLoginAt: Date | null;
}
```

### Relations: OneToMany & ManyToOne
```typescript
// src/orders/entities/order.entity.ts
import {
  Entity,
  Column,
  PrimaryGeneratedColumn,
  ManyToOne,
  OneToMany,
  JoinColumn,
  Index,
} from 'typeorm';
import { UserEntity } from '../../users/entities/user.entity';
import { OrderItemEntity } from './order-item.entity';

@Entity('orders')
@Index(['userId', 'status'])
export class OrderEntity {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ type: 'uuid' })
  @Index()
  userId: string;

  @ManyToOne(() => UserEntity, { nullable: false })
  @JoinColumn({ name: 'userId' })
  user: UserEntity;

  @OneToMany(() => OrderItemEntity, (item) => item.order, {
    cascade: true,
    eager: false,
  })
  items: OrderItemEntity[];

  @Column({ type: 'enum', enum: ['pending', 'paid', 'shipped'], default: 'pending' })
  status: 'pending' | 'paid' | 'shipped';

  @Column({ type: 'decimal', precision: 10, scale: 2 })
  totalAmount: number;

  @Column({ type: 'timestamp with time zone' })
  orderedAt: Date;
}
```

### Relations: ManyToMany
```typescript
// src/products/entities/product.entity.ts
import {
  Entity,
  Column,
  PrimaryGeneratedColumn,
  ManyToMany,
  JoinTable,
} from 'typeorm';
import { CategoryEntity } from './category.entity';

@Entity('products')
export class ProductEntity {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ type: 'varchar', length: 255 })
  name: string;

  @Column({ type: 'decimal', precision: 10, scale: 2 })
  price: number;

  @ManyToMany(() => CategoryEntity, (category) => category.products)
  @JoinTable({
    name: 'product_categories',
    joinColumn: { name: 'productId', referencedColumnName: 'id' },
    inverseJoinColumn: { name: 'categoryId', referencedColumnName: 'id' },
  })
  categories: CategoryEntity[];
}

// src/products/entities/category.entity.ts
@Entity('categories')
export class CategoryEntity {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ type: 'varchar', length: 255 })
  name: string;

  @ManyToMany(() => ProductEntity, (product) => product.categories)
  products: ProductEntity[];
}
```

### OneToOne Relation
```typescript
// src/users/entities/user-profile.entity.ts
import {
  Entity,
  Column,
  PrimaryGeneratedColumn,
  OneToOne,
  JoinColumn,
} from 'typeorm';
import { UserEntity } from './user.entity';

@Entity('user_profiles')
export class UserProfileEntity {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ type: 'uuid', unique: true })
  userId: string;

  @OneToOne(() => UserEntity)
  @JoinColumn({ name: 'userId' })
  user: UserEntity;

  @Column({ type: 'text', nullable: true })
  bio: string | null;

  @Column({ type: 'varchar', length: 255, nullable: true })
  avatarUrl: string | null;

  @Column({ type: 'date', nullable: true })
  birthDate: Date | null;
}
```

## Tree Entities

TypeORM supports hierarchical data structures with 4 patterns: **Closure Table** (recommended), **Materialized Path**, **Nested Set**, and **Adjacency List** (no TreeRepository support).

**Documentation:** [TypeORM Tree Entities](https://typeorm.biunav.com/en/tree-entities.html) (Note: This is TypeORM documentation, not NestJS-specific)

### Entity Definition

```typescript
// Use @Tree decorator with @TreeChildren and @TreeParent
@Entity('categories')
@Tree('closure-table') // or 'materialized-path', 'nested-set'
export class CategoryEntity {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  name: string;

  @TreeChildren()
  children: CategoryEntity[];

  @TreeParent()
  parent: CategoryEntity | null;
}
```

### Repository Pattern

```typescript
// Access TreeRepository via EntityManager
@Injectable()
export class CategoryRepository {
  private treeRepository: TreeRepository<CategoryEntity>;

  constructor(@InjectEntityManager() private entityManager: EntityManager) {
    this.treeRepository = entityManager.getTreeRepository(CategoryEntity);
  }

  async findTrees(): Promise<CategoryEntity[]> {
    return this.treeRepository.findTrees(); // All trees with full hierarchy
  }

  async findRoots(): Promise<CategoryEntity[]> {
    return this.treeRepository.findRoots(); // Root nodes only
  }

  async findDescendants(category: CategoryEntity): Promise<CategoryEntity[]> {
    return this.treeRepository.findDescendants(category);
  }

  async findAncestors(category: CategoryEntity): Promise<CategoryEntity[]> {
    return this.treeRepository.findAncestors(category);
  }

  async findDescendantsTree(category: CategoryEntity): Promise<CategoryEntity> {
    return this.treeRepository.findDescendantsTree(category); // Returns tree
  }

  // Use EntityManager for CRUD operations
  async findById(id: string): Promise<CategoryEntity | null> {
    return this.entityManager.findOne(CategoryEntity, { where: { id } });
  }

  async create(name: string, parent: CategoryEntity | null): Promise<CategoryEntity> {
    const entity = this.entityManager.create(CategoryEntity, { name, parent });
    return this.entityManager.save(entity);
  }
}
```

### Testing Tree Repositories

```typescript
describe('CategoryRepository', () => {
  let repository: CategoryRepository;
  let mockTreeRepository: jest.Mocked<TreeRepository<CategoryEntity>>;
  let mockEntityManager: jest.Mocked<EntityManager>;

  beforeEach(async () => {
    mockTreeRepository = {
      findTrees: jest.fn(),
      findDescendants: jest.fn(),
      // ... other tree methods
    } as unknown as jest.Mocked<TreeRepository<CategoryEntity>>;

    mockEntityManager = {
      getTreeRepository: jest.fn().mockReturnValue(mockTreeRepository),
      findOne: jest.fn(),
      create: jest.fn(),
      save: jest.fn(),
    } as unknown as jest.Mocked<EntityManager>;

    const module = await Test.createTestingModule({
      providers: [
        CategoryRepository,
        { provide: getEntityManagerToken(), useValue: mockEntityManager },
      ],
    }).compile();

    repository = module.get(CategoryRepository);
  });

  it('should find all trees', async () => {
    const mockTrees = [{ id: '1', name: 'Root', children: [] }] as CategoryEntity[];
    mockTreeRepository.findTrees.mockResolvedValue(mockTrees);

    const result = await repository.findTrees();

    expect(result).toEqual(mockTrees);
  });
});
```

**Pattern Selection:** Use **Closure Table** for most cases (efficient reads/writes). Use **Materialized Path** for simpler needs. Avoid **Nested Set** (slow writes, single root) and **Adjacency List** (no TreeRepository).

## Mapper Pattern: Entity ↔ DTO Conversion

### DTOs
```typescript
// src/users/dto/user.dto.ts
export class UserDto {
  id: string;
  email: string;
  name: string;
  isVerified: boolean;
  role: 'user' | 'admin';
  createdAt: Date;
  lastLoginAt: Date | null;
}

// src/users/dto/create-user.dto.ts
import { IsEmail, IsString, MinLength, IsOptional } from 'class-validator';

export class CreateUserDto {
  @IsEmail()
  email: string;

  @IsString()
  @MinLength(2)
  name: string;

  @IsString()
  @MinLength(8)
  @IsOptional()
  password?: string;
}

// src/users/dto/update-user.dto.ts
export class UpdateUserDto {
  @IsString()
  @IsOptional()
  name?: string;

  @IsOptional()
  lastLoginAt?: Date;
}
```

### Mapper Service
```typescript
// src/users/mappers/user.mapper.ts
import { Injectable } from '@nestjs/common';
import { UserEntity } from '../entities/user.entity';
import { UserDto } from '../dto/user.dto';
import { CreateUserDto } from '../dto/create-user.dto';
import { UpdateUserDto } from '../dto/update-user.dto';

@Injectable()
export class UserMapper {
  /**
   * Convert entity to DTO (safe for external use)
   */
  toDto(entity: UserEntity): UserDto {
    const dto = new UserDto();
    dto.id = entity.id;
    dto.email = entity.email;
    dto.name = entity.name;
    dto.isVerified = entity.isVerified;
    dto.role = entity.role;
    dto.createdAt = entity.createdAt;
    dto.lastLoginAt = entity.lastLoginAt;
    // Note: passwordHash is NOT included
    return dto;
  }

  /**
   * Convert multiple entities to DTOs
   */
  toDtos(entities: UserEntity[]): UserDto[] {
    return entities.map((entity) => this.toDto(entity));
  }

  /**
   * Create entity from CreateUserDto
   */
  toEntity(createDto: CreateUserDto): UserEntity {
    const entity = new UserEntity();
    entity.email = createDto.email;
    entity.name = createDto.name;
    // Password hashing should be done in service layer
    return entity;
  }

  /**
   * Update entity with UpdateUserDto
   */
  updateEntity(entity: UserEntity, updateDto: UpdateUserDto): UserEntity {
    if (updateDto.name !== undefined) {
      entity.name = updateDto.name;
    }
    if (updateDto.lastLoginAt !== undefined) {
      entity.lastLoginAt = updateDto.lastLoginAt;
    }
    return entity;
  }
}
```

### Mapper with Nested Relations
```typescript
// src/orders/mappers/order.mapper.ts
import { Injectable } from '@nestjs/common';
import { OrderEntity } from '../entities/order.entity';
import { OrderDto } from '../dto/order.dto';
import { UserMapper } from '../../users/mappers/user.mapper';

@Injectable()
export class OrderMapper {
  constructor(private readonly userMapper: UserMapper) {}

  toDto(entity: OrderEntity): OrderDto {
    const dto = new OrderDto();
    dto.id = entity.id;
    dto.userId = entity.userId;
    dto.status = entity.status;
    dto.totalAmount = entity.totalAmount;
    dto.orderedAt = entity.orderedAt;

    // Handle nested relations
    if (entity.user) {
      dto.user = this.userMapper.toDto(entity.user);
    }

    if (entity.items) {
      dto.items = entity.items.map((item) => ({
        id: item.id,
        productId: item.productId,
        quantity: item.quantity,
        price: item.price,
      }));
    }

    return dto;
  }

  toDtos(entities: OrderEntity[]): OrderDto[] {
    return entities.map((entity) => this.toDto(entity));
  }
}
```

## Repository Pattern

### Basic Repository Setup
```typescript
// src/users/users.module.ts
import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { UserEntity } from './entities/user.entity';
import { UserRepository } from './repositories/user.repository';
import { UserService } from './services/user.service';
import { UserMapper } from './mappers/user.mapper';

@Module({
  imports: [TypeOrmModule.forFeature([UserEntity])],
  providers: [UserRepository, UserService, UserMapper],
  exports: [UserService],
})
export class UsersModule {}
```

### Custom Repository
```typescript
// src/users/repositories/user.repository.ts
import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository, FindOptionsWhere } from 'typeorm';
import { UserEntity } from '../entities/user.entity';

@Injectable()
export class UserRepository {
  constructor(
    @InjectRepository(UserEntity)
    private readonly repository: Repository<UserEntity>,
  ) {}

  /**
   * Find user by ID
   */
  async findById(id: string): Promise<UserEntity | null> {
    return this.repository.findOne({ where: { id } });
  }

  /**
   * Find user by email
   */
  async findByEmail(email: string): Promise<UserEntity | null> {
    return this.repository.findOne({ where: { email } });
  }

  /**
   * Create and save new user
   */
  async create(data: Partial<UserEntity>): Promise<UserEntity> {
    const user = this.repository.create(data);
    return this.repository.save(user);
  }

  /**
   * Update existing user
   */
  async update(id: string, data: Partial<UserEntity>): Promise<UserEntity> {
    await this.repository.update(id, data);
    const updated = await this.findById(id);
    if (!updated) {
      throw new Error('User not found after update');
    }
    return updated;
  }

  /**
   * Delete user
   */
  async delete(id: string): Promise<void> {
    await this.repository.delete(id);
  }

  /**
   * Find users with pagination
   */
  async findWithPagination(
    page: number,
    limit: number,
    filters?: { role?: string; isVerified?: boolean },
  ): Promise<{ data: UserEntity[]; total: number }> {
    const query = this.repository.createQueryBuilder('user');

    if (filters?.role) {
      query.andWhere('user.role = :role', { role: filters.role });
    }

    if (filters?.isVerified !== undefined) {
      query.andWhere('user.isVerified = :isVerified', {
        isVerified: filters.isVerified,
      });
    }

    query.skip((page - 1) * limit).take(limit);

    const [data, total] = await query.getManyAndCount();

    return { data, total };
  }

  /**
   * Find users by IDs (batch operation)
   */
  async findByIds(ids: string[]): Promise<UserEntity[]> {
    return this.repository
      .createQueryBuilder('user')
      .whereInIds(ids)
      .getMany();
  }

  /**
   * Check if email exists
   */
  async existsByEmail(email: string): Promise<boolean> {
    const count = await this.repository.count({ where: { email } });
    return count > 0;
  }
}
```

### Repository with Relations
```typescript
// src/orders/repositories/order.repository.ts
import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { OrderEntity } from '../entities/order.entity';

@Injectable()
export class OrderRepository {
  constructor(
    @InjectRepository(OrderEntity)
    private readonly repository: Repository<OrderEntity>,
  ) {}

  /**
   * Find order with user relation
   */
  async findByIdWithUser(id: string): Promise<OrderEntity | null> {
    return this.repository.findOne({
      where: { id },
      relations: ['user'],
    });
  }

  /**
   * Find order with all relations
   */
  async findByIdWithRelations(id: string): Promise<OrderEntity | null> {
    return this.repository.findOne({
      where: { id },
      relations: ['user', 'items'],
    });
  }

  /**
   * Find user orders with query builder (avoid N+1)
   */
  async findUserOrders(userId: string): Promise<OrderEntity[]> {
    return this.repository
      .createQueryBuilder('order')
      .leftJoinAndSelect('order.items', 'item')
      .where('order.userId = :userId', { userId })
      .orderBy('order.orderedAt', 'DESC')
      .getMany();
  }

  /**
   * Find orders with filtering and sorting
   */
  async findOrders(filters: {
    userId?: string;
    status?: string;
    fromDate?: Date;
    toDate?: Date;
  }): Promise<OrderEntity[]> {
    const query = this.repository.createQueryBuilder('order');

    if (filters.userId) {
      query.andWhere('order.userId = :userId', { userId: filters.userId });
    }

    if (filters.status) {
      query.andWhere('order.status = :status', { status: filters.status });
    }

    if (filters.fromDate) {
      query.andWhere('order.orderedAt >= :fromDate', {
        fromDate: filters.fromDate,
      });
    }

    if (filters.toDate) {
      query.andWhere('order.orderedAt <= :toDate', {
        toDate: filters.toDate,
      });
    }

    return query
      .orderBy('order.orderedAt', 'DESC')
      .getMany();
  }

  /**
   * Create order with items (use in transaction)
   */
  async create(data: Partial<OrderEntity>): Promise<OrderEntity> {
    const order = this.repository.create(data);
    return this.repository.save(order);
  }

  /**
   * Update order status
   */
  async updateStatus(
    id: string,
    status: 'pending' | 'paid' | 'shipped',
  ): Promise<void> {
    await this.repository.update(id, { status });
  }
}
```

## Transaction Management with @Transactional

### Service with @Transactional
```typescript
// src/orders/services/order.service.ts
import { Injectable, NotFoundException } from '@nestjs/common';
import { Transactional, Propagation, IsolationLevel } from 'typeorm-transactional';
import { OrderRepository } from '../repositories/order.repository';
import { OrderItemRepository } from '../repositories/order-item.repository';
import { UserRepository } from '../../users/repositories/user.repository';
import { OrderMapper } from '../mappers/order.mapper';
import { CreateOrderDto } from '../dto/create-order.dto';
import { OrderDto } from '../dto/order.dto';

@Injectable()
export class OrderService {
  constructor(
    private readonly orderRepository: OrderRepository,
    private readonly orderItemRepository: OrderItemRepository,
    private readonly userRepository: UserRepository,
    private readonly orderMapper: OrderMapper,
  ) {}

  /**
   * Create order with items in a transaction
   * REQUIRED: Creates new transaction
   */
  @Transactional({
    propagation: Propagation.REQUIRED,
    isolationLevel: IsolationLevel.READ_COMMITTED,
  })
  async createOrder(dto: CreateOrderDto): Promise<OrderDto> {
    // Verify user exists
    const user = await this.userRepository.findById(dto.userId);
    if (!user) {
      throw new NotFoundException('User not found');
    }

    // Calculate total
    const totalAmount = dto.items.reduce(
      (sum, item) => sum + item.price * item.quantity,
      0,
    );

    // Create order
    const order = await this.orderRepository.create({
      userId: dto.userId,
      totalAmount,
      status: 'pending',
      orderedAt: new Date(),
    });

    // Create order items
    for (const itemDto of dto.items) {
      await this.orderItemRepository.create({
        orderId: order.id,
        productId: itemDto.productId,
        quantity: itemDto.quantity,
        price: itemDto.price,
      });
    }

    // Fetch complete order with relations
    const savedOrder = await this.orderRepository.findByIdWithRelations(
      order.id,
    );
    if (!savedOrder) {
      throw new Error('Order not found after creation');
    }

    return this.orderMapper.toDto(savedOrder);
  }

  /**
   * Update order status with transaction
   */
  @Transactional()
  async updateOrderStatus(
    orderId: string,
    status: 'pending' | 'paid' | 'shipped',
  ): Promise<OrderDto> {
    const order = await this.orderRepository.findByIdWithRelations(orderId);
    if (!order) {
      throw new NotFoundException('Order not found');
    }

    await this.orderRepository.updateStatus(orderId, status);

    // Additional business logic based on status
    if (status === 'paid') {
      // Send confirmation email, etc.
      await this.handlePaymentConfirmation(orderId);
    }

    const updated = await this.orderRepository.findByIdWithRelations(orderId);
    if (!updated) {
      throw new Error('Order not found after update');
    }

    return this.orderMapper.toDto(updated);
  }

  /**
   * REQUIRES_NEW: Always creates new transaction
   */
  @Transactional({ propagation: Propagation.REQUIRES_NEW })
  async handlePaymentConfirmation(orderId: string): Promise<void> {
    // This runs in a separate transaction
    // Even if outer transaction fails, this commit is independent
    // Use for logging, auditing, notifications
  }

  /**
   * Get order (read-only, no transaction needed)
   */
  async getOrder(orderId: string): Promise<OrderDto> {
    const order = await this.orderRepository.findByIdWithRelations(orderId);
    if (!order) {
      throw new NotFoundException('Order not found');
    }
    return this.orderMapper.toDto(order);
  }

  /**
   * Cancel order and refund (complex transaction)
   */
  @Transactional()
  async cancelOrder(orderId: string): Promise<void> {
    const order = await this.orderRepository.findByIdWithRelations(orderId);
    if (!order) {
      throw new NotFoundException('Order not found');
    }

    if (order.status !== 'pending') {
      throw new Error('Only pending orders can be cancelled');
    }

    // Update status
    await this.orderRepository.updateStatus(orderId, 'cancelled');

    // Restore inventory
    for (const item of order.items) {
      await this.restoreInventory(item.productId, item.quantity);
    }

    // Process refund if payment was made
    if (order.status === 'paid') {
      await this.processRefund(orderId);
    }
  }

  private async restoreInventory(
    productId: string,
    quantity: number,
  ): Promise<void> {
    // Implementation details
  }

  private async processRefund(orderId: string): Promise<void> {
    // Implementation details
  }
}
```

### Transaction Propagation Modes
```typescript
// src/examples/transaction-propagation.service.ts
import { Injectable } from '@nestjs/common';
import { Transactional, Propagation } from 'typeorm-transactional';

@Injectable()
export class TransactionExampleService {
  /**
   * REQUIRED (default): Use existing transaction or create new one
   */
  @Transactional({ propagation: Propagation.REQUIRED })
  async methodRequired(): Promise<void> {
    // Most common use case
  }

  /**
   * REQUIRES_NEW: Always create new transaction
   * Useful for logging, auditing that should succeed even if main transaction fails
   */
  @Transactional({ propagation: Propagation.REQUIRES_NEW })
  async methodRequiresNew(): Promise<void> {
    // Independent transaction
  }

  /**
   * SUPPORTS: Use existing transaction if available, otherwise non-transactional
   */
  @Transactional({ propagation: Propagation.SUPPORTS })
  async methodSupports(): Promise<void> {
    // Flexible approach
  }

  /**
   * MANDATORY: Must be called within existing transaction, throw error otherwise
   */
  @Transactional({ propagation: Propagation.MANDATORY })
  async methodMandatory(): Promise<void> {
    // Enforces transaction usage
  }

  /**
   * NEVER: Must NOT be in transaction, throw error if transaction exists
   */
  @Transactional({ propagation: Propagation.NEVER })
  async methodNever(): Promise<void> {
    // Explicitly non-transactional
  }
}
```

### Transaction Isolation Levels
```typescript
// src/examples/transaction-isolation.service.ts
import { Injectable } from '@nestjs/common';
import { Transactional, IsolationLevel } from 'typeorm-transactional';

@Injectable()
export class IsolationExampleService {
  /**
   * READ_UNCOMMITTED: Lowest isolation, allows dirty reads
   */
  @Transactional({ isolationLevel: IsolationLevel.READ_UNCOMMITTED })
  async readUncommitted(): Promise<void> {
    // Can read uncommitted changes from other transactions
  }

  /**
   * READ_COMMITTED: Prevents dirty reads (default for most DBs)
   */
  @Transactional({ isolationLevel: IsolationLevel.READ_COMMITTED })
  async readCommitted(): Promise<void> {
    // Only reads committed data
  }

  /**
   * REPEATABLE_READ: Prevents dirty and non-repeatable reads
   */
  @Transactional({ isolationLevel: IsolationLevel.REPEATABLE_READ })
  async repeatableRead(): Promise<void> {
    // Same query returns same results within transaction
  }

  /**
   * SERIALIZABLE: Highest isolation, prevents all anomalies
   */
  @Transactional({ isolationLevel: IsolationLevel.SERIALIZABLE })
  async serializable(): Promise<void> {
    // Complete isolation, potential for deadlocks
  }
}
```

## Testing TypeORM Services

### Unit Testing with @Transactional
```typescript
// src/orders/services/__tests__/order.service.spec.ts
import { Test, TestingModule } from '@nestjs/testing';
import { OrderService } from '../order.service';
import { OrderRepository } from '../../repositories/order.repository';
import { OrderItemRepository } from '../../repositories/order-item.repository';
import { UserRepository } from '../../../users/repositories/user.repository';
import { OrderMapper } from '../../mappers/order.mapper';
import { NotFoundException } from '@nestjs/common';

describe('OrderService', () => {
  let service: OrderService;
  let orderRepository: jest.Mocked<OrderRepository>;
  let orderItemRepository: jest.Mocked<OrderItemRepository>;
  let userRepository: jest.Mocked<UserRepository>;
  let orderMapper: jest.Mocked<OrderMapper>;

  beforeEach(async () => {
    // Create mocks
    const mockOrderRepository = {
      create: jest.fn(),
      findByIdWithRelations: jest.fn(),
      updateStatus: jest.fn(),
    };

    const mockOrderItemRepository = {
      create: jest.fn(),
    };

    const mockUserRepository = {
      findById: jest.fn(),
    };

    const mockOrderMapper = {
      toDto: jest.fn(),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        OrderService,
        {
          provide: OrderRepository,
          useValue: mockOrderRepository,
        },
        {
          provide: OrderItemRepository,
          useValue: mockOrderItemRepository,
        },
        {
          provide: UserRepository,
          useValue: mockUserRepository,
        },
        {
          provide: OrderMapper,
          useValue: mockOrderMapper,
        },
      ],
    }).compile();

    service = module.get<OrderService>(OrderService);
    orderRepository = module.get(OrderRepository);
    orderItemRepository = module.get(OrderItemRepository);
    userRepository = module.get(UserRepository);
    orderMapper = module.get(OrderMapper);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('createOrder', () => {
    it('should create order with items successfully', async () => {
      // Arrange
      const createOrderDto = {
        userId: 'user-123',
        items: [
          { productId: 'prod-1', quantity: 2, price: 10 },
          { productId: 'prod-2', quantity: 1, price: 20 },
        ],
      };

      const mockUser = { id: 'user-123', email: 'test@test.com' };
      const mockOrder = {
        id: 'order-123',
        userId: 'user-123',
        totalAmount: 40,
        status: 'pending',
      };
      const mockOrderDto = { ...mockOrder };

      userRepository.findById.mockResolvedValue(mockUser as any);
      orderRepository.create.mockResolvedValue(mockOrder as any);
      orderItemRepository.create.mockResolvedValue({} as any);
      orderRepository.findByIdWithRelations.mockResolvedValue(mockOrder as any);
      orderMapper.toDto.mockReturnValue(mockOrderDto as any);

      // Act
      const result = await service.createOrder(createOrderDto);

      // Assert
      expect(userRepository.findById).toHaveBeenCalledWith('user-123');
      expect(orderRepository.create).toHaveBeenCalledWith({
        userId: 'user-123',
        totalAmount: 40,
        status: 'pending',
        orderedAt: expect.any(Date),
      });
      expect(orderItemRepository.create).toHaveBeenCalledTimes(2);
      expect(result).toEqual(mockOrderDto);
    });

    it('should throw NotFoundException if user not found', async () => {
      // Arrange
      const createOrderDto = {
        userId: 'user-123',
        items: [{ productId: 'prod-1', quantity: 1, price: 10 }],
      };

      userRepository.findById.mockResolvedValue(null);

      // Act & Assert
      await expect(service.createOrder(createOrderDto)).rejects.toThrow(
        NotFoundException,
      );
      expect(orderRepository.create).not.toHaveBeenCalled();
    });

    it('should rollback transaction on error', async () => {
      // Arrange
      const createOrderDto = {
        userId: 'user-123',
        items: [{ productId: 'prod-1', quantity: 1, price: 10 }],
      };

      userRepository.findById.mockResolvedValue({ id: 'user-123' } as any);
      orderRepository.create.mockResolvedValue({ id: 'order-123' } as any);
      orderItemRepository.create.mockRejectedValue(new Error('Database error'));

      // Act & Assert
      await expect(service.createOrder(createOrderDto)).rejects.toThrow(
        'Database error',
      );

      // In real scenario with @Transactional, the transaction would rollback
      // In unit tests, we just verify the error propagates
    });
  });

  describe('updateOrderStatus', () => {
    it('should update order status successfully', async () => {
      // Arrange
      const orderId = 'order-123';
      const mockOrder = {
        id: orderId,
        status: 'pending',
        items: [],
      };
      const mockUpdatedOrder = { ...mockOrder, status: 'paid' };

      orderRepository.findByIdWithRelations
        .mockResolvedValueOnce(mockOrder as any)
        .mockResolvedValueOnce(mockUpdatedOrder as any);
      orderRepository.updateStatus.mockResolvedValue(undefined);
      orderMapper.toDto.mockReturnValue(mockUpdatedOrder as any);

      // Mock the private method (if needed for testing)
      jest.spyOn(service as any, 'handlePaymentConfirmation')
        .mockResolvedValue(undefined);

      // Act
      const result = await service.updateOrderStatus(orderId, 'paid');

      // Assert
      expect(orderRepository.updateStatus).toHaveBeenCalledWith(
        orderId,
        'paid',
      );
      expect(result.status).toBe('paid');
    });

    it('should throw NotFoundException if order not found', async () => {
      // Arrange
      orderRepository.findByIdWithRelations.mockResolvedValue(null);

      // Act & Assert
      await expect(
        service.updateOrderStatus('invalid-id', 'paid'),
      ).rejects.toThrow(NotFoundException);
    });
  });
});
```

### Repository Unit Testing
```typescript
// src/users/repositories/__tests__/user.repository.spec.ts
import { Test, TestingModule } from '@nestjs/testing';
import { getRepositoryToken } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { UserRepository } from '../user.repository';
import { UserEntity } from '../../entities/user.entity';

describe('UserRepository', () => {
  let userRepository: UserRepository;
  let mockRepository: jest.Mocked<Repository<UserEntity>>;

  beforeEach(async () => {
    // Create mock TypeORM repository
    const mockTypeOrmRepository = {
      findOne: jest.fn(),
      create: jest.fn(),
      save: jest.fn(),
      update: jest.fn(),
      delete: jest.fn(),
      count: jest.fn(),
      createQueryBuilder: jest.fn(() => ({
        whereInIds: jest.fn().mockReturnThis(),
        andWhere: jest.fn().mockReturnThis(),
        skip: jest.fn().mockReturnThis(),
        take: jest.fn().mockReturnThis(),
        getMany: jest.fn(),
        getManyAndCount: jest.fn(),
      })),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        UserRepository,
        {
          provide: getRepositoryToken(UserEntity),
          useValue: mockTypeOrmRepository,
        },
      ],
    }).compile();

    userRepository = module.get<UserRepository>(UserRepository);
    mockRepository = module.get(getRepositoryToken(UserEntity));
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('findById', () => {
    it('should find user by id', async () => {
      // Arrange
      const userId = 'user-123';
      const mockUser = {
        id: userId,
        email: 'test@test.com',
        name: 'Test User',
      };
      mockRepository.findOne.mockResolvedValue(mockUser as UserEntity);

      // Act
      const result = await userRepository.findById(userId);

      // Assert
      expect(mockRepository.findOne).toHaveBeenCalledWith({
        where: { id: userId },
      });
      expect(result).toEqual(mockUser);
    });

    it('should return null if user not found', async () => {
      // Arrange
      mockRepository.findOne.mockResolvedValue(null);

      // Act
      const result = await userRepository.findById('invalid-id');

      // Assert
      expect(result).toBeNull();
    });
  });

  describe('create', () => {
    it('should create and save user', async () => {
      // Arrange
      const userData = { email: 'test@test.com', name: 'Test User' };
      const savedUser = { id: 'user-123', ...userData };

      mockRepository.create.mockReturnValue(userData as UserEntity);
      mockRepository.save.mockResolvedValue(savedUser as UserEntity);

      // Act
      const result = await userRepository.create(userData);

      // Assert
      expect(mockRepository.create).toHaveBeenCalledWith(userData);
      expect(mockRepository.save).toHaveBeenCalledWith(userData);
      expect(result).toEqual(savedUser);
    });
  });

  describe('findWithPagination', () => {
    it('should return paginated users', async () => {
      // Arrange
      const mockUsers = [
        { id: '1', email: 'user1@test.com', name: 'User 1' },
        { id: '2', email: 'user2@test.com', name: 'User 2' },
      ];
      const total = 10;

      const queryBuilder = mockRepository.createQueryBuilder();
      (queryBuilder.getManyAndCount as jest.Mock).mockResolvedValue([
        mockUsers,
        total,
      ]);

      // Act
      const result = await userRepository.findWithPagination(1, 10);

      // Assert
      expect(result.data).toEqual(mockUsers);
      expect(result.total).toBe(total);
      expect(queryBuilder.skip).toHaveBeenCalledWith(0);
      expect(queryBuilder.take).toHaveBeenCalledWith(10);
    });

    it('should apply filters', async () => {
      // Arrange
      const queryBuilder = mockRepository.createQueryBuilder();
      (queryBuilder.getManyAndCount as jest.Mock).mockResolvedValue([[], 0]);

      // Act
      await userRepository.findWithPagination(1, 10, {
        role: 'admin',
        isVerified: true,
      });

      // Assert
      expect(queryBuilder.andWhere).toHaveBeenCalledWith(
        'user.role = :role',
        { role: 'admin' },
      );
      expect(queryBuilder.andWhere).toHaveBeenCalledWith(
        'user.isVerified = :isVerified',
        { isVerified: true },
      );
    });
  });

  describe('existsByEmail', () => {
    it('should return true if email exists', async () => {
      // Arrange
      mockRepository.count.mockResolvedValue(1);

      // Act
      const result = await userRepository.existsByEmail('test@test.com');

      // Assert
      expect(result).toBe(true);
      expect(mockRepository.count).toHaveBeenCalledWith({
        where: { email: 'test@test.com' },
      });
    });

    it('should return false if email does not exist', async () => {
      // Arrange
      mockRepository.count.mockResolvedValue(0);

      // Act
      const result = await userRepository.existsByEmail('test@test.com');

      // Assert
      expect(result).toBe(false);
    });
  });
});
```

### Integration Testing with Real Database
```typescript
// src/orders/services/__tests__/order.service.integration.spec.ts
import { Test, TestingModule } from '@nestjs/testing';
import { TypeOrmModule } from '@nestjs/typeorm';
import { DataSource } from 'typeorm';
import { OrderService } from '../order.service';
import { OrderRepository } from '../../repositories/order.repository';
import { OrderItemRepository } from '../../repositories/order-item.repository';
import { UserRepository } from '../../../users/repositories/user.repository';
import { OrderMapper } from '../../mappers/order.mapper';
import { UserMapper } from '../../../users/mappers/user.mapper';
import { OrderEntity } from '../../entities/order.entity';
import { OrderItemEntity } from '../../entities/order-item.entity';
import { UserEntity } from '../../../users/entities/user.entity';
import { initializeTransactionalContext } from 'typeorm-transactional';

describe('OrderService Integration', () => {
  let module: TestingModule;
  let service: OrderService;
  let dataSource: DataSource;
  let userRepository: UserRepository;

  beforeAll(async () => {
    // Initialize transaction context
    initializeTransactionalContext();

    module = await Test.createTestingModule({
      imports: [
        TypeOrmModule.forRoot({
          type: 'postgres',
          host: 'localhost',
          port: 5433,
          username: 'test',
          password: 'test',
          database: 'test_db',
          entities: [UserEntity, OrderEntity, OrderItemEntity],
          synchronize: true, // OK for tests
          dropSchema: true, // Clean slate for each test run
        }),
        TypeOrmModule.forFeature([UserEntity, OrderEntity, OrderItemEntity]),
      ],
      providers: [
        OrderService,
        OrderRepository,
        OrderItemRepository,
        UserRepository,
        OrderMapper,
        UserMapper,
      ],
    }).compile();

    service = module.get<OrderService>(OrderService);
    userRepository = module.get<UserRepository>(UserRepository);
    dataSource = module.get<DataSource>(DataSource);
  });

  afterAll(async () => {
    await dataSource.destroy();
    await module.close();
  });

  beforeEach(async () => {
    // Clean database before each test
    await dataSource.query('TRUNCATE TABLE orders CASCADE');
    await dataSource.query('TRUNCATE TABLE users CASCADE');
  });

  describe('createOrder', () => {
    it('should create order with transaction', async () => {
      // Arrange
      const user = await userRepository.create({
        email: 'test@test.com',
        name: 'Test User',
        role: 'user',
      });

      const createOrderDto = {
        userId: user.id,
        items: [
          { productId: 'prod-1', quantity: 2, price: 10 },
          { productId: 'prod-2', quantity: 1, price: 20 },
        ],
      };

      // Act
      const result = await service.createOrder(createOrderDto);

      // Assert
      expect(result.id).toBeDefined();
      expect(result.userId).toBe(user.id);
      expect(result.totalAmount).toBe(40);
      expect(result.items).toHaveLength(2);
    });

    it('should rollback transaction on error', async () => {
      // Arrange
      const createOrderDto = {
        userId: 'non-existent-user',
        items: [{ productId: 'prod-1', quantity: 1, price: 10 }],
      };

      // Act & Assert
      await expect(service.createOrder(createOrderDto)).rejects.toThrow();

      // Verify no order was created
      const orders = await dataSource.query('SELECT * FROM orders');
      expect(orders).toHaveLength(0);
    });
  });
});
```

## Advanced Query Patterns

### Query Builder Examples
```typescript
// src/users/repositories/user.repository.advanced.ts
export class UserRepositoryAdvanced {
  /**
   * Complex search with multiple conditions
   */
  async searchUsers(criteria: {
    searchTerm?: string;
    roles?: string[];
    isVerified?: boolean;
    createdAfter?: Date;
  }): Promise<UserEntity[]> {
    const query = this.repository.createQueryBuilder('user');

    if (criteria.searchTerm) {
      query.andWhere(
        '(user.email ILIKE :search OR user.name ILIKE :search)',
        { search: `%${criteria.searchTerm}%` },
      );
    }

    if (criteria.roles && criteria.roles.length > 0) {
      query.andWhere('user.role IN (:...roles)', { roles: criteria.roles });
    }

    if (criteria.isVerified !== undefined) {
      query.andWhere('user.isVerified = :isVerified', {
        isVerified: criteria.isVerified,
      });
    }

    if (criteria.createdAfter) {
      query.andWhere('user.createdAt > :createdAfter', {
        createdAfter: criteria.createdAfter,
      });
    }

    return query
      .orderBy('user.createdAt', 'DESC')
      .getMany();
  }

  /**
   * Prevent N+1 with proper joins
   */
  async findUsersWithOrders(): Promise<UserEntity[]> {
    return this.repository
      .createQueryBuilder('user')
      .leftJoinAndSelect('user.orders', 'order')
      .leftJoinAndSelect('order.items', 'item')
      .where('order.status != :status', { status: 'cancelled' })
      .orderBy('user.createdAt', 'DESC')
      .addOrderBy('order.orderedAt', 'DESC')
      .getMany();
  }

  /**
   * Aggregate queries
   */
  async getUserStatistics(): Promise<{
    totalUsers: number;
    verifiedUsers: number;
    adminCount: number;
    avgOrdersPerUser: number;
  }> {
    const result = await this.repository
      .createQueryBuilder('user')
      .select('COUNT(DISTINCT user.id)', 'totalUsers')
      .addSelect(
        'COUNT(DISTINCT CASE WHEN user.isVerified = true THEN user.id END)',
        'verifiedUsers',
      )
      .addSelect(
        'COUNT(DISTINCT CASE WHEN user.role = :role THEN user.id END)',
        'adminCount',
      )
      .setParameter('role', 'admin')
      .getRawOne();

    return {
      totalUsers: parseInt(result.totalUsers),
      verifiedUsers: parseInt(result.verifiedUsers),
      adminCount: parseInt(result.adminCount),
      avgOrdersPerUser: 0, // Additional query needed
    };
  }

  /**
   * Raw SQL for complex queries
   */
  async findUsersByComplexCriteria(): Promise<UserEntity[]> {
    return this.repository.query(
      `
      SELECT u.*
      FROM users u
      WHERE u.id IN (
        SELECT DISTINCT o.user_id
        FROM orders o
        WHERE o.total_amount > $1
        GROUP BY o.user_id
        HAVING COUNT(*) > $2
      )
      ORDER BY u.created_at DESC
      `,
      [1000, 5],
    );
  }
}
```

### Optimistic Locking
```typescript
// src/products/entities/product.entity.ts
import { Entity, Column, PrimaryGeneratedColumn, VersionColumn } from 'typeorm';

@Entity('products')
export class ProductEntity {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ type: 'varchar', length: 255 })
  name: string;

  @Column({ type: 'integer' })
  stock: number;

  @VersionColumn()
  version: number; // Automatically managed by TypeORM
}

// src/products/services/product.service.ts
@Injectable()
export class ProductService {
  @Transactional()
  async updateStock(productId: string, quantity: number): Promise<void> {
    const product = await this.productRepository.findById(productId);
    if (!product) {
      throw new NotFoundException('Product not found');
    }

    if (product.stock < quantity) {
      throw new Error('Insufficient stock');
    }

    product.stock -= quantity;

    try {
      await this.productRepository.save(product);
    } catch (error) {
      // Handle optimistic locking failure
      if (error.message.includes('version')) {
        throw new Error('Product was modified by another process');
      }
      throw error;
    }
  }
}
```

## Best Practices

### 1. Always Use Repository Pattern (Data Mapper)
```typescript
// ✅ GOOD: Repository pattern
@Injectable()
export class UserService {
  constructor(private readonly userRepository: UserRepository) {}

  async createUser(dto: CreateUserDto): Promise<UserDto> {
    const entity = await this.userRepository.create(dto);
    return this.userMapper.toDto(entity);
  }
}

// ❌ BAD: Direct entity manipulation (Active Record)
const user = new UserEntity();
user.email = 'test@test.com';
await user.save(); // Avoid this pattern
```

### 2. Use @Transactional for Write Operations
```typescript
// ✅ GOOD: Explicit transaction management
@Transactional()
async createOrder(dto: CreateOrderDto): Promise<OrderDto> {
  const order = await this.orderRepository.create(dto);
  await this.orderItemRepository.createMany(order.id, dto.items);
  return this.orderMapper.toDto(order);
}

// ❌ BAD: No transaction for multi-step operations
async createOrder(dto: CreateOrderDto): Promise<OrderDto> {
  const order = await this.orderRepository.create(dto);
  // If this fails, order already created (data inconsistency)
  await this.orderItemRepository.createMany(order.id, dto.items);
}
```

### 3. Prevent N+1 Queries
```typescript
// ✅ GOOD: Eager loading with proper joins
async findOrders(): Promise<OrderEntity[]> {
  return this.repository
    .createQueryBuilder('order')
    .leftJoinAndSelect('order.items', 'item')
    .leftJoinAndSelect('order.user', 'user')
    .getMany();
}

// ❌ BAD: Lazy loading causing N+1
async findOrders(): Promise<OrderEntity[]> {
  const orders = await this.repository.find(); // 1 query
  for (const order of orders) {
    order.items = await this.itemRepository.findByOrderId(order.id); // N queries
  }
  return orders;
}
```

### 4. Always Map Entities to DTOs
```typescript
// ✅ GOOD: Return DTOs, not entities
async getUser(id: string): Promise<UserDto> {
  const entity = await this.userRepository.findById(id);
  if (!entity) throw new NotFoundException();
  return this.userMapper.toDto(entity); // Excludes sensitive fields
}

// ❌ BAD: Exposing entities directly
async getUser(id: string): Promise<UserEntity> {
  return this.userRepository.findById(id); // Exposes passwordHash, etc.
}
```

### 5. Use Pagination for Large Datasets
```typescript
// ✅ GOOD: Pagination with limits
async findUsers(page: number, limit: number): Promise<PaginatedResult<UserDto>> {
  const { data, total } = await this.userRepository.findWithPagination(page, limit);
  return {
    data: this.userMapper.toDtos(data),
    total,
    page,
    pageCount: Math.ceil(total / limit),
  };
}

// ❌ BAD: Loading all records
async findUsers(): Promise<UserDto[]> {
  const users = await this.userRepository.find(); // Could be millions of records
  return this.userMapper.toDtos(users);
}
```

### 6. Index Foreign Keys and Frequently Queried Columns
```typescript
@Entity('orders')
@Index(['userId', 'status']) // Composite index for common queries
export class OrderEntity {
  @Column({ type: 'uuid' })
  @Index() // Single column index
  userId: string;

  @Column({ type: 'varchar' })
  @Index() // Index for filtering
  status: string;
}
```

### 7. Use Connection Pooling
```typescript
export const getTypeOrmConfig = (): DataSourceOptions => ({
  type: 'postgres',
  // ...other options
  poolSize: 10, // Max connections
  extra: {
    max: 10,
    min: 2,
    idleTimeoutMillis: 30000,
    connectionTimeoutMillis: 2000,
  },
});
```

### 8. Handle Errors Properly
```typescript
@Transactional()
async updateUser(id: string, dto: UpdateUserDto): Promise<UserDto> {
  try {
    const user = await this.userRepository.findById(id);
    if (!user) {
      throw new NotFoundException('User not found');
    }

    const updated = this.userMapper.updateEntity(user, dto);
    const saved = await this.userRepository.save(updated);
    return this.userMapper.toDto(saved);
  } catch (error) {
    if (error.code === '23505') { // Postgres unique violation
      throw new ConflictException('Email already exists');
    }
    throw error;
  }
}
```

## Common Issues & Solutions

### Issue: Transaction Context Not Initialized
```typescript
// ❌ Problem: Error "No transaction context found"
// Solution: Initialize in main.ts BEFORE creating app

// src/main.ts
import { initializeTransactionalContext } from 'typeorm-transactional';

async function bootstrap() {
  initializeTransactionalContext(); // MUST be first!
  const app = await NestFactory.create(AppModule);
  await app.listen(3000);
}
```

### Issue: Repository Not Found
```typescript
// ❌ Problem: Error "Cannot find repository"
// Solution: Import TypeOrmModule.forFeature in module

@Module({
  imports: [
    TypeOrmModule.forFeature([UserEntity]), // Must include entity
  ],
  providers: [UserRepository, UserService],
})
export class UsersModule {}
```

### Issue: Circular Dependencies with Mappers
```typescript
// ❌ Problem: Circular dependency between mappers
// Solution: Use forwardRef or restructure dependencies

@Injectable()
export class OrderMapper {
  constructor(
    @Inject(forwardRef(() => UserMapper))
    private readonly userMapper: UserMapper,
  ) {}
}
```

### Issue: N+1 Query Problem
```typescript
// ❌ Problem: Multiple queries for related data
// Solution: Use joins or eager loading

// Instead of:
const orders = await this.orderRepository.find();
for (const order of orders) {
  order.user = await this.userRepository.findById(order.userId);
}

// Do:
const orders = await this.orderRepository
  .createQueryBuilder('order')
  .leftJoinAndSelect('order.user', 'user')
  .getMany();
```

### Issue: Transaction Not Rolling Back
```typescript
// ❌ Problem: Changes persist despite error
// Solution: Ensure @Transactional decorator is applied

// Check: Is initializeTransactionalContext() called in main.ts?
// Check: Is dataSource properly registered with addTransactionalDataSource()?
// Check: Are you using async/await properly?
```

### Issue: Slow Queries
```typescript
// ❌ Problem: Queries taking too long
// Solution: Add indexes, use query builder, enable query logging

// Add logging to find slow queries
TypeOrmModule.forRoot({
  logging: true,
  maxQueryExecutionTime: 1000, // Log queries > 1 second
});

// Add appropriate indexes
@Entity('orders')
@Index(['userId', 'createdAt']) // For common query patterns
export class OrderEntity {}
```

### Issue: Memory Leaks with Large Result Sets
```typescript
// ❌ Problem: Loading too much data into memory
// Solution: Use streaming or pagination

// Instead of:
const allOrders = await this.orderRepository.find();

// Do:
async *streamOrders(): AsyncGenerator<OrderEntity> {
  const stream = await this.repository
    .createQueryBuilder('order')
    .stream();

  for await (const order of stream) {
    yield order;
  }
}
```

## Documentation
- [TypeORM Official Docs](https://typeorm.io/)
- [NestJS TypeORM Integration](https://docs.nestjs.com/techniques/database)
- [typeorm-transactional GitHub](https://github.com/Aliheym/typeorm-transactional)
- [TypeORM Migrations](https://typeorm.io/migrations)
- [TypeORM Query Builder](https://typeorm.io/select-query-builder)

**Use for**: TypeORM Data Mapper pattern, repository implementation, transaction management with @Transactional, entity-DTO mapping, query optimization, N+1 prevention, unit and integration testing with TypeORM, production-ready database operations.
