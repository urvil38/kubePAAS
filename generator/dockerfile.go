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

func GenerateDockerFile(appConfig latest.KubepaasConfig) error {
	timeout := 10 * time.Second
	client := client.NewHTTPClient(&timeout)

	b, err := json.Marshal(appConfig)
	if err != nil {
		return fmt.Errorf("Couldn't marshal registration details: %v", err.Error())
	}

	res, err := client.Post(fmt.Sprintf(generatorEndPoint, "dockerfile"), "application/json", bytes.NewReader(b))
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
		err = ioutil.WriteFile(filepath.Join(config.KubeConfig.ProjectRoot, "Dockerfile"), b, 0644)
		if err != nil {
			return fmt.Errorf("Unable to create Dockerfile , %v", err)
		}
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("Runtime is not support right now")
	default:
		return fmt.Errorf("Inernal server Error")
	}

}
