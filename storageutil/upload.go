package storageutil

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/urvil38/spinner"
)

type uploadObject struct {
	source      string
	destination string
	bucketName  string
}

func CreateUploadObject(sourcePath, destinationPath, bucketName string) *uploadObject {
	return &uploadObject{
		source:      sourcePath,
		destination: destinationPath,
		bucketName:  bucketName,
	}
}

func (u *uploadObject) Upload() error {
	client, err := GetStorageClient()
	if err != nil {
		return err
	}

	writer := client.Bucket(u.bucketName).Object(u.destination).NewWriter(context.Background())

	reader, err := os.Open(u.source)
	if err != nil {
		return err
	}
	defer reader.Close()
	if err != nil {
		return err
	}

	fi, err := reader.Stat()
	if err != nil {
		return err
	}
	b := fi.Size()

	s := spinner.New(fmt.Sprintf("Uploding: %.2f KB of tar file ", float64(b)/1024))
	s.Start()
	_, err = io.Copy(writer, reader)
	if err != nil {
		s.Stop()
		return err
	}
	err = writer.Close()
	if err != nil {
		s.Stop()
		return fmt.Errorf(" \x1b[31mPlease check your internet connection ℹ\x1b[0m")
	}
	fmt.Printf("Successfully uploaded file : %v\n", u.source)
	err = os.Remove(u.source)
	if err != nil {
		return err
	}
	return nil
}
