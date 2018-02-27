package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
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
	session := getSessionDetails()
	fmt.Printf("This session home directory: %v\n", session.UserDir)
	if len(os.Args) < 3 {
		log.Fatalf("USAGE: noteit -<n/a/v> <details>")
	}
	switch os.Args[1] {
	case "-n":
		session.createNewNotebook(os.Args[2])
	case "-a":
		session.addNote(os.Args[2:])
	case "-v":
		session.viewNote(os.Args[2])
	case "-e":
		session.editNote(os.Args[2])
	}
}

func (s *NoteItSession) addNote(input []string) {
	fmt.Printf("Add to note: %v\n", strings.Join(input[:], " "))
	// open correct notebook
	// append to correct file within notebook
}

func (s *NoteItSession) viewNote(input string) {
	fmt.Printf("User would like to view all notes")
	// print contents of note to screen
}

func (s *NoteItSession) createNewNotebook(input string) {
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
	var filename bytes.Buffer
	filename.WriteString(s.UserDir)
	filename.WriteString(input)
	os.Mkdir(filename.String(), 644)
}

func (s *NoteItSession) editNote(input string) {
	fmt.Printf("User would like to edit notebook: %v\n", input)
	// open requested note in vim (or default editor)
}

func getSessionDetails() *NoteItSession {
	p := new(NoteItSession)
	p.UserDir = os.Getenv("HOME")
	var buffer bytes.Buffer
	buffer.WriteString(p.UserDir)
	buffer.WriteString("/noteit/")
	p.UserDir = buffer.String()
	return p
}
