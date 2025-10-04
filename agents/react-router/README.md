# React Router v7 Expert Agents

Specialized subagents for React Router v7, designed for use with Claude Code. Each agent focuses on specific React Router domains and provides production-ready, type-safe solutions.

## Available Agents

### Core Routing
- **react-router-core-expert**: Routing modes, route configuration, nested routes, URL parameters, navigation, and file-based routing conventions

### Data Management
- **react-router-data-loading-expert**: Loader functions, useLoaderData, defer, streaming, revalidation, prefetching, and data fetching strategies
- **react-router-forms-actions-expert**: Form component, action functions, validation, progressive enhancement, optimistic UI, and useFetcher

### Developer Tools
- **react-router-hooks-expert**: All routing hooks including useLoaderData, useActionData, useNavigation, useParams, useLocation, useSearchParams, useFetcher, and navigation state management
- **react-router-error-handling-expert**: Error boundaries, route error handling, useRouteError, error responses, and recovery patterns

### Framework Features
- **react-router-framework-expert**: Framework Mode setup, file-based routing, SSR, build configuration, deployment, meta tags, and route modules
- **react-router-testing-expert**: Testing loaders, actions, components with routing context, integration tests, and mocking strategies

## Agent Selection Guidelines

**When building new applications:**
- Use `react-router-framework-expert` for Framework Mode setup and configuration
- Use `react-router-core-expert` for routing architecture and file structure

**When implementing data features:**
- Use `react-router-data-loading-expert` for data fetching, loaders, and prefetching
- Use `react-router-forms-actions-expert` for forms, mutations, and validation

**When working with routing state:**
- Use `react-router-hooks-expert` for accessing params, location, navigation state
- Use `react-router-error-handling-expert` for error boundaries and error recovery

**When writing tests:**
- Use `react-router-testing-expert` for loader/action tests, component tests, and integration tests

## Usage with Claude Code

Add these agents to your Claude Code workflow:

```markdown
Before starting work, select appropriate React Router agent(s):

- react-router-core-expert: Route setup and navigation
- react-router-data-loading-expert: Data fetching and loaders
- react-router-forms-actions-expert: Forms and mutations
- react-router-hooks-expert: Hook usage and routing state
- react-router-error-handling-expert: Error boundaries
- react-router-framework-expert: Framework setup and SSR
- react-router-testing-expert: Testing
```

## Key Features

### Production-Ready Patterns
- Type-safe implementations with TypeScript
- Best practices and common patterns
- Error handling and edge cases
- Performance optimization

### Framework Mode Focus
- File-based routing conventions
- Server-side rendering (SSR)
- Progressive enhancement
- Meta tags and SEO

### Comprehensive Coverage
- All three routing modes (Framework, Data, Declarative)
- Complete hook reference
- Testing strategies
- Deployment patterns

## Examples

### Setting up a new route with data loading
```typescript
// Use: react-router-core-expert + react-router-data-loading-expert

// app/routes/users.$id.tsx
export async function loader({ params }: Route.LoaderArgs) {
  const user = await fetchUser(params.id);
  if (!user) {
    throw new Response("Not found", { status: 404 });
  }
  return { user };
}

export default function User() {
  const { user } = useLoaderData<typeof loader>();
  return <h1>{user.name}</h1>;
}
```

### Creating a form with validation
```typescript
// Use: react-router-forms-actions-expert

export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  const errors = validate(formData);

  if (errors) {
    return { errors };
  }

  await createUser(formData);
  return redirect("/users");
}

export default function NewUser() {
  const actionData = useActionData<typeof action>();

  return (
    <Form method="post">
      <input name="email" />
      {actionData?.errors?.email && <p>{actionData.errors.email}</p>}
      <button type="submit">Create</button>
    </Form>
  );
}
```

### Testing a route
```typescript
// Use: react-router-testing-expert

import { createMemoryRouter, RouterProvider } from "react-router";

it("displays user from loader", async () => {
  const router = createMemoryRouter([
    {
      path: "/users/:id",
      loader: () => ({ user: { name: "John" } }),
      Component: User,
    },
  ], {
    initialEntries: ["/users/123"],
  });

  render(<RouterProvider router={router} />);

  expect(await screen.findByText("John")).toBeInTheDocument();
});
```

## Documentation References

- [React Router v7 Official Docs](https://reactrouter.com)
- [Framework Mode Guide](https://reactrouter.com/start/framework)
- [Data Loading](https://reactrouter.com/start/framework/data-loading)
- [Actions & Forms](https://reactrouter.com/start/framework/actions)
- [Hooks Reference](https://reactrouter.com/hooks)

## Design Principles

1. **Precision**: Concise, focused expertise for Claude Code efficiency
2. **Type Safety**: TypeScript-first with proper type inference
3. **Production Ready**: Real-world patterns, not just examples
4. **Progressive Enhancement**: Works without JavaScript when possible
5. **Performance**: Optimized data loading and rendering strategies
6. **Framework Mode**: Emphasizes modern React Router v7 approach

## Contributing

These agents are designed to be:
- **Lean**: No unnecessary explanations
- **Precise**: Exactly what's needed, nothing more
- **Practical**: Production-ready code patterns
- **Current**: React Router v7 specific

---

**Created for Claude Code** - Optimized for AI-assisted development workflows
