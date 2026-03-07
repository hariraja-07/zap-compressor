package main

import (
	"fmt"
	"os"

	"github.com/zap-tool/zap/internal/cli"
)

var Version = "dev"

func main() {
	fmt.Printf("ZAP v%s - Compression Tool\n\n", Version)
	if err := cli.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
