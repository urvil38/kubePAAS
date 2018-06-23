package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//Tarit tar the folder given by source path and place it at target path
func Tarit(source, target string) (path string, err error) {

	if _, err := os.Stat(source); err != nil {
		return "", fmt.Errorf("Unable to tar file - %v", err.Error())
	}

	if filepath.Ext(target) != "tgz" || filepath.Ext(target) == "" {
		target = target + ".tgz"
	}
	writer, err := os.Create(target)
	if err != nil {
		return "", fmt.Errorf("Unable to tar file - %v", err.Error())
	}

	gzw := gzip.NewWriter(writer)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		if fi.IsDir() && fi.Name()[0] == '.' || fi.IsDir() && fi.Name() == "node_modules" {
			return filepath.SkipDir
		}

		header.Name = strings.TrimPrefix(strings.Replace(file, source, "", -1), string(filepath.Separator))

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		f, err := os.Open(file)
		defer f.Close()
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
		return nil

	})
	return target, nil
}
