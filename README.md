# scientiaCLI

A command line interface for [scientia](https://scientia.doc.ic.ac.uk) for Imperial DoC students to download materials faster and easier.

Look at the [cronjob docs](https://github.com/goDoCer/scientiaCLI#cron-job) with which you never have to worry about downloading materials ever again!

## Installation

Right now we only support automatic Linux installation, but the code should work for Windows and MacOS.

### Linux

```bash
wget https://raw.githubusercontent.com/goDoCer/scientiaCLI/main/linux-installer.sh
sh install.sh
```

### Windows/MacOS

Coming soon (Make an issue if you want this really quick)

## How to use

1. Login to Scientia. `scientia-cli login`
2. Create a directory to save all your files to `mkdir <save-directory>`
3. Set the save directory, `scientia-cli save-dir <save-directory>`
4. Download all the materials `scientia-cli download all`

Note that the `download` command only downloads files that have been updated on scienttia. It will only download a file if the file on scientia is newer than the file present on your machine.

To only download the materials for a particular course, say COMP40009 Computing Practical, run `scientia-cli download 40009`.

## Cron Job

The [cron](https://en.wikipedia.org/wiki/Cron) command-line utility is a job scheduler on Unix-like operating systems. Schedule the download command to run everyday.

### Daily Cron Job

Make sure you have logged in and have set a save directory using the `login` and the `save-dir` commands respectively.

1. Create the daily cron job file `sudo touch /etc/cron.daily/scientia-cli`.
2. Open the file using vim (or any other editor) - `sudo vim /etc/cron.daily/scientia-cli`.

   ```sh
   #! /bin/bash

    # NOTE: confirm your scientia-cli installation location by running `which scientia-cli` accordingly
   /usr/local/bin/scientia-cli/scientia-cli download all
   ```

## Cron Job without sudo access/Custom cronjob

1. Run `crontab -u $USER -e`.
2. Add the following line at the end of the file to run the command at 5:30 am everyday:

   ```sh
    30 5 * * *  /usr/local/bin/scientia-cli/scientia-cli download all
   ```

   You can use this [website](https://crontab.guru/#30_5_*_*_*) to figure out what `30 5 * * *` stand for and configure thema ccordingly.

To request features/report bugs, feel free to raise an issue or even create a PR :)
