package generator

import (
	"path/filepath"
	"os"
	"io/ioutil"
	"net/http"
	"bytes"
	"time"
	"encoding/json"
	"fmt"
	"github.com/urvil38/kubepaas/http/client"
	"github.com/urvil38/kubepaas/config"
)

func GenerateKubernetesConfig(appConfig config.AppConfig,projectMetadata config.ProjectMetaData) error{

	kubernetes := config.Kubernetes{
		ProjectName: appConfig.ProjectName,
		CurrentVersion: projectMetadata.CurrentVersion,
		Port: appConfig.Port,
	}

	timeout := 10 * time.Second
	client := client.NewHTTPClient(&timeout)

	b, err := json.Marshal(kubernetes)
	if err != nil {
		return fmt.Errorf("Couldn't marshal registration details: %v", err.Error())
	}

	res, err := client.Post(fmt.Sprintf(generatorEndPoint, "kubernetes"), "application/json", bytes.NewReader(b))
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
		projectRoot, _ := os.Getwd()
		err = os.MkdirAll(filepath.Join(projectRoot,"kubernetes"),0777)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(projectRoot,"kubernetes","kubernetes.yaml"), b, 0777)
		if err != nil {
			return fmt.Errorf("Unable to create kubernetes.yaml , %v", err)
		}
		return nil
	default:
		return fmt.Errorf("Inernal server Error")
	}
}