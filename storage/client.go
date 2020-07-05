package storage

import (
	"context"
	"github.com/urvil38/kubepaas/util"
	"io/ioutil"
	"log"
	"path/filepath"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func GetStorageClient() (storageClient *storage.Client, err error) {

	data, err := ioutil.ReadFile(filepath.Join(util.GetConfigFolderPath(), "kubepaas-storage.json"))
	if err != nil {
		log.Fatal(err)
	}
	config, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/devstorage.read_write")
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
