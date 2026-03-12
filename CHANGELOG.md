# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-03-13

### Added
- Initial Go-based release with cross-platform support
- Compression modes:
  - `fast` - tar.zst (zstd fastest)
  - `normal` - tar.gz (zstd default)
  - `ultra` - tar.zst (zstd best compression)
  - `zip` - zip format
- Auto-detect format on extraction (magic bytes detection)
- Progress bars during compression
- Symlink preservation in archives
- Auto-cleanup (deletes source after compression)
- Cross-platform build script (Linux, macOS, Windows)
- CI/CD workflow for automated releases

### Changed
- Migrated from bash to Go binary
- No external dependencies required

### Fixed
- Zstd output corruption on large files
- Symlinks being dereferenced in tar.gz
- Progress bar showing 0% on Windows
- Empty archive extraction panic
- Windows binary naming (.exe)

---

## [0.1.0] - 2024-04-10

### Added
- Initial bash-based compression tool
- Fast compression (zstd)
- Ultra compression (7z)
- Auto-cleanup feature