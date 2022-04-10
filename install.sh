go build
sudo mkdir -p /usr/local/bin/scientiaCLI
sudo mv scientiaCLI /usr/local/bin/scientiaCLI/

sudo touch /usr/local/bin/scientiaCLI/token.txt  
sudo chmod -R 777 /usr/local/bin/scientiaCLI/token.txt

sudo cp ./default-config.json /usr/local/bin/scientiaCLI/config.json
sudo chmod -R 777 /usr/local/bin/scientiaCLI/config.json

sudo cp autocomplete.sh /etc/bash_completion.d/scientiaCLI
sudo chmod +x /etc/bash_completion.d/scientiaCLI

# NOT SURE IF WE NEED THIS?
# echo "[ -r /usr/share/bash-completion/bash_completion ] && . /usr/share/bash-completion/bash_completion" >> ~/.bashrc

echo "Scientia installed successfully"
echo "Scientia is now available in /usr/local/bin/scientiaCLI"
# echo "export PATH=\$PATH:/usr/local/bin/scientiaCLI" >> ~/.zshrc

# UNCOMMENT FOR NORMAL BASH 
# echo "export PATH=\$PATH:/usr/local/bin/scientiaCLI" >> ~/.bashrc

# UNCOMMENT FOR ZSH 
# echo "export PATH=\$PATH:/usr/local/bin/scientiaCLI" >> ~/.zshrc

# UNCOMMENT FOR FISH 
# set -U fish_user_paths /usr/local/go/bin $fish_user_paths

# sudo wget -O /etc/bash_completion.d/scientiaCLI https://raw.githubusercontent.com/urfave/cli/master/autocomplete/bash_autocomplete && sudo chmod +x /etc/bash_completion.d/scientiaCLI