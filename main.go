package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Context struct {
	Memory      map[string]string
	Bot         map[string]string
	Topic       string
	aimlRoot    *AIMLRoot
	LastRecv    string
	LastSent    string
	ThatMatches []string
}

func InitContext() Context {
	var context Context
	context.Memory = make(map[string]string)
	context.Bot = make(map[string]string)

	context.Bot["name"] = "StupidBot"

	return context
}

func main() {
	var err error

	context := InitContext()

	context.aimlRoot, err = Parse("sample/alice.aiml")
	if err != nil {
		log.Fatal(err)
	}

	inputFile, err := os.Open("sample/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()

		output, err := context.run(line)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("User %s\n>> %s\n", line, strings.TrimSpace(output))
	}
}
