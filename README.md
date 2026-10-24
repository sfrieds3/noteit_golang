# noteit

noteit is a small CLI note-taking tool written in Go. Notes are kept in
**notebooks**, each of which is a plain markdown file under the noteit
directory (default `~/noteit/`):

```
~/noteit/
├── default.md
├── work.md
└── ideas.md
```

Each note is appended as a markdown list item, so notebooks stay readable and
easy to edit by hand.

## Install

```bash
go install github.com/sfrieds3/noteit_golang@latest
```

Or build from a checkout with `go build -o noteit .`.

## Usage

```bash
noteit add buy milk                 # append "- buy milk" to the default notebook
noteit add -n work call bob         # append to the "work" notebook (created if needed)
noteit view work                    # print a notebook
noteit view                         # print the default notebook
noteit edit work                    # open a notebook in your editor
noteit list                         # list all notebooks
noteit version
```

Flags must come before the note text: `noteit add -n work call bob`, not
`noteit add call bob -n work`.

## Configuration

noteit reads an optional JSON config from `~/.config/noteit/config.json`
(or the platform equivalent; override with `$NOTEIT_CONFIG`):

```json
{
    "path": "~/noteit",
    "editor": "nvim",
    "defaultNotebook": "default"
}
```

All fields are optional. The editor is chosen from the config `editor` value,
then `$VISUAL`, then `$EDITOR`, falling back to `vi`. Editor values may include
arguments, e.g. `"code --wait"`.

## Development

```bash
go test ./...
go vet ./...
```

## Ideas

- tags on notes and notebooks
- `noteit add` with no text opens the editor
- search across notebooks
