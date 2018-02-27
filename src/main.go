package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
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
	Notebook     Notebook
	NoteBody     string
	DateAdded    time.Time
	LastModified time.Time
}

type Notebook struct {
	Name     string
	NumNotes int
	Tags     string
}

func main() {
	//session := startNoteItSession()
	if len(os.Args) < 3 {
		log.Fatalf("USAGE: noteit -<n/a/v> <details>")
	}
	switch os.Args[1] {
	case "-n":
		createNewNotebook(os.Args[2])
	case "-a":
		addNote(os.Args[2:])
	case "-v":
		viewNote(os.Args[2])
	case "-e":
		editNote(os.Args[2])
	}
}

func addNote(input []string) {
	fmt.Printf("Add to note: %v\n", strings.Join(input[:], " "))
	// open correct notebook
	// append to correct file within notebook
}

func viewNote(input string) {
	fmt.Printf("User would like to view all notes")
	// print contents of note to screen
}

func createNewNotebook(input string) {
	fmt.Printf("User would like to create a new notebook: %v\n", input)
	notebook := new(Notebook)
	notebook.Name = input
	notebook.NumNotes = 0
	// TODO: get tags from user input
	notebook.Tags = ""
	// create new Notebook struct
	// prompt user for new notebook name and tags
	// create new folder
	// allow user to add README?
}

func editNote(input string) {
	fmt.Printf("User would like to edit notebook: %v\n", input)
	// open requested note in vim (or default editor)
}

func startNoteItSession() *NoteItSession {
	p := new(NoteItSession)
	if usr, err := user.Current(); err != nil {
		p.UserDir = usr.HomeDir
		fmt.Println("User directory: %v", p.UserDir)
	} else {
		log.Fatal(err)
	}
	return p
}
