package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !Login {
			fmt.Println("Login or Signup in order to show profile")
			return
		}
		err := getProfile()
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func getProfile() error {
	fmt.Printf("UserID : %s\nName   : %v\nEmail  : %v\n", ConfigValue.ID, ConfigValue.Name, ConfigValue.Email)
	return nil
}

func init() {
	rootCmd.AddCommand(profileCmd)
}
