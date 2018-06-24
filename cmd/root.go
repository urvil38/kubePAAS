package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
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
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

func checkConfigFileExists() bool {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Couldn't Find current working directory beacause of : %v\n", err)
	}

	if _, err := os.Stat(filepath.Join(wd, "app.yaml")); err != nil {
		if _, err := os.Stat(filepath.Join(wd, "app.yml")); err != nil {
			fmt.Printf("No app.yaml file exist. Make sure you have app.yaml file in current project\n")
			return false
		}
	}

	return true
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

	exists := checkConfigFileExists()
	if !exists {
		return
	}

	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubepaas.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
