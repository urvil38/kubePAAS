package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/urvil38/kubepaas/util"
	"github.com/urvil38/kubepaas/storageutil"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the application to kubepaas platform",
	Long: `Using deploy commnad you can deploy your code to kubepaas platform.
It require app.yaml file to be in your current directory where you running kubepaas deploy command.`,
	Run: func(cmd *cobra.Command, args []string) {
		exists := checkConfigFileExists()
		if !exists {
			return
		}
		tarFilePath, err := generateTarFolder()
		if err != nil {
			fmt.Printf("Unable to create zip folder :%v", err.Error())
		}
		err = uploadFile(tarFilePath)
		if err != nil {
			fmt.Printf("Error while Uploding File:%v", err.Error())
		}
	},
}

func uploadFile(source string) error {
	bucketName := "staging-kubepaas-ml"
	destination := filepath.Base(source)
	uploadObject := storageutil.NewUploadObject(source,destination,bucketName)
	return uploadObject.Upload()
}

func checkConfigFileExists() bool {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Couldn't Find current working directory beacause of : %v\n", err)
	}

	if _, err := os.Stat(filepath.Join(wd, "app.yaml")); err != nil {
		if _, err := os.Stat(filepath.Join(wd, "app.yml")); err != nil {
			fmt.Printf("No app.yaml file exist. Make sure you have app.yaml file in current project\n")
			return false
		}
	}
	return true
}

func generateTarFolder() (path string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Couldn't Find current working directory becauser of %v\n", err)
	}
	temp := os.TempDir()
	temptar := filepath.Join(temp, filepath.Base(wd))
	targetPath, err := util.Tarit(wd, temptar)
	if err != nil {
		return "", err
	}
	return targetPath, nil
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
