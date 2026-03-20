#!/bin/bash

set -euo pipefail

APP_NAME="media-suite"
INSTALL_DIR="/usr/bin"
SRC_DIR="./src"
BUILD_DIR="./build"
GO_MOD_NAME="media-suite"

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}==========================================${NC}"
<<<<<<< HEAD
echo -e "${BLUE}           media-suite  Installer       ${NC}"
=======
echo -e "${BLUE}    media-suite | Catppuccin Installer       ${NC}"
>>>>>>> 384a0a11dc78403475f1423e4ae5f67d7082ac0f
echo -e "${BLUE}==========================================${NC}"

echo -e "${GREEN}[+] Checking System Dependencies...${NC}"

if [[ -f /etc/debian_version ]]; then
    echo "Detected Debian/Ubuntu based system."
    sudo apt-get update -qq
    echo "Installing GCC, graphics libs, yt-dlp, ffmpeg, Go..."
    sudo apt-get install -y golang git gcc libgl1-mesa-dev xorg-dev yt-dlp ffmpeg
elif [[ -f /etc/arch-release ]]; then
    echo "Detected Arch Linux."
    sudo pacman -S --noconfirm go git base-devel libgl xorg-server yt-dlp ffmpeg
else
    echo -e "${RED}[!] Unsupported distro for auto-install.${NC}"
    echo "Ensure Go, GCC, libgl1-mesa-dev, xorg-dev, yt-dlp, ffmpeg are installed."
    read -p "Press Enter to continue anyway..."
fi

echo -e "${GREEN}[+] Setting up Go Module...${NC}"
if [[ ! -f go.mod ]]; then
    go mod init "$GO_MOD_NAME"
fi

echo "Downloading dependencies..."
go mod tidy
go get fyne.io/fyne/v2

echo -e "${GREEN}[+] Building Application...${NC}"
mkdir -p "$BUILD_DIR"
go build -ldflags "-s -w" -o "$BUILD_DIR/$APP_NAME" "$SRC_DIR/main.go"

echo -e "${GREEN}[+] Installing to $INSTALL_DIR...${NC}"
if [[ -f "$BUILD_DIR/$APP_NAME" ]]; then
    sudo mv "$BUILD_DIR/$APP_NAME" "$INSTALL_DIR/$APP_NAME"
    sudo chmod +x "$INSTALL_DIR/$APP_NAME"
    echo -e "${GREEN}Successfully installed to $INSTALL_DIR/$APP_NAME${NC}"
else
    echo -e "${RED}Binary not found – build failed.${NC}"
    exit 1
fi

echo -e "${BLUE}==========================================${NC}"
echo -e "${BLUE}Done! Run the app by typing: $APP_NAME${NC}"
<<<<<<< HEAD
echo -e "${BLUE}==========================================${NC}"
=======
echo -e "${BLUE}==========================================${NC}"
>>>>>>> 384a0a11dc78403475f1423e4ae5f67d7082ac0f
