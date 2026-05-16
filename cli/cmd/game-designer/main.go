package main

import (
	"os"

	"github.com/example/game-designer-cli/internal/commands"
)

func main() {
	root := commands.NewRootCmd()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
