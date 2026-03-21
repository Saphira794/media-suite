package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	AppTitle       = "Catppuccin Media Suite"
	AppVersion     = "2.0.3-Optimized"
	DefaultPadding = 10
)

type CatppuccinPalette struct {
	Name      string
	Rosewater string
	Flamingo  string
	Pink      string
	Mauve     string
	Red       string
	Maroon    string
	Peach     string
	Yellow    string
	Green     string
	Teal      string
	Sky       string
	Sapphire  string
	Blue      string
	Lavender  string
	Text      string
	Subtext1  string
	Subtext0  string
	Overlay2  string
	Overlay1  string
	Overlay0  string
	Surface2  string
	Surface1  string
	Surface0  string
	Base      string
	Mantle    string
	Crust     string
}

type ThemeManager struct {
	CurrentPalette CatppuccinPalette
	Variant        fyne.ThemeVariant
}

func NewThemeManager(flavor string) *ThemeManager {
	tm := &ThemeManager{}
	tm.SetFlavor(flavor)
	return tm
}

func (t *ThemeManager) SetFlavor(flavor string) {
	switch strings.ToLower(flavor) {
	case "latte":
		t.CurrentPalette = PaletteLatte
		t.Variant = theme.VariantLight
	case "frappe":
		t.CurrentPalette = PaletteFrappe
		t.Variant = theme.VariantDark
	case "macchiato":
		t.CurrentPalette = PaletteMacchiato
		t.Variant = theme.VariantDark
	default: 
		t.CurrentPalette = PaletteMocha
		t.Variant = theme.VariantDark
	}
}

func (t *ThemeManager) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	p := t.CurrentPalette

	switch n {
	case theme.ColorNameBackground:
		return hexToColor(p.Base)
	case theme.ColorNameForeground:
		return hexToColor(p.Text)
	case theme.ColorNameButton:
		return hexToColor(p.Surface0)
	case theme.ColorNameDisabledButton:
		return hexToColor(p.Surface1)
	case theme.ColorNameDisabled:
		return hexToColor(p.Overlay0)
	case theme.ColorNameError:
		return hexToColor(p.Red)
	case theme.ColorNameFocus:
		return hexToColor(p.Mauve)
	case theme.ColorNameHover:
		return hexToColor(p.Surface1)
	case theme.ColorNameInputBackground:
		return hexToColor(p.Mantle)
	case theme.ColorNamePlaceHolder:
		return hexToColor(p.Overlay1)
	case theme.ColorNamePrimary:
		return hexToColor(p.Mauve)
	case theme.ColorNameScrollBar:
		return hexToColor(p.Surface1)
	case theme.ColorNameShadow:
		return hexToColor(p.Crust)
	case theme.ColorNameMenuBackground:
		return hexToColor(p.Mantle)
	case theme.ColorNameSeparator:
		return hexToColor(p.Surface1)
	default:
		return theme.DefaultTheme().Color(n, v)
	}
}

func (t *ThemeManager) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}
func (t *ThemeManager) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}
func (t *ThemeManager) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}

