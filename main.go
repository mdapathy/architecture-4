package main

import (
	"bufio"
	"fmt"
	"github.com/mdapathy/architecture-4/commands"
	"github.com/mdapathy/architecture-4/engine"
	"log"
	"os"
)

func processCLParameters() string {

	if len(os.Args) < 2 {
		log.Fatalf("Not enough parameters specified, required 1, got %d", len(os.Args)-1)
	} else if dir, err := os.Stat(os.Args[1]); os.IsNotExist(err) || dir.IsDir() {
		log.Fatalf("Improper path to file %s ", os.Args[1])

	} else if len(os.Args) > 2 {
		fmt.Println("Ignoring all parameters except " + os.Args[1])
	}

	return os.Args[1]

}

func main() {
	inputFile := processCLParameters()

	eventLoop := new(engine.EventLoop)

	eventLoop.Start()
	input, err := os.Open(inputFile)

	if err != nil {
		log.Fatalf("Error occured while reading file: " + err.Error())
	}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		scanner.Text()
		commandLine := scanner.Text()
		cmd := commands.Parse(commandLine) // parse the line to get an instance of Command cmd
		eventLoop.Post(cmd)
	}

	if err := input.Close(); err != nil {
		log.Fatalf("Error while closing the file: " + err.Error())
	}

	eventLoop.AwaitFinish()
}
