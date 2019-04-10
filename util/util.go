package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetConfigFilePath() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("HOME env is not found")
	}
	return filepath.Join(home, ".config", "kubepaas", "config"), nil
}

func GetConfigFolderPath() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".config", "kubepaas")
}

func ConfigFileExists() bool {
	path, err := GetConfigFilePath()
	if err != nil {
		return false
	}
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
