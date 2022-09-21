#!/bin/bash
set -e

if [ -z "${SCIENTIA_DEV}" ]; then
    # Fetch latest release form github
    curl -s https://api.github.com/repos/goDoCer/scientiaCLI/releases/latest \
    | grep "browser_download_url" \
    | cut -d : -f 2,3 \
    | tr -d \" \
    | wget -qO scientia-cli -qi -
    
    # Fetch default config file
    CFG=$(wget -qO - https://raw.githubusercontent.com/goDoCer/scientiaCLI/main/default-config.json | cat)
else
    CFG=$(cat default-config.json)
    go build -o scientia-cli
fi

chmod +x scientia-cli

INSTALL_DIR="$HOME/.scientia-cli"
TOKEN_FILE="$INSTALL_DIR/token.txt"
CFG_FILE="$INSTALL_DIR/config.json"

mkdir -p "$INSTALL_DIR"
mv scientia-cli "$INSTALL_DIR"

touch "$TOKEN_FILE"
chmod -R 777 "$TOKEN_FILE"

if [ ! -f "$CFG_FILE" ]; then
    touch "$CFG_FILE"
    chmod -R 777 "$CFG_FILE"
    
    echo "$CFG" | tee "$CFG_FILE" > /dev/null
fi

set +e

cat << EndOfMessage
======================================================================================================
Add scientia-cli to your path. You can do so by running the command depending upon which shell you use

# BASH
echo "export PATH=\\\$PATH:$HOME/.scientia-cli" >> ~/.bashrc && source ~/.bashrc

# ZSH
echo "export PATH=\\\$PATH:$HOME/.scientia-cli" >> ~/.zshrc && source ~/.zshrc

# FISH
set -U fish_user_paths $HOME/.scientia-cli \$fish_user_paths

======================================================================================================
EndOfMessage

echo "scientia-cli installed successfully"
echo "scientia-cli is now installed at $INSTALL_DIR"

echo "to uninstall run rm -r $INSTALL_DIR at any time"
