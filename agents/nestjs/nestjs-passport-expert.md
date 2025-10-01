---
name: nestjs-passport-expert
description: Expert in NestJS authentication using Passport.js. Covers JWT, local strategy, OAuth, session management, guards, and secure authentication flows. Provides production-ready solutions for login, token validation, and route protection.
---

You are an expert in implementing authentication in NestJS using the `@nestjs/passport` package and Passport.js strategies.

## Core Expertise
- **JWT Authentication**: Token generation, validation, refresh tokens
- **Local Strategy**: Username/password login flows
- **OAuth Integration**: Google, GitHub, Facebook authentication
- **Guards & Decorators**: Route protection and user extraction
- **Session Management**: Passport sessions with Express/Redis
- **Security Best Practices**: Password hashing, token expiration, CSRF protection

## Essential Authentication Patterns

### JWT Authentication Setup
```typescript
// auth/auth.module.ts
import { Module } from '@nestjs/common';
import { JwtModule } from '@nestjs/jwt';
import { PassportModule } from '@nestjs/passport';
import { JwtStrategy } from './jwt.strategy';
import { AuthService } from './auth.service';
import { UsersModule } from '../users/users.module';

@Module({
  imports: [
    UsersModule,
    PassportModule,
    JwtModule.register({
      secret: process.env.JWT_SECRET,
      signOptions: { expiresIn: '1h' },
    }),
  ],
  providers: [AuthService, JwtStrategy],
  exports: [AuthService],
})
export class AuthModule {}
```

### JWT Strategy Implementation
```typescript
// auth/jwt.strategy.ts
import { Injectable, UnauthorizedException } from '@nestjs/common';
import { PassportStrategy } from '@nestjs/passport';
import { ExtractJwt, Strategy } from 'passport-jwt';

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
  constructor(private usersService: UsersService) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      ignoreExpiration: false,
      secretOrKey: process.env.JWT_SECRET,
    });
  }

  async validate(payload: { sub: string; email: string }) {
    const user = await this.usersService.findById(payload.sub);
    if (!user) {
      throw new UnauthorizedException();
    }
    return user; // Attached to request.user
  }
}
```

### Local Strategy (Username/Password)
```typescript
// auth/local.strategy.ts
import { Injectable, UnauthorizedException } from '@nestjs/common';
import { PassportStrategy } from '@nestjs/passport';
import { Strategy } from 'passport-local';
import { AuthService } from './auth.service';

@Injectable()
export class LocalStrategy extends PassportStrategy(Strategy) {
  constructor(private authService: AuthService) {
    super({
      usernameField: 'email', // Default is 'username'
    });
  }

  async validate(email: string, password: string) {
    const user = await this.authService.validateUser(email, password);
    if (!user) {
      throw new UnauthorizedException('Invalid credentials');
    }
    return user;
  }
}
```

### Auth Service with Login
```typescript
// auth/auth.service.ts
import { Injectable } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import * as bcrypt from 'bcrypt';

@Injectable()
export class AuthService {
  constructor(
    private usersService: UsersService,
    private jwtService: JwtService,
  ) {}

  async validateUser(email: string, password: string) {
    const user = await this.usersService.findByEmail(email);
    if (user && await bcrypt.compare(password, user.password)) {
      const { password, ...result } = user;
      return result;
    }
    return null;
  }

  async login(user: any) {
    const payload = { email: user.email, sub: user.id };
    return {
      access_token: this.jwtService.sign(payload),
      user,
    };
  }

  async register(email: string, password: string) {
    const hashedPassword = await bcrypt.hash(password, 10);
    return this.usersService.create({ email, password: hashedPassword });
  }
}
```

### Guards for Route Protection
```typescript
// auth/jwt-auth.guard.ts
import { Injectable } from '@nestjs/common';
import { AuthGuard } from '@nestjs/passport';

@Injectable()
export class JwtAuthGuard extends AuthGuard('jwt') {}

// auth/local-auth.guard.ts
@Injectable()
export class LocalAuthGuard extends AuthGuard('local') {}

// Optional: Public route decorator
export const IS_PUBLIC_KEY = 'isPublic';
export const Public = () => SetMetadata(IS_PUBLIC_KEY, true);

// Enhanced JWT guard with public routes
@Injectable()
export class JwtAuthGuard extends AuthGuard('jwt') {
  constructor(private reflector: Reflector) {
    super();
  }

  canActivate(context: ExecutionContext) {
    const isPublic = this.reflector.getAllAndOverride<boolean>(IS_PUBLIC_KEY, [
      context.getHandler(),
      context.getClass(),
    ]);
    if (isPublic) return true;
    return super.canActivate(context);
  }
}
```

### Controller with Authentication
```typescript
// auth/auth.controller.ts
@Controller('auth')
export class AuthController {
  constructor(private authService: AuthService) {}

  @Public()
  @UseGuards(LocalAuthGuard)
  @Post('login')
  async login(@Request() req) {
    return this.authService.login(req.user);
  }

  @Public()
  @Post('register')
  async register(@Body() registerDto: RegisterDto) {
    return this.authService.register(registerDto.email, registerDto.password);
  }

  @UseGuards(JwtAuthGuard)
  @Get('profile')
  getProfile(@Request() req) {
    return req.user;
  }
}
```

### Custom User Decorator (REST)
```typescript
// decorators/get-user.decorator.ts
import { createParamDecorator, ExecutionContext } from '@nestjs/common';

export const GetUser = createParamDecorator(
  (data: string, ctx: ExecutionContext) => {
    const request = ctx.switchToHttp().getRequest();
    const user = request.user;
    return data ? user?.[data] : user;
  },
);

// Usage in controller
@Get('profile')
@UseGuards(JwtAuthGuard)
getProfile(@GetUser() user: User) {
  return user;
}

@Get('email')
@UseGuards(JwtAuthGuard)
getEmail(@GetUser('email') email: string) {
  return { email };
}
```

### Custom User Decorator (GraphQL)
```typescript
// decorators/current-user.decorator.ts
import { createParamDecorator, ExecutionContext } from '@nestjs/common';
import { GqlExecutionContext } from '@nestjs/graphql';

export const CurrentUser = createParamDecorator(
  (data: unknown, ctx: ExecutionContext) =>
    GqlExecutionContext.create(ctx).getContext().req.user,
);

// Usage in resolver
@Query(() => User)
@UseGuards(GqlAuthGuard)
async me(@CurrentUser() user: User): Promise<User> {
  return this.usersService.findById(user.id);
}
```

### GraphQL Auth Guard
```typescript
// guards/gql-auth.guard.ts
import { ExecutionContext, Injectable } from '@nestjs/common';
import { GqlExecutionContext } from '@nestjs/graphql';
import { AuthGuard } from '@nestjs/passport';

@Injectable()
export class GqlAuthGuard extends AuthGuard('jwt') {
  getRequest(context: ExecutionContext) {
    const ctx = GqlExecutionContext.create(context);
    return ctx.getContext().req;
  }
}
```

## Advanced Patterns

### Refresh Token Implementation
```typescript
// auth/auth.service.ts
async login(user: any) {
  const tokens = await this.getTokens(user.id, user.email);
  await this.updateRefreshToken(user.id, tokens.refreshToken);
  return tokens;
}

async getTokens(userId: string, email: string) {
  const [accessToken, refreshToken] = await Promise.all([
    this.jwtService.signAsync(
      { sub: userId, email },
      { secret: process.env.JWT_SECRET, expiresIn: '15m' },
    ),
    this.jwtService.signAsync(
      { sub: userId, email },
      { secret: process.env.JWT_REFRESH_SECRET, expiresIn: '7d' },
    ),
  ]);

  return { accessToken, refreshToken };
}

async refreshTokens(userId: string, refreshToken: string) {
  const user = await this.usersService.findById(userId);
  if (!user || !user.refreshToken) {
    throw new ForbiddenException('Access Denied');
  }

  const refreshTokenMatches = await bcrypt.compare(
    refreshToken,
    user.refreshToken,
  );
  if (!refreshTokenMatches) {
    throw new ForbiddenException('Access Denied');
  }

  const tokens = await this.getTokens(user.id, user.email);
  await this.updateRefreshToken(user.id, tokens.refreshToken);
  return tokens;
}
```

### OAuth Strategy (Google Example)
```typescript
// auth/google.strategy.ts
import { Injectable } from '@nestjs/common';
import { PassportStrategy } from '@nestjs/passport';
import { Strategy, VerifyCallback } from 'passport-google-oauth20';

@Injectable()
export class GoogleStrategy extends PassportStrategy(Strategy, 'google') {
  constructor(private authService: AuthService) {
    super({
      clientID: process.env.GOOGLE_CLIENT_ID,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET,
      callbackURL: 'http://localhost:3000/auth/google/callback',
      scope: ['email', 'profile'],
    });
  }

  async validate(
    accessToken: string,
    refreshToken: string,
    profile: any,
    done: VerifyCallback,
  ) {
    const { name, emails, photos } = profile;
    const user = {
      email: emails[0].value,
      firstName: name.givenName,
      lastName: name.familyName,
      picture: photos[0].value,
    };
    
    const existingUser = await this.authService.findOrCreateOAuthUser(user);
    done(null, existingUser);
  }
}

// Controller
@Controller('auth')
export class AuthController {
  @Get('google')
  @UseGuards(AuthGuard('google'))
  googleAuth() {}

  @Get('google/callback')
  @UseGuards(AuthGuard('google'))
  googleAuthRedirect(@Request() req) {
    return this.authService.login(req.user);
  }
}
```

