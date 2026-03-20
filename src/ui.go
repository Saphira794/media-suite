package main

import (
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Debounced logging to reduce UI update frequency and CPU usage.
var (
	logQueue          []string
	logFlushMu        sync.Mutex
	logFlushScheduled bool
	logFlushInterval  = 250 * time.Millisecond
)

// runOnMain is a small helper that executes a function intended for the UI thread.
// In this simplified setup it just calls f() directly so the code compiles on
// older Fyne versions that do not provide App.Schedule or Driver().RunOnMain.
func runOnMain(f func()) {
	if f != nil {
		f()
	}
}

// logSystem logs a message to the console with debounced UI updates.
func logSystem(msg string) {
	// CRITICAL FIX: Ensure the window is ready before proceeding.
	if !isWindowReady || consoleLog == nil {
		return
	}

	ts := time.Now().Format("15:04:05")
	formatted := fmt.Sprintf("[%s] %s", ts, msg)

	// Add to the queue and schedule a flush if not already scheduled.
	logFlushMu.Lock()
	logQueue = append(logQueue, formatted)
	if !logFlushScheduled {
		logFlushScheduled = true
		go func() {
			time.Sleep(logFlushInterval)
			logFlushMu.Lock()
			batch := strings.Join(logQueue, "\n")
			logQueue = nil
			logFlushScheduled = false
			logFlushMu.Unlock()
			// Schedule UI update on the main Fyne thread
			runOnMain(func() {
				if consoleLog != nil {
					consoleLog.SetText(consoleLog.Text + batch + "\n")
					consoleLog.CursorRow = len(strings.Split(consoleLog.Text, "\n"))
					consoleLog.Refresh()
				}
			})
		}()
	}
	logFlushMu.Unlock()
}

// addToHistory adds an item to the history and refreshes the history list.
func addToHistory(item string) {
	historyData = append(historyData, item)

	// Always schedule list refresh
	runOnMain(func() {
		historyList.Refresh()
	})
}

// parseURL parses a URL string and returns a URL pointer.
func parseURL(urlStr string) *url.URL {
	u, _ := url.Parse(urlStr)
	return u
}

// UI Builders

func buildDownloaderTab() fyne.CanvasObject {
	// -- Header --
	title := canvas.NewText("Media Downloader", hexToColor(globalTheme.CurrentPalette.Mauve))
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := widget.NewLabel("Powered by yt-dlp")

	// -- Inputs --
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Paste URL from YouTube, Twitch, Twitter, etc...")
	urlEntry.OnChanged = func(s string) { globalState.DownloadURL = s }

	// -- Config Grid --

	// Format
	downloadFormatSelect = widget.NewSelect([]string{
		"Video (MP4)", "Video (MKV)", "Video (WebM)",
		"Audio (MP3)", "Audio (M4A)", "Audio (WAV)", "Audio (FLAC)",
		"Thumbnail Only (JPG)",
	}, func(s string) {
		globalState.DownloadFormat = s
	})

	// Quality
	downloadQualitySelect = widget.NewSelect([]string{"Best", "4K", "1440p", "1080p", "720p", "480p", "Worst"}, func(s string) {
		globalState.Quality = s
	})

	// Options
	checkMeta := widget.NewCheck("Embed Metadata", func(b bool) { globalState.EmbedMetadata = b })
	checkMeta.SetChecked(true)

	checkThumb := widget.NewCheck("Embed Thumbnail (Video)", func(b bool) { globalState.EmbedThumbnail = b })
	checkThumb.SetChecked(true)

	checkSubs := widget.NewCheck("Download Subtitles", func(b bool) { globalState.EmbedSubs = b })

	// Path
	pathEntry := widget.NewEntry()
	pathEntry.SetText(globalState.DownloadPath)
	pathBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				pathEntry.SetText(uri.Path())
				globalState.DownloadPath = uri.Path()
			}
		}, globalWindow)
	})

	// -- Action Button --
	dlBtn := widget.NewButtonWithIcon("START DOWNLOAD", theme.DownloadIcon(), nil)
	dlBtn.Importance = widget.HighImportance
	dlBtn.OnTapped = func() {
		if globalState.IsBusy {
			return
		}
		if globalState.DownloadURL == "" {
			dialog.ShowError(fmt.Errorf("URL cannot be empty"), globalWindow)
			return
		}

		globalState.IsBusy = true
		dlBtn.Disable()
		progressBar.SetValue(0)
		statusLabel.SetText("Initializing download...")

		go func() {
			err := performDownload()
			globalState.IsBusy = false

			// Always schedule UI updates
			runOnMain(func() {
				dlBtn.Enable()
				if err != nil {
					logSystem("ERROR: " + err.Error())
					statusLabel.SetText("Download Failed")
				} else {
					logSystem("Download sequence completed.")
					statusLabel.SetText("Download Complete")
					progressBar.SetValue(1.0)
					addToHistory(globalState.DownloadURL + " (" + globalState.DownloadFormat + ")")
				}
			})
		}()
	}

	// Layout
	formContainer := container.NewVBox(
		widget.NewLabel("Source URL"),
		urlEntry,
		widget.NewSeparator(),
		widget.NewLabel("Configuration"),
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabel("Format"), downloadFormatSelect),
			container.NewVBox(widget.NewLabel("Quality"), downloadQualitySelect),
		),
		widget.NewSeparator(),
		widget.NewLabel("Post-Processing"),
		container.NewHBox(checkMeta, checkThumb, checkSubs),
		widget.NewSeparator(),
		widget.NewLabel("Output Directory"),
		container.NewBorder(nil, nil, nil, pathBtn, pathEntry),
		layout.NewSpacer(),
		dlBtn,
	)

	return container.NewPadded(container.NewBorder(
		container.NewVBox(title, subtitle, widget.NewSeparator()),
		nil, nil, nil,
		formContainer,
	))
}

