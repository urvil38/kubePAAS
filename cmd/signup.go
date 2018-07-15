package cmd

import (
	"fmt"
	"github.com/urvil38/kubepaas/authservice"
	"os"

	"github.com/spf13/cobra"
	"github.com/urvil38/kubepaas/questions"
	"github.com/urvil38/kubepaas/types"
	"gopkg.in/AlecAivazis/survey.v1"
)

// signupCmd represents the signup command
var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := prompForRegisterUser()
		if err != nil {
			fmt.Print(err)
			os.Exit(0)
		}
	},
}

func prompForRegisterUser() error {
	var user types.UserInfo
	if err := survey.Ask(questions.RegisterUserQuestion, &user); err != nil {
		return err
	}
	err := authservice.RegisterUser(user)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(signupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
