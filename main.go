package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
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

// Constants for NoteIt.
const (
	// Current goserver version
	Version = "0.0.1"
	Path    = "./" // should be able to be changed
)

// NoteItSession stores session data for NoteIt.
type NoteItSession struct {
	// TODO: this should be pulled from JSON
	UserDir         string
	NotebookPath    string
	DefaultNotebook string
	DefaultEditor   string
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
	var useNotebook = flag.String("n", "", "flag to specify creation of new notebook")
	var addNote = flag.String("a", "", "quick add using command line arg")
	var editNote = flag.String("e", "", "open note in default editor")

	flag.Parse()

	session := getSessionDetails()

	if len(os.Args) < 3 {
		flag.Usage()
		log.Fatalf("USAGE: noteit -<n/a/e> <details>")
	}

	if *useNotebook != "" {
		session.setNotebookPath(*useNotebook)
		session.getNotebook(*useNotebook)
	}

	if *addNote != "" {
		session.addNote(*addNote)
	}

	if *editNote != "" {
		session.editNote(*editNote)
	}
}

// setNotebookPath sets path of notebook in NoteItSession struct
func (s *NoteItSession) setNotebookPath(n string) {
	notebookPath := new(strings.Builder)

	if _, err := notebookPath.WriteString(s.UserDir); err != nil {
		log.Fatalf("error writing directory name, %s\n", s.UserDir)
	}

	if _, err := notebookPath.WriteString(n); err != nil {
		log.Fatalf("error writing directory name, %s\n", n)
	}

	if _, err := notebookPath.WriteString(".md"); err != nil {
		log.Fatalf("error writing directory name, %s\n", n)
	}

	s.NotebookPath = notebookPath.String()
}

// getNotebook ensures the notebook (i.e. folder) is available
// and will create new folder if folder has not been created yet
func (s *NoteItSession) getNotebook(n string) {

	_, err := os.Stat(s.NotebookPath)

	if os.IsNotExist(err) {
		// create directory
		f, err := os.Create(s.NotebookPath)
		defer f.Close()
		if err != nil {
			log.Fatalf("Unable to create notebook: %s\n", s.NotebookPath)
		}

		_, err = f.WriteString("# ")
		if err != nil {
			log.Fatalf("error writing to newly created notebook %v, %v\n", s.NotebookPath, err)
		}

		_, err = f.WriteString(n)
		if err != nil {
			log.Fatalf("error writing to newly created notebook %v, %v\n", s.NotebookPath, err)
		}

		_, err = f.WriteString("\n\n")
		if err != nil {
			log.Fatalf("error writing to newly created notebook %v, %v\n", s.NotebookPath, err)
		}
	}
}

// addNote adds a note to the selected notebook
func (s *NoteItSession) addNote(n string) {
	// appned to s.NotebookPath
	// if NotebookPath == nil
	// add to defualt notebook path

	if s.NotebookPath == "" {
		// append to notebook
		fmt.Printf("No notebook specified, will add to default notebook\n")
		s.setNotebookPath("default")
		s.getNotebook("default")
	}

	f, err := os.OpenFile(s.NotebookPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer f.Close()

	_, err = f.WriteString("- ")
	if err != nil {
		log.Fatalf("error writing to file: %v. Error: %v\n", s.NotebookPath, err)
	}

	_, err = f.WriteString(n)
	if err != nil {
		log.Fatalf("error writing to file: %v. Error: %v\n", s.NotebookPath, err)
	}

	_, err = f.WriteString("\n")
	if err != nil {
		log.Fatalf("error writing new line file: %v. Error: %v\n", s.NotebookPath, err)
	}

	// alert user note has been updated
	fmt.Printf("%v note updated!\n", s.NotebookPath)
}

// editNote opens specified note in vim
func (s *NoteItSession) editNote(n string) {
	s.setNotebookPath(n)

	cmd := exec.Command("nvim", s.NotebookPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running cmd: %v\n", cmd.Args[:])
	}
}

func getSessionDetails() *NoteItSession {
	// TODO: make this read json file
	p := new(NoteItSession)
	p.UserDir = os.Getenv("HOME")
	var buffer bytes.Buffer
	buffer.WriteString(p.UserDir)
	buffer.WriteString("/noteit/")
	p.UserDir = buffer.String()
	return p
}
