package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
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

// Constants for NoteIt.
const (
	// Current goserver version
	Version = "0.0.1"
	Path    = "./" // should be able to be changed
)

// NoteItSession stores session data for NoteIt.
type NoteItSession struct {
	UserDir string
}

// Note struct contains details about the note to be saved.
type Note struct {
	Notebook     Notebook
	NoteBody     string
	DateAdded    time.Time
	LastModified time.Time
}

// Notebook struct that contains details about notebook.
type Notebook struct {
	Name     string
	NumNotes int
	Tags     string
}

func main() {
	var createNotebook = flag.String("n", "", "flag to specify creation of new notebook")
	var addNote = flag.String("a", "", "flag to add new note")
	var editNote = flag.String("e", "", "flag to edit note by opening note in default editor")
	//var note = flag.String("n", "", "contents of note")

	flag.Parse()

	fmt.Printf("FLAGS: createNotebook: %v, addNote: %v, editNote: %v\n", *createNotebook, *addNote, *editNote)

	session := getSessionDetails()
	fmt.Printf("This session home directory: %v\n", session.UserDir)
	if len(os.Args) < 3 {
		log.Fatalf("USAGE: noteit -<n/a/v> <details>")
	}
	if *createNotebook != "" {
		session.createNewNotebook(*createNotebook)
	}

	if *addNote != "" {
		session.addNote(*addNote)
	}

	if *editNote != "" {
		session.editNote(*editNote)
	}
}

func (s *NoteItSession) addNote(input string) {
	fmt.Printf("Add to note: %v\n", input)

	// buffer to create notebook path
	nBook := s.getNotebookPath(input)
	fmt.Printf("notebook to add note to: %v\n", nBook)
	// append to correct file within notebook

	notePath := s.getNotePath(nBook, input)
	fmt.Printf("note path: %v\n", notePath)

	// TODO: fix permissions
	f, err := os.OpenFile(notePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 755)
	defer f.Close()
	if err != nil {
		log.Fatalf("Notebook %s does not exist. Please create with -n and try again.", nBook)
	}

	f.WriteString("\n")
	//f.WriteString(strings.Split(input[2:], " "))
	f.WriteString("This is a test")
}

func (s *NoteItSession) viewNote(input string) {
	fmt.Printf("User would like to view note")
	// print contents of note to screen
}

func (s *NoteItSession) createNewNotebook(input string) {
	fmt.Printf("User would like to create a new notebook: %v\n", input)
	notebook := new(Notebook)
	notebook.Name = input
	notebook.NumNotes = 0

	// TODO: get tags from user input
	notebook.Tags = ""

	// get filename for notebook
	filename := s.getNotebookPath(input)
	// TODO: fix permissions
	os.Mkdir(filename, 755)
}

func (s *NoteItSession) getNotePath(notebookPath, note string) string {
	var noteFile bytes.Buffer
	noteFile.WriteString(notebookPath)
	noteFile.WriteString(note)
	return noteFile.String()
}

func (s *NoteItSession) getNotebookPath(nBook string) string {
	var notebookPath bytes.Buffer
	notebookPath.WriteString(s.UserDir)
	notebookPath.WriteString(nBook)
	notebookPath.WriteString("/")
	return notebookPath.String()
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
