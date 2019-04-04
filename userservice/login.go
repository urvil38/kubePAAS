package userservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/questions"
	"github.com/urvil38/spinner"
	"github.com/urvil38/kubepaas/http/client"
)

func Login(auth questions.AuthCredential) error {
	timeout := 15 * time.Second
	client := client.NewHTTPClient(&timeout)

	b, err := json.Marshal(auth)
	if err != nil {
		return fmt.Errorf("Unable to marshal struct :%v", err.Error())
	}

	s := spinner.New("Loging you in")
	s.Start()
	res, err := client.Post(fmt.Sprintf(userserviceEndpoint, "login"), "application/json", bytes.NewReader(b))
	if err != nil {
		s.Stop()
		return fmt.Errorf("Unable to Login.Connection Timeout ⏱")
	}

	if res.TLS == nil {
		s.Stop()
		fmt.Println("WARNING! Communication is not secure, please consider using HTTPS. Letsencrypt.org offers free SSL/TLS certificates.")
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		s.Stop()
		fmt.Println("Opss,User not found!!")
	case http.StatusUnauthorized:
		s.Stop()
		fmt.Println("Username or Password is incorrect!")
	case http.StatusOK:
		token, err := getToken(res)
		if err != nil {
			return fmt.Errorf("Cound't get user configuration details: %v", err.Error())
		}
		var conf config.Config
		conf.Token, conf.Email = token, auth.Email
		userConf, err := getUserProfile(conf)
		if err != nil {
			s.Stop()
			return fmt.Errorf("Cound't get profile of user: %v", err.Error())
		}
		conf.ID, conf.Name = userConf.ID, userConf.Name
		err = config.CreateConfigFile(conf)
		if err != nil {
			s.Stop()
			return fmt.Errorf("Cound't write user configuration details: %v", err.Error())
		}
		s.Stop()
		fmt.Println("Successfully logged you in!!")
	default:
		fmt.Println("Server Error!!")
	}

	return nil
}

func getToken(res *http.Response) (token string, err error) {

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Unable to read response body: %v", err.Error())
	}
	var c config.AuthToken

	err = json.Unmarshal(b, &c)
	if err != nil {
		return "", fmt.Errorf("Unable to unmarshal response: %v", err.Error())
	}
	return c.Token, nil
}
