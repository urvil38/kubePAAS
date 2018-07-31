package config

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/urvil38/kubepaas/util"
	"go.opencensus.io/trace"
)

func CreateConfigFile(ctx context.Context, c *Config) error {
	ctx, span := trace.StartSpan(ctx, "createconfigfile")
	defer span.End()
	buffer := new(bytes.Buffer)

	buffer.WriteString(c.Token + "\n" + c.Email + "\n" + c.ID + "\n" + c.Name)
	err := ioutil.WriteFile(util.GetConfigFilePath(), buffer.Bytes(), 0777)
	if err != nil {
		return fmt.Errorf("Cound't Write to config file: %v", err.Error())
	}
	return nil
}
