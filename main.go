// Command noteit is a small CLI for taking notes in markdown notebooks.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const version = "0.1.0"

const usage = `usage: noteit <command> [arguments]

commands:
  add [-n notebook] <note...>   append a note to a notebook
  edit [notebook]               open a notebook in your editor
  view [notebook]               print a notebook
  list                          list all notebooks
  version                       print the version

The default notebook is used when no notebook is given.
`

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "noteit: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		fmt.Fprint(stdout, usage)
		return nil
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
		return nil
	case "version", "--version":
		fmt.Fprintf(stdout, "noteit %s\n", version)
		return nil
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	store := Store{Dir: cfg.Path}

	switch cmd {
	case "add":
		return cmdAdd(store, cfg, args, stdout)
	case "edit":
		return cmdEdit(store, cfg, args)
	case "view":
		return cmdView(store, cfg, args, stdout)
	case "list", "ls":
		return cmdList(store, stdout)
	default:
		return fmt.Errorf("unknown command %q\n\n%s", cmd, usage)
	}
}

func cmdAdd(store Store, cfg Config, args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	notebook := fs.String("n", cfg.DefaultNotebook, "notebook to add the note to")
	if err := fs.Parse(args); err == flag.ErrHelp {
		return nil
	} else if err != nil {
		return err
	}

	note := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if note == "" {
		return fmt.Errorf("nothing to add\n\nusage: noteit add [-n notebook] <note...>")
	}

	path, err := store.Add(*notebook, note)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "added to %s\n", path)
	return nil
}

func cmdEdit(store Store, cfg Config, args []string) error {
	name, err := notebookArg(cfg, args)
	if err != nil {
		return err
	}
	path, err := store.Ensure(name)
	if err != nil {
		return err
	}

	editor := cfg.editorCommand()
	cmd := exec.Command(editor[0], append(editor[1:], path)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running editor %q: %w", strings.Join(editor, " "), err)
	}
	return nil
}

func cmdView(store Store, cfg Config, args []string, stdout io.Writer) error {
	name, err := notebookArg(cfg, args)
	if err != nil {
		return err
	}
	data, err := store.Read(name)
	if err != nil {
		return err
	}
	_, err = stdout.Write(data)
	return err
}

func cmdList(store Store, stdout io.Writer) error {
	names, err := store.List()
	if err != nil {
		return err
	}
	for _, n := range names {
		fmt.Fprintln(stdout, n)
	}
	return nil
}

// notebookArg returns the single optional notebook argument, or the
// configured default notebook when none is given.
func notebookArg(cfg Config, args []string) (string, error) {
	switch len(args) {
	case 0:
		return cfg.DefaultNotebook, nil
	case 1:
		return args[0], nil
	default:
		return "", fmt.Errorf("expected at most one notebook, got %d", len(args))
	}
}
