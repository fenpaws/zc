package cmd

import (
	"bufio"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fenpaws/zc/internal"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(decompressCmd)
}

// decompressCmd represents the decompress command.
var decompressCmd = &cobra.Command{
	Use:   "decompress [source file] [destination folder]",
	Short: "Decompress a zlib compressed file or directory",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceFileName := args[0]
		destFolder := args[1]

		return decompress(sourceFileName, destFolder)
	},
}

func decompress(sourceFileName, destFolder string) error {
	sourceFile, err := os.Open(sourceFileName)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	reader, err := zlib.NewReader(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to create zlib reader: %v", err)
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	bar := utils.CreateProgressBar(0, "decompressing") // Use 0 or estimate if size is known

	if err := processDecompression(reader, scanner, destFolder, bar); err != nil {
		return fmt.Errorf("decompression failed: %v", err)
	}

	fmt.Println("Decompression completed")
	return nil
}

func processDecompression(reader io.Reader, scanner *bufio.Scanner, destFolder string, bar *progressbar.ProgressBar) error {
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "file:") {
			relPath := strings.TrimPrefix(line, "file:")
			destPath := filepath.Join(destFolder, relPath)
			destDir := filepath.Dir(destPath)

			if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", destDir, err)
			}

			destFile, err := os.Create(destPath)
			if err != nil {
				return fmt.Errorf("failed to create destination file: %v", err)
			}

			_, err = io.Copy(io.MultiWriter(destFile, bar), reader)
			destFile.Close()

			if err != nil {
				return fmt.Errorf("failed to copy decompressed data: %v", err)
			}
		}
	}

	return scanner.Err()
}
