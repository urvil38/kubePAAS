package authservice

import (
	"go.opencensus.io/trace"
	"context"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/types"
	"github.com/urvil38/spinner"
)

func Login(ctx context.Context,auth types.AuthCredential) error {
	ctx,span := trace.StartSpan(ctx,"login")
	defer span.End()

	b, err := json.Marshal(auth)
	if err != nil {
		return fmt.Errorf("Unable to marshal struct :%v", err.Error())
	}

	timeout := 20 * time.Second
	c := newHTTPClient(&timeout)
	s := spinner.New("Loging you in")
	s.Start()
	res, err := c.Client.Post(fmt.Sprintf(userserviceAPI, "login"), "application/json", bytes.NewReader(b))
	if err != nil {
		s.Stop()
		return errors.New("Unable to logged in you.Please check your internet connection")
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.TLS == nil {
		s.Stop()
		fmt.Println("WARNING! Communication is not secure, please consider using HTTPS. Letsencrypt.org offers free SSL/TLS certificates.")
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
		var c config.Config
		c.Token,c.Email = token, auth.Email
		err = getProfile(ctx,&c)
		if err != nil {
			s.Stop()
			return fmt.Errorf("Cound't get profile of user: %v", err.Error())
		}
		err = config.CreateConfigFile(ctx,&c)
		if err != nil {
			s.Stop()
			return fmt.Errorf("Cound't write user configuration details: %v", err.Error())
		}
		s.Stop()
		fmt.Println("Successfully logged you in!!")
	default:
		fmt.Println("Something unexpected happend!!")
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
