# Assign arguments to variables
EMAIL="$1"
USER="$2"
VERSION="$3"

ARCH=$(uname -m)
OS=$(uname -os)

OS_VAL=$(echo "$OS" | cut -d ' ' -f 1)

# Validate input arguments
if [ -z "${EMAIL}" ] || [ -z "${USER}" ]; then
    echo "Usage: $0 <EMAIL> <USER> <VERSION>"
    exit 1
fi

if sudo -nl &>/dev/null || [ "$(id -u)" -ne 0 ]; then
    echo "User has sudo privileges without a password."

    # Check for wget or curl availability
    DOWNLOADER=""
    if command -v wget &>/dev/null; then
        DOWNLOADER="wget"
    elif command -v curl &>/dev/null; then
        DOWNLOADER="curl -LO"
    else
        echo "Neither wget nor curl is available. Please install either of them."
        exit 1
    fi

    if [[ "$OS_VAL" == *"Linux"* && "$ARCH" == *"x86_64"* ]]; then
        $DOWNLOADER https://github.com/devlup-labs/spok/releases/download/${VERSION}/verifier_${VERSION}_linux_amd64
        mv verifier_${VERSION}_linux_amd64 verifier
    elif [[ "$OS_VAL" == *"Linux"* && "$ARCH" == *"arm64"* ]]; then
        $DOWNLOADER https://github.com/devlup-labs/spok/releases/download/${VERSION}/verifier_${VERSION}_linux_arm64
        mv verifier_${VERSION}_linux_arm64 verifier
    elif [[ "$OS_VAL" == *"Darwin"* && "$ARCH" == *"x86_64"* ]]; then
        $DOWNLOADER https://github.com/devlup-labs/spok/releases/download/${VERSION}/verifier_${VERSION}_darwin_amd64
        mv verifier_${VERSION}_darwin_amd64 verifier
    elif [[ "$OS_VAL" == *"Darwin"* && "$ARCH" == *"arm64"* ]]; then
        $DOWNLOADER https://github.com/devlup-labs/spok/releases/download/${VERSION}/verifier_${VERSION}_darwin_arm64
        mv verifier_${VERSION}_darwin_arm64 verifier
    else
        echo "This OS: $OS_VAL and ARCH: $ARCH is not supported please contact the developers for help :)"
    fi

    # Create the default spok directory
    sudo mkdir -p /etc/spok

    # Install the spok verifier
    sudo mv verifier /etc/spok/
    sudo chown root /etc/spok/verifier
    sudo chmod 700 /etc/spok/verifier

    # Create root policy yaml file which is a yaml file
    sudo touch /etc/spok/policy.yml
    sudo chown root /etc/spok/policy.yml
    sudo chmod 600 /etc/spok/policy.yml

    sudo /etc/spok/verifier add "${EMAIL}" "${USER}"

    # Comment out existing AuthorizedKeysCommand configuration
    # TODO: How do these other AuthorizedKeysCommands interact with our own?
    if [[ "$OS_VAL" == *"Linux"* ]]; then
        sudo sed -i '/^AuthorizedKeysCommand /s/^/#/' "/etc/ssh/sshd_config"
        sudo sed -i '/^AuthorizedKeysCommandUser /s/^/#/' "/etc/ssh/sshd_config"
    elif [[ "$OS_VAL" == *"Darwin"* ]]; then
        sudo sed -i '' '/^AuthorizedKeysCommand /s/^/#/' "/etc/ssh/sshd_config"
        sudo sed -i '' '/^AuthorizedKeysCommandUser /s/^/#/' "/etc/ssh/sshd_config"
    else
        echo "This OS: $OS_VAL is not supported please contact the developers for help :)"
    fi

    # Add our AuthorizedKeysCommand line so that the spok verifier is called when ssh-ing in
    sudo tee -a /etc/ssh/sshd_config >/dev/null <<EOT
AuthorizedKeysCommand /etc/spok/verifier verify %u %k %t
AuthorizedKeysCommandUser root
EOT

    if sudo systemctl restart sshd &>/dev/null; then
        echo "SSH service restarted successfully."
    elif sudo systemctl restart ssh &>/dev/null; then
        echo "SSH service restarted successfully."
    else
        exit 0
    fi

    exit 0
else
    echo "The user is not root and does not have sudo access; adding resources to user's home directory"

    mkdir ~/.spok

    # Check for wget or curl availability
    DOWNLOADER=""
    if command -v wget &>/dev/null; then
        DOWNLOADER="wget"
    elif command -v curl &>/dev/null; then
        DOWNLOADER="curl -LO"
    else
        echo "Neither wget nor curl is available. Please install either of them."
        exit 1
    fi

    if [[ "$OS_VAL" == *"Linux"* && "$ARCH" == *"x86_64"* ]]; then
        $DOWNLOADER https://github.com/devlup-labs/spok/releases/download/${VERSION}/verifier_${VERSION}_linux_amd64
        mv verifier_${VERSION}_linux_amd64 verifier
    elif [[ "$OS_VAL" == *"Linux"* && "$ARCH" == *"arm64"* ]]; then
        $DOWNLOADER https://github.com/devlup-labs/spok/releases/download/${VERSION}/verifier_${VERSION}_linux_arm64
        mv verifier_${VERSION}_linux_arm64 verifier
    elif [[ "$OS_VAL" == *"Darwin"* && "$ARCH" == *"x86_64"* ]]; then
        $DOWNLOADER https://github.com/devlup-labs/spok/releases/download/${VERSION}/verifier_${VERSION}_darwin_amd64
        mv verifier_${VERSION}_darwin_amd64 verifier
    elif [[ "$OS_VAL" == *"Darwin"* && "$ARCH" == *"arm64"* ]]; then
        $DOWNLOADER https://github.com/devlup-labs/spok/releases/download/${VERSION}/verifier_${VERSION}_darwin_arm64
        mv verifier_${VERSION}_darwin_arm64 verifier
    else
        echo "This OS: $OS_VAL and ARCH: $ARCH is not supported please contact the developers for help :)"
    fi

    mv verifier ~/.spok/verifier
    # Install the verifier client to user's home directory
    chmod 700 ~/.spok/verifier

    # Create a personal policy yaml file in user's home directory
    touch ~/.spok/policy.yml
    chmod 600 ~/.spok/policy.yml

    ~/.spok/verifier add "${EMAIL}" "${USER}"

    exit 0
fi
