package cloudbuild

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/urvil38/kubepaas/config"
	"google.golang.org/api/cloudbuild/v1"
)

const (
	maxRetries    = 60 * 60
	retryInterval = 1 * time.Second
)

type builder struct {
	ProjectName       string
	Build             cloudbuild.Build
	CloudbuildService cloudbuild.Service
}

var errTimeout = errors.New("Operation %s timed out. This operation may still be underway")

func newBuilder(projectName, buildType string) (builder, error) {
	var build cloudbuild.Build
	if buildType == "docker" {
		cb, err := ioutil.ReadFile(filepath.Join(config.KubeConfig.KubepaasRoot, "docker-cloudbuild.json"))
		if err != nil {
			return builder{}, err
		}
		err = json.Unmarshal(cb, &build)
		if err != nil {
			return builder{}, err
		}
	}

	if buildType == "kubernetes" {
		cb, err := ioutil.ReadFile(filepath.Join(config.KubeConfig.KubepaasRoot, "kubernetes-cloudbuild.json"))
		if err != nil {
			return builder{}, err
		}
		err = json.Unmarshal(cb, &build)
		if err != nil {
			return builder{}, err
		}
	}

	cloudbuildService, _ := getCloudBuildClient()
	return builder{
		ProjectName:       projectName,
		Build:             build,
		CloudbuildService: *cloudbuildService,
	}, nil
}

func CreateNewBuild(projectName string, buildType string) error {
	builder, err := newBuilder(projectName, buildType)
	if err != nil {
		return err
	}

	ProjectBuildService := cloudbuild.NewProjectsBuildsService(&builder.CloudbuildService)
	ProjectBuildCreateCall := ProjectBuildService.Create(builder.ProjectName, &builder.Build)

	buildOp, err := ProjectBuildCreateCall.Do()
	if err != nil {
		return err
	}

	err = waitAndStreamLogs(buildOp)
	if err != nil {
		return err
	}
	return nil
}

type logConfig struct {
	buildID   string
	logBucket string
}

func waitAndStreamLogs(buildOp *cloudbuild.Operation) error {
	m, _ := buildOp.Metadata.MarshalJSON()

	var cb customBuild
	err := json.Unmarshal(m, &cb)
	if err != nil {
		fmt.Print(err)
	}

	var gcsPrefix = "gs://"
	if strings.Contains(cb.Build.LogsBucket, gcsPrefix) {
		cb.Build.LogsBucket = strings.TrimPrefix(cb.Build.LogsBucket, gcsPrefix)
	}

	lc := logConfig{
		buildID:   cb.Build.ID,
		logBucket: cb.Build.LogsBucket,
	}

	//fmt.Printf("Started cloud build [%v]\n", lc.buildID)
	cloudbuildLogfileFmtString := `log-%s.txt`

	if lc.logBucket != "" {
		logObject := fmt.Sprintf(cloudbuildLogfileFmtString, lc.buildID)
		logTailer := NewLogTailer(lc.logBucket, logObject)
		var callback func(bool)
		if logTailer != nil {
			callback = logTailer.Poll
		}
		op, err := waitForOperation(buildOp, callback)
		if err == errTimeout {
			return errTimeout
		}

		if logTailer != nil {
			logTailer.Poll(true)
		}

		finalStatus := getStatusFromOperation(op)
		if finalStatus != "SUCCESS" {
			return fmt.Errorf("cloud build failed. Faliure status : %v", finalStatus)
		}
	}
	return nil
}

func waitForOperation(buildOp *cloudbuild.Operation, callback func(bool)) (cloudbuild.Operation, error) {
	completedOperation, err := pollUntilDone(buildOp, callback)
	if err != nil {
		return cloudbuild.Operation{}, err
	}
	return completedOperation, nil
}

func getStatusFromOperation(op cloudbuild.Operation) string {
	m, _ := op.Metadata.MarshalJSON()
	var cb customBuild
	err := json.Unmarshal(m, &cb)
	if err != nil {
		fmt.Printf("Unable to unmarshal Operation struct : %v\n", err)
	}
	if cb.Build.Status != "" {
		return cb.Build.Status
	}
	return "UNKNOWN"
}

func pollUntilDone(op *cloudbuild.Operation, callback func(bool)) (cloudbuild.Operation, error) {
	if op.Done {
		return *op, nil
	}

	cloudService, _ := getCloudBuildClient()
	operationGetCall := cloudService.Operations.Get(op.Name)

	for i := 0; i < maxRetries; i++ {
		operation, err := operationGetCall.Do()
		if err != nil {
			return cloudbuild.Operation{}, err
		}
		if operation.Done {
			//fmt.Printf("Operation %v complete\n", operation.Name)
			return *operation, nil
		}
		//fmt.Printf("Operation %v not complete. Waiting %v\n",operation.Name,RETRY_INTERVAL.String())
		time.Sleep(retryInterval)
		if callback != nil {
			callback(false)
		}
	}
	return cloudbuild.Operation{}, errors.New("Unexpected errors occured in pollUntilDone() function")
}
