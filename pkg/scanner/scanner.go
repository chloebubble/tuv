package scanner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// UVProject represents a uv Python project
type UVProject struct {
	Name          string
	Path          string
	PythonVersion string
	Size          int64
	LastModified  time.Time
	HasVenv       bool
	HasLock       bool
}

// Scanner scans directories for uv projects
type Scanner struct {
	ParentDir string
}

// NewScanner creates a new scanner for the given parent directory
func NewScanner(parentDir string) *Scanner {
	return &Scanner{
		ParentDir: parentDir,
	}
}

// ScanProjects scans the parent directory for uv projects
func (s *Scanner) ScanProjects() ([]UVProject, error) {
	var projects []UVProject

	// Check if parent directory exists
	if _, err := os.Stat(s.ParentDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("parent directory does not exist: %s", s.ParentDir)
	}

	// List all directories in the parent directory
	entries, err := os.ReadDir(s.ParentDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectPath := filepath.Join(s.ParentDir, entry.Name())

		// Check for .python-version file
		pythonVersionPath := filepath.Join(projectPath, ".python-version")
		_, hasPythonVersion := os.Stat(pythonVersionPath)

		// Check for .venv directory
		venvPath := filepath.Join(projectPath, ".venv")
		_, hasVenv := os.Stat(venvPath)

		// Check for uv.lock file
		uvLockPath := filepath.Join(projectPath, "uv.lock")
		_, hasUVLock := os.Stat(uvLockPath)

		// Check for pyproject.toml file
		pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
		_, hasPyproject := os.Stat(pyprojectPath)

		// If it has at least one of these files, consider it a uv project
		if !os.IsNotExist(hasPythonVersion) || !os.IsNotExist(hasPyproject) || (!os.IsNotExist(hasVenv) && !os.IsNotExist(hasUVLock)) {
			info, _ := entry.Info()

			// Get Python version
			pythonVersion := "unknown"
			if !os.IsNotExist(hasPythonVersion) {
				if versionBytes, err := os.ReadFile(pythonVersionPath); err == nil {
					pythonVersion = strings.TrimSpace(string(versionBytes))
				}
			}

			// Calculate directory size
			size, _ := getDirSize(projectPath)

			project := UVProject{
				Name:          entry.Name(),
				Path:          projectPath,
				PythonVersion: pythonVersion,
				Size:          size,
				LastModified:  info.ModTime(),
				HasVenv:       !os.IsNotExist(hasVenv),
				HasLock:       !os.IsNotExist(hasUVLock),
			}

			projects = append(projects, project)
		}
	}

	return projects, nil
}

// getDirSize calculates the total size of a directory in bytes
func getDirSize(path string) (int64, error) {
	var size int64

	// Use du command for efficiency on Unix systems
	cmd := exec.Command("du", "-sk", path)
	output, err := cmd.Output()
	if err == nil {
		parts := strings.Fields(string(output))
		if len(parts) > 0 {
			kbSize, err := strconv.ParseInt(parts[0], 10, 64)
			if err == nil {
				return kbSize * 1024, nil
			}
		}
	}

	// Fallback to manual calculation if du fails
	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

// FormatSize formats a file size in bytes to a human-readable string
func FormatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// RunUVCommand runs a uv command in the specified project directory
func RunUVCommand(projectPath string, args ...string) (string, error) {
	cmd := exec.Command("uv", args...)
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	return string(output), err
}
