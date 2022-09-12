#!/bin/bash
set -e

if [ -z "${SCIENTIA_DEV}" ]; then
    # Fetch latest release form github
    curl -s https://api.github.com/repos/goDoCer/scientiaCLI/releases/latest \
    | grep "browser_download_url" \
    | cut -d : -f 2,3 \
    | tr -d \" \
    | wget -O scientia-cli -qi -
    
    # Fetch default config file
    CFG=$(wget -O - https://raw.githubusercontent.com/goDoCer/scientiaCLI/main/default-config.json | cat)
    
else
    CFG=$(cat default-config.json)
    go build -o scientia-cli
fi

chmod +x scientia-cli

INSTALL_DIR="/usr/local/bin/scientia-cli"
TOKEN_FILE="$INSTALL_DIR/token.txt"
CFG_FILE="$INSTALL_DIR/config.json"

sudo mkdir -p $INSTALL_DIR
sudo mv scientia-cli $INSTALL_DIR

sudo touch $TOKEN_FILE
sudo chmod -R 777 $TOKEN_FILE

if [ ! -f $CFG_FILE ]; then
    sudo touch $CFG_FILE
    printf "%s" "$CFG" > $CFG_FILE
    sudo chmod -R 777 $CFG_FILE
fi

set +e

cat << EndOfMessage
======================================================================================================
Add scientia-cli to your path. You can do so by running the command depending upon which shell you use

# BASH
echo "export PATH=\\\$PATH:/usr/local/bin/scientia-cli" >> ~/.bashrc

# ZSH
echo "export PATH=\\\$PATH:/usr/local/bin/scientia-cli" >> ~/.zshrc

# FISH
set -U fish_user_paths /usr/local/bin/scientia-cli \$fish_user_paths

Run source ~/.bashrc afterwards
======================================================================================================
EndOfMessage

echo "Scientia installed successfully"
echo "Scientia is now available in /usr/local/bin/scientia-cli"
