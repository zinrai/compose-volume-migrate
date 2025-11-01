package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "export":
		if err := runExport(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "import":
		if err := runImport(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`compose-volume-migrate %s

Usage:
  compose-volume-migrate <export|import|help>

Commands:
  export    Export external volumes to tar.gz files
  import    Import external volumes from tar.gz files
  help      Show this help message

Description:
  Migrate Docker Compose external volumes between hosts.
  Only volumes with "external: true" are processed.

`, version)
}
