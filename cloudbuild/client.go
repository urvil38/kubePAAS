package cloudbuild

import (
	"context"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/urvil38/kubepaas/util"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/option"
)

func getCloudBuildClient() (*cloudbuild.Service, error) {

	data, err := ioutil.ReadFile(filepath.Join(util.GetConfigFolderPath(), "kubepaas-cloudbuild.json"))
	if err != nil {
		log.Fatal(err)
	}

	config, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Fatal(err)
	}

	httpClient := config.Client(oauth2.NoContext)
	cloudbuildClinet, err := cloudbuild.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return cloudbuildClinet, nil
}

func getCloudBuildLogStorageClient() (storageClient *storage.Client, err error) {

	data, err := ioutil.ReadFile(filepath.Join(util.GetConfigFolderPath(), "kubepaas-storage.json"))
	if err != nil {
		log.Fatal(err)
	}
	config, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/devstorage.read_only")
	if err != nil {
		log.Fatal(err)
	}

	httpClient := config.Client(oauth2.NoContext)
	storageClient, err = storage.NewClient(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return storageClient, nil
}
