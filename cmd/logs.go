package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/http/client"
)

type LogOptions struct {
	Follow        bool
	Name          string
	tail          int
	since         time.Duration
	containerName string
}

func newLogOptions() *LogOptions {
	return &LogOptions{}
}

func newLogCmd() *cobra.Command {
	o := newLogOptions()

	logCmd := &cobra.Command{
		Use:   "logs",
		Short: "Print the logs of current running containers",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := o.GetLogs(os.Stdout)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		},
	}

	o.AddFlags(logCmd)

	return logCmd
}

func (o *LogOptions) GetLogs(w io.Writer) error {
	c := client.NewHTTPClient(nil)

	o.Name = config.KAppConfig.Metadata.Name

	if o.Name == "" {
		return fmt.Errorf("project name doesn't provided, Please make sure that app.yml file exists on the root of the project")
	}

	u, err := url.Parse(config.CLIConf.GeneratorEndpoint + "/logs" + "/" + o.Name)
	if err != nil {
		return fmt.Errorf("unable to parse generator URL, %v", err)
	}

	q := u.Query()

	if o.tail > 0 {
		q.Set("tail", strconv.Itoa(o.tail))
	}

	if o.containerName != "" {
		q.Set("container_name", o.containerName)
	}

	if o.Follow {
		q.Set("stream", strconv.FormatBool(o.Follow))
	}

	if o.since != time.Duration(0) {
		q.Set("since", o.since.String())
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Timeout: Please check your internet connection!")
	}

	_, err = io.Copy(w, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func (o *LogOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&o.tail, "tail", "t", 30, "Lines of recent log file to display. Defaults to -1 with no selector, showing all log lines otherwise 30")
	cmd.Flags().BoolVarP(&o.Follow, "follow", "f", false, "Specify if the logs should be streamed.")
	cmd.Flags().DurationVarP(&o.since, "since", "s", time.Duration(0), "Only return logs newer than a relative duration like 5s, 2m, or 3h. Defaults to all logs.")
	cmd.Flags().StringVarP(&o.containerName, "container", "c", "", "Specify container name")
}

func init() {
	rootCmd.AddCommand(newLogCmd())
}
