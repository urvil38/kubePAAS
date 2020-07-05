package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/urvil38/kubepaas/banner"
	"github.com/urvil38/kubepaas/cloudbuild"
	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/generator"

	"github.com/spf13/cobra"
)

// rolloutCmd represents the rollout command
var rolloutCmd = &cobra.Command{
	Use:   "rollout",
	Short: "RollBack to Older or Newer Version",
	Run: func(cmd *cobra.Command, args []string) {
		if !Login {
			fmt.Println("Login or Signup in order to deploy your app to kubepaas")
			return
		}
		exists := config.CheckAppConfigFileExists() && config.ProjectMetaDataFileExist()
		if !exists {
			os.Exit(0)
		}

		appConfig, err := config.ParseAppConfigFile()
		if err != nil {
			fmt.Printf("Error while deploying application: %v", err)
			os.Exit(0)
		}

		forward, err := cmd.Flags().GetBool("forward")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		backward, err := cmd.Flags().GetBool("backward")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		history, err := cmd.Flags().GetBool("history")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		var projectMetaData config.ProjectMetaData
		if config.ProjectMetaDataFileExist() {
			f, _ := os.Open(filepath.Join(config.KubeConfig.KubepaasRoot, ".project.json"))
			defer f.Close()
			b, _ := ioutil.ReadAll(f)
			_ = json.Unmarshal(b, &projectMetaData)
		}

		if history {
			printVersionHistory(projectMetaData)
			os.Exit(0)
		}

		if len(projectMetaData.Versions) <= 1 {
			fmt.Println("There are no version history to rollback")
			os.Exit(0)
		}

		if forward && backward {
			fmt.Println("Please specify only one flag.i.e either --forward of --backward")
		}

		if !forward && !backward {
			fmt.Println("Please specify either --forward or --backward flag")
			os.Exit(0)
		}

		var currentIndex int
		for i := range projectMetaData.Versions {
			if projectMetaData.Versions[i] == projectMetaData.CurrentVersion {
				currentIndex = i
			}
		}

		var updatedIndex int
		if backward {
			updatedIndex = currentIndex - 1
			if updatedIndex >= 0 {
				projectMetaData.CurrentVersion = projectMetaData.Versions[updatedIndex]
			} else {
				fmt.Println("Unable to rollback backward as there are no version available.")
				printVersionHistory(projectMetaData)
				os.Exit(0)
			}
		}

		if forward {
			updatedIndex = currentIndex + 1
			if updatedIndex <= len(projectMetaData.Versions)-1 {
				projectMetaData.CurrentVersion = projectMetaData.Versions[updatedIndex]
			} else {
				fmt.Println("Unable to rollback forward as there are no version available.")
				printVersionHistory(projectMetaData)
				os.Exit(0)
			}
		}

		err = writeToProjectMetaDataFile(projectMetaData)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		fmt.Println(banner.PrintKubernetesUpdateMessage())

		err = generator.GenerateKubernetesCloudBuildFile(projectMetaData)
		if err != nil {
			fmt.Print(err)
			os.Exit(0)
		}
		fmt.Println(banner.SuccessUpdateKubernetesCloudbuildMessage())

		err = generator.GenerateKubernetesConfig(*appConfig, projectMetaData)
		if err != nil {
			fmt.Print(err)
			os.Exit(0)
		}
		fmt.Println(banner.SuccessUpdateKubernetesMessage())

		err = cloudbuild.CreateNewBuild("kubepaas-261611", "kubernetes")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	},
}

const (
	connector = "︎︎⬇︎"
)

func printVersionHistory(projectMetaData config.ProjectMetaData) {
	fmt.Println("          VERSION HISTORY        ")
	fmt.Println("          ═══════════════        ")
	fmt.Printf("\t\t%s\n", connector)
	for i, version := range projectMetaData.Versions {
		if i == len(projectMetaData.Versions)-1 {
			if version == projectMetaData.CurrentVersion {
				fmt.Printf("%d ⟹   %s\n", i+1, green(version))
			} else {
				fmt.Printf("%d ⟹   %s\n", i+1, version)
			}
		} else {
			if version == projectMetaData.CurrentVersion {
				fmt.Printf("%d ⟹   %s\n\t\t%s\n", i+1, green(version), connector)
			} else {
				fmt.Printf("%d ⟹   %s\n\t\t%s\n", i+1, version, connector)
			}
		}
	}
}

func green(s string) string {
	return color.HiGreenString(s+"\t%s", "⬅︎  CURRENT-VERSION")
}

func init() {
	rootCmd.AddCommand(rolloutCmd)
	rolloutCmd.Flags().BoolP("forward", "f", false, "Increment Current version")
	rolloutCmd.Flags().BoolP("backward", "b", false, "Decrement Current Version")
	rolloutCmd.Flags().BoolP("history", "", false, "Shows version history")
}
