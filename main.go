package main

import (
	"os"

	"github.com/acs-dl/unverified-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
