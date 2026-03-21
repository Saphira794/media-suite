$APP_NAME = "media-suite"
$INSTALL_DIR = "C:\Program Files\${APP_NAME}"
$SRC_DIR = ".\src"
$BUILD_DIR = ".\build"
$GO_MOD_NAME = "media-suite"

function WriteColor($Text, $Color = "White") {
    Write-Host $Text -ForegroundColor $Color
}

WriteColor "==========================================" "Blue"
WriteColor "         media-suite Installer           " "Blue"
WriteColor "==========================================" "Blue"

WriteColor "[+] Checking System Dependencies..." "Green"

function Ensure-Package($Package, $ScoopName, $WingetName) {
    if (-not (Get-Command $Package -ErrorAction SilentlyContinue)) {
        WriteColor "[*] $Package not found. Attempting to install..." "Yellow"

        if (Get-Command scoop -ErrorAction SilentlyContinue) {
            WriteColor "[*] Installing $Package via Scoop..." "Blue"
            scoop install $ScoopName
        }
        elseif (Get-Command winget -ErrorAction SilentlyContinue) {
            WriteColor "[*] Installing $Package via Winget..." "Blue"
            winget install --id=$WingetName -e --silent
        }
        else {
            WriteColor "[!] Neither Scoop nor Winget is installed. Please install $Package manually." "Red"
            exit 1
        }
    }
    else {
        WriteColor "[+] $Package found." "Green"
    }
}

Ensure-Package "go" "golang" "GoLang.Go"
Ensure-Package "git" "git" "Git.Git"
Ensure-Package "gcc" "mingw" "GCC.GCC"
Ensure-Package "ffmpeg" "ffmpeg" "FFmpeg.FFmpeg"
Ensure-Package "yt-dlp" "yt-dlp" "yt-dlp.yt-dlp"

WriteColor "[+] Setting up Go Module..." "Green"

if (-not (Test-Path "go.mod")) {
    go mod init $GO_MOD_NAME
}

WriteColor "Downloading dependencies..." "Green"
go mod tidy
go get fyne.io/fyne/v2

WriteColor "[+] Building Application..." "Green"

if (-not (Test-Path $BUILD_DIR)) {
    New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null
}

go build -ldflags "-s -w" -o "$BUILD_DIR\$APP_NAME.exe" "$SRC_DIR\main.go"

WriteColor "[+] Installing to $INSTALL_DIR..." "Green"

if (-not (Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
}

$TargetPath = Join-Path $INSTALL_DIR "$APP_NAME.exe"

if (Test-Path "$BUILD_DIR\$APP_NAME.exe") {
    Copy-Item "$BUILD_DIR\$APP_NAME.exe" $TargetPath -Force
    WriteColor "Successfully installed to $TargetPath" "Green"
}
else {
    WriteColor "Binary not found – build failed." "Red"
    exit 1
}

Write-Color "==========================================" "Blue"
Write-Color "Done! Run the app by typing: $TargetPath" "Blue"
Write-Color "==========================================" "Blue"