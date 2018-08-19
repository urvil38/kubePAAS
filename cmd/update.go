package cmd

import (
	"fmt"

	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/questions"
	"github.com/urvil38/kubepaas/userservice"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Used for changing values of diffrent configuration",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "change password",
	Run: func(cmd *cobra.Command, args []string) {
		if !Login {
			fmt.Println("Login or Signup in order to change password")
			return
		}
		err := changePassword()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func changePassword() error {
	var pass config.ChangePassword
	err := survey.Ask(questions.ChangePassword, &pass)
	if err != nil {
		return err
	}

	authToken, email := ConfigValue.AuthToken.Token, ConfigValue.Email

	err = userservice.ChangePassword(pass, authToken, email)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.AddCommand(passwordCmd)
}
