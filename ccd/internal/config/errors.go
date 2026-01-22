package config

import "fmt"

// ConfigExistsError is returned when attempting to create a config file
// that already exists without the force flag.
type ConfigExistsError struct {
	Path string
}

func (e *ConfigExistsError) Error() string {
	return fmt.Sprintf("config file already exists: %s (use --force to overwrite)", e.Path)
}

// ConfigWriteError is returned when writing the config file fails.
type ConfigWriteError struct {
	Path  string
	Cause error
}

func (e *ConfigWriteError) Error() string {
	return fmt.Sprintf("failed to write config file %s: %v", e.Path, e.Cause)
}

func (e *ConfigWriteError) Unwrap() error {
	return e.Cause
}

// ExecutablePathError is returned when the executable path cannot be determined.
type ExecutablePathError struct {
	Cause error
}

func (e *ExecutablePathError) Error() string {
	return fmt.Sprintf("failed to determine executable path: %v", e.Cause)
}

func (e *ExecutablePathError) Unwrap() error {
	return e.Cause
}
