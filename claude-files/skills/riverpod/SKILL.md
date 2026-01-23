---
name: riverpod
version: 1.0.0
description: Flutter state management with Riverpod 3.x and flutter_hooks. This skill should be used when the user asks to "add Riverpod", "create provider", "manage Flutter state", "use hooks in Flutter", "setup HookConsumerWidget", "create notifier", "handle async state", or needs guidance on provider patterns, ref.watch vs ref.read, AsyncValue handling, or Flutter state management best practices.
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
```

```dart
// main.dart
void main() {
  runApp(ProviderScope(child: MyApp()));
}
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

## Provider Types (Codegen - Preferred)

```dart
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'user_provider.g.dart';  // REQUIRED

// Simple provider
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
}

// Async Notifier
@riverpod
class UserNotifier extends _$UserNotifier {
  @override
  Future<User> build() async => await _fetchUser();

  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(_fetchUser);
  }

  Future<User> _fetchUser() async {
    return ref.read(apiProvider).fetchUser();
  }
}
```

For code generation setup and troubleshooting, see `references/codegen.md`.

## ref Methods

| Method | Use Case | Rebuilds? |
|--------|----------|-----------|
| `ref.watch()` | In build method, reactive | Yes |
| `ref.read()` | In callbacks, one-time | No |
| `ref.listen()` | Side effects (snackbar, nav) | No |

```dart
@override
Widget build(BuildContext context, WidgetRef ref) {
  // watch in build - reactive
  final user = ref.watch(userProvider);

  // listen for side effects
  ref.listen(authProvider, (prev, next) {
    if (next == null) context.go('/login');
  });

  return ElevatedButton(
    onPressed: () {
      // read in callback - one-time
      ref.read(counterProvider.notifier).increment();
    },
    child: Text(user.name),
  );
}
```

## AsyncValue Handling

```dart
final userAsync = ref.watch(userProvider);

return userAsync.when(
  data: (user) => Text(user.name),
  loading: () => const CircularProgressIndicator(),
  error: (error, stack) => Text('Error: $error'),
);
```

## Common Hooks

```dart
// State
final counter = useState(0);

// Controllers (auto-disposed)
final textController = useTextEditingController();
final scrollController = useScrollController();

// Side effects
useEffect(() {
  // Runs on mount
  return () => print('Disposed');  // Cleanup
}, []);  // Empty = run once

// Memoization
final expensive = useMemoized(() => compute(value), [value]);
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
```

## Best Practices

1. **Use `HookConsumerWidget`** as default widget base
2. **Use code generation** (`@riverpod`) for new providers
3. **`ref.watch()` in build**, `ref.read()` in callbacks
4. **Never call `ref.watch()` in callbacks** - causes rebuilds
5. **Use `AsyncValue.when()`** for async state
6. **Keep providers small and focused**
7. **Use `ref.listen()` for navigation/snackbars**

For common mistakes and how to avoid them, see `references/common-mistakes.md`.
For testing patterns, see `references/testing.md`.
