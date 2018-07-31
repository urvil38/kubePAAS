package cmd

import (
	"os"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/urvil38/kubepaas/util"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from kubepaas platform",
	Run: func(cmd *cobra.Command, args []string) {
		err := logOut()
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func logOut() error {
	configFilePath := util.GetConfigFilePath()
	if _,err := os.Stat(configFilePath) ; err != nil {
		return fmt.Errorf("You are not logged in")
	}

	err := os.Remove(configFilePath)
	if err != nil {
		return fmt.Errorf("Unable to logged you out")
	}
	fmt.Println("Successfully logged out")
	return nil
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
