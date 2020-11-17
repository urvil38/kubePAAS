package authservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/http/client"
)

func getUserProfile(conf config.AuthConfig) (userConf config.UserConfig, err error) {
	var userConfig config.UserConfig

	timeout := 15 * time.Second
	client := client.NewHTTPClient(&timeout)

	req, err := http.NewRequest("GET", config.CLIConf.AuthEndpoint+"/user/"+conf.Email, nil)
	req.Header.Add("x-access-token", conf.AuthToken.Token)
	res, err := client.Do(req)
	if err != nil {
		return userConfig, errors.New("Unable to get Profile.Check your internet connection")
	}

	switch res.StatusCode {
	case http.StatusUnauthorized:
		return userConfig, fmt.Errorf("You are not authorized")
	case http.StatusInternalServerError:
		return userConfig, fmt.Errorf("Server is not available right now, Try again in 30 seconds. Sorry for inconvience")
	case http.StatusOK:
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return userConfig, fmt.Errorf("Unable to read response body: %v", err.Error())
		}
		res.Body.Close()
		err = json.Unmarshal(b, &userConfig)
		if err != nil {
			return userConfig, fmt.Errorf("Unable to Unmarshal profile response: %v", err.Error())
		}
	}
	return userConfig, nil
}
