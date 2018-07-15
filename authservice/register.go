package authservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/urvil38/kubepaas/types"
	"github.com/urvil38/spinner"
)

func RegisterUser(user types.UserInfo) error {

	timeout := 15 * time.Second
	c := newHTTPClient(&timeout)

	b, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("Couldn't marshal data : %v", err.Error())
	}

	s := spinner.New("Registering You")
	s.Start()
	res, err := c.Client.Post(fmt.Sprintf(userserviceAPI, "user"), "application/json", bytes.NewReader(b))
	if err != nil {
		s.Stop()
		return errors.New("Unable to register you.Check your internet connection")
	}
	if res.Body != nil {
		res.Body.Close()
	}
	s.Stop()

	switch res.StatusCode {
	case http.StatusCreated:
		fmt.Println("Welcome,Thank you for registering with kubepaas")
	case http.StatusConflict:
		fmt.Println("User is already exists!")
	default:
		fmt.Println("Something wronge is happaned")
	}
	return nil
}