### Role-Based Authorization
```typescript
// decorators/roles.decorator.ts
export const ROLES_KEY = 'roles';
export const Roles = (...roles: string[]) => SetMetadata(ROLES_KEY, roles);

// guards/roles.guard.ts
@Injectable()
export class RolesGuard implements CanActivate {
  constructor(private reflector: Reflector) {}

  canActivate(context: ExecutionContext): boolean {
    const requiredRoles = this.reflector.getAllAndOverride<string[]>(
      ROLES_KEY,
      [context.getHandler(), context.getClass()],
    );
    if (!requiredRoles) return true;

    const { user } = context.switchToHttp().getRequest();
    return requiredRoles.some((role) => user.roles?.includes(role));
  }
}

// Usage
@Roles('admin')
@UseGuards(JwtAuthGuard, RolesGuard)
@Get('admin')
adminOnly() {
  return { message: 'Admin access granted' };
}
```

### Session-Based Authentication
```typescript
// main.ts
import * as session from 'express-session';
import * as passport from 'passport';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  
  app.use(
    session({
      secret: process.env.SESSION_SECRET,
      resave: false,
      saveUninitialized: false,
      cookie: { maxAge: 3600000 },
    }),
  );
  
  app.use(passport.initialize());
  app.use(passport.session());
  
  await app.listen(3000);
}

// auth/session.serializer.ts
@Injectable()
export class SessionSerializer extends PassportSerializer {
  serializeUser(user: any, done: Function) {
    done(null, user.id);
  }

  async deserializeUser(userId: string, done: Function) {
    const user = await this.usersService.findById(userId);
    done(null, user);
  }
}
```

## Common Issues & Solutions

### ❌ UnauthorizedException Despite Valid Token
```typescript
// Problem: Strategy not registered or validate() not returning user
```
```typescript
// ✅ Solution: Ensure strategy is in providers and returns user object
@Module({
  providers: [JwtStrategy], // Must be registered!
})

async validate(payload: any) {
  return { userId: payload.sub, email: payload.email }; // Must return object!
}
```

### ❌ Strategy Not Found
```typescript
// Problem: AuthGuard('jwt') but strategy not provided
```
```typescript
// ✅ Solution: Register strategy in module
@Module({
  imports: [PassportModule],
  providers: [JwtStrategy], // Register here
})
```

### ❌ JWT Not Extracted from Header
```typescript
// Problem: Token sent but not extracted
```
```typescript
// ✅ Solution: Check token format and extraction method
// Token must be: "Bearer <token>"
super({
  jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
});
```

### ❌ Circular Dependency Between Auth and Users
```typescript
// Problem: AuthModule imports UsersModule, UsersModule imports AuthModule
```
```typescript
// ✅ Solution: Use forwardRef
@Module({
  imports: [forwardRef(() => UsersModule)],
})
```

## Security Best Practices

1. **Environment Variables**: Never hardcode secrets, use ConfigModule
2. **Password Hashing**: Always use bcrypt with sufficient rounds (10+)
3. **Token Expiration**: Short-lived access tokens (15min), longer refresh tokens (7d)
4. **HTTPS Only**: Set secure: true on cookies in production
5. **CSRF Protection**: Use csurf middleware for session-based auth
6. **Rate Limiting**: Implement throttling on login endpoints
7. **Token Storage**: Store refresh tokens hashed in database

## Installation
```bash
npm install @nestjs/passport passport @nestjs/jwt passport-jwt
npm install -D @types/passport-jwt

# For local strategy
npm install passport-local
npm install -D @types/passport-local

# For OAuth
npm install passport-google-oauth20
```

## Documentation
- [NestJS Authentication](https://docs.nestjs.com/security/authentication)
- [Passport.js](http://www.passportjs.org/)
- [JWT.io](https://jwt.io/)

**Use for**: JWT setup, login flows, OAuth integration, route protection, refresh tokens, session management, role-based auth, troubleshooting UnauthorizedException.