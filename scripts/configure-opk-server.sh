# Assign arguments to variables
EMAIL="$1"
USER="$2"

# Validate input arguments
if [ -z "${EMAIL}" ] || [ -z "${USER}" ]; then
    echo "Usage: $0 <EMAIL> <USER>"
    exit 1
fi

if sudo -nl &> /dev/null || $(id -u) -ne 0; then
    echo "User has sudo privileges without a password."
    # Create the default opk directory
    sudo mkdir -p /etc/opk

    # Install the opk verifier
    sudo mv opk-ssh /etc/opk
    sudo chown root /etc/opk/opk-ssh
    sudo chmod 700 /etc/opk/opk-ssh

    # Create root policy yaml file which is a yaml file
    sudo touch /etc/opk/policy.yml
    sudo chown root /etc/opk/policy.yml
    sudo chmod 600 /etc/opk/policy.yml

    sudo /etc/opk/opk-ssh add "${EMAIL}" "${USER}"

    # Comment out existing AuthorizedKeysCommand configuration
    # TODO: How do these other AuthorizedKeysCommands interact with our own?
    sudo sed -i '/^AuthorizedKeysCommand /s/^/#/' "/etc/ssh/sshd_config"
    sudo sed -i '/^AuthorizedKeysCommandUser /s/^/#/' "/etc/ssh/sshd_config"

    # Add our AuthorizedKeysCommand line so that the opk verifier is called when ssh-ing in
    sudo tee -a /etc/ssh/sshd_config > /dev/null <<EOT
AuthorizedKeysCommand /etc/opk/opk-ssh verify %u %k %t
AuthorizedKeysCommandUser root
EOT

    sudo systemctl restart sshd || sudo systemctl restart ssh
    exit 0
else
    echo "The user is not root and does not have sudo access; adding resources to user's home directory"

    mkdir ~/.opk

    # Install the opk-ssh client to user's home directory
    mv opk-ssh ~/.opk
    chmod 700 ~/.opk/opk-ssh

    # Create a personal policy yaml file in user's home directory
    touch ~/.opk/policy.yml
    chmod 600 ~/.opk/policy.yml

    ~/.opk/opk-ssh add "${EMAIL}" "${USER}"
    exit 0
fi
