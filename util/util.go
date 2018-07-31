package util

import (
	"path/filepath"
	"os"
)

func GetConfigFilePath() string {
	home := os.Getenv("HOME")
	return filepath.Join(home,".config","kubepaas","config")
}

func GetConfigFolderPath() string {
	home := os.Getenv("HOME")
	return filepath.Join(home,".config","kubepaas")
}

func ConfigFileExists() bool {
	path := GetConfigFilePath()
	if _,err := os.Stat(path); err != nil {
		return false
	}
	return true
}