package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urvil38/kubepaas/schema/latest"
	"github.com/urvil38/kubepaas/util"

	"sigs.k8s.io/yaml"
)

var (
	KubeConfig KubepaasConfig
	KAppConfig latest.KubepaasConfig
)

const kubepaasAppConfigFile = `app.yml`

func CreateConfigFile(c AuthConfig) error {
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

	if _, err := os.Stat(filepath.Join(wd, kubepaasAppConfigFile)); err != nil {
		fmt.Println("\x1b[31m✗ No app.yml file exist. Make sure you have app.yml file in current project ℹ\x1b[0m")
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
		return filepath.Join(wd, kubepaasAppConfigFile), nil
	}
	return "", fmt.Errorf("Coun't find app.yml file")
}

func ParseAppConfigFile() (*latest.KubepaasConfig, error) {

	path, err := getAppConfigPath()
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, &KAppConfig)
	if err != nil {
		return nil, err
	}

	return &KAppConfig, nil
}

func ProjectMetaDataFileExist() bool {
	if _, err := os.Stat(filepath.Join(KubeConfig.KubepaasRoot, ".project.json")); err != nil {
		return false
	}
	return true
}
