package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urvil38/kubepaas/util"
)

func CreateConfigFile(c Config) error {
	buffer := new(bytes.Buffer)

	buffer.WriteString(c.Token + "\n" + c.Email + "\n" + c.ID + "\n" + c.Name)
	confFileName, err := util.GetConfigFilePath()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(confFileName, buffer.Bytes(), 0400)
	if err != nil {
		return fmt.Errorf("Cound't Write to config file: %v", err.Error())
	}
	return nil
}

func CheckAppConfigFileExists() bool {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Couldn't Find current working directory beacause of : %v\n", err)
	}

	if _, err := os.Stat(filepath.Join(wd, "app.json")); err != nil {
		fmt.Println("\x1b[31m✗ No app.json file exist. Make sure you have app.json file in current project ℹ\x1b[0m")
		return false
	}

	return true
}

func getAppConfigPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Couldn't Find current working directory beacause of : %v\n", err)
	}

	if CheckAppConfigFileExists() {
		return filepath.Join(wd, "app.json"), nil
	}
	return "", fmt.Errorf("Coun't find app.json file")
}

func ParseAppConfigFile() (AppConfig, error) {
	var appConfig AppConfig
	path, err := getAppConfigPath()
	if err != nil {
		return appConfig, err
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return appConfig, err
	}
	err = json.Unmarshal(b, &appConfig)
	if err != nil {
		return appConfig, err
	}
	return appConfig, nil
}

func ProjectMetaDataFileExist() bool {
	wd,_ := os.Getwd()
	if _, err := os.Stat(filepath.Join(wd, ".project.json")); err != nil {
		return false
	}
	return true
}
