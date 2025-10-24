package cmd

import (
	"compress/zlib"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fenpaws/zc/internal"
	"github.com/spf13/cobra"
)

var compressionLevel int

func init() {
	rootCmd.AddCommand(compressCmd)
	compressCmd.Flags().IntVarP(&compressionLevel, "level", "l", zlib.DefaultCompression, "Compression level (0-9)")
}

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:   "compress [source file/folder] [destination file]",
	Short: "Compress a file or directory using zlib",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourcePath := args[0]
		destFileName := args[1]

		return compress(sourcePath, destFileName)
	},
}

func compress(sourcePath, destFileName string) error {
	destFile, err := os.Create(destFileName)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	writer, err := zlib.NewWriterLevel(destFile, compressionLevel)
	if err != nil {
		return fmt.Errorf("failed to create zlib writer: %v", err)
	}
	defer writer.Close()

	fileCount, totalSize, err := utils.CountFilesAndSize(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to compute size: %v", err)
	}

	bar := utils.CreateProgressBar(totalSize, "compressing")

	if err := filepath.Walk(sourcePath, utils.CreateWalkFunction(writer, bar, sourcePath)); err != nil {
		return fmt.Errorf("failed to walk source path: %v", err)
	}

	fmt.Printf("Successfully compressed %d files into %s\n", fileCount, destFileName)
	return nil
}
