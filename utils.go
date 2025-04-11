package gv

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CopyFile(src, dst string) error {
	// Create directory structure for the destination file if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy data
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	// Ensure data is written to disk
	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to flush to disk: %w", err)
	}

	return nil
}

// ExecCommand executes a command with arguments and returns the combined output
func ExecCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	err := cmd.Run()
	if err != nil {
		return output.String(), fmt.Errorf("command failed: %w", err)
	}

	return output.String(), nil
}

// ExecCommandWithCallback executes a command and calls the callback with the output
func ExecCommandWithCallback(command string, args []string, callback func(output string)) error {
	output, err := ExecCommand(command, args)
	callback(output)
	return err
}

// ExecStringCommand executes a command from a string with space-separated arguments
func ExecStringCommand(cmdString string) (string, error) {
	parts := strings.Fields(cmdString)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	return ExecCommand(parts[0], parts[1:])
}
