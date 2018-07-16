# noteit
CLI note taking tool, written in golang

---

## Usage

```bash
$ noteit -a <noteName> -n <notebookName> <note> # add note
$ noteit -e <noteName> # edit note - will open in default editor
$ noteit -n <notebookName> # create new notebook
$ noteit -v <noteName> # view note
$ noteit -va <notebookName> # view all notes in notebook
$ noteit -vn # view all notebooks
```


## Overview

- notes separated into directories
- each directory can have multiple notes
- each note and directory can have tags
- command line args:
        - a: add new or append to existing note
        - n: specify notebook to use (create new notebook if necessary)
        - v: view note [TODO]
        - va view all notes in directory [TODO]
        - vd: view all directory names [TODO]
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
