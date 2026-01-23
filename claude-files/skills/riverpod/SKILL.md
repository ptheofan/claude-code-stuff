---
name: flutter-riverpod
description: Flutter state management with Riverpod 3.x and flutter_hooks. Use when building Flutter apps, managing state, or working with providers and hooks.
---

# Flutter Riverpod + Hooks Patterns

Use **context7** for Riverpod/Flutter API docs. This skill defines OUR conventions.

## Setup

```yaml
# pubspec.yaml
dependencies:
  flutter_hooks: ^0.20.0
  hooks_riverpod: ^3.2.0
  riverpod_annotation: ^4.0.1

dev_dependencies:
  build_runner:
  riverpod_generator: ^4.0.2
  riverpod_lint:  # Optional but recommended
```

```dart
// main.dart
void main() {
  runApp(
    ProviderScope(
      child: MyApp(),
    ),
  );
}
```

```yaml
# analysis_options.yaml (optional - enables riverpod lint rules)
analyzer:
  plugins:
    - riverpod_lint
```

## Widget Base Class

**Always use `HookConsumerWidget`** - combines hooks + Riverpod:

```dart
class MyWidget extends HookConsumerWidget {
  const MyWidget({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // Hooks
    final counter = useState(0);
    final controller = useTextEditingController();
    
    // Riverpod
    final user = ref.watch(userProvider);
    
    return Text('${user.name}: ${counter.value}');
  }
}
```

For stateful with lifecycle methods, use `StatefulHookConsumerWidget`.

## Provider Types (Riverpod 3.x)

### With Code Generation (Preferred)

```dart
// providers/user_provider.dart
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'user_provider.g.dart';  // REQUIRED!

// Simple provider (read-only)
@riverpod
String greeting(Ref ref) => 'Hello';

// Async provider
@riverpod
Future<User> user(Ref ref) async {
  return await ref.watch(apiProvider).fetchUser();
}

// Notifier (mutable state)
@riverpod
class Counter extends _$Counter {
  @override
  int build() => 0;

  void increment() => state++;
  void decrement() => state--;
}

// Async Notifier
@riverpod
class UserNotifier extends _$UserNotifier {
  @override
  Future<User> build() async {
    return await _fetchUser();
  }

  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(_fetchUser);
  }

  Future<User> _fetchUser() async {
    return ref.read(apiProvider).fetchUser();
  }
}
```

## Code Generation

### Setup

```yaml
# pubspec.yaml
dependencies:
  riverpod_annotation: ^4.0.1

dev_dependencies:
  build_runner:
  riverpod_generator: ^4.0.2
```

### File Structure

Every file using `@riverpod` needs a `part` directive:

```dart
// user_provider.dart
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'user_provider.g.dart';  // REQUIRED - must match filename

@riverpod
String greeting(Ref ref) => 'Hello';
```

### Commands

```bash
# One-time build
dart run build_runner build

# Watch mode (rebuilds on save) - RECOMMENDED during development
dart run build_runner watch

# Fix conflicts (when .g.dart is out of sync)
dart run build_runner build --delete-conflicting-outputs
```

### Common Errors

| Error | Cause | Fix |
|-------|-------|-----|
| `Could not find generator for @riverpod` | Missing `riverpod_generator` | Add to dev_dependencies |
| `part 'x.g.dart' not found` | Generator not run | Run `build_runner build` |
| `Conflicting outputs` | Stale generated files | Add `--delete-conflicting-outputs` |
| `The name '_$ClassName' isn't defined` | Missing part directive | Add `part 'filename.g.dart';` |

### Without Code Generation

```dart
// Simple provider
final greetingProvider = Provider<String>((ref) => 'Hello');

// State provider (simple mutable)
final counterProvider = StateProvider<int>((ref) => 0);

// Future provider
final userProvider = FutureProvider<User>((ref) async {
  return ref.watch(apiProvider).fetchUser();
});

// Stream provider
final messagesProvider = StreamProvider<List<Message>>((ref) {
  return ref.watch(chatProvider).messageStream;
});

// StateNotifier (complex mutable - legacy but still works)
final todosProvider = StateNotifierProvider<TodosNotifier, List<Todo>>((ref) {
  return TodosNotifier();
});
```

## ref Methods

| Method | Use Case | Rebuilds? |
|--------|----------|-----------|
| `ref.watch()` | In build method, reactive | Yes |
| `ref.read()` | In callbacks, one-time | No |
| `ref.listen()` | Side effects (snackbar, nav) | No |

```dart
class MyWidget extends HookConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // ✅ watch in build - reactive
    final user = ref.watch(userProvider);
    
    // ✅ listen for side effects
    ref.listen(authProvider, (prev, next) {
      if (next == null) {
        context.go('/login');
      }
    });

    return ElevatedButton(
      onPressed: () {
        // ✅ read in callback - one-time
        ref.read(counterProvider.notifier).increment();
      },
      child: Text(user.name),
    );
  }
}
```

## AsyncValue Handling

```dart
class UserWidget extends HookConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final userAsync = ref.watch(userProvider);

    return userAsync.when(
      data: (user) => Text(user.name),
      loading: () => const CircularProgressIndicator(),
      error: (error, stack) => Text('Error: $error'),
    );
  }
}

// Skip loading on refresh (keep showing old data)
return userAsync.when(
  skipLoadingOnRefresh: true,
  data: (user) => Text(user.name),
  loading: () => const CircularProgressIndicator(),
  error: (error, stack) => Text('Error: $error'),
);
```

## Common Hooks

