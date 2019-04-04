package cloudbuild

import (
	"strings"
	"time"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
	"google.golang.org/api/cloudbuild/v1"
)

const (
	MAX_RETRIES = 60 * 60
	RETRY_INTERVAL = 1 * time.Second
)

type builder struct {
	ProjectName string
	Build *cloudbuild.Build
	CloudbuildService *cloudbuild.Service
}

var errTimeout = errors.New("Operation %s timed out. This operation may still be underway")

func newBuilder(projectName string) (*builder,error) {
	var projectRoot , _ = os.Getwd()
	cb,err := ioutil.ReadFile(filepath.Join(projectRoot,"cloudbuild.json"))
	if err != nil {
		return &builder{},err
	}
	var build cloudbuild.Build
	err = json.Unmarshal(cb,&build)
	if err != nil {
		return &builder{},err
	}

	cloudbuildService, _ := getCloudBuildClient() 
	return &builder{
		ProjectName : projectName,
		Build : &build,
		CloudbuildService : cloudbuildService,
	},nil
}

func CreateNewBuild(projectName string) error{
	builder,err := newBuilder(projectName)
	if err != nil {
		return err
	}

	ProjectBuildService := cloudbuild.NewProjectsBuildsService(builder.CloudbuildService)
	ProjectBuildCreateCall := ProjectBuildService.Create(builder.ProjectName,builder.Build)

	buildOp,err := ProjectBuildCreateCall.Do()
	if err != nil {
		return err
	}

	waitAndStreamLogs(buildOp)
	return nil
}

type logConfig struct {
	buildID string
	logURL string
	logBucket string
}

func waitAndStreamLogs(buildOp *cloudbuild.Operation) {
	m ,_ := buildOp.Metadata.MarshalJSON()
	
	var cb customBuild
	err := json.Unmarshal(m,&cb)
	if err != nil {
		fmt.Print(err)
	}

	var gcsPrefix = "gs://"
	if strings.Contains(cb.Build.LogsBucket,gcsPrefix) {
		cb.Build.LogsBucket = strings.TrimPrefix(cb.Build.LogsBucket,gcsPrefix)
	}

	log := logConfig{
		buildID:cb.Build.ID,
		logURL:cb.Build.LogURL,
		logBucket:cb.Build.LogsBucket,
	}

	fmt.Printf("Started cloud build [%v]\n",log.buildID)
	cloudbuildLogfileFmtString := `log-%s.txt`

	if log.logBucket != "" {
		logObject := fmt.Sprintf(cloudbuildLogfileFmtString,log.buildID)
		logTailer := NewLogTailer(log.logBucket,logObject)
		var callback func(bool) 
		if logTailer != nil {
			callback = logTailer.Poll
		}
		op ,err := waitForOperation(buildOp,callback)
		if err == errTimeout {
			fmt.Println("Error: Cloud build timed out.")
			return 
		}

		if logTailer != nil {
			logTailer.Poll(true)
		}

		finalStatus := getStatusFromOperation(op)
		if finalStatus != "SUCCESS" {
			fmt.Printf("Cloud build failed. Faliure status : %v\n",finalStatus)
		}
	}
}

func waitForOperation(buildOp *cloudbuild.Operation,callback func(bool)) (cloudbuild.Operation,error){
	completedOperation, err := pollUntilDone(buildOp,callback)
	if err != nil {
		return cloudbuild.Operation{},err
	}
	return completedOperation,nil
}

func getStatusFromOperation(op cloudbuild.Operation) string {
	m ,_ := op.Metadata.MarshalJSON()
	var cb customBuild
	err := json.Unmarshal(m,&cb)
	if err != nil {
		fmt.Printf("Unable to unmarshal Operation struct : %v\n",err)
	}
	if cb.Build.Status != "" {
		return cb.Build.Status
	}
	return "UNKNOWN"
}

func pollUntilDone(op *cloudbuild.Operation,callback func(bool)) (cloudbuild.Operation,error) {
	if op.Done {
		return *op,nil
	}

	cloudService,_ := getCloudBuildClient()
	operationGetCall := cloudService.Operations.Get(op.Name)

	for i := 0 ; i < MAX_RETRIES ; i++ {
		operation,err := operationGetCall.Do()
		if err != nil {
			return cloudbuild.Operation{},err
		}
		if operation.Done {
			fmt.Printf("Operation %v complete\n",operation.Name)
			return *operation,nil
		}
		//fmt.Printf("Operation %v not complete. Waiting %v\n",operation.Name,RETRY_INTERVAL.String())
		time.Sleep(RETRY_INTERVAL)
		if callback != 
		nil {
			callback(false)
		}
	}
	return cloudbuild.Operation{},errors.New("Unexpected errors occured in pollUntilDone() function")
}



