go build
sudo mkdir -p /usr/local/bin/scientia-cli
sudo mv scientia-cli /usr/local/bin/scientia-cli/

sudo touch /usr/local/bin/scientia-cli/token.txt  
sudo chmod -R 777 /usr/local/bin/scientia-cli/token.txt

echo "Scientia installed successfully"
echo "Scientia is now available in /usr/local/bin/scientia-cli"
# echo "export PATH=\$PATH:/usr/local/bin/scientia-cli" >> ~/.bashrc

# UNCOMMENT FOR NORMAL BASH 
# echo "export PATH=\$PATH:/usr/local/bin/scientia-cli" >> ~/.bashrc

# UNCOMMENT FOR FISH 
# set -U fish_user_paths /usr/local/go/bin $fish_user_paths

