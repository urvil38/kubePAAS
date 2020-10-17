package authservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/http/client"
	"github.com/urvil38/spinner"
)

func ChangePassword(pass config.ChangePassword, authToken, email string) error {
	timeout := 15 * time.Second
	client := client.NewHTTPClient(&timeout)

	b, err := json.Marshal(pass)
	if err != nil {
		return fmt.Errorf("Unable to marshal password Struct: %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf(authserviceEndpoint, "user/"+email+"/password"), bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-access-token", authToken)

	s := spinner.New("Changing Password")
	s.Start()

	res, err := client.Do(req)
	if err != nil {
		s.Stop()
		return fmt.Errorf("Unable to Change Password.Check Internet Connection")
	}

	if res.Body != nil {
		s.Stop()
		res.Body.Close()
	}

	s.Stop()
	switch res.StatusCode {
	case http.StatusUnauthorized:
		fmt.Println("Invalid Current Password")
	case http.StatusBadRequest:
		fmt.Println("Invalid Request")
	case http.StatusOK:
		fmt.Println("Password changed successfully")
	default:
		fmt.Println("Server Error!!")
	}
	return nil
}
