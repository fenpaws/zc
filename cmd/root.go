package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zc",
	Short: "ZC: Zlib Compression/Decompression CLI Tool",
	Long: `ZC (Zlib Compression) is a command-line tool designed for compressing and decompressing files 
using the zlib format. It supports various compression levels and provides real-time progress.
Examples:
  # Compress a file with default settings
  zc compress -l 6 source.txt compressed.zlib

  # Decompress a file
  zc decompress compressed.zlib decompressed.txt
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// No persistent or local flags to set for the root command
}
