package storage

import (
	"context"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
)

func GetStorageClient() (storageClient *storage.Client, err error) {

	config := jwt.Config{
		PrivateKeyID: "ff5a14b807fa7c77e1110b3ee9f373ce6dfea0af",
		Email:        "storagekubepaas@kubepaas.iam.gserviceaccount.com",
		PrivateKey:   []byte(""),//private key
		Scopes:       []string{"https://www.googleapis.com/auth/devstorage.read_write"},
		TokenURL:     google.JWTTokenURL,
	}

	httpClient := config.Client(oauth2.NoContext)
	storageClient, err = storage.NewClient(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return storageClient, nil
}
