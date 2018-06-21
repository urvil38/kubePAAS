package cmd

import (
	"google.golang.org/api/option"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
	"cloud.google.com/go/storage"
	"github.com/urvil38/kubepaas/utils"
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
		tarFilePath,err := generateTarFolder()
		if err != nil {
			fmt.Printf("Unable to create zip folder :%v",err.Error())
		}
		err = uploadFile(tarFilePath)
		if err != nil {
			fmt.Printf("Error while Uploding File:%v", err.Error())
		}
	},
}

func uploadFile(source string) error {
	bucketName := "staging-kubepaas-ml"
	
	ctx := context.Background()
	httpClient := utils.CreateStorageClient()
	client,err := storage.NewClient(ctx,option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	writer := client.Bucket(bucketName).Object(filepath.Base(source)).NewWriter(ctx)
	
	reader,err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()
	if err != nil {
		return err
	}
	b,err := io.Copy(writer,reader)
	if err != nil {
		return err
	}
	fmt.Print("Bytes writtern:",b)
	err = writer.Close()
	if err != nil {
		fmt.Printf("Successfully uploaded file at path:%v",source)
		err := os.Remove(source)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkConfigFileExists() bool{
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

func generateTarFolder() (path string,err error) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Couldn't Find current working directory becauser of %v\n", err)
	}
	temp := os.TempDir()
	temptar := filepath.Join(temp,filepath.Base(wd))
	targetPath,err := utils.Tarit(wd,temptar)
	if err != nil {
		return "",err
	}
	return targetPath,nil
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
