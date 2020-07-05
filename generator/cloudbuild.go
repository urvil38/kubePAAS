package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/urvil38/kubepaas/schema/latest"

	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/http/client"
)

func GenerateDockerCloudBuildFile(projectMetaData config.ProjectMetaData, appMetaData latest.KubepaasConfig) error {
	timeout := 10 * time.Second
	client := client.NewHTTPClient(&timeout)

	type cloudBuild struct {
		ProjectName    string `json:"project_name"`
		CurrentVersion string `json:"current_version"`
		DockerFilePath string `json:"dockerfilePath"`
	}

	cb := cloudBuild{
		ProjectName:    projectMetaData.ProjectName,
		CurrentVersion: projectMetaData.CurrentVersion,
		DockerFilePath: appMetaData.Deploy.DockerfilePath,
	}

	b, err := json.Marshal(cb)
	if err != nil {
		return fmt.Errorf("Couldn't marshal registration details: %v", err.Error())
	}

	res, err := client.Post(fmt.Sprintf(generatorEndPoint, "cloudbuild/t/docker"), "application/json", bytes.NewReader(b))
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
		err = ioutil.WriteFile(filepath.Join(config.KubeConfig.KubepaasRoot, "docker-cloudbuild.json"), b, 0644)
		if err != nil {
			return fmt.Errorf("Unable to create docker-cloudbuild.json file: , %v", err)
		}
		return nil
	default:
		return fmt.Errorf("Inernal server Error")
	}
}

func GenerateKubernetesCloudBuildFile(projectMetaData config.ProjectMetaData) error {
	timeout := 10 * time.Second
	client := client.NewHTTPClient(&timeout)

	type cloudBuild struct {
		ProjectName    string `json:"project_name"`
		CurrentVersion string `json:"current_version"`
	}

	cb := cloudBuild{
		ProjectName:    projectMetaData.ProjectName,
		CurrentVersion: projectMetaData.CurrentVersion,
	}

	b, err := json.Marshal(cb)
	if err != nil {
		return fmt.Errorf("Couldn't marshal registration details: %v", err.Error())
	}

	res, err := client.Post(fmt.Sprintf(generatorEndPoint, "cloudbuild/t/kubernetes"), "application/json", bytes.NewReader(b))
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
		err = ioutil.WriteFile(filepath.Join(config.KubeConfig.KubepaasRoot, "kubernetes-cloudbuild.json"), b, 0644)
		if err != nil {
			return fmt.Errorf("Unable to create kubernetes-cloudbuild.json file: , %v", err)
		}
		return nil
	default:
		return fmt.Errorf("Inernal server Error")
	}
}
