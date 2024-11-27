package main

import (
	"log"
	"os/exec"
	"strings"
)

func main() {
	log.Println("Starting Lute converter...")
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Found ffmpeg at %q\n", ffmpegPath)
	}
	ffmpegCommand := exec.Command(ffmpegPath, "--help")
	var ffmpegOutput strings.Builder
	ffmpegCommand.Stdout = &ffmpegOutput
	ffmpegCommand.Run()
	log.Printf("Executing '%s'\n %s\n", ffmpegCommand, ffmpegOutput.String())
}
