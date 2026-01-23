# Riverpod Code Generation

## Setup

```yaml
# pubspec.yaml
dependencies:
  riverpod_annotation: ^4.0.1

dev_dependencies:
  build_runner:
  riverpod_generator: ^4.0.2
```

## File Structure

Every file using `@riverpod` needs a `part` directive:

```dart
// user_provider.dart
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'user_provider.g.dart';  // REQUIRED - must match filename

@riverpod
String greeting(Ref ref) => 'Hello';
```

## Commands

```bash
# One-time build
dart run build_runner build

# Watch mode (rebuilds on save) - RECOMMENDED during development
dart run build_runner watch

# Fix conflicts (when .g.dart is out of sync)
dart run build_runner build --delete-conflicting-outputs
```

## Common Errors

| Error | Cause | Fix |
|-------|-------|-----|
| `Could not find generator for @riverpod` | Missing `riverpod_generator` | Add to dev_dependencies |
| `part 'x.g.dart' not found` | Generator not run | Run `build_runner build` |
| `Conflicting outputs` | Stale generated files | Add `--delete-conflicting-outputs` |
| `The name '_$ClassName' isn't defined` | Missing part directive | Add `part 'filename.g.dart';` |
