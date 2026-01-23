---
name: react-router
description: React Router v7 patterns including routing, data loading, forms, actions, and error handling. Use when building routing, navigation, or data-fetching in React applications.
---

# React Router Patterns

Use **context7** for React Router API docs. This skill defines OUR conventions.

## Route Configuration

### File-Based Routes (Framework Mode)
```
app/
├── routes/
│   ├── _index.tsx           # /
│   ├── about.tsx            # /about
│   ├── users.tsx            # /users (layout)
│   ├── users._index.tsx     # /users
│   ├── users.$id.tsx        # /users/:id
│   └── users.$id.edit.tsx   # /users/:id/edit
```

### Route Module Structure
```typescript
// app/routes/users.$id.tsx
import type { Route } from './+types/users.$id';

// Loader: Fetch data
export async function loader({ params }: Route.LoaderArgs) {
  const user = await getUser(params.id);
  if (!user) {
    throw new Response('Not found', { status: 404 });
  }
  return { user };
}

// Action: Handle mutations
export async function action({ request, params }: Route.ActionArgs) {
  const formData = await request.formData();
  await updateUser(params.id, Object.fromEntries(formData));
  return redirect(`/users/${params.id}`);
}

// Component
export default function UserPage({ loaderData }: Route.ComponentProps) {
  const { user } = loaderData;
  return <div>{user.name}</div>;
}

// Error Boundary
export function ErrorBoundary() {
  const error = useRouteError();
  return <ErrorDisplay error={error} />;
}
```

## Data Loading

### Basic Loader
```typescript
export async function loader({ params, request }: Route.LoaderArgs) {
  const url = new URL(request.url);
  const query = url.searchParams.get('q');
  
  const users = await searchUsers(query);
  return { users, query };
}

function UsersPage({ loaderData }: Route.ComponentProps) {
  const { users, query } = loaderData;
  // ...
}
```

### Deferred Data
```typescript
import { defer, Await } from 'react-router';
import { Suspense } from 'react';

export async function loader({ params }: Route.LoaderArgs) {
  // Critical data - awaited
  const user = await getUser(params.id);
  
  // Non-critical - deferred
  const postsPromise = getUserPosts(params.id);
  
  return defer({
    user,
    posts: postsPromise,
  });
}

function UserPage({ loaderData }: Route.ComponentProps) {
  const { user, posts } = loaderData;
  
  return (
    <div>
      <h1>{user.name}</h1>
      <Suspense fallback={<PostsSkeleton />}>
        <Await resolve={posts}>
          {(resolvedPosts) => <PostsList posts={resolvedPosts} />}
        </Await>
      </Suspense>
    </div>
  );
}
```

## Forms & Actions

### Basic Form
```typescript
import { Form, useNavigation } from 'react-router';

function CreateUserForm() {
  const navigation = useNavigation();
  const isSubmitting = navigation.state === 'submitting';

  return (
    <Form method="post">
      <input name="name" required />
      <input name="email" type="email" required />
      <button type="submit" disabled={isSubmitting}>
        {isSubmitting ? 'Creating...' : 'Create User'}
      </button>
    </Form>
  );
}

export async function action({ request }: Route.ActionArgs) {
  const formData = await request.formData();
  
  const errors = validate(formData);
  if (errors) {
    return { errors };
  }
  
  const user = await createUser(Object.fromEntries(formData));
  return redirect(`/users/${user.id}`);
}
```

### useFetcher (No Navigation)
```typescript
function LikeButton({ postId }: { postId: string }) {
  const fetcher = useFetcher();
  const isLiking = fetcher.state !== 'idle';

  return (
    <fetcher.Form method="post" action={`/posts/${postId}/like`}>
      <button disabled={isLiking}>
        {isLiking ? '...' : '❤️'}
      </button>
    </fetcher.Form>
  );
}
```

### Optimistic UI
```typescript
function TodoItem({ todo }: { todo: Todo }) {
  const fetcher = useFetcher();
  
  // Use optimistic value if submitting
  const isComplete = fetcher.formData
    ? fetcher.formData.get('complete') === 'true'
    : todo.complete;

  return (
    <fetcher.Form method="post">
      <input
        type="checkbox"
        name="complete"
        value="true"
        checked={isComplete}
        onChange={(e) => fetcher.submit(e.target.form)}
      />
      <span className={isComplete ? 'line-through' : ''}>
        {todo.title}
      </span>
    </fetcher.Form>
  );
}
```

## Navigation & Hooks

```typescript
import {
  useNavigate,
  useParams,
  useSearchParams,
  useLocation,
  useNavigation,
} from 'react-router';

function Component() {
  const navigate = useNavigate();
  const params = useParams();           // { id: '123' }
  const [searchParams, setSearchParams] = useSearchParams();
  const location = useLocation();       // { pathname, search, hash }
  const navigation = useNavigation();   // { state, location, formData }

  // Programmatic navigation
  navigate('/users');
  navigate(-1); // Back
  navigate('/users', { replace: true });

  // Update search params
  setSearchParams({ page: '2', sort: 'name' });
}
```

## Error Handling

### Route Error Boundary
```typescript
import { useRouteError, isRouteErrorResponse } from 'react-router';

export function ErrorBoundary() {
  const error = useRouteError();

  if (isRouteErrorResponse(error)) {
    return (
      <div>
        <h1>{error.status} {error.statusText}</h1>
        <p>{error.data}</p>
      </div>
    );
  }

  return (
    <div>
      <h1>Unexpected Error</h1>
      <p>{error instanceof Error ? error.message : 'Unknown error'}</p>
    </div>
  );
}
```

### Throwing Responses
```typescript
export async function loader({ params }: Route.LoaderArgs) {
  const user = await getUser(params.id);
  
  if (!user) {
    throw new Response('User not found', { status: 404 });
  }
  
  if (!canViewUser(user)) {
    throw new Response('Forbidden', { status: 403 });
  }
  
  return { user };
}
```

## Testing

```typescript
import { createMemoryRouter, RouterProvider } from 'react-router';
import { render, screen, waitFor } from '@testing-library/react';

describe('UserPage', () => {
  it('renders user data', async () => {
    const routes = [
      {
        path: '/users/:id',
        element: <UserPage />,
        loader: () => ({ user: { id: '1', name: 'Test User' } }),
      },
    ];

    const router = createMemoryRouter(routes, {
      initialEntries: ['/users/1'],
    });

    render(<RouterProvider router={router} />);

    await waitFor(() => {
      expect(screen.getByText('Test User')).toBeInTheDocument();
    });
  });
});
```
