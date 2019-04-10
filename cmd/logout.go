package cmd

import (
	"fmt"
	"os"

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
	configFilePath, err := util.GetConfigFilePath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(configFilePath); err != nil {
		return fmt.Errorf("You are not logged in")
	}

	err = os.Remove(configFilePath)
	if err != nil {
		return fmt.Errorf("Unable to logged you out")
	}
	fmt.Println("Successfully logged out")
	return nil
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