func buildConverterTab() fyne.CanvasObject {
	title := canvas.NewText("File Converter", hexToColor(globalTheme.CurrentPalette.Green))
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Input File
	inputLabel := widget.NewLabel("No file selected")
	inputLabel.Wrapping = fyne.TextWrapBreak

	selectFileBtn := widget.NewButton("Select Input File", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				globalState.ConvertSourceFile = reader.URI().Path()
				inputLabel.SetText("Selected: " + filepath.Base(globalState.ConvertSourceFile))
			}
		}, globalWindow)
		fd.Show()
	})

	// Target Format
	convertFormatSelect = widget.NewSelect([]string{"mp3", "wav", "aac", "flac", "mp4", "mkv", "avi", "gif", "webm"}, func(s string) {
		globalState.ConvertDestFormat = s
	})

	// Convert Button
	convertBtn := widget.NewButtonWithIcon("CONVERT NOW", theme.MediaPlayIcon(), nil)
	convertBtn.Importance = widget.HighImportance
	convertBtn.OnTapped = func() {
		if globalState.IsBusy {
			return
		}
		if globalState.ConvertSourceFile == "" {
			dialog.ShowError(fmt.Errorf("Please select a file first"), globalWindow)
			return
		}

		globalState.IsBusy = true

		// Schedule UI updates
		runOnMain(func() {
			convertBtn.Disable()
			progressBar.SetValue(0)
			statusLabel.SetText("Converting...")
		})

		go func() {
			err := performConversion()
			globalState.IsBusy = false

			// Always schedule UI updates
			runOnMain(func() {
				convertBtn.Enable()
				if err != nil {
					logSystem("Conversion Error: " + err.Error())
					statusLabel.SetText("Conversion Failed")
				} else {
					statusLabel.SetText("Conversion Complete")
					progressBar.SetValue(1.0)
					addToHistory("Converted: " + filepath.Base(globalState.ConvertSourceFile) + " -> " + globalState.ConvertDestFormat)
				}
			})
		}()
	}

	return container.NewPadded(container.NewVBox(
		title,
		widget.NewSeparator(),
		widget.NewLabel("1. Input File"),
		selectFileBtn,
		inputLabel,
		widget.NewSeparator(),
		widget.NewLabel("2. Target Format"),
		convertFormatSelect,
		widget.NewSeparator(),
		layout.NewSpacer(),
		convertBtn,
	))
}

func buildHistoryTab() fyne.CanvasObject {
	historyData = []string{}
	historyList = widget.NewList(
		func() int { return len(historyData) },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.FileIcon()), widget.NewLabel("Template"))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			label := box.Objects[1].(*widget.Label)
			label.SetText(historyData[i])
		},
	)

	clearBtn := widget.NewButtonWithIcon("Clear History", theme.DeleteIcon(), func() {
		historyData = []string{}
		historyList.Refresh()
	})

	return container.NewBorder(
		widget.NewLabelWithStyle("Session History", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewPadded(clearBtn),
		nil, nil,
		historyList,
	)
}

func buildSettingsTab() fyne.CanvasObject {
	// Theme Selector
	themeSelect = widget.NewSelect([]string{"Latte", "Frappe", "Macchiato", "Mocha"}, func(s string) {
		globalTheme.SetFlavor(s)
		globalApp.Settings().SetTheme(globalTheme)

		// Refreshing the window content ensures colors update everywhere
		if globalWindow.Content() != nil {
			globalWindow.Content().Refresh()
		}
		logSystem("Theme switched to: " + s)
	})

	// Binary Check
	checkBtn := widget.NewButton("Check Dependencies", func() {
		_, errY := exec.LookPath("yt-dlp")
		_, errF := exec.LookPath("ffmpeg")
		_, errA := exec.LookPath("aria2c") // Check aria2c as well

		msg := "Dependencies Status:\n"
		if errY == nil {
			msg += "yt-dlp found\n"
		} else {
			msg += "yt-dlp NOT found\n"
		}
		if errF == nil {
			msg += "ffmpeg found\n"
		} else {
			msg += "ffmpeg NOT found\n"
		}
		if errA == nil {
			msg += "aria2c found (Multi-thread acceleration available)\n"
		} else {
			msg += "aria2c NOT found (Single-thread download)\n"
		}

		dialog.ShowInformation("System Check", msg, globalWindow)
	})

	return container.NewPadded(container.NewVBox(
		widget.NewLabelWithStyle("Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewLabel("Appearance"),
		container.NewGridWithColumns(2, widget.NewLabel("Theme Flavor"), themeSelect),
		widget.NewSeparator(),
		widget.NewLabel("System"),
		checkBtn,
		widget.NewSeparator(),
		widget.NewLabel("About"),
		widget.NewLabel("Catppuccin Downloader "+AppVersion),
		widget.NewHyperlink("Visit Catppuccin", parseURL("https://github.com/catppuccin/catppuccin")),
	))
}
