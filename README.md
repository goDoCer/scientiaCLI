# scientiaCLI

A command line interface for [scientia](https://scientia.doc.ic.ac.uk) for Imperial DoC students to download materials faster and easier.

## Installation

Right now we only support automatic Linux installation, but the code should work for Windows and MacOS.

### Linux

1. Make sure you have go installed ([instructions](https://go.dev/doc/install))
2. Clone this repository: `git clone git@github.com:goDoCer/scientiaCLI.git`
3. Run the installer script in the root of the repo: `cd scientiaCLI && sh install.sh`
4. Feel free to remove the repository: `rm -rf scientiaCLI`

### Windows/MacOS

Coming soon (Make an issue if you want this really quick)

## How to use

1. Login to Scientia. `scientia-cli login`
2. Create a directory to save all your files to `mkdir <save-directory>`
3. Set the save directory, `scientia-cli save-dir <save-directory>`
4. Download all the materials `scientia-cli download all`

To only download the materials for a particular course, say COMP40009 Computing Practical, run `scientia-cli download 40009`.

## How to install Go
1. Get the latest tar for your OS from [here](https://go.dev/dl/) by running `wget <link you just copied>`
2. Run `rm -rf /usr/local/go && tar -C /usr/local -xzf <go1.XX.X.os-arch.tar.gz>`
3. Add go to your PATH by running `echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc`
4. Run `source ~/.bashrc` to update your PATH

To request features/report bugs, feel free to raise an issue or even create a PR :)
