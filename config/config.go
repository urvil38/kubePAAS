package config

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/urvil38/kubepaas/util"
)

func CreateConfigFile(c Config) error {
	buffer := new(bytes.Buffer)

	buffer.WriteString(c.Token + "\n" + c.Email + "\n" + c.ID + "\n" + c.Name)
	confFileName ,err := util.GetConfigFilePath()
	if err  != nil {
		return err
	}
	err = ioutil.WriteFile(confFileName, buffer.Bytes(), 0400)
	if err != nil {
		return fmt.Errorf("Cound't Write to config file: %v", err.Error())
	}
	return nil
}
