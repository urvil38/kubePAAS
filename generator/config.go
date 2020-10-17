package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/urvil38/kubepaas/config"

	"github.com/urvil38/kubepaas/http/client"
)

const (
	generatorEndPoint = "https://generator.kubepaas.ml/%s"
)

func KmanagerConf(metadata *config.ProjectMetaData) error {
	timeout := 10 * time.Second
	client := client.NewHTTPClient(&timeout)

	res, err := client.Get(fmt.Sprintf(generatorEndPoint, "kubepaas/config"))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, metadata)
	if err != nil {
		return err
	}
	return nil
}
