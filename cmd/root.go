package cmd

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/urvil38/kubepaas/config"

	"github.com/spf13/cobra"
	"github.com/urvil38/kubepaas/util"
)

const (
	Banner = `

██╗  ██╗██╗   ██╗██████╗ ███████╗██████╗  █████╗  █████╗ ███████╗
██║ ██╔╝██║   ██║██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔══██╗██╔════╝
█████╔╝ ██║   ██║██████╔╝█████╗  ██████╔╝███████║███████║███████╗
██╔═██╗ ██║   ██║██╔══██╗██╔══╝  ██╔═══╝ ██╔══██║██╔══██║╚════██║
██║  ██╗╚██████╔╝██████╔╝███████╗██║     ██║  ██║██║  ██║███████║
╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚══════╝╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝

`
)

var (
	ConfigValue config.AuthConfig
	Login       bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubepaas",
	Short: "A CLI for interacting with kubepaas platform",
	Long: `A tool for interacting with kubepaas platform 
and used for all kind of command that This plateform will support`,
	Run: func(cmd *cobra.Command, args []string) {
		printBanner()
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func printBanner() {
	rand.Seed(time.Now().UnixNano())
	colorCounter := rand.Intn(7)
	fmt.Printf("\x1b[3%dm%v\x1b[0m", colorCounter+1, Banner)
}

func init() {

	err := os.MkdirAll(util.GetConfigFolderPath(), 0777)
	if err != nil {
		fmt.Printf("Unable to create config Folder: %v", err.Error())
		os.Exit(1)
	}

	configFilePath, _ := util.GetConfigFilePath()

	if _, err := os.Stat(configFilePath); err == nil {
		Login = true
	}

	if util.ConfigFileExists() {
		confFileName, _ := util.GetConfigFilePath()

		b, err := ioutil.ReadFile(confFileName)
		if err != nil {
			return
		}
		str := strings.Split(string(b), "\n")

		ConfigValue.AuthToken.Token = str[0]
		ConfigValue.UserConfig.Email = str[1]
		ConfigValue.UserConfig.ID = str[2]
		ConfigValue.UserConfig.Name = str[3]
	}
}
