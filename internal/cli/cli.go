package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zap-tool/zap/internal/compressor"
	"github.com/zap-tool/zap/internal/extractor"
)

var (
	mode        string
	autoCleanup bool
)

func Run() error {
	rootCmd := &cobra.Command{
		Use:   "zap <file|dir>",
		Short: "Fast compression tool",
		Long:  "One-command compression and extraction with auto-cleanup",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			if _, err := os.Stat(target); os.IsNotExist(err) {
				return fmt.Errorf("target not found: %s", target)
			}
			return compress(target)
		},
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "extract <archive>",
		Short: "Extract archive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			return extractor.Extract(target)
		},
	})

	return rootCmd.Execute()
}

func compress(target string) error {
	var modeVal compressor.CompressionMode
	switch mode {
	case "ultra":
		modeVal = compressor.ModeUltra
	case "fast":
		modeVal = compressor.ModeFast
	case "zip":
		modeVal = compressor.ModeZip
	default:
		modeVal = compressor.ModeNormal
	}

	c := compressor.NewCompressor(modeVal)
	archive, err := c.CompressFile(target)
	if err != nil {
		return fmt.Errorf("compression failed: %w", err)
	}
	fmt.Printf("Created: %s\n", archive)
	return nil
}

var _ = cobra.Command{}
