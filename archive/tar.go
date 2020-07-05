package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	isDir = 1 << iota
	isSymLink
	isHardLink
)

//Tarit tar the folder given by source path and place it at target path
func MakeTarBall(source, target string) (path string, err error) {

	if _, err := os.Stat(source); err != nil {
		return "", fmt.Errorf("Unable to tar file - %v", err.Error())
	}

	if filepath.Ext(target) != "tgz" || filepath.Ext(target) == "" {
		target = target + ".tgz"
	}
	toFile, err := os.Create(target)
	if err != nil {
		return "", fmt.Errorf("Unable to tar file - %v", err.Error())
	}

	//gzip writer
	gzw := gzip.NewWriter(toFile)
	defer gzw.Close()

	//tar writer
	tgzw := tar.NewWriter(gzw)
	defer tgzw.Close()

	filepath.Walk(source, func(file string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fileInfo.IsDir() && fileInfo.Name() == ".git" || fileInfo.IsDir() && fileInfo.Name() == "node_modules" || fileInfo.IsDir() && fileInfo.Name() == "kubepaas" {
			return filepath.SkipDir
		}

		header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("cannot convert stat value to syscall.Stat_t")
		}

		//headerType for store type of header
		//   001          010            100     <- byteOrder
		//	isDir	   isSymLink     isHardLink
		//    |            |              |
		// Directory  SymbolicLink     HardLink
		var headerType byte

		// Check if file is Directory.
		if fileInfo.IsDir() {
			headerType |= isDir
		}

		// True if the file is a symlink.
		if fileInfo.Mode()&os.ModeSymlink != 0 {
			headerType |= isSymLink
		}

		// Check if file have hardlink.
		nlink := uint32(stat.Nlink)
		if nlink > 1 {
			headerType |= isHardLink
		}

		header.Name = strings.TrimPrefix(strings.Replace(file, source, "", -1), string(filepath.Separator))

		//setting header's Typeflag according to headerType
		if headerType&isSymLink == isSymLink {
			header.Typeflag = tar.TypeSymlink
		}
		if headerType&isHardLink == isHardLink {
			header.Typeflag = tar.TypeLink
		}
		if headerType&isDir == isDir {
			header.Name += "/"
			header.Typeflag = tar.TypeDir
		}

		//writing header information to tar-gzip writer
		if err := tgzw.WriteHeader(header); err != nil {
			return err
		}

		//if file have symlink or is directory we just return
		//because we can't open that file
		if headerType&isSymLink == isSymLink || headerType&isDir == isDir {
			return nil
		}

		//Open file for copy to tar-gzip writer
		f, err := os.Open(file)
		defer f.Close()
		if err != nil {
			return err
		}

		//Copy content of file to tar-gzip writer
		//This step do tar and gzip of given file f
		if _, err := io.Copy(tgzw, f); err != nil {
			return err
		}
		return nil

	})
	return target, nil
}
