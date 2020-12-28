package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/urvil38/kubepaas/archive"
	"github.com/urvil38/kubepaas/banner"
	"github.com/urvil38/kubepaas/cloudbuild"
	"github.com/urvil38/kubepaas/config"
	"github.com/urvil38/kubepaas/generator"
	"github.com/urvil38/kubepaas/storage"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the application to kubepaas platform",
	Long: `Using deploy commnad you can deploy your code to kubepaas platform.
It require app.yaml file to be in your current directory where you running kubepaas deploy command.`,
	Run: func(cmd *cobra.Command, args []string) {
		var canUpdate bool

		// if !Login {
		// 	fmt.Println("Login or Signup in order to deploy your app to kubepaas")
		// 	return
		// }
		exists, err := config.CheckAppConfigFileExists()
		if !exists {
			fmt.Println(err)
			os.Exit(0)
		}

		err = os.MkdirAll(config.KubeConfig.KubepaasRoot, 0777)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		canUpdateFlag, err := cmd.Flags().GetBool("update")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		if canUpdateFlag {
			canUpdate = true
		}

		fmt.Println(banner.PrintDeployingMessage())

		appConfig, err := config.ParseAppConfigFile()
		if err != nil {
			fmt.Printf("Error while deploying application: %v", err)
			os.Exit(0)
		}

		var currentVersion string

		currentVersion, err = generateNewVersionNumber()
		if err != nil {
			log.Fatalf("coun't generate new version number from uuid : %v", err)
		}

		var projectMetaData config.ProjectMetaData

		err = generator.KmanagerConf(&projectMetaData)
		if err != nil {
			fmt.Printf("Couldn't Get kubepaas cluster data: %v \n", err)
			os.Exit(1)
		}
		if config.ProjectMetaDataFileExist() {
			f, _ := os.Open(filepath.Join(config.KubeConfig.KubepaasRoot, ".project.json"))
			defer f.Close()
			b, _ := ioutil.ReadAll(f)
			_ = json.Unmarshal(b, &projectMetaData)
		}

		if !config.ProjectMetaDataFileExist() {
			_, err := os.Create(filepath.Join(config.KubeConfig.KubepaasRoot, ".project.json"))
			if err != nil {
				fmt.Printf("Coun't create project.json file : %v\n", err)
				os.Exit(1)
			}

			projectMetaData.ProjectName = appConfig.Metadata.Name
			projectMetaData.Versions = append(projectMetaData.Versions, currentVersion)
			if len(projectMetaData.Versions) == 0 {
				projectMetaData.CurrentVersion = currentVersion
			} else {
				projectMetaData.CurrentVersion = projectMetaData.Versions[len(projectMetaData.Versions)-1]
			}
			err = writeToProjectMetaDataFile(projectMetaData)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			if canUpdate {

				projectMetaData.Versions = append(projectMetaData.Versions, currentVersion)
				if len(projectMetaData.Versions) == 0 {
					projectMetaData.CurrentVersion = currentVersion
				} else {
					projectMetaData.CurrentVersion = projectMetaData.Versions[len(projectMetaData.Versions)-1]
				}
				err := writeToProjectMetaDataFile(projectMetaData)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}
		fmt.Println(banner.PrintProjectInfo(*appConfig, projectMetaData))

		fmt.Println(banner.PrintDockerfileMessage())
		err = generator.GenerateDockerFile(*appConfig)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		fmt.Println(banner.SuccessDockerfileMessage())

		fmt.Println(banner.PrintCloudBuildMessage())
		err = generator.GenerateDockerCloudBuildFile(projectMetaData, *appConfig)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		fmt.Println(banner.SuccessDockerCloudbuildMessage())

		err = generator.GenerateKubernetesCloudBuildFile(projectMetaData)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		fmt.Println(banner.SuccessKubernetesCloudbuildMessage())

		fmt.Println(banner.PrintUploadSourceCodeMessage())
		tarFilePath, err := generateTarBallFromSourceCode(currentVersion)
		if err != nil {
			fmt.Printf("Unable to create tar folder :%v", err.Error())
		}
		err = uploadSourceCodeToGCS(projectMetaData.SourceCodeBucket, tarFilePath, projectMetaData.ProjectName, currentVersion)
		if err != nil {
			fmt.Printf("Error while Uploding File :%v\n", err.Error())
		}

		err = cloudbuild.CreateNewBuild(projectMetaData.GCPProject, "docker")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(banner.PrintKubernetesMessage())
		err = generator.GenerateKubernetesConfig(*appConfig, projectMetaData)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		fmt.Println(banner.SuccessKubernetesMessage())

		fmt.Println(banner.PrintUploadKubernetesMessage())
		tarFilePath, err = generateTarBallFromKubernetes(currentVersion)
		if err != nil {
			fmt.Printf("Unable to create tar folder: %v", err)
		}
		err = uploadKubernetesConfigToGCS(projectMetaData.SourceCodeBucket, tarFilePath, projectMetaData.ProjectName, currentVersion)
		if err != nil {
			fmt.Printf("Error while Uploding File :%v\n", err.Error())
		}

		err = cloudbuild.CreateNewBuild(projectMetaData.GCPProject, "kubernetes")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func writeToProjectMetaDataFile(projectMetaData interface{}) error {
	b, err := json.Marshal(projectMetaData)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(config.KubeConfig.KubepaasRoot, ".project.json"), b, 600)
	if err != nil {
		return err
	}
	return nil
}

func uploadSourceCodeToGCS(bucket string, source string, projectName string, currentVersion string) error {
	uploadObject := storage.NewUploadObject(source, projectName+"/"+projectName+"-"+currentVersion+".tgz", bucket)
	return uploadObject.UploadTarBallToGCS()
}

func uploadKubernetesConfigToGCS(bucket string, source string, projectName string, currentVersion string) error {
	uploadObjectFormatString := `%s/kubernetes-%s-%s.tgz`
	uploadObject := storage.NewUploadObject(source, fmt.Sprintf(uploadObjectFormatString, projectName, projectName, currentVersion), bucket)
	return uploadObject.UploadTarBallToGCS()
}

func generateTarBallFromSourceCode(currentVersion string) (path string, err error) {
	temp := os.TempDir()
	tempTarBallPath := filepath.Join(temp, filepath.Base(config.KubeConfig.ProjectRoot))

	tempTarBallPath = tempTarBallPath + "-" + currentVersion
	targetPath, err := archive.MakeTarBall(config.KubeConfig.ProjectRoot, tempTarBallPath)
	if err != nil {
		return "", err
	}
	return targetPath, nil
}

func generateTarBallFromKubernetes(currentVersion string) (path string, err error) {
	if _, err := os.Stat(filepath.Join(config.KubeConfig.KubepaasRoot, "kubernetes", "kubernetes.yaml")); os.IsNotExist(err) {
		return "", err
	}
	temp := os.TempDir()
	tempTarBallPath := filepath.Join(temp, filepath.Base(config.KubeConfig.ProjectRoot))
	tempTarBallPath = tempTarBallPath + "-" + "kubernetes" + "-" + currentVersion
	targetPath, err := archive.MakeTarBall(filepath.Join(config.KubeConfig.KubepaasRoot, "kubernetes"), tempTarBallPath)
	if err != nil {
		return "", err
	}
	return targetPath, nil
}

func generateNewVersionNumber() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().Bool("update", false, "specify that kubepaas will generate new tarBall with newly generated version number")
}
