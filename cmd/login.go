package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/urvil38/kubepaas/questions"
	"github.com/urvil38/kubepaas/userservice"
	"github.com/urvil38/kubepaas/util"
	"gopkg.in/AlecAivazis/survey.v1"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to kubepaas platform",
	Run: func(cmd *cobra.Command, args []string) {

		//checking if user is alreay logged in
		if util.ConfigFileExists() {
			fmt.Println("You are already loggin to kubepaas platform")
			os.Exit(0)
		}

		//there are two ways user can provide auth credentials
		//1.using --email and --password flags
		//2.using survey prompt (using survey package)
		//if email && password are not empty than we handle using method => (1) otherwise method => (2).
		var auth questions.AuthCredential

		emailFlagValue := cmd.Flags().Lookup("email").Value.String()
		passwordFlagValue := cmd.Flags().Lookup("password").Value.String()
		if emailFlagValue != "" || passwordFlagValue != "" {
			//method 1
			err := parseFlags(emailFlagValue, passwordFlagValue)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			auth = questions.AuthCredential{Email: emailFlagValue, Password: passwordFlagValue}
		} else {
			//method 2
			var err error
			auth, err = prompForUserLogin()
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		}

		//contact with userservice in order to login
		err := userservice.Login(auth)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func prompForUserLogin() (basicAuth questions.AuthCredential, err error) {
	if err := survey.Ask(questions.LoginUser, &basicAuth); err != nil {
		return questions.AuthCredential{}, err
	}
	return basicAuth, nil
}

func parseFlags(emailFlagValue, passwordFlagValue string) (err error) {
	if emailFlagValue == "" {
		return fmt.Errorf("You must Have to provide email for login.\nHELP: kubepaas login --help")
	}
	if passwordFlagValue == "" {
		return fmt.Errorf("You must Have to provide password for login.\nHELP: kubepaas login --help")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(loginCmd)

	//attach two flags on login cmd.
	//i.e. kubepaas login --email=<email> --password=<password>
	loginCmd.Flags().String("email", "", "Email address of accout you want to login")
	loginCmd.Flags().String("password", "", "Password of accout you want to login")
}
