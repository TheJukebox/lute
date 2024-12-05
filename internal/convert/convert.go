package convert

import (
	"log"           // providing logging for us - not necessary, but I like to start logging ASAP.
	"os"            // for using Stat to check if files exist
	"os/exec"       // this is a builtin for executing commands on the host OS
	"path/filepath" //Looks for the path of a file
	"strings"       // we use this the strings library to collect output from ffmpeg
)

//Currently throws file out next to wherever it's run from.
//Also returns a string of the absolute path of the generated file

// Expects a string. In this case the full path of the file to be converted. Returns output.
func ConvertFile(filePath string) (string, error) {

	log.Println("Starting Lute converter...")

	//Error handling for existence of ffmpeg
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Printf("Something went wrong retrieving FFMPEG: %q\n", err)
		return "", err
	} else {
		log.Printf("Found ffmpeg at %q\n", ffmpegPath)
	}

	//Error handling for existence of provided file path
	_, err = os.Stat(filePath)
	if err != nil {
		log.Printf("Something went wrong retrieving filepath: %q\n", err)
		return "", err
	} else {
		log.Printf("Found file at %q\n", filePath)
	}

	//Creates a slice filled with the relevant arguments to run the FFMPEG exec.Command()
	//Also logs the literal output of the slice.
	manualArg := []string{"-i", filePath, "-c:a", "aac", "-f", "adts", "output.aac"}
	log.Printf("Literal manualArg slice output: %#v\n", manualArg)

	// exec.Command lets us compile a command as an object before we execute it.
	// that way we can programatically construct them!
	ffmpegCommand := exec.Command("ffmpeg", manualArg...)

	//Use strings.Builder since it handles memory efficiently
	var ffmpegOutput strings.Builder

	ffmpegCommand.Stdout = &ffmpegOutput

	//Will check for errors in command. Like if the file already exists.
	//So make sure to do something with it each run.
	err = ffmpegCommand.Run()
	if err != nil {
		log.Printf("Something went wrong executing FFMPEG command: %q\n", err)
		return "", err
	}

	// Logs output with format string
	log.Printf("Executing '%s' %s \n", ffmpegCommand, ffmpegOutput.String())

	//Sets return
	pathOutput, _ := filepath.Abs("output.aac")
	log.Println("Output file absolute path:", pathOutput)

	return pathOutput, nil
}
