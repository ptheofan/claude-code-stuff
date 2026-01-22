package sync

import (
	"errors"
	"testing"
)

func TestMappingSourceNotFoundError_Error(t *testing.T) {
	err := &MappingSourceNotFoundError{
		Source:  "nonexistent.md",
		Mapping: "nonexistent.md -> target.md",
	}

	msg := err.Error()

	if msg == "" {
		t.Error("expected non-empty error message")
	}
	if !contains(msg, "nonexistent.md") {
		t.Errorf("expected error to contain source path, got: %s", msg)
	}
}

func TestMappingSourceNotFoundError_ImplementsError(t *testing.T) {
	var err error = &MappingSourceNotFoundError{
		Source:  "test.md",
		Mapping: "test.md -> test.md",
	}

	if err == nil {
		t.Error("expected non-nil error")
	}
}

func TestMappingOverlapError_Error(t *testing.T) {
	err := &MappingOverlapError{
		TargetPath: "CLAUDE.md",
		Mapping1:   "source1.md -> CLAUDE.md",
		Mapping2:   "source2.md -> CLAUDE.md",
	}

	msg := err.Error()

	if msg == "" {
		t.Error("expected non-empty error message")
	}
	if !contains(msg, "CLAUDE.md") {
		t.Errorf("expected error to contain target path, got: %s", msg)
	}
	if !contains(msg, "source1.md") || !contains(msg, "source2.md") {
		t.Errorf("expected error to contain both mapping references, got: %s", msg)
	}
}

func TestMappingOverlapError_ImplementsError(t *testing.T) {
	var err error = &MappingOverlapError{
		TargetPath: "test.md",
		Mapping1:   "a.md -> test.md",
		Mapping2:   "b.md -> test.md",
	}

	if err == nil {
		t.Error("expected non-nil error")
	}
}

func TestInvalidMappingError_Error(t *testing.T) {
	err := &InvalidMappingError{
		Mapping: " -> target.md",
		Reason:  "empty source",
	}

	msg := err.Error()

	if msg == "" {
		t.Error("expected non-empty error message")
	}
	if !contains(msg, "empty source") {
		t.Errorf("expected error to contain reason, got: %s", msg)
	}
}

func TestInvalidMappingError_ImplementsError(t *testing.T) {
	var err error = &InvalidMappingError{
		Mapping: "test",
		Reason:  "test reason",
	}

	if err == nil {
		t.Error("expected non-nil error")
	}
}

func TestErrors_CanBeCheckedWithErrorsAs(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "MappingSourceNotFoundError",
			err:  &MappingSourceNotFoundError{Source: "x", Mapping: "x -> y"},
		},
		{
			name: "MappingOverlapError",
			err:  &MappingOverlapError{TargetPath: "x", Mapping1: "a", Mapping2: "b"},
		},
		{
			name: "InvalidMappingError",
			err:  &InvalidMappingError{Mapping: "x", Reason: "y"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "MappingSourceNotFoundError":
				var target *MappingSourceNotFoundError
				if !errors.As(tt.err, &target) {
					t.Errorf("expected errors.As to match %s", tt.name)
				}
			case "MappingOverlapError":
				var target *MappingOverlapError
				if !errors.As(tt.err, &target) {
					t.Errorf("expected errors.As to match %s", tt.name)
				}
			case "InvalidMappingError":
				var target *InvalidMappingError
				if !errors.As(tt.err, &target) {
					t.Errorf("expected errors.As to match %s", tt.name)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
