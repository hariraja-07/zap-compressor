package extractor

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/gzip"
)

var magicBytes = []map[string][]byte{
	{".zst": []byte{0x28, 0xb5, 0x2f, 0xfd}}, // zstd magic
	{".gz": []byte{0x1f, 0x8b}},              // gzip magic
}

func detectFormat(data []byte) string {
	for _, m := range magicBytes {
		for ext, magic := range m {
			if len(data) >= len(magic) && bytes.Equal(data[:len(magic)], magic) {
				return ext
			}
		}
	}
	return ""
}

func Extract(archive string) error {
	f, err := os.Open(archive)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer f.Close()

	buf := make([]byte, 8)
	n, err := f.Read(buf)
	if err != nil || n == 0 {
		return fmt.Errorf("empty or unreadable archive")
	}

	format := detectFormat(buf)
	if format == "" {
		if len(archive) > 4 && archive[len(archive)-4:] == ".zip" {
			format = ".zip"
		} else {
			return fmt.Errorf("unsupported archive format")
		}
	}

	f.Seek(0, 0)

	var reader io.Reader
	switch format {
	case ".zst", ".gz":
		gr, err := gzip.NewReader(f)
		if err != nil {
			return fmt.Errorf("gzip decode: %w", err)
		}
		defer gr.Close()
		reader = gr
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	tr := tar.NewReader(reader)

	extractedDir := filepath.Base(archive)
	for _, ext := range []string{".tar.zst", ".tar.gz", ".tar.xz", ".zip"} {
		if len(extractedDir) >= len(ext) {
			extractedDir = extractedDir[:len(extractedDir)-len(ext)]
		}
	}

	hasContent := false
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read tar header: %w", err)
		}
		hasContent = true

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

	// Handle empty archive - was causing nil pointer panic before
	if !hasContent {
		return fmt.Errorf("archive is empty")
	}

	os.Remove(archive)
	fmt.Printf("Extracted to: %s\n", extractedDir)
	return nil
}
