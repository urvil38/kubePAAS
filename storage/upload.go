package storage

import (
	"path/filepath"
	"time"
	"github.com/briandowns/spinner"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

//Size of maximum upload file
//i.e. upload object's size should not be more than 20MB.
const (
	KB      = 1 << 10
	MB      = 1 << 20
	maxSize = 20 * MB
)

type uploadObject struct {
	source      string
	destination string
	bucketName  string
}

func NewUploadObject(sourcePath, destinationPath, bucketName string) *uploadObject {
	return &uploadObject{
		source:      sourcePath,
		destination: destinationPath,
		bucketName:  bucketName,
	}
}

func (u *uploadObject) UploadTarBallToGCS() error {
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

	fileInfo, err := reader.Stat()
	if err != nil {
		return err
	}

	b := fileInfo.Size()
	if b > maxSize {
		_ = removeSourceFile(u.source)
		return fmt.Errorf("You can't upload file which is more than 20MB,Your object size is %.2fMB", float64(b)/MB)
	}
	
	s := spinner.New(spinner.CharSets[11],100*time.Millisecond)
	s.Color("yellow","bold")
	s.Suffix = fmt.Sprintf(" Uploading: %.2f KB of source code file ", float64(b)/KB)
	s.Start()

	_, err = io.Copy(writer, reader)
	if err != nil {
		s.Stop()
		return err
	}
	err = writer.Close()
	if err != nil {
		s.Stop()
		_ = removeSourceFile(u.source)
		return fmt.Errorf(" \x1b[31mPlease check your internet connection ℹ\x1b[0m")
	}
	s.Stop()
	fmt.Println(color.HiGreenString("Successfully uploaded %s file ✔︎",filepath.Base(u.source)))
	err = removeSourceFile(u.source)
	if err != nil {
		return err
	}
	return nil
}

func removeSourceFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
