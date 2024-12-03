package convert

import (
	"log"           // providing logging for us - not necessary, but I like to start logging ASAP.
	"os"            // for using Stat to check if files exist
	"os/exec"       // this is a builtin for executing commands on the host OS
	"path/filepath" //Looks for the path of a file
	"strings"       // we use this the strings library to collect output from ffmpeg
)

//Currently spits file out next to wherever it's run from.
//Also returns a string of the absolute path of the generated file

//FFMPEG fails to parse properly when introducing whitespace into path.

// Expects a string. In this case the full path of the file to be converted.
func ConvertFile(filePath string) (returnPath string) {

	log.Println("Starting Lute converter...")

	//Error handling for existence of ffmpeg.
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		log.Printf("Unable to find ffmpeg. Terminating...")
		log.Fatal(err)
	} else {
		log.Printf("Found ffmpeg at %q\n", ffmpegPath)
	}

	inputPath, err := os.Stat(filePath)
	if err != nil {
		log.Printf("File not located. Please check your path. Terminating...")
		log.Fatal(err)
	} else {
		log.Printf("Path for file %q is valid", inputPath.Name())
	}

	//Puts together the ffmpeg command with a string. Fields also screams with path whitespace
	commandArg := "-i " + filePath + " -c:a aac output.aac"
	//commandArg := fmt.Sprintf("-i %s -c:a aac output.aac", filePath)
	argParts := strings.Fields(commandArg)

	// exec.Command lets us compile a command as an object before we execute it.
	// that way we can programatically construct them!
	ffmpegCommand := exec.Command(ffmpegPath, argParts...)

	//Use strings.Builder since it handles memory efficiently
	var ffmpegOutput strings.Builder

	ffmpegCommand.Stdout = &ffmpegOutput
	ffmpegCommand.Run()

	// Logs output with format string
	log.Printf("Executing '%s' %s \n", ffmpegCommand, ffmpegOutput.String())

	//Will need to be changed when we specify where we actually want
	//to output converted files later.

	pathOutput, err := os.Stat("output.aac")
	if err != nil {
		log.Fatal(pathOutput, err)
	} else {
		pathOutput, _ := filepath.Abs("output.aac")
		log.Print("Output file absolute path: ", pathOutput)
		returnPath = pathOutput

	}
	return returnPath
}
