package main

import (
	"fmt"
	"os"
)

/*
- notes separated into directories
- each directory can have multiple notes
- each note can have tags
- directories can have tags???
- command line args:
	-
*/

const (
	// Current goserver version
	Version = "0.0.1"
)

func main() {
	userArgs := os.Args[0:]
	switch userArgs[1] {
	case "n":
		newNote(userArgs[2])
	case "a":
		addNote(userArgs[2])
	case "d":
		deleteNote(userArgs[2])
	}
}

func deleteNote(input string) {
	fmt.Printf("Delete a note: %v\n", input)
}

func addNote(input string) {
	fmt.Printf("Add to note:%v\n", input)
}

func newNote(input string) {
	fmt.Printf("User would like to create a new notebook%v\n", input)
}
