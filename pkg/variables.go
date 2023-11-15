package pkg

import (
	"log"
	"os"
	"path/filepath"
)

// Version is the application version, set by linker flags during build time.
var Version = "0.0.0"

// Directory is the root for all nogfx files.
var Directory = "/tmp/nogfx"

func init() {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed acquiring home directory: %s", err)
	}

	dir = filepath.Join(dir, "nogfx")

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed creating directory %q: %s", dir, err)
	}

	Directory = dir
}
