package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/util"
)

const (
	setCmdUsage   = "set [PROPERTY] [VALUE]"
	unsetCmdUsage = "unset [PROPERTY]"
)

var (
	setCmdErr   = fmt.Sprintf("expected %s", setCmdUsage)
	unsetCmdErr = fmt.Sprintf("expected %s", unsetCmdUsage)
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "view and edit Kubepaas CLI properties",
	Long: `The kubepaas config command group lets you set, view and unset properties
	used by kubepaas CLI.

	A configuration is a set of properties that govern the behavior of kubepaas CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var setConfigCmd = &cobra.Command{
	Use:   setCmdUsage,
	Short: "Set a Kubepaas CLI property",
	Run: func(cmd *cobra.Command, args []string) {
		err := validate(args, 2, setCmdErr)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		err = config.CLIConf.Set(args[0], args[1])
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

var unsetConfigCmd = &cobra.Command{
	Use:   unsetCmdUsage,
	Short: "Unset a Kubepaas CLI property",
	Run: func(cmd *cobra.Command, args []string) {
		err := validate(args, 1, unsetCmdErr)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		err = config.CLIConf.Unset(args[0])
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

var listConfigCmd = &cobra.Command{
	Use:   "list",
	Short: "list a Kubepaas CLI config",
	Run: func(cmd *cobra.Command, args []string) {
		fp, err := util.GetConfigFilePath()
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		b, err := ioutil.ReadFile(fp)
		if errors.Is(err, os.ErrNotExist) {
			cmd.PrintErrln(".config.json file not found!")
		}

		fmt.Println(string(b))
	},
}

func validate(args []string, argCount int, usageErr string) error {
	if len(args) != argCount {
		return errors.New(usageErr)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setConfigCmd)
	configCmd.AddCommand(unsetConfigCmd)
	configCmd.AddCommand(listConfigCmd)
}
