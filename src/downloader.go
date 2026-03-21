package main

import (
	"os/exec"
	"path/filepath"
	"strings"
)

func performDownload() error {
	logSystem("Preparing download for: " + globalState.DownloadURL)

	args := []string{"--newline", "--no-colors"}

	_, errAria := exec.LookPath("aria2c")
	if errAria == nil {
		logSystem("aria2c found. Enabling multi-threaded download acceleration (-x16 -k1M).")
		args = append(args, "--external-downloader", "aria2c")
		args = append(args, "--external-downloader-args", "-x16 -k1M")
	}

	outTmpl := filepath.Join(globalState.DownloadPath, "%(title)s.%(ext)s")
	args = append(args, "-o", outTmpl)

	if strings.Contains(globalState.DownloadFormat, "Audio") {
		args = append(args, "-x") 

		audioFmt := "mp3" 
		if strings.Contains(globalState.DownloadFormat, "M4A") {
			audioFmt = "m4a"
		}
		if strings.Contains(globalState.DownloadFormat, "WAV") {
			audioFmt = "wav"
		}
		if strings.Contains(globalState.DownloadFormat, "FLAC") {
			audioFmt = "flac"
		}

		args = append(args, "--audio-format", audioFmt)
		args = append(args, "--audio-quality", "0") 

	} else if strings.Contains(globalState.DownloadFormat, "Thumbnail") {
		args = append(args, "--write-thumbnail", "--skip-download")
		args = append(args, "--convert-thumbnails", "jpg")
	} else {
		args = append(args, "--merge-output-format", "mp4") 

		qualityArg := "bestvideo+bestaudio/best"
		switch globalState.Quality {
		case "4K":
			qualityArg = "bestvideo[height<=2160]+bestaudio/best[height<=2160]"
		case "1440p":
			qualityArg = "bestvideo[height<=1440]+bestaudio/best[height<=1440]"
		case "1080p":
			qualityArg = "bestvideo[height<=1080]+bestaudio/best[height<=1080]"
		case "720p":
			qualityArg = "bestvideo[height<=720]+bestaudio/best[height<=720]"
		case "480p":
			qualityArg = "bestvideo[height<=480]+bestaudio/best[height<=480]"
		case "Worst":
			qualityArg = "worst"
		}
		args = append(args, "-f", qualityArg)
	}

	if globalState.EmbedMetadata {
		args = append(args, "--add-metadata")
	}
	if globalState.EmbedThumbnail && !strings.Contains(globalState.DownloadFormat, "Thumbnail") {
		args = append(args, "--embed-thumbnail")
	}
	if globalState.EmbedSubs {
		args = append(args, "--write-auto-sub", "--sub-lang", "en", "--embed-subs")
	}

	args = append(args, globalState.DownloadURL)

	return runCommandWithProgress("yt-dlp", args...)
}

