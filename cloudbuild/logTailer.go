package cloudbuild

import (
	"io/ioutil"
	"strings"
	"context"
	"fmt"
)

type logTailer struct {
	logBucket string
	logObject string
}

var cursor int64

const (
	gcsURLPattern = `https://www.googleapis.com/storage/v1/b/%s/o/%s?alt=media`
)

func NewLogTailer(logBucket, logObject string) *logTailer {
	return &logTailer{
		logBucket: logBucket,
		logObject: logObject,
	}
}

func (l logTailer) Poll(isLast bool) {
	client,err := getCloudBuildLogStorageClient()
	if err != nil {
		fmt.Print(err)
	}

	res,err := client.Bucket(l.logBucket).Object(l.logObject).NewRangeReader(context.Background(),cursor,-1)
	if err != nil {
		return
	}

	if res != nil {
		b,err := ioutil.ReadAll(res)
		if err != nil {
			fmt.Print(err)
		}

		if cursor == 0 {
			printFirstLine()
		}
		cursor += int64(len(b))
		fmt.Print(strings.TrimRight(string(b),"\n"))
		if isLast {
			printLastLine()
		}
	}
}

func printLastLine() {
	fmt.Println(" START OF CLOUDBUILD ")
}

func printFirstLine() {
	fmt.Println(" END OF CLOUDBUILD ")
}