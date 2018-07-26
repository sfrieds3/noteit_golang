# noteit
noteit is a CLI note taking tool, written in golang. Notes are filed in
notebooks, which are organized as directories under the noteit directory
(default is ~/noteit/). Usage examples are below:

---

## Usage

file system:

        noteit
        /     \
       /       \
      /         \
  notebook1   notebook2
  /     \        /   \
 /       \      /     \
 note1  note2  note3  note4

```bash
$ noteit -a <noteName> -n <notebookName> <note> # append to note in specified notebook
$ noteit -e <noteName> # edit note in default notebook - will open in default editor
$ noteit -n <notebookName> # check if notebook exists, create if not
$ noteit -v <noteName> -n <notebookName> # view note in notebook
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
