package generator

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/util"

	"github.com/urvil38/kubepaas/http/client"
)

func KmanagerConf(metadata *config.ProjectMetaData) error {
	timeout := 10 * time.Second
	client := client.NewHTTPClient(&timeout)

	res, err := client.Get(config.CLIConf.GeneratorEndpoint + "/kubepaas/config")
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

	if metadata.CloudBuildSecret != "" {
		cloudBuildSecretFilePath := filepath.Join(util.GetConfigFolderPath(), "kubepaas-cloudbuild.json")
		err = ioutil.WriteFile(cloudBuildSecretFilePath, []byte(metadata.CloudBuildSecret), 0600)
		if err != nil {
			return err
		}
	}

	if metadata.CloudStorageSecret != "" {
		cloudStorageSecretFilePath := filepath.Join(util.GetConfigFolderPath(), "kubepaas-storage.json")
		err = ioutil.WriteFile(cloudStorageSecretFilePath, []byte(metadata.CloudStorageSecret), 0600)
		if err != nil {
			return err
		}
	}

	return nil
}
