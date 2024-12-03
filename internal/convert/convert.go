package convert

import (
	"log"           // providing logging for us - not necessary, but I like to start logging ASAP.
	"os"            // for using Stat to check if files exist
	"os/exec"       // this is a builtin for executing commands on the host OS
	"path/filepath" //Looks for the path of a file
	"strings"       // we use this the strings library to collect output from ffmpeg
)

//Currently spits it out next to wherever it's run from.
//Still also fails to parse properly when introducing whitespace into path.

// Expects a string. In this case the full path of the file to be converted.
func ConvertFile(filePath string) string {

	log.Println("Starting Lute converter from INTERNAL...")

	//Error handling for existence of ffmpeg.
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Printf("Unable to find ffmpeg. Terminating...")
		log.Fatal(err)
	} else {
		log.Printf("Found ffmpeg at %q\n", ffmpegPath)
	}

	//Will probably have to do something if path has whitespace.
	inputPath, err := os.Stat(filePath)
	if err != nil {
		log.Printf("File not located. Please check your path. Terminating...")
		log.Fatal(err)
	} else {
		log.Printf("Path located at %q is valid\n", inputPath.Name())
	}

	//Puts together the ffmpeg command with a string. Fields also screams with path whitespace
	commandArg := "-i " + filePath + " -c:a aac output.aac"
	argParts := strings.Fields(commandArg)

	// exec.Command lets us compile a command as an object before we execute it.
	// that way we can programatically construct them!
	ffmpegComm := exec.Command(ffmpegPath, argParts...)

	//Use strings.Builder since it handles memory efficiently
	var ffmpegOutput strings.Builder

	ffmpegComm.Stdout = &ffmpegOutput
	ffmpegComm.Run()

	// Logs output with format string
	log.Printf("Executing '%s'\n %s\n", ffmpegComm, ffmpegOutput.String())

	outputPath, err := filepath.Abs("output.aac")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("Output file absolute path: ", outputPath)
	}

	return outputPath
}
