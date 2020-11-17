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
	return filepath.Join(home, ".kubepaas", "config.json"), nil
}

func GetAuthConfigFilePath() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("HOME env is not found")
	}
	return filepath.Join(home, ".kubepaas", "auth"), nil
}

func GetConfigFolderPath() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".kubepaas")
}

func AuthConfigFileExists() bool {
	path, err := GetAuthConfigFilePath()
	if err != nil {
		return false
	}
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
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