func hexToColor(hex string) color.Color {
	format := "%02x%02x%02x"
	var r, g, b uint8
	if len(hex) == 7 {
		fmt.Sscanf(hex, "#"+format, &r, &g, &b)
	} else {
		return color.RGBA{128, 128, 128, 255}
	}
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

var PaletteMocha = CatppuccinPalette{
	Name: "Mocha", Rosewater: "#f5e0dc", Flamingo: "#f2cdcd", Pink: "#f5c2e7", Mauve: "#cba6f7", Red: "#f38ba8", Maroon: "#eba0ac", Peach: "#fab387", Yellow: "#f9e2af", Green: "#a6e3a1", Teal: "#94e2d5", Sky: "#89dceb", Sapphire: "#74c7ec", Blue: "#89b4fa", Lavender: "#b4befe", Text: "#cdd6f4", Subtext1: "#bac2de", Subtext0: "#a6adc8", Overlay2: "#9399b2", Overlay1: "#7f849c", Overlay0: "#6c7086", Surface2: "#585b70", Surface1: "#45475a", Surface0: "#313244", Base: "#1e1e2e", Mantle: "#181825", Crust: "#11111b",
}

var PaletteMacchiato = CatppuccinPalette{
	Name: "Macchiato", Rosewater: "#f4dbd6", Flamingo: "#f0c6c6", Pink: "#f5bde6", Mauve: "#c6a0f6", Red: "#ed8796", Maroon: "#ee99a0", Peach: "#f5a97f", Yellow: "#eed49f", Green: "#a6da95", Teal: "#8bd5ca", Sky: "#91d7e3", Sapphire: "#7dc4e4", Blue: "#8aadf4", Lavender: "#b7bdf8", Text: "#cad3f5", Subtext1: "#b8c0e0", Subtext0: "#a6da95", Overlay2: "#939ab7", Overlay1: "#8087a2", Overlay0: "#6e738d", Surface2: "#5b6078", Surface1: "#494d64", Surface0: "#363a4f", Base: "#24273a", Mantle: "#1e2030", Crust: "#181926",
}

var PaletteFrappe = CatppuccinPalette{
	Name: "Frappe", Rosewater: "#f2d5cf", Flamingo: "#eebebe", Pink: "#f4b8e4", Mauve: "#ca9ee6", Red: "#e78284", Maroon: "#ea999c", Peach: "#ef9f76", Yellow: "#e5c890", Green: "#a6d189", Teal: "#81c8be", Sky: "#99d1db", Sapphire: "#85c1dc", Blue: "#8caaee", Lavender: "#babbf1", Text: "#c6d0f5", Subtext1: "#b5bfe2", Subtext0: "#a5adce", Overlay2: "#949cbb", Overlay1: "#838ba7", Overlay0: "#737994", Surface2: "#626880", Surface1: "#51576d", Surface0: "#414559", Base: "#303446", Mantle: "#292c3c", Crust: "#232634",
}

var PaletteLatte = CatppuccinPalette{
	Name: "Latte", Rosewater: "#dc8a78", Flamingo: "#dd7878", Pink: "#ea76cb", Mauve: "#8839ef", Red: "#d20f39", Maroon: "#e64553", Peach: "#fe640b", Yellow: "#df8e1d", Green: "#40a02b", Teal: "#179299", Sky: "#04a5e5", Sapphire: "#209fb5", Blue: "#1e66f5", Lavender: "#7287fd", Text: "#4c4f69", Subtext1: "#5c5f77", Subtext0: "#6c6f85", Overlay2: "#7c7f93", Overlay1: "#8c8fa1", Overlay0: "#9ca0b0", Surface2: "#acb0be", Surface1: "#bcc0cc", Surface0: "#ccd0da", Base: "#eff1f5", Mantle: "#e6e9ef", Crust: "#dce0e8",
}


type AppState struct {
	DownloadURL    string
	DownloadFormat string 
	Quality        string 
	DownloadPath   string
	IsBusy         bool
	EmbedMetadata  bool
	EmbedThumbnail bool
	EmbedSubs      bool
	CustomArgs     string

	ConvertSourceFile string
	ConvertDestFormat string 
	ConvertStatus     string

	CurrentThemeName string
	LogHistory       []string
}

var globalState AppState
var globalApp fyne.App
var globalWindow fyne.Window
var globalTheme *ThemeManager

var (
	consoleLog  *widget.Entry
	progressBar *widget.ProgressBar
	statusLabel *widget.Label
	historyList *widget.List
	historyData []string

	downloadFormatSelect  *widget.Select
	downloadQualitySelect *widget.Select
	convertFormatSelect   *widget.Select
	themeSelect           *widget.Select

	isWindowReady bool
)


func main() {
	globalApp = app.New()
	globalWindow = globalApp.NewWindow(fmt.Sprintf("%s v%s", AppTitle, AppVersion))
	globalWindow.Resize(fyne.NewSize(1000, 750))

	globalTheme = NewThemeManager("Mocha")
	globalApp.Settings().SetTheme(globalTheme)

	wd, _ := os.Getwd()
	globalState = AppState{
		DownloadFormat:    "Video (MP4)",
		Quality:           "Best",
		DownloadPath:      wd,
		EmbedMetadata:     true,
		EmbedThumbnail:    true,
		CurrentThemeName:  "Mocha",
		ConvertDestFormat: "mp3",
	}

	mainTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Downloader", theme.DownloadIcon(), buildDownloaderTab()),
		container.NewTabItemWithIcon("Converter", theme.MediaPlayIcon(), buildConverterTab()),
		container.NewTabItemWithIcon("History", theme.HistoryIcon(), buildHistoryTab()),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), buildSettingsTab()),
	)
	mainTabs.SetTabLocation(container.TabLocationLeading)

	consoleLog = widget.NewMultiLineEntry()
	consoleLog.Disable()
	consoleLog.TextStyle = fyne.TextStyle{Monospace: true}
	consoleLog.SetPlaceHolder("System logs will appear here...")

	progressBar = widget.NewProgressBar()
	progressBar.SetValue(0)

	statusLabel = widget.NewLabel("Ready")
	statusLabel.Alignment = fyne.TextAlignCenter

	logScroll := container.NewScroll(consoleLog)
	logScroll.SetMinSize(fyne.NewSize(0, 150))

	finalLayout := container.NewBorder(
		nil,
		container.NewVBox(progressBar, statusLabel, logScroll),
		nil,
		nil,
		mainTabs,
	)

	globalWindow.SetContent(finalLayout)
	isWindowReady = true 

	if themeSelect != nil {
		runOnMain(func() { themeSelect.SetSelected(globalState.CurrentThemeName) })
	}
	if downloadFormatSelect != nil {
		runOnMain(func() { downloadFormatSelect.SetSelected(globalState.DownloadFormat) })
	}
	if downloadQualitySelect != nil {
		runOnMain(func() { downloadQualitySelect.SetSelected(globalState.Quality) })
	}
	if convertFormatSelect != nil {
		runOnMain(func() { convertFormatSelect.SetSelected(globalState.ConvertDestFormat) })
	}

	logSystem("Welcome to " + AppTitle)
	logSystem("Engine initialized. Ready for operations.")

	globalWindow.ShowAndRun()
}


