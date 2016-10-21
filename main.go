package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Context struct {
	Memory   map[string]string
	Topic    string
	aimlRoot *AIMLRoot
	LastSent string
}

func main() {
	var err error
	var context Context
	context.Memory = make(map[string]string)

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
