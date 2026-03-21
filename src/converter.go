package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func performConversion() error {
	inputFile := globalState.ConvertSourceFile
	ext := filepath.Ext(inputFile)
	outputFile := strings.TrimSuffix(inputFile, ext) + "_converted." + globalState.ConvertDestFormat

	logSystem(fmt.Sprintf("Converting %s -> %s", inputFile, outputFile))

	args := []string{"-i", inputFile, "-y"} 

	switch globalState.ConvertDestFormat {
	case "mp3":
		args = append(args, "-q:a", "0", "-map", "a")
	case "gif":
		args = append(args, "-vf", "fps=10,scale=320:-1:flags=lanczos", "-c:v", "gif")
	case "mp4":
		args = append(args, "-c:v", "libx264", "-c:a", "aac", "-strict", "experimental")
	}

	args = append(args, outputFile)

	return runCommandWithProgress("ffmpeg", args...)
}

