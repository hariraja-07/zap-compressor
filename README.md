# ZAP - Smart Compression Tool

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)  
![Go](https://img.shields.io/badge/Go-1.24-blue?logo=go)

> **One-command compression and extraction with auto-cleanup.**  
> Lightweight, fast, and smart – just `zap` it!

---

## Installation

### Download Pre-built Binary

```bash
# Linux
curl -L https://github.com/zap-tool/zap/releases/download/v1.0.0/zap-1.0.0-linux-amd64 -o zap
chmod +x zap

# macOS
curl -L https://github.com/zap-tool/zap/releases/download/v1.0.0/zap-1.0.0-darwin-amd64 -o zap
chmod +x zap

# Windows
curl -L https://github.com/zap-tool/zap/releases/download/v1.0.0/zap-1.0.0-windows-amd64.exe -o zap.exe
```

### Build from Source

```bash
git clone https://github.com/zap-tool/zap.git
cd zap
./build.sh
```

---

## Usage

### Compression

```bash
zap <file|folder>         # Fast compression (tar.zst)
zap -m fast <file>        # Fast compression
zap -m normal <file>     # Normal compression (tar.gz)
zap -m ultra <file>      # Ultra compression (zstd max)
zap -m zip <file>       # Zip format
```

> Creates archive and **deletes the original**.

---

### Extraction

```bash
zap <archive.tar.zst>    # Auto-detect and extract
zap extract <archive>
```

> Restores files and **deletes the archive**.

---

## Options

| Flag | Description | Default |
|------|------------|---------|
| `-m, --mode` | Compression mode | `fast` |
| `-h, --help` | Show help | - |

---

## Compression Modes

| Mode | Format | Speed | Ratio |
|------|-------|-------|-------|
| `fast` | `.tar.zst` | Fastest | Good |
| `normal` | `.tar.gz` | Normal | Better |
| `ultra` | `.tar.zst` | Slow | Best |
| `zip` | `.zip` | Normal | Standard |

---

## Features

- ✅ **Self-contained** (no external dependencies)
- ✅ **Cross-platform** (Linux, macOS, Windows)
- ✅ **Auto-detect format** on extraction
- ✅ **Progress bars** during compression
- ✅ **Preserves symlinks**
- ✅ **Auto-cleanup** (deletes source after operation)

---

## License

MIT License - see [LICENSE](LICENSE)