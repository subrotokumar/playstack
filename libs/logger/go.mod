module gitlab.com/subrotokumar/glitchr/pkg/logger

go 1.25.5

// Correct local replacement path to the core lib (was incorrectly pointing to ../libs/core
// which resolves to libs/libs/core). This should point to ../core from this directory.
replace gitlab.com/subrotokumar/glitchr/pkg/core => ../core

require gitlab.com/subrotokumar/glitchr/pkg/core v0.0.0-00010101000000-000000000000

require (
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
