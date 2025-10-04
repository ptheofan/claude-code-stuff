# Apollo GraphQL React Expert Agents

Specialized subagents for Apollo Client React, designed for use with Claude Code. Each agent focuses on specific Apollo Client domains and provides production-ready, type-safe solutions.

## Available Agents

### Core Setup & Configuration
- **apollo-client-setup-expert**: ApolloClient instantiation, ApolloProvider, link configuration, cache setup, authentication, and error handling

### Data Operations
- **apollo-queries-expert**: useQuery hook, useLazyQuery, fetch policies, polling, refetching, and query patterns
- **apollo-mutations-expert**: useMutation hook, cache updates, optimistic responses, refetching queries, and mutation patterns
- **apollo-subscriptions-expert**: useSubscription hook, WebSocket setup, real-time updates, and subscription patterns

### Advanced Features
- **apollo-caching-expert**: InMemoryCache configuration, normalization, type policies, field policies, cache reading/writing, and persistence
- **apollo-testing-expert**: MockedProvider, query/mutation testing, component testing, integration tests, and mocking strategies

## Agent Selection Guidelines

**When setting up Apollo Client:**
- Use `apollo-client-setup-expert` for client configuration, authentication, links, and error handling

**When fetching data:**
- Use `apollo-queries-expert` for data fetching, polling, refetching, and query optimization
- Use `apollo-subscriptions-expert` for real-time data and WebSocket connections

**When modifying data:**
- Use `apollo-mutations-expert` for mutations, cache updates, and optimistic UI

**When working with cache:**
- Use `apollo-caching-expert` for cache configuration, normalization, field policies, and cache manipulation

**When writing tests:**
- Use `apollo-testing-expert` for testing components with Apollo Client

## Usage with Claude Code

Add these agents to your Claude Code workflow:

```markdown
Before starting work, select appropriate Apollo Client agent(s):

- apollo-client-setup-expert: Client setup and configuration
- apollo-queries-expert: Data fetching and queries
- apollo-mutations-expert: Data mutations
- apollo-caching-expert: Cache management
- apollo-subscriptions-expert: Real-time subscriptions
- apollo-testing-expert: Testing
```

## Key Features

### Production-Ready Patterns
- Type-safe implementations with TypeScript
- Best practices and common patterns
- Error handling and edge cases
- Performance optimization

### Complete Coverage
- Client setup and authentication
- Queries, mutations, and subscriptions
- Advanced caching strategies
- Comprehensive testing patterns

### Real-World Examples
- Authentication flows
- Optimistic UI updates
- Real-time data synchronization
- Offline support

## Examples

### Setting up Apollo Client with authentication
```typescript
// Use: apollo-client-setup-expert

import { ApolloClient, InMemoryCache, HttpLink, from } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import { onError } from '@apollo/client/link/error';

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
    graphQLErrors.forEach(({ extensions }) => {
      if (extensions?.code === 'UNAUTHENTICATED') {
        window.location.href = '/login';
      }
    });
  }
});

const httpLink = new HttpLink({
  uri: 'https://api.example.com/graphql',
});

const client = new ApolloClient({
  link: from([errorLink, authLink, httpLink]),
  cache: new InMemoryCache(),
});
```

### Fetching data with polling
```typescript
// Use: apollo-queries-expert

const { data, startPolling, stopPolling } = useQuery(GET_LIVE_DATA, {
  pollInterval: 5000, // Poll every 5 seconds
});

useEffect(() => {
  if (shouldPoll) {
    startPolling(5000);
  } else {
    stopPolling();
  }
}, [shouldPoll]);
```

### Mutation with optimistic UI
```typescript
// Use: apollo-mutations-expert

const [toggleTodo] = useMutation(TOGGLE_TODO, {
  optimisticResponse: {
    toggleTodo: {
      __typename: 'Todo',
      id: todoId,
      completed: !currentCompleted,
    },
  },
});
```

### Real-time subscriptions
```typescript
// Use: apollo-subscriptions-expert

const { data } = useSubscription(COMMENT_SUBSCRIPTION, {
  variables: { postId },
  onData: ({ client, data }) => {
    client.cache.modify({
      fields: {
        comments(existingComments = []) {
          const newCommentRef = client.cache.writeFragment({
            data: data.data.commentAdded,
            fragment: COMMENT_FRAGMENT,
          });
          return [...existingComments, newCommentRef];
        },
      },
    });
  },
});
```

### Advanced caching
```typescript
// Use: apollo-caching-expert

const cache = new InMemoryCache({
  typePolicies: {
    Query: {
      fields: {
        posts: {
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
  },
});
```

### Testing components
```typescript
// Use: apollo-testing-expert

const mocks = [
  {
    request: {
      query: GET_DOG,
      variables: { breed: 'bulldog' },
    },
    result: {
      data: {
        dog: {
          __typename: 'Dog',
          id: '1',
          breed: 'Bulldog',
        },
      },
    },
  },
];

it('renders dog', async () => {
  render(
    <MockedProvider mocks={mocks}>
      <Dog breed="bulldog" />
    </MockedProvider>
  );

  expect(await screen.findByText('Bulldog')).toBeInTheDocument();
});
```

## Documentation References

- [Apollo Client React Docs](https://www.apollographql.com/docs/react)
- [Getting Started](https://www.apollographql.com/docs/react/get-started)
- [Queries](https://www.apollographql.com/docs/react/data/queries)
- [Mutations](https://www.apollographql.com/docs/react/data/mutations)
- [Caching](https://www.apollographql.com/docs/react/caching/overview)
- [Subscriptions](https://www.apollographql.com/docs/react/data/subscriptions)
- [Testing](https://www.apollographql.com/docs/react/development-testing/testing)

## Design Principles

1. **Precision**: Concise, focused expertise for Claude Code efficiency
2. **Type Safety**: TypeScript-first with proper type inference
3. **Production Ready**: Real-world patterns, not just examples
4. **Performance**: Optimized queries, mutations, and caching
5. **Testing**: Comprehensive testing patterns with MockedProvider
6. **Real-Time**: WebSocket subscriptions for live data

## Contributing

These agents are designed to be:
- **Lean**: No unnecessary explanations
- **Precise**: Exactly what's needed, nothing more
- **Practical**: Production-ready code patterns
- **Current**: Apollo Client 3.x specific

---

**Created for Claude Code** - Optimized for AI-assisted development workflows
