package cmd

import (
	"log"
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"encoding/json"

	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/urvil38/kubepaas/storage"
	"github.com/urvil38/kubepaas/archive"
	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/generator"
	"github.com/urvil38/kubepaas/cloudbuild"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the application to kubepaas platform",
	Long: `Using deploy commnad you can deploy your code to kubepaas platform.
It require app.yaml file to be in your current directory where you running kubepaas deploy command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !Login {
			fmt.Println("Login or Signup in order to deploy your app to kubepaas")
			return
		}
		exists := config.CheckAppConfigFileExists()
		if !exists {
			os.Exit(0)
		}

		var appConfig config.AppConfig
		appConfig,err := config.ParseAppConfigFile()
		if err != nil {
			fmt.Println("Error while deploying application: %v",err)
			os.Exit(0)
		}

		var currentVersion string

		currentVersion,err = generateNewVersionNumber()
		if err != nil {
			log.Fatalf("coun't generate new version number from uuid : %v",err)
		}

		var projectMetaData config.ProjectMetaData
		f,_ := os.Open(filepath.Join(PROJECT_ROOT,".project.json"))
		defer f.Close()
		b,_ := ioutil.ReadAll(f)
		_ = json.Unmarshal(b,&projectMetaData)

		if !projectMetaDataFileExist() {
			_, err := os.Create(filepath.Join(PROJECT_ROOT,".project.json"))
			if err != nil {
				fmt.Printf("Coun't create project.json file : %v\n",err)
				os.Exit(0)
			}
			
			projectMetaData.ProjectName = appConfig.ProjectName
			projectMetaData.Versions = append(projectMetaData.Versions,currentVersion)
			if len(projectMetaData.Versions) == 0 {
				projectMetaData.CurrentVersion = currentVersion
			}else{
				projectMetaData.CurrentVersion = projectMetaData.Versions[len(projectMetaData.Versions)-1]
			}
			b,err := json.Marshal(projectMetaData)
			if err != nil {
				fmt.Printf("coun't serealize config.ProjectMetaData : %v\n",err)
				os.Exit(0)
			}
			err = ioutil.WriteFile(filepath.Join(PROJECT_ROOT,".project.json"),b,600)
			if err != nil {
				fmt.Printf("coun't write to project file : %v\n",err)
				os.Exit(0)
			}
		}else{
			if CAN_UPDATE_VERSION {
				
				f,_ := os.Open(filepath.Join(PROJECT_ROOT,".project.json"))
				defer f.Close()
				b,_ := ioutil.ReadAll(f)
				_ = json.Unmarshal(b,&projectMetaData)
				
				projectMetaData.Versions = append(projectMetaData.Versions,currentVersion)
				if len(projectMetaData.Versions) == 0 {
					projectMetaData.CurrentVersion = currentVersion
				}else{
					projectMetaData.CurrentVersion = projectMetaData.Versions[len(projectMetaData.Versions)-1]
				}
				b,err := json.Marshal(projectMetaData)
				if err != nil {
					fmt.Printf("coun't serealize config.ProjectMetaData : %v\n",err)
					os.Exit(0)
				}
				err = ioutil.WriteFile(filepath.Join(PROJECT_ROOT,".project.json"),b,600)
				if err != nil {
					fmt.Printf("coun't write to project file : %v\n",err)
					os.Exit(0)
				}
			}
		}


		err = generator.GenerateDockerFile(appConfig)
		if err != nil {
			fmt.Print(err)
		}

		err = generator.GenerateCloudBuildFile(projectMetaData)
		if err != nil {
			fmt.Print(err)
		}
		tarFilePath, err := generateTarFolder(currentVersion)
		if err != nil {
			fmt.Printf("Unable to create tar folder :%v", err.Error())
		}
		err = uploadObjectToGCS(tarFilePath,projectMetaData.ProjectName,currentVersion)
		if err != nil {
			fmt.Printf("Error while Uploding File :%v\n", err.Error())
		}
		err = cloudbuild.CreateNewBuild("kubepaas")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	},
}

func uploadObjectToGCS(source string,projectName string,currentVersion string) error {
	
	bucketName := "staging-kubepaas-ml"
	
	uploadObject := storage.NewUploadObject(source, projectName+"/"+projectName+"-"+currentVersion+".tgz", bucketName)
	return uploadObject.UploadTarBallToGCS()
}

func generateTarFolder(currentVersion string) (path string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Couldn't Find current working directory becauser of %v\n", err)
	}
	temp := os.TempDir()
	temptar := filepath.Join(temp, filepath.Base(wd))
	
	temptar = temptar + "-" + currentVersion
	targetPath, err := archive.MakeTarBall(wd, temptar)
	if err != nil {
		return "", err
	}
	return targetPath, nil
}

func projectMetaDataFileExist() bool {
	if _, err := os.Stat(filepath.Join(PROJECT_ROOT,".project.json")); err != nil {
		return false
	}
	return true
}

func generateNewVersionNumber() (string,error) {
	uuid,err := uuid.NewV4()
	if err != nil {
		return "",err
	}
	return uuid.String(),nil
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
