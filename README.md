# scientiaCLI

A command line interface for [scientia](https://scientia.doc.ic.ac.uk) for Imperial DoC students to download materials faster and easier.

## Installation

Right now we only support automatic Linux installation, but the code should work for Windows and MacOS.

### Linux

```bash
wget https://raw.githubusercontent.com/goDoCer/scientiaCLI/main/linux-installer.sh
sh linux-installer.sh
```

### Windows/MacOS

Coming soon (Make an issue if you want this really quick)

## How to use

1. Login to Scientia. `scientia-cli login`
2. Create a directory to save all your files to `mkdir <save-directory>`
3. Set the save directory, `scientia-cli save-dir <save-directory>`
4. Download all the materials `scientia-cli download all`


**Note that the `download` command only downloads files that have been updated on scientia**. It will only download a file, if that file on scientia is newer than the file present on your machine. **If you take notes directly in the provided slides/pdfs** please run the download command with the `--new-only` flag so that you do not lose your work (`scientia-cli download --new-only all`).

To only download the materials for a particular course, say COMP40009 Computing Practical, run `scientia-cli download 40009`.


To request features/report bugs, feel free to raise an issue or even create a PR :)
