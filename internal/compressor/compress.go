package compressor

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

type CompressionMode int

const (
	ModeFast CompressionMode = iota
	ModeNormal
	ModeUltra
)

type Compressor struct {
	Mode CompressionMode
}

func NewCompressor(mode CompressionMode) *Compressor {
	return &Compressor{Mode: mode}
}

func (c *Compressor) CompressFile(source string) (string, error) {
	info, err := os.Stat(source)
	if err != nil {
		return "", fmt.Errorf("stat source: %w", err)
	}

	var output string
	var outputFile *os.File
	var writer io.WriteCloser

	switch c.Mode {
	case ModeFast, ModeUltra:
		output = source + ".tar.zst"
		outputFile, err = os.Create(output)
		if err != nil {
			return "", fmt.Errorf("create archive: %w", err)
		}
		level := zstd.SpeedFastest
		if c.Mode == ModeUltra {
			level = zstd.SpeedBetterCompression
		}
		encoder, err := zstd.NewWriter(outputFile, zstd.WithEncoderLevel(level))
		if err != nil {
			outputFile.Close()
			os.Remove(output)
			return "", fmt.Errorf("create zstd encoder: %w", err)
		}
		writer = encoder
	case ModeNormal:
		output = source + ".tar.gz"
		outputFile, err = os.Create(output)
		if err != nil {
			return "", fmt.Errorf("create archive: %w", err)
		}
		encoder, err := zstd.NewWriter(outputFile, zstd.WithEncoderLevel(zstd.SpeedDefault))
		if err != nil {
			outputFile.Close()
			os.Remove(output)
			return "", fmt.Errorf("create gzip encoder: %w", err)
		}
		writer = encoder
	}

	tw := tar.NewWriter(writer)

	baseDir := filepath.Base(source)
	err = filepath.Walk(source, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = filepath.Join(baseDir, walkPath)
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if !info.IsDir() {
			f, err := os.Open(walkPath)
			if err != nil {
				return err
			}
			defer f.Close()
			io.Copy(tw, f)
		}
		return nil
	})
	if err != nil {
		os.Remove(output)
		return "", fmt.Errorf("walk source: %w", err)
	}

	if err := tw.Close(); err != nil {
		os.Remove(output)
		return "", fmt.Errorf("close tar: %w", err)
	}

	writer.Close()
	outputFile.Close()

	if info.IsDir() {
		os.RemoveAll(source)
	} else {
		os.Remove(source)
	}

	return output, nil
}