func runCommandWithProgress(bin string, args ...string) error {
	logSystem(fmt.Sprintf("CMD: %s %v", bin, args))

	cmd := exec.Command(bin, args...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	lastProgressUpdate := time.Now()

	reProgress := regexp.MustCompile(`\[download\]\s+(\d+\.\d+)%`)
	reFFmpegProgress := regexp.MustCompile(`(frame=|size=|time=|bitrate=|speed=)`)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			matches := reProgress.FindStringSubmatch(line)
			if len(matches) > 1 {
				p, _ := strconv.ParseFloat(matches[1], 64)

				if time.Since(lastProgressUpdate) > 500*time.Millisecond {
					lastProgressUpdate = time.Now()

					runOnMain(func() {
						progressBar.SetValue(p / 100.0)
						statusLabel.SetText(fmt.Sprintf("Downloading... %.1f%%", p))
					})
				}
			} else {
				if !strings.HasPrefix(line, "[download]") {
					logSystem(line)
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()

			if reFFmpegProgress.MatchString(line) {
				continue
			}

			logSystem("[stderr] " + line)
		}
	}()

	wg.Wait()
	return cmd.Wait()
}

type HistoryFile struct {
	Items []string `json:"items"`
}

func saveHistory() {
	f, err := os.Create("history.json")
	if err != nil {
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(HistoryFile{Items: historyData})
}

