package authservice

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/urvil38/kubepaas/types"
	"github.com/urvil38/spinner"
)

func Login(auth types.AuthCredential) error {
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
		return err
	}
	s.Stop()
	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		fmt.Println("Opss,User not found!!")
	case http.StatusUnauthorized:
		fmt.Println("Username or Password is incorrect!")
	case http.StatusOK:
		fmt.Println("Successfully logged you in!!")
	default:
		fmt.Println("Something unexpected happend!!")
	}
	
	return nil
}
