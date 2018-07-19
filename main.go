package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Constants for NoteIt.
// Constants include: current version of noteit
const (
	Version = "0.0.2"
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

func main() {
	var useNotebook = flag.String("n", "", "flag to specify creation of new notebook")
	var addNote = flag.String("a", "", "quick add using command line arg")
	var editNote = flag.String("e", "", "open note in default editor")
	//var viewNote = flag.String("v", "", "view specified note")
	//var viewNotebooks = flag.String("vn", "", "view all notebook names")

	flag.Parse()

	session := getSessionDetails()

	if len(os.Args) < 3 && len(os.Args) > 1 {
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
	if s.NotebookPath == "" {
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

func (s *NoteItSession) printNote(n string) {
	// show all notes in notebook n
}

// getSessionDetails sets up the current NoteItSession.
// This holds details such as user's $HOME, and default note.
// These can be overridden by the noteit.json file.
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
