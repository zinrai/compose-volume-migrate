package main

import (
	"fmt"
	"os"
)

func runExport() error {
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

	// Check for existing tar.gz files
	for _, volumeName := range volumes {
		tarFile := volumeName + ".tar.gz"
		if _, err := os.Stat(tarFile); err == nil {
			return fmt.Errorf("%s already exists", tarFile)
		}
	}

	// Export each volume
	for _, volumeName := range volumes {
		if !volumeExists(volumeName) {
			fmt.Printf("Skip %s (not found)\n", volumeName)
			continue
		}

		// Check if volume is in use
		inUse, err := checkVolumeInUse(volumeName)
		if err != nil {
			return err
		}
		if inUse {
			return fmt.Errorf("volume %s is in use by running container", volumeName)
		}

		fmt.Printf("Exporting %s\n", volumeName)

		tarFile := volumeName + ".tar.gz"
		if err := exportVolume(volumeName, tarFile); err != nil {
			return fmt.Errorf("export failed for %s: %w", volumeName, err)
		}
	}

	return nil
}
