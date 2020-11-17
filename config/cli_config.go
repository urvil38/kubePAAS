package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/urvil38/kubepaas/util"
)

type CLIConfig struct {
	GeneratorEndpoint string `json:"-"`
	AuthEndpoint      string `json:"-"`
	kv                map[string]string
}

var CLIConfigKeys = []string{"generator-endpoint", "auth-endpoint"}

func (c *CLIConfig) ValidKey(key string) bool {
	for _, k := range CLIConfigKeys {
		if k == key {
			return true
		}
	}
	return false
}

func NewCLIConfig() *CLIConfig {
	return &CLIConfig{
		kv: make(map[string]string),
	}
}

func (c *CLIConfig) Write() error {
	fp, err := util.GetConfigFilePath()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(c.kv, "", "   ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fp, b, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (c *CLIConfig) Read() error {
	fp, err := util.GetConfigFilePath()
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &c.kv)
	if err != nil {
		return err
	}

	val, err := c.Get("generator-endpoint")
	if err != nil {
		c.GeneratorEndpoint = ""
	} else {
		c.GeneratorEndpoint = strings.TrimSuffix(val, "/")
	}

	val, err = c.Get("auth-endpoint")
	if err != nil {
		c.AuthEndpoint = ""
	} else {
		c.AuthEndpoint = strings.TrimSuffix(val, "/")
	}

	return nil
}

func (c *CLIConfig) Get(key string) (string, error) {
	value, ok := c.kv[key]
	if !ok {
		return "", fmt.Errorf("no property found of name: %s", key)
	}
	return value, nil
}

func (c *CLIConfig) Set(key, value string) error {
	if !c.ValidKey(key) {
		return fmt.Errorf("Unknown property \"%s\"", key)
	}
	c.kv[key] = value
	err := c.Write()
	if err != nil {
		return err
	}
	return nil
}

func (c *CLIConfig) Unset(key string) error {
	_, ok := c.kv[key]
	if !ok {
		return fmt.Errorf("no property found of name: %s", key)
	}
	delete(c.kv, key)
	err := c.Write()
	if err != nil {
		return err
	}
	return nil
}
