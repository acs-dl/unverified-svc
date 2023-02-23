package main

import (
	"os"

	"gitlab.com/distributed_lab/acs/unverified-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
