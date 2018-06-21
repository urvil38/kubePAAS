package utils

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"io"
)

//Zipit zip the folder given by source path and place it at target path
func Zipit(source, target string) (path string ,err error) {
	if filepath.Ext(target) != "zip" || filepath.Ext(target) == "" {
		target = target + ".zip"
	}
	zipfile, err := os.Create(target)
	if err != nil {
		return "",err
	}
	defer zipfile.Close()

	zw := zip.NewWriter(zipfile)
	defer zw.Close()

	filepath.Walk(source, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name()[0] == '.' || info.IsDir() && info.Name() == "node_modules" {
			return filepath.SkipDir
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		
		header.Name = strings.TrimPrefix(strings.Replace(file,source,"",-1),string(filepath.Separator))

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		w, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(w, f)
		return err
	})

	return target,nil
}