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
    # Create the default sos directory
    sudo mkdir -p /etc/sos

    # Install the sos verifier
    sudo mv verifier /etc/sos
    sudo chown root /etc/sos/verifier
    sudo chmod 700 /etc/sos/verifier

    # Create root policy yaml file which is a yaml file
    sudo touch /etc/sos/policy.yml
    sudo chown root /etc/sos/policy.yml
    sudo chmod 600 /etc/sos/policy.yml

    sudo /etc/sos/verifier add "${EMAIL}" "${USER}"

    # Comment out existing AuthorizedKeysCommand configuration
    # TODO: How do these other AuthorizedKeysCommands interact with our own?
    sudo sed -i '/^AuthorizedKeysCommand /s/^/#/' "/etc/ssh/sshd_config"
    sudo sed -i '/^AuthorizedKeysCommandUser /s/^/#/' "/etc/ssh/sshd_config"

    # Add our AuthorizedKeysCommand line so that the sos verifier is called when ssh-ing in
    sudo tee -a /etc/ssh/sshd_config > /dev/null <<EOT
AuthorizedKeysCommand /etc/sos/verifier verify %u %k %t
AuthorizedKeysCommandUser root
EOT

    sudo systemctl restart sshd || sudo systemctl restart ssh
    exit 0
else
    echo "The user is not root and does not have sudo access; adding resources to user's home directory"

    mkdir ~/.sos

    # Install the verifier client to user's home directory
    mv verifier ~/.sos
    chmod 700 ~/.sos/verifier

    # Create a personal policy yaml file in user's home directory
    touch ~/.sos/policy.yml
    chmod 600 ~/.sos/policy.yml

    ~/.sos/verifier add "${EMAIL}" "${USER}"
    exit 0
fi
