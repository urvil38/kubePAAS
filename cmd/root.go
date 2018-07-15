package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	banner = `

	██╗  ██╗██╗   ██╗██████╗ ███████╗██████╗  █████╗  █████╗ ███████╗
	██║ ██╔╝██║   ██║██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔══██╗██╔════╝
	█████╔╝ ██║   ██║██████╔╝█████╗  ██████╔╝███████║███████║███████╗
	██╔═██╗ ██║   ██║██╔══██╗██╔══╝  ██╔═══╝ ██╔══██║██╔══██║╚════██║
	██║  ██╗╚██████╔╝██████╔╝███████╗██║     ██║  ██║██║  ██║███████║
	╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚══════╝╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝

`
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubepaas",
	Short: "A CLI for interacting with kubepaas platform",
	Long: `A tool for interacting with kubepaas platform 
and used for all kind of command that This plateform will support`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rand.Seed(time.Now().UnixNano())
	colorCounter := rand.Intn(7)
	fmt.Printf("\x1b[3%dm%v\x1b[0m", colorCounter+1, banner)

	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubepaas.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
