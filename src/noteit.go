package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"time"
)

/*
- notes separated into directories
- each directory can have multiple notes
- each note and directory can have tags
- command line args:
	- a: add new or append to existing note
	- n: specify notebook to use (create new notebook if necessary)
	- v: view note
	- va view all notes in directory
- for notes:
	- quick add, append to existing note, or add to new note
	- for more in depth editing, can open editor of your choice
- process for adding note:
	- noteit -a <noteName> -n <notebookName> <note>
		- if <note> not specified, open vim or other editor
	- default to noteName = notebookName if not specified
	- must specify notebook to add to
	- note name is optional
	- default action is to append to noteName notebookName
*/

const (
	// Current goserver version
	Version = "0.0.1"
	Path    = "./" // should be able to be changed
)

type NoteItSession struct {
	UserDir string
}

type Note struct {
	notebook     Notebook
	noteBody     string
	dateAdded    time.Time
	lastModified time.Time
}

type Notebook struct {
	numNotes    int
	dateAdded   time.Time
	lasModified time.Time
	tags        string
}

func main() {
	session := startNoteItSession()
	userArgs := os.Args[0:]
	switch userArgs[1] {
	case "n":
		createNewNotebook(userArgs[2])
	case "a":
		addNote(userArgs[2])
	case "v":
		viewNote(userArgs[2])
	}
}

func addNote(input string) {
	fmt.Printf("Add to note: %v\n", input)
}

func viewNote(input string) {
	fmt.Printf("User would like to view all notes")
}

func createNewNotebook(input string) {
	fmt.Printf("User would like to create a new notebook: %v\n", input)
	// create new Notebook struct
	// prompt user for new notebook name and tags
	// create new folder
	// allow user to add README?
}

func startNoteItSession() NoteItSession {
	p := NoteItSession
	if usr, err := user.Current(); err != nil {
		p.UserDir = usr.HomeDir
		fmt.Println("User directory: %v", p.UserDir)
	} else {
		log.Fatal(err)
	}
	return p
}
