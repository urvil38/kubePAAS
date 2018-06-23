package storageutil

import (
	"io"
	"os"
	"path/filepath"
	"context"
	"fmt"
)

type uploadObject struct{
	source string
	destination string
	bucketName string
}

func NewUploadObject(sourcePath,destinationPath,bucketName string) *uploadObject {
	return &uploadObject{
		source:sourcePath,
		destination: destinationPath,
		bucketName: bucketName,
	}
}

func (u *uploadObject) Upload() error {
	client, err := GetStorageClient()
	if err != nil {
		return err
	}

	writer := client.Bucket(u.bucketName).Object(filepath.Base(u.destination)).NewWriter(context.Background())

	reader, err := os.Open(u.source)
	if err != nil {
		return err
	}
	defer reader.Close()
	if err != nil {
		return err
	}
	b, err := io.Copy(writer, reader)
	if err != nil {
		return err
	}
	fmt.Printf("Bytes writtern: %.2f KB\n", float64(b)/1024)
	err = writer.Close()
	if err != nil {
		if err != nil {
			return err
		}
	}
	fmt.Printf("Successfully uploaded file at path: %v\n", u.source)
	err = os.Remove(u.source)
	if err != nil {
		return err
	}
	return nil
}