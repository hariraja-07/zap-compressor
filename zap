#!/bin/bash

# Color and emoji setup
GREEN='\033[0;32m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color
CHECK="✅"
ZAP="⚡"

# Generate a safe archive name
get_safe_name() {
    local path="${1%/}"  # Remove trailing slash
    local base=$(basename "$path")
    # If filename would be invalid, use timestamp
    if [[ "$base" =~ [/\\] ]]; then
        base="zap_$(date +%s)"
    fi
    echo "${base}"
}

show_help() {
    echo -e "${CYAN}Usage:${NC}"
    echo -e "  ${GREEN}zap <file/folder>${NC}         # Fast compression (zstd)"
    echo -e "  ${GREEN}zap -u <file/folder>${NC}      # Ultra compression (7z)"
    echo -e "  ${GREEN}zap <archive.zst/.7z>${NC}     # Auto-extract"
}

# Main execution
if [ $# -eq 0 ]; then
    show_help
    exit 1
fi

# Mode detection
if [ "$1" == "-u" ]; then
    MODE="ultra"
    TARGET="$2"
else
    MODE="fast"
    TARGET="$1"
fi

# Normalize target path
TARGET="${TARGET%/}"
if [ ! -e "$TARGET" ]; then
    echo -e "${RED}❌ Error: '${TARGET}' not found!${NC}"
    exit 1
fi

# ---- Compression ----
if [ -d "$TARGET" ] || [ -f "$TARGET" ] && [[ "$TARGET" != *.zst && "$TARGET" != *.7z ]]; then
    echo -e "${ZAP} ${CYAN}Compressing '${TARGET}'...${NC}"
    
    SAFE_NAME=$(get_safe_name "$TARGET")
    
    case "$MODE" in
        fast)
            ARCHIVE="${SAFE_NAME}.zst"
            if tar --zstd --preserve-permissions --xattrs -cvf "$ARCHIVE" "$TARGET"; then
                echo -e "${CHECK} ${GREEN}Compressed to: ${ARCHIVE}${NC}"
                rm -rf "$TARGET" && echo -e "${ZAP} Deleted original: ${TARGET}"
            else
                echo -e "${RED}❌ Compression failed!${NC}"
                [ -f "$ARCHIVE" ] && rm -f "$ARCHIVE"
                exit 1
            fi
            ;;
        ultra)
            ARCHIVE="${SAFE_NAME}.7z"
            if 7z a -t7z -mx=9 -snl -spf "$ARCHIVE" "$TARGET"; then
                echo -e "${CHECK} ${GREEN}Compressed to: ${ARCHIVE}${NC}"
                rm -rf "$TARGET" && echo -e "${ZAP} Deleted original: ${TARGET}"
            else
                echo -e "${RED}❌ Compression failed!${NC}"
                [ -f "$ARCHIVE" ] && rm -f "$ARCHIVE"
                exit 1
            fi
            ;;
    esac

# ---- Extraction ----
elif [[ "$TARGET" == *.zst ]]; then
    echo -e "${ZAP} ${CYAN}Extracting '${TARGET}'...${NC}"
    EXTRACT_DIR="${TARGET%.zst}"
    
    if tar --zstd --preserve-permissions --xattrs -xvf "$TARGET"; then
        echo -e "${CHECK} ${GREEN}Extracted to: ${EXTRACT_DIR}${NC}"
        rm -f "$TARGET" && echo -e "${ZAP} Deleted archive: ${TARGET}"
    else
        echo -e "${RED}❌ Extraction failed!${NC}"
        exit 1
    fi

elif [[ "$TARGET" == *.7z ]]; then
    echo -e "${ZAP} ${CYAN}Extracting '${TARGET}'...${NC}"
    EXTRACT_DIR="${TARGET%.7z}"
    
    if 7z x -snl "$TARGET"; then
        echo -e "${CHECK} ${GREEN}Extracted to: ${EXTRACT_DIR}${NC}"
        rm -f "$TARGET" && echo -e "${ZAP} Deleted archive: ${TARGET}"
    else
        echo -e "${RED}❌ Extraction failed!${NC}"
        exit 1
    fi

else
    echo -e "${RED}❌ Error: Unsupported file type!${NC}"
    show_help
    exit 1
fi

