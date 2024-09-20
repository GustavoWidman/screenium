# Check if user running is root
if [[ $EUID -ne 0 ]]; then
    echo "This script must be run as root"
    exit 1
fi

# Check if x86_64 architecture is detected
if [[ $(uname -m) != "x86_64" ]]; then
    echo "This script is only compatible with x86_64 architecture"
    exit 1
fi

# Check if Arch Linux is detected
# If it is, use AUR method (screenium-bin)
# TODO

# if [[ $(uname -r) == *arch* ]]; then
#     echo "Arch Linux detected"
# fi

echo "Looking for the latest release..."

# Find out the latest release tag
url=$(curl -s https://api.github.com/repos/GustavoWidman/screenium/releases/latest | grep "releases/tag" | awk '{print $2}' | sed 's|[\"\,]*||g')
tag=$(echo $url | sed 's|https://github.com/GustavoWidman/screenium/releases/tag/||g')
download_link="https://github.com/GustavoWidman/screenium/releases/download/$tag/screenium"

echo "Downloading $download_link..."
tmpfile=$(mktemp)
curl -sL $download_link -o $tmpfile

echo "Installing screenium"
sudo cp $tmpfile /usr/local/bin/screenium
sudo chmod +x /usr/local/bin/screenium

echo "Cleaning up"
rm $tmpfile

echo "Done!"
echo "Run \"screenium\" to start using it."