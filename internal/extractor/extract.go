package extractor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"archive/tar"
	"github.com/klauspost/compress/zstd"
)

func Extract(archive string) error {
	info, err := os.Stat(archive)
	if err != nil {
		return fmt.Errorf("stat archive: %w", err)
	}

	extractedDir := filepath.Base(archive)
	for _, ext := range []string{".tar.zst", ".tar.gz", ".tar.xz"} {
		if len(extractedDir) > len(ext) {
			extractedDir = extractedDir[:len(extractedDir)-len(ext)]
		}
	}

	if err := os.MkdirAll(extractedDir, info.Mode()); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	f, err := os.Open(archive)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer f.Close()

	var reader io.Reader
	switch {
	case len(archive) > 5 && archive[len(archive)-5:] == ".zst":
		zstdReader, err := zstd.NewReader(f)
		if err != nil {
			return fmt.Errorf("zstd decode: %w", err)
		}
		reader = zstdReader
	default:
		return fmt.Errorf("unsupported format")
	}

	tr := tar.NewReader(reader)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read tar: %w", err)
		}

		path := filepath.Join(extractedDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("create dir: %w", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return fmt.Errorf("create parent dir: %w", err)
			}
			outFile, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("create file: %w", err)
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("write file: %w", err)
			}
			outFile.Close()
		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, path); err != nil {
				return fmt.Errorf("create symlink: %w", err)
			}
		}
	}

	os.Remove(archive)
	fmt.Printf("Extracted to: %s\n", extractedDir)
	return nil
}
