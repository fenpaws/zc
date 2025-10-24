package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

// CountFilesAndSize calculates the number of files and total size in a directory.
func CountFilesAndSize(sourcePath string) (int, int64, error) {
	fileCount := 0
	totalSize := int64(0)
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fileCount++
		totalSize += info.Size()
		return nil
	})
	return fileCount, totalSize, err
}

// CreateWalkFunction returns a WalkFunc for compressing files.
func CreateWalkFunction(writer io.Writer, bar *progressbar.ProgressBar, basePath string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		headerPath := strings.ReplaceAll(relPath, string(filepath.Separator), "/")

		fmt.Fprintf(writer, "file:%s\n", headerPath)

		_, err = io.Copy(io.MultiWriter(writer, bar), file)
		if err != nil {
			return err
		}
		fmt.Fprintln(writer)
		return nil
	}
}

// CreateProgressBar initializes a progress bar for a given size.
func CreateProgressBar(totalSize int64, description string) *progressbar.ProgressBar {
	return progressbar.DefaultBytes(totalSize, description)
}
