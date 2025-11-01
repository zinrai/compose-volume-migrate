package main

import (
	"fmt"
	"os"
)

func runImport() error {
	// Check Docker availability
	if err := checkDocker(); err != nil {
		return err
	}

	// Find and parse compose file
	composePath, err := findComposeFile()
	if err != nil {
		return err
	}

	compose, err := parseComposeFile(composePath)
	if err != nil {
		return err
	}

	// Get external volumes
	volumes := getExternalVolumes(compose)
	if len(volumes) == 0 {
		fmt.Println("No external volumes found")
		return nil
	}

	// Check for missing tar.gz files
	for _, volumeName := range volumes {
		tarFile := volumeName + ".tar.gz"
		if _, err := os.Stat(tarFile); os.IsNotExist(err) {
			return fmt.Errorf("%s not found", tarFile)
		}
	}

	// Import each volume
	for _, volumeName := range volumes {
		fmt.Printf("Importing %s\n", volumeName)

		// Create volume if not exists
		if !volumeExists(volumeName) {
			if err := createVolume(volumeName); err != nil {
				return err
			}
		}

		tarFile := volumeName + ".tar.gz"
		if err := importVolume(volumeName, tarFile); err != nil {
			return fmt.Errorf("import failed for %s: %w", volumeName, err)
		}
	}

	return nil
}
