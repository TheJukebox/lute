package main

import (
	"log"     // providing logging for us - not necessary, but I like to start logging ASAP.
	"os/exec" // this is a builtin for executing commands on the host OS
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
	// exec.Command lets us compile a command as an object before we execute it.
	// that way we can programatically construct them!
	ffmpegCommand := exec.Command(ffmpegPath, "--help")

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
}
