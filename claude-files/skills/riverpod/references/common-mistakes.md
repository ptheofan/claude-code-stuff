# Common Riverpod Mistakes

## Modifying Provider During Build

```dart
// WRONG - modifying during build
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
// CORRECT - use useEffect or callback
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

## Using ref.watch() in Callbacks

```dart
// WRONG - watch in callback causes unnecessary rebuilds
onPressed: () {
  final user = ref.watch(userProvider);  // DON'T
  doSomething(user);
}
```

```dart
// CORRECT - use ref.read() in callbacks
onPressed: () {
  final user = ref.read(userProvider);  // DO
  doSomething(user);
}
```

## Modifying State in initState/dispose

```dart
// WRONG - modifying in lifecycle methods
class MyWidget extends ConsumerStatefulWidget {
  @override
  void initState() {
    super.initState();
    ref.read(counterProvider.notifier).increment();  // THROWS
  }
}
```

```dart
// CORRECT - use WidgetsBinding or Future
class MyWidget extends ConsumerStatefulWidget {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(counterProvider.notifier).increment();
    });
  }
}

// BETTER - use HookConsumerWidget with useEffect
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

## Calling async Methods Without Guard

```dart
// WRONG - unhandled errors crash the app
Future<void> loadData() async {
  state = const AsyncLoading();
  state = AsyncData(await fetchData());  // Throws on error!
}
```

```dart
// CORRECT - use AsyncValue.guard
Future<void> loadData() async {
  state = const AsyncLoading();
  state = await AsyncValue.guard(() => fetchData());
}
```
