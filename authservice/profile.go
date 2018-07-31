package authservice

import (
	"go.opencensus.io/trace"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"fmt"
	"net/http"
	"time"
	"github.com/urvil38/kubepaas/config"
)

func getProfile(ctx context.Context,conf *config.Config) error {
	ctx,span := trace.StartSpan(ctx,"getProfile")
	defer span.End()
	
	timeout := 20 * time.Second
	c := newHTTPClient(&timeout)
	
	req,err := http.NewRequest("GET",fmt.Sprintf(userserviceAPI,"user"+"/"+conf.Email),nil)
	req.Header.Add("x-access-token",conf.AuthToken.Token)
	res,err := c.Client.Do(req)
	if err != nil {
		return errors.New("Unable to register you.Check your internet connection")
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("You are not authorized")
	case http.StatusInternalServerError:
		return fmt.Errorf("Server is not available right now, Try again in 30 seconds. Sorry for inconvience")
	case http.StatusOK:
		b,err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Unable to read response body: %v",err.Error())
		}
		err = json.Unmarshal(b,&conf.UserConfig)
		if err != nil {
			return fmt.Errorf("Unable to Unmarshal profile response: %v",err.Error())
		}
	}
	return nil
}