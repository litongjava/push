package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fileName := "push.txt"

	fmt.Println("Starting push execution...")

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		fmt.Printf("Executing line %d: %s\n", lineNumber, line)

		// Split command and arguments
		parts := strings.Fields(line)
		cmdName := parts[0]
		args := parts[1:]

		cmd := exec.Command(cmdName, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error executing command on line %d: %v\n", lineNumber, err)
		} else {
			fmt.Printf("Command on line %d executed successfully\n", lineNumber)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	fmt.Println("Push execution finished.")
}
