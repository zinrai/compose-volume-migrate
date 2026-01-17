package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const busyboxImage = "busybox:stable-glibc"

// checkDocker checks if docker command is available
func checkDocker() error {
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker not available: %w", err)
	}
	return nil
}

// checkVolumeInUse checks if volume is being used by any running container
func checkVolumeInUse(volumeName string) (bool, error) {
	cmd := exec.Command("docker", "ps", "-q", "--filter", fmt.Sprintf("volume=%s", volumeName))
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check volume usage: %w", err)
	}
	return len(bytes.TrimSpace(output)) > 0, nil
}

// volumeExists checks if docker volume exists
func volumeExists(volumeName string) bool {
	cmd := exec.Command("docker", "volume", "inspect", volumeName)
	return cmd.Run() == nil
}

// createVolume creates docker volume
func createVolume(volumeName string) error {
	cmd := exec.Command("docker", "volume", "create", volumeName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create volume %s: %w", volumeName, err)
	}
	return nil
}

// exportVolume exports volume to tar.gz
func exportVolume(volumeName, tarFile string) error {
	cwd, err := getCurrentDir()
	if err != nil {
		return err
	}

	cmd := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/data:ro", volumeName),
		"-v", fmt.Sprintf("%s:/backup", cwd),
		busyboxImage,
		"tar", "czf", fmt.Sprintf("/backup/%s", tarFile), "-C", "/data", ".")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to export volume: %w\n%s", err, stderr.String())
	}

	return nil
}

// importVolume imports volume from tar.gz
func importVolume(volumeName, tarFile string) error {
	cwd, err := getCurrentDir()
	if err != nil {
		return err
	}

	cmd := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/data", volumeName),
		"-v", fmt.Sprintf("%s:/backup", cwd),
		busyboxImage,
		"tar", "xzf", fmt.Sprintf("/backup/%s", tarFile), "-C", "/data")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to import volume: %w\n%s", err, stderr.String())
	}

	return nil
}

// getCurrentDir returns current directory absolute path
func getCurrentDir() (string, error) {
	cmd := exec.Command("pwd")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
