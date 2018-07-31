package cmd

import (
	"context"
	"go.opencensus.io/trace"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/urvil38/kubepaas/authservice"
	"github.com/urvil38/kubepaas/questions"
	"github.com/urvil38/kubepaas/types"
	"gopkg.in/AlecAivazis/survey.v1"
	"github.com/urvil38/kubepaas/util"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to kubepaas platform",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if util.ConfigFileExists() {
			fmt.Println("You are already loggin to kubepaas platform")
			os.Exit(0)
		}
		auth := types.AuthCredential{}
		if cmd.Flags().Lookup("email").Value.String() != "" || cmd.Flags().Lookup("password").Value.String() != "" {
			var err error
			auth,err = parseFlags(cmd)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		}else{
			var err error
			auth,err = prompForUserLogin()
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		}

		ctx,span := trace.StartSpan(context.Background(),"login")
		defer span.End()

		err := authservice.Login(ctx,auth)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func prompForUserLogin() (basicAuth types.AuthCredential,err error) {
	if err := survey.Ask(questions.LoginUserQuestion, &basicAuth); err != nil {
		return types.AuthCredential{},err
	}
	return basicAuth,nil
}

func parseFlags(cmd *cobra.Command) (basicAuth types.AuthCredential,err error) {
		email := cmd.Flags().Lookup("email").Value.String()
		if email == "" {
			return basicAuth,fmt.Errorf("You must Have to provide email for login")
		}
		password := cmd.Flags().Lookup("password").Value.String()
		if password == "" {
			return basicAuth,fmt.Errorf("You must Have to provide password for login")
		}
		return types.AuthCredential{Email:email,Password:password},nil
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().String("email", "", "Email address of accout you want to login")
	loginCmd.Flags().String("password","","Password of accout you want to login")
}
