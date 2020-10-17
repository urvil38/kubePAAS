package authservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urvil38/spinner"
	"net/http"

	"github.com/urvil38/kubepaas/questions"
)

func RegistrationInit(client *http.Client, signupInfo questions.UserInfo) error {
	b, err := json.Marshal(signupInfo)
	if err != nil {
		return fmt.Errorf("Couldn't marshal registration details: %v", err.Error())
	}

	res, err := client.Post(fmt.Sprintf(authserviceEndpoint, "user"), "application/json", bytes.NewReader(b))
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
		fmt.Printf("Check %v for authectication code\n", signupInfo.Email)
	case http.StatusConflict:
		return fmt.Errorf("User is already exists!")
	default:
		return fmt.Errorf("Inernal server Error!!")
	}

	return nil
}

func RegistrationFinish(client *http.Client, signupInfo questions.UserInfo) error {
	b, err := json.Marshal(signupInfo)
	if err != nil {
		return fmt.Errorf("Couldn't marshal registration details: %v", err.Error())
	}

	s := spinner.New("Registering You")
	s.Start()
	res, err := client.Post(fmt.Sprintf(authserviceEndpoint, "user"), "application/json", bytes.NewReader(b))
	if err != nil {
		s.Stop()
		return fmt.Errorf("Unable to Signup.Check Internet Connection")
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.TLS == nil {
		s.Stop()
		fmt.Println("WARNING! Communication is not secure, please consider using HTTPS. Letsencrypt.org offers free SSL/TLS certificates.")
	}

	s.Stop()
	switch res.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("Email is not valid!!")
	}
	return nil
}
