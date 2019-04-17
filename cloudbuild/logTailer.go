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
	cursor int
}

func NewLogTailer(logBucket, logObject string) *logTailer {
	return &logTailer{
		logBucket: logBucket,
		logObject: logObject,
		cursor: 0,
	}
}

func (l *logTailer) Poll(isLast bool) {

	client, err := getCloudBuildLogStorageClient()
	if err != nil {
		fmt.Print(err)
	}

	res, err := client.Bucket(l.logBucket).Object(l.logObject).NewRangeReader(context.Background(), int64(l.cursor), -1)
	if err != nil {
		if isLast {
			printLastLine()
		}
		return
	}

	if res != nil {
		b, err := ioutil.ReadAll(res)
		defer res.Close()
		if err != nil {
			fmt.Print(err)
		}
		if l.cursor == 0 && !isLast && len(b) > 0 {
			printFirstLine()
		}
		l.cursor += len(b)
		fmt.Println(strings.TrimRight(string(b), "\n"))
		if isLast && len(b) > 0{
			printLastLine()
			return
		}
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
