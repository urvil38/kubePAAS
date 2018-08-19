package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/urvil38/kubepaas/questions"
	"github.com/urvil38/kubepaas/types"
	"github.com/urvil38/kubepaas/userservice"
	"gopkg.in/AlecAivazis/survey.v1"
)

// signupCmd represents the signup command
var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Sign up for kubepaas platform",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := prompForRegisterUser()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	},
}

func prompForRegisterUser() error {
	var user types.UserInfo
	if err := survey.Ask(questions.RegisterUser, &user); err != nil {
		return err
	}
	err := userservice.RegisterUser(user)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(signupCmd)
}
