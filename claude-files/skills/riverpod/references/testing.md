# Riverpod Testing

## Unit Testing Providers

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
