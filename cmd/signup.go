package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/urvil38/kubepaas/http/client"
	"github.com/urvil38/kubepaas/questions"
	"github.com/urvil38/kubepaas/userservice"
	"gopkg.in/AlecAivazis/survey.v1"
)

// signupCmd represents the signup command
var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Sign up for kubepaas platform",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var signupInfo questions.UserInfo

		timeout := 10 * time.Second
		client := client.NewHTTPClient(&timeout)

		err := promptForRegisterInit(client, &signupInfo)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		err = promptForRegisterFinish(client, &signupInfo)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	},
}

func promptForRegisterInit(client *http.Client, signupInfo *questions.UserInfo) error {

	if err := survey.Ask(questions.RegisterUserInit, signupInfo); err != nil {
		return err
	}

	err := userservice.RegistrationInit(client, *signupInfo)
	if err != nil {
		return err
	}
	return nil
}

func promptForRegisterFinish(client *http.Client, signupInfo *questions.UserInfo) error {

	if err := survey.Ask(questions.RegisterUserFinish, signupInfo); err != nil {
		return err
	}
	err := userservice.RegistrationFinish(client, *signupInfo)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(signupCmd)
}
