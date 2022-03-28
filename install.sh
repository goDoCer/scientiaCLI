go build
sudo mkdir -p /usr/local/bin/scientia-cli
sudo mv scientia-cli /usr/local/bin/scientia-cli/

sudo touch /usr/local/bin/scientia-cli/token.txt  
sudo chmod -R 777 /usr/local/bin/scientia-cli/token.txt

sudo cp ./config.json /usr/local/bin/scientia-cli/config.json

sudo cp autocomplete.sh /etc/bash_completion.d/scientia-cli
sudo chmod +x /etc/bash_completion.d/scientia-cli

# NOT SURE IF WE NEED THIS?
# echo "[ -r /usr/share/bash-completion/bash_completion ] && . /usr/share/bash-completion/bash_completion" >> ~/.bashrc

echo "Scientia installed successfully"
echo "Scientia is now available in /usr/local/bin/scientia-cli"
# echo "export PATH=\$PATH:/usr/local/bin/scientia-cli" >> ~/.zshrc

# UNCOMMENT FOR NORMAL BASH 
# echo "export PATH=\$PATH:/usr/local/bin/scientia-cli" >> ~/.bashrc

# UNCOMMENT FOR ZSH 
# echo "export PATH=\$PATH:/usr/local/bin/scientia-cli" >> ~/.zshrc

# UNCOMMENT FOR FISH 
# set -U fish_user_paths /usr/local/go/bin $fish_user_paths

# sudo wget -O /etc/bash_completion.d/scientia-cli https://raw.githubusercontent.com/urfave/cli/master/autocomplete/bash_autocomplete && sudo chmod +x /etc/bash_completion.d/scientia-cli