---
name: apollo-graphql
version: 1.0.0
description: Apollo Client patterns for React applications including setup, queries, mutations, caching, subscriptions, and testing. This skill should be used when the user asks to "setup Apollo Client", "write GraphQL query", "create mutation", "configure Apollo cache", "add GraphQL subscription", "test Apollo components", "use useQuery", "use useMutation", or needs guidance on Apollo error handling, optimistic updates, or cache normalization in React/React Native.
---

# Apollo GraphQL Patterns

Use **context7** for Apollo Client API docs. This skill defines OUR conventions.

## Client Setup

```typescript
import { ApolloClient, InMemoryCache, HttpLink, from } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import { onError } from '@apollo/client/link/error';

const httpLink = new HttpLink({ uri: '/graphql' });

const authLink = setContext((_, { headers }) => {
  const token = localStorage.getItem('token');
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : '',
    },
  };
});

const errorLink = onError(({ graphQLErrors, networkError }) => {
  if (graphQLErrors) {
    graphQLErrors.forEach(({ message, locations, path }) => {
      console.error(`[GraphQL error]: ${message}`);
    });
  }
  if (networkError) {
    console.error(`[Network error]: ${networkError}`);
  }
});

export const client = new ApolloClient({
  link: from([errorLink, authLink, httpLink]),
  cache: new InMemoryCache(),
});
```

## Queries

```typescript
// Define query
const GET_USER = gql`
  query GetUser($id: ID!) {
    user(id: $id) {
      id
      name
      email
    }
  }
`;

// Use in component
function UserProfile({ userId }: { userId: string }) {
  const { data, loading, error } = useQuery(GET_USER, {
    variables: { id: userId },
  });

  if (loading) return <Spinner />;
  if (error) return <ErrorMessage error={error} />;

  return <div>{data.user.name}</div>;
}
```

### Query Patterns
```typescript
// Skip until ready
const { data } = useQuery(GET_USER, {
  variables: { id: userId },
  skip: !userId,
});

// Polling
const { data } = useQuery(GET_NOTIFICATIONS, {
  pollInterval: 30000, // 30 seconds
});

// Lazy query (manual trigger)
const [getUser, { data, loading }] = useLazyQuery(GET_USER);
```

## Mutations

```typescript
const CREATE_USER = gql`
  mutation CreateUser($input: CreateUserInput!) {
    createUser(input: $input) {
      id
      name
    }
  }
`;

function CreateUserForm() {
  const [createUser, { loading }] = useMutation(CREATE_USER, {
    // Update cache after mutation
    update(cache, { data: { createUser } }) {
      cache.modify({
        fields: {
          users(existingUsers = []) {
            const newUserRef = cache.writeFragment({
              data: createUser,
              fragment: gql`
                fragment NewUser on User {
                  id
                  name
                }
              `,
            });
            return [...existingUsers, newUserRef];
          },
        },
      });
    },
    // Optimistic response
    optimisticResponse: {
      createUser: {
        __typename: 'User',
        id: 'temp-id',
        name: input.name,
      },
    },
  });

  const handleSubmit = async (input: CreateUserInput) => {
    await createUser({ variables: { input } });
  };
}
```

## Caching

### Type Policies
```typescript
const cache = new InMemoryCache({
  typePolicies: {
    Query: {
      fields: {
        // Pagination
        users: {
          keyArgs: ['filter'],
          merge(existing = [], incoming, { args }) {
            const offset = args?.offset ?? 0;
            const merged = existing.slice(0);
            for (let i = 0; i < incoming.length; i++) {
              merged[offset + i] = incoming[i];
            }
            return merged;
          },
        },
      },
    },
    User: {
      // Custom cache key
      keyFields: ['id'],
    },
  },
});
```

### Cache Operations
```typescript
// Read from cache
const user = client.readFragment({
  id: `User:${userId}`,
  fragment: USER_FRAGMENT,
});

// Write to cache
client.writeFragment({
  id: `User:${userId}`,
  fragment: USER_FRAGMENT,
  data: updatedUser,
});

// Evict from cache
client.cache.evict({ id: `User:${userId}` });
client.cache.gc(); // Garbage collect
```

## Subscriptions

```typescript
const MESSAGES_SUBSCRIPTION = gql`
  subscription OnNewMessage($channelId: ID!) {
    messageAdded(channelId: $channelId) {
      id
      content
      sender {
        id
        name
      }
    }
  }
`;

function Messages({ channelId }: { channelId: string }) {
  const { data: queryData } = useQuery(GET_MESSAGES, {
    variables: { channelId },
  });

  useSubscription(MESSAGES_SUBSCRIPTION, {
    variables: { channelId },
    onData: ({ client, data }) => {
      // Update cache with new message
      client.cache.modify({
        fields: {
          messages(existingMessages = []) {
            return [...existingMessages, data.data.messageAdded];
          },
        },
      });
    },
  });
}
```

## Testing

```typescript
import { MockedProvider } from '@apollo/client/testing';
import { render, screen, waitFor } from '@testing-library/react';

const mocks = [
  {
    request: {
      query: GET_USER,
      variables: { id: '1' },
    },
    result: {
      data: {
        user: {
          __typename: 'User',
          id: '1',
          name: 'Test User',
          email: 'test@example.com',
        },
      },
    },
  },
];

describe('UserProfile', () => {
  it('renders user data', async () => {
    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <UserProfile userId="1" />
      </MockedProvider>
    );

    expect(screen.getByText('Loading...')).toBeInTheDocument();
    
    await waitFor(() => {
      expect(screen.getByText('Test User')).toBeInTheDocument();
    });
  });

  it('handles error', async () => {
    const errorMocks = [
      {
        request: { query: GET_USER, variables: { id: '1' } },
        error: new Error('User not found'),
      },
    ];

    render(
      <MockedProvider mocks={errorMocks}>
        <UserProfile userId="1" />
      </MockedProvider>
    );

    await waitFor(() => {
      expect(screen.getByText(/error/i)).toBeInTheDocument();
    });
  });
});
```

## Error Handling

```typescript
// Global error link (see setup above)
// Component-level error handling
const { error } = useQuery(GET_USER);

if (error) {
  if (error.graphQLErrors.some(e => e.extensions?.code === 'UNAUTHENTICATED')) {
    return <Redirect to="/login" />;
  }
  return <ErrorBoundary error={error} />;
}
```
