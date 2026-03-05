package compressor

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
)

type CompressionMode int

const (
	ModeFast CompressionMode = iota
	ModeNormal
	ModeUltra
	ModeZip
)

type Compressor struct {
	Mode     CompressionMode
	Progress func(current, total int64)
}

func NewCompressor(mode CompressionMode) *Compressor {
	return &Compressor{Mode: mode}
}

func (c *Compressor) CompressFile(source string) (string, error) {
	info, err := os.Lstat(source)
	if err != nil {
		return "", fmt.Errorf("stat source: %w", err)
	}

	switch c.Mode {
	case ModeZip:
		return c.compressZip(source, info)
	default:
		return c.compressTar(source, info)
	}
}

func (c *Compressor) compressTar(source string, sourceInfo os.FileInfo) (string, error) {
	var output string
	var outputFile *os.File
	var writer io.WriteCloser
	var err error

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
	_ = baseDir
	err = filepath.Walk(source, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var link string
		if info.Mode()&os.ModeSymlink != 0 {
			link, err = os.Readlink(walkPath)
			if err != nil {
				return err
			}
		}
		header, err := tar.FileInfoHeader(info, link)
		if err != nil {
			return err
		}
		header.Name = walkPath
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if !info.IsDir() && link == "" {
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

	if sourceInfo.IsDir() {
		os.RemoveAll(source)
	} else {
		os.Remove(source)
	}

	return output, nil
}

func (c *Compressor) compressZip(source string, sourceInfo os.FileInfo) (string, error) {
	output := source + ".zip"
	outputFile, err := os.Create(output)
	if err != nil {
		return "", fmt.Errorf("create archive: %w", err)
	}
	defer outputFile.Close()

	gw, err := gzip.NewWriterLevel(outputFile, gzip.BestCompression)
	if err != nil {
		outputFile.Close()
		os.Remove(output)
		return "", fmt.Errorf("create gzip encoder: %w", err)
	}
	defer gw.Close()

	tw := tar.NewWriter(gw)

	baseDir := filepath.Base(source)
	_ = baseDir
	err = filepath.Walk(source, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var link string
		if info.Mode()&os.ModeSymlink != 0 {
			link, err = os.Readlink(walkPath)
			if err != nil {
				return err
			}
		}
		header, err := tar.FileInfoHeader(info, link)
		if err != nil {
			return err
		}
		header.Name = walkPath
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if !info.IsDir() && link == "" {
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

	if sourceInfo.IsDir() {
		os.RemoveAll(source)
	} else {
		os.Remove(source)
	}

	return output, nil
}
