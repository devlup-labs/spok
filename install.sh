# Check if script is run as root
if [ "$(id -u)" -ne 0 ]; then
    echo "This script must be run as root. Restarting as root..."
    sudo "$0" "$@"
    exit $?
fi

# Define directories
binDir="bin"
configDir="scripts"
installDir="/usr/bin"
configInstallDir="/etc/spok"
readmeDir="/usr/share/doc/spok"
licenseDir="/usr/share/licenses/spok"

# Create installation directories if they don't exist
mkdir -p "$configInstallDir"
mkdir -p "$readmeDir"
mkdir -p "$licenseDir"

# Copy binaries
cp "$binDir/spok" "$installDir"
cp "$configDir/configure-spok-server.sh" "$configInstallDir"
cp "README.md" "$readmeDir"
cp "LICENSE" "$licenseDir"

echo "Installation of SPoK completed successfully!"
