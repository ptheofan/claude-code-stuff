package config

// Mapping defines a source-to-target path mapping.
// Source is relative to the working directory.
// Target is relative to the target directory.
type Mapping struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}
