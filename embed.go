package getwellsoon

import (
	"embed"
	_ "embed"
)

// DefaultConfig holds the default configuration.
//
//go:embed config/config.default.yaml
var DefaultConfig []byte

//go:embed migrations
var Migrations embed.FS
