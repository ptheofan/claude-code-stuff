package sync

import "fmt"

// MappingSourceNotFoundError is returned when a mapping's source path does not exist.
type MappingSourceNotFoundError struct {
	Source  string
	Mapping string
}

func (e *MappingSourceNotFoundError) Error() string {
	return fmt.Sprintf("mapping source not found: %s (mapping: %s)", e.Source, e.Mapping)
}

// MappingOverlapError is returned when two mappings target the same path.
type MappingOverlapError struct {
	TargetPath string
	Mapping1   string
	Mapping2   string
}

func (e *MappingOverlapError) Error() string {
	return fmt.Sprintf("mapping overlap at target %q: %s conflicts with %s",
		e.TargetPath, e.Mapping1, e.Mapping2)
}

// InvalidMappingError is returned when a mapping has empty source or target.
type InvalidMappingError struct {
	Mapping string
	Reason  string
}

func (e *InvalidMappingError) Error() string {
	return fmt.Sprintf("invalid mapping %q: %s", e.Mapping, e.Reason)
}
