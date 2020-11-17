package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/http/client"
	"github.com/urvil38/kubepaas/schema/latest"
)

func GenerateKubernetesConfig(appConfig latest.KubepaasConfig, projectMetadata config.ProjectMetaData) error {

	timeout := 10 * time.Second
	client := client.NewHTTPClient(&timeout)

	kubernetes := config.Kubernetes{
		ProjectName:    projectMetadata.ProjectName,
		CurrentVersion: projectMetadata.CurrentVersion,
		Spec:           appConfig,
	}

	b, err := json.Marshal(kubernetes)
	if err != nil {
		return fmt.Errorf("Couldn't marshal registration details: %v", err.Error())
	}

	res, err := client.Post(config.CLIConf.GeneratorEndpoint+"/kubernetes", "application/json", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("Unable to Signup.Check Internet Connection")
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	// if res.TLS == nil {
	// 	fmt.Println("WARNING! Communication is not secure, please consider using HTTPS. Letsencrypt.org offers free SSL/TLS certificates.")
	// }

	switch res.StatusCode {
	case http.StatusOK:
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Coun't read body of response , %v", err)
		}

		err = os.MkdirAll(filepath.Join(config.KubeConfig.KubepaasRoot, "kubernetes"), 0777)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(config.KubeConfig.KubepaasRoot, "kubernetes", "kubernetes.yaml"), b, 0777)
		if err != nil {
			return fmt.Errorf("Unable to create kubernetes.yaml , %v", err)
		}
		return nil
	default:
		return fmt.Errorf("Inernal server Error")
	}
}
