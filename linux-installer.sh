#!/bin/bash
set -e
curl -s https://api.github.com/repos/goDoCer/scientiaCLI/releases/latest \
    | grep "browser_download_url" \
    | cut -d : -f 2,3 \
    | tr -d \" \
    | wget -O scientia-cli -qi -
    
sudo mkdir -p /usr/local/bin/scientia-cli
sudo mv scientia-cli /usr/local/bin/scientia-cli/

sudo touch /usr/local/bin/scientia-cli/token.txt  
sudo chmod -R 777 /usr/local/bin/scientia-cli/token.txt

sudo cp ./default-config.json /usr/local/bin/scientia-cli/config.json
sudo chmod -R 777 /usr/local/bin/scientia-cli/config.json
set +e

cat << EndOfMessage
======================================================================================================
Add scientia-cli to your path. You can do so by running the command depending upon which shell you use

# BASH
echo "export PATH=\\\$PATH:/usr/local/bin/scientia-cli" >> ~/.bashrc

# ZSH 
echo "export PATH=\\\$PATH:/usr/local/bin/scientia-cli" >> ~/.zshrc

# FISH 
set -U fish_user_paths /usr/local/go/bin \$fish_user_paths

Run source ~/.bashrc afterwards
======================================================================================================
EndOfMessage

echo "Scientia installed successfully"
echo "Scientia is now available in /usr/local/bin/scientia-cli"