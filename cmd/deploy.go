package cmd

import (
	"strings"
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
	fileName := filepath.Base(source)
	folderName := strings.Split(fileName,".")[0]
	uploadObject := storageutil.CreateUploadObject(source,folderName+"/"+fileName,bucketName)
	return uploadObject.Upload()
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
