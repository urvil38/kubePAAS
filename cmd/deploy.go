package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/urvil38/kubepaas/storageutil"
	"github.com/urvil38/kubepaas/util"
)

var login bool

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the application to kubepaas platform",
	Long: `Using deploy commnad you can deploy your code to kubepaas platform.
It require app.yaml file to be in your current directory where you running kubepaas deploy command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !login {
			fmt.Println("Login or Signup in order to deploy your app to kubepaas")
			return
		}
		exists := checkConfigFileExists()
		if !exists {
			os.Exit(0)
		}
		tarFilePath, err := generateTarFolder()
		if err != nil {
			fmt.Printf("Unable to create tar folder :%v", err.Error())
		}
		err = uploadFile(tarFilePath)
		if err != nil {
			fmt.Printf("Error while Uploding File :%v\n", err.Error())
		}
	},
}

func uploadFile(source string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	bucketName := "staging-kubepaas-ml"
	fileName := filepath.Base(source)
	folderName := filepath.Base(wd)
	uploadObject := storageutil.NewUploadObject(source, folderName+"/"+fileName, bucketName)
	return uploadObject.Upload()
}

func generateTarFolder() (path string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Couldn't Find current working directory becauser of %v\n", err)
	}
	temp := os.TempDir()
	temptar := filepath.Join(temp, filepath.Base(wd))
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	temptar = temptar + "-" + id.String()
	targetPath, err := util.Tarit(wd, temptar)
	if err != nil {
		return "", err
	}
	return targetPath, nil
}

func checkConfigFileExists() bool {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Couldn't Find current working directory beacause of : %v\n", err)
	}

	if _, err := os.Stat(filepath.Join(wd, "app.yaml")); err != nil {
		fmt.Println("\x1b[31m✗ No app.yaml file exist. Make sure you have app.yaml file in current project ℹ\x1b[0m")
		return false
	}

	return true
}

func init() {
	configFilePath := util.GetConfigFilePath()
	if _,err := os.Stat(configFilePath) ; err == nil {
		login = true
	}
	rootCmd.AddCommand(deployCmd)
}
