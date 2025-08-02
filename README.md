# ⚡ ZAP - Smart Compression Tool

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)  
![Shell Script](https://img.shields.io/badge/Shell_Script-121011?logo=gnu-bash&logoColor=white)

> **One-command compression and extraction with auto-cleanup.**  
> Lightweight, fast, and smart – just `zap` it!

---

## 🚀 Quick Start

### 🔧 Install Dependencies

```bash
sudo apt install zstd p7zip-full
```

### 📥 Download and Install

```bash
git clone https://github.com/yourusername/zap-compressor.git
cd zap-compressor
sudo cp zap /usr/local/bin/
```

 ---

## 💻 Usage

### ⚙️ Compression

```bash
zap <file/folder>         # Fast compression using zstd
zap -u <file/folder>      # Ultra compression using 7z
```

> 🗃️ Creates `archive.zst` or `archive.7z` and **deletes the original**.

---

### 📦 Extraction

```bash
zap <archive.zst/.7z>     # Auto-detect and extract
```

> 🔄 Restores original files/folders and **deletes the archive**.

---

## ✨ Features

- ✅ **Automatic format detection**
- ♻️ **Self-cleaning** (removes source after operation)
- 🚀 **Two compression modes**: Fast (zstd) and Ultra (7z)
- 🔐 **Preserves file permissions and attributes**
- 💡 **Simple & minimal interface**

---

## 🤝 Contributing & Support

As a new GitHub user, I welcome:

- 🐛 Bug reports  
- 💡 Feature requests  
- 🛠️ Pull requests  
- 📚 Documentation improvements

> Please open an [issue](https://github.com/yourusername/zap-compressor/issues) for suggestions or problems!

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).