```dart
class MyWidget extends HookConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // State
    final counter = useState(0);
    final isLoading = useState(false);

    // Controllers (auto-disposed)
    final textController = useTextEditingController();
    final scrollController = useScrollController();
    final focusNode = useFocusNode();

    // Animation (auto-disposed)
    final animController = useAnimationController(
      duration: const Duration(milliseconds: 300),
    );

    // Side effects
    useEffect(() {
      // Runs on mount (and when dependencies change)
      animController.forward();
      
      // Return cleanup function
      return () => print('Disposed');
    }, []); // Empty = run once

    // Memoization
    final expensive = useMemoized(() => computeExpensive(counter.value), [counter.value]);

    return Column(
      children: [
        TextField(controller: textController),
        Text('Count: ${counter.value}'),
        ElevatedButton(
          onPressed: () => counter.value++,
          child: const Text('Increment'),
        ),
      ],
    );
  }
}
```

## Provider Modifiers

```dart
// Family - parameterized providers
@riverpod
Future<User> userById(Ref ref, String id) async {
  return ref.watch(apiProvider).fetchUser(id);
}

// Usage
final user = ref.watch(userByIdProvider('123'));

// AutoDispose (default in codegen, explicit without)
final dataProvider = FutureProvider.autoDispose<Data>((ref) async {
  // Auto-disposed when no longer listened to
  return fetchData();
});

// Keep alive temporarily
@riverpod
Future<Data> cachedData(Ref ref) async {
  final link = ref.keepAlive();
  
  // Dispose after 30 seconds of no listeners
  final timer = Timer(Duration(seconds: 30), link.close);
  ref.onDispose(timer.cancel);
  
  return fetchData();
}
```

## Testing Providers

```dart
void main() {
  test('counter increments', () {
    final container = ProviderContainer();
    addTearDown(container.dispose);

    expect(container.read(counterProvider), 0);
    
    container.read(counterProvider.notifier).increment();
    
    expect(container.read(counterProvider), 1);
  });

  test('user provider fetches data', () async {
    final container = ProviderContainer(
      overrides: [
        // Mock the API
        apiProvider.overrideWithValue(MockApi()),
      ],
    );
    addTearDown(container.dispose);

    final user = await container.read(userProvider.future);
    expect(user.name, 'Test User');
  });
}
```

## Widget Testing with Hooks

```dart
testWidgets('counter widget', (tester) async {
  await tester.pumpWidget(
    ProviderScope(
      overrides: [
        counterProvider.overrideWith((ref) => 10),
      ],
      child: const MaterialApp(home: CounterWidget()),
    ),
  );

  expect(find.text('10'), findsOneWidget);
  
  await tester.tap(find.byType(ElevatedButton));
  await tester.pump();
  
  expect(find.text('11'), findsOneWidget);
});
```

## Project Structure

```
lib/
├── main.dart
├── core/
│   ├── providers/           # Global providers (api, auth, etc.)
│   └── hooks/               # Custom hooks
├── features/
│   └── users/
│       ├── providers/
│       │   ├── user_provider.dart
│       │   └── user_provider.g.dart
│       ├── models/
│       └── widgets/
```

## Common Mistakes

### ❌ Modifying Provider During Build

```dart
// ❌ WRONG - modifying during build
class MyWidget extends HookConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // This throws: "Tried to modify a provider while the widget tree was building"
    ref.read(friendsProvider.notifier).markAsRead();
    
    return Container();
  }
}
```

```dart
// ✅ CORRECT - use useEffect or callback
class MyWidget extends HookConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // Option 1: useEffect (runs after build)
    useEffect(() {
      ref.read(friendsProvider.notifier).markAsRead();
      return null;
    }, []);

    // Option 2: Future.microtask (delays until after build)
    useEffect(() {
      Future.microtask(() {
        ref.read(friendsProvider.notifier).markAsRead();
      });
      return null;
    }, []);

    return Container();
  }
}
```

### ❌ Using ref.watch() in Callbacks

```dart
// ❌ WRONG - watch in callback causes unnecessary rebuilds
onPressed: () {
  final user = ref.watch(userProvider);  // DON'T
  doSomething(user);
}
```

```dart
// ✅ CORRECT - use ref.read() in callbacks
onPressed: () {
  final user = ref.read(userProvider);  // DO
  doSomething(user);
}
```

### ❌ Modifying State in initState/dispose

```dart
// ❌ WRONG - modifying in lifecycle methods
class MyWidget extends ConsumerStatefulWidget {
  @override
  void initState() {
    super.initState();
    ref.read(counterProvider.notifier).increment();  // THROWS
  }
}
```

```dart
// ✅ CORRECT - use WidgetsBinding or Future
class MyWidget extends ConsumerStatefulWidget {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(counterProvider.notifier).increment();
    });
  }
}

// ✅ BETTER - use HookConsumerWidget with useEffect
class MyWidget extends HookConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    useEffect(() {
      ref.read(counterProvider.notifier).increment();
      return null;
    }, []);
    
    return Container();
  }
}
```

### ❌ Calling async Methods Without Guard

```dart
// ❌ WRONG - unhandled errors crash the app
Future<void> loadData() async {
  state = const AsyncLoading();
  state = AsyncData(await fetchData());  // Throws on error!
}
```

```dart
// ✅ CORRECT - use AsyncValue.guard
Future<void> loadData() async {
  state = const AsyncLoading();
  state = await AsyncValue.guard(() => fetchData());
}
```

## Best Practices

1. **Use `HookConsumerWidget`** as default widget base
2. **Use code generation** (`@riverpod`) for new providers
3. **`ref.watch()` in build**, `ref.read()` in callbacks
4. **Never call `ref.watch()` in callbacks** - causes rebuilds
5. **Use `AsyncValue.when()`** for async state
6. **Keep providers small and focused**
7. **Use `ref.listen()` for navigation/snackbars** - not `ref.watch()`
