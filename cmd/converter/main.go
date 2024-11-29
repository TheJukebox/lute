package main

import (
	"fmt"
	"log"     // providing logging for us - not necessary, but I like to start logging ASAP.
	"os"      // for using Stat to check if files exist
	"os/exec" // this is a builtin for executing commands on the host OS
	"path/filepath"

	//Looks for the path of a file
	"strings" // we use this the strings library to collect output from ffmpeg
)

func main() {
	log.Println("Starting Lute converter...")
	/* most functions in Go return two values
	you can usually expect to receive the intended output and nil
	or if there's an exception, nil and the returned error

	in this case we use `:=` as a shorthand for initialising and declaring
	these two variables.
	*/
	ffmpegPath, err := exec.LookPath("ffmpeg")
	// if the err variable is NOT empty
	if err != nil {
		log.Fatal(err) // log the error and call os.Exit(1)
	} else {
		log.Printf("Found ffmpeg at %q\n", ffmpegPath) // use a format string to log the path of ffmpeg on the host
	}

	//Splits command for exec.Command to recognise as well as obtaining input from user.
	var comInput string

	//Needs to be able to accept spaces and have error handling for it
	fmt.Print("Please provide file name and extension(no spaces): ")
	fmt.Scanln(&comInput)

	//Checks if file exists
	inputFile, err := os.Stat(comInput)
	if err != nil {
		log.Fatal(inputFile, err) // log the error and call os.Exit(1)
	}
	//Obtains file path
	comArg := "-i " + comInput + " -f hls -c:a aac output.m3u8"
	argParts := strings.Fields(comArg)

	// exec.Command lets us compile a command as an object before we execute it.
	// that way we can programatically construct them!

	ffmpegCommand := exec.Command(ffmpegPath, argParts...)

	// here we're declaring a var but not setting a value for it
	var ffmpegOutput strings.Builder

	// here we say that stdout (output) of the command should be written
	// to our string. we use strings.Builder because it handles memory very
	// efficiently for writing strings.
	// the '&' operator in this case is returning the memory address of the string.
	// we're telling go to write the output to this location in memory.
	ffmpegCommand.Stdout = &ffmpegOutput
	ffmpegCommand.Run()

	// now we log the output with another format string.
	log.Printf("Executing '%s'\n %s\n", ffmpegCommand, ffmpegOutput.String())

	//Outputs absolute path of the converted file
	outputPath, err := filepath.Abs("output.m3u8")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("Output file absolute path: ", outputPath)
	}
}
