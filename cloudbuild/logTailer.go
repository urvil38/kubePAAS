package cloudbuild

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/urvil38/kubepaas/banner"
)

type logTailer struct {
	logBucket string
	logObject string
}

func NewLogTailer(logBucket, logObject string) *logTailer {
	return &logTailer{
		logBucket: logBucket,
		logObject: logObject,
	}
}
var cursor = 0
func (l logTailer) Poll(isLast bool) {
	if isLast {
		cursor = 0
	}

	client, err := getCloudBuildLogStorageClient()
	if err != nil {
		fmt.Print(err)
	}

	res, err := client.Bucket(l.logBucket).Object(l.logObject).NewRangeReader(context.Background(), int64(cursor), -1)
	if err != nil {
		return
	}

	if res != nil {
		b, err := ioutil.ReadAll(res)
		defer res.Close()
		if err != nil {
			fmt.Print(err)
		}

		if cursor == 0 {
			printFirstLine()
		}
		cursor += len(b)
		fmt.Print(strings.TrimRight(string(b), "\n"))
	}

	if isLast {
		printLastLine()
	}
}

func printFirstLine() {
	fmt.Println(banner.StartCloudBuildLog())
}

func printLastLine() {
	fmt.Println(banner.EndCloudBuildLog())
}
