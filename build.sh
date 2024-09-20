# Check if user running is root
if [[ $EUID -ne 0 ]]; then
    echo "This script must be run as root"
    exit 1
fi

# check if "Go" is installed
if ! command -v go &> /dev/null
then
    echo "Go is not installed. Please install it first."
    exit 1
fi

# Check if x86_64 architecture is detected
if [[ $(uname -m) != "x86_64" ]]; then
    echo "This code has only been tested on x86_64 architecture, you are on your own."
fi


# Check if Arch Linux is detected
# If it is, use AUR method (screenium-git)
# TODO

# if [[ $(uname -r) == *arch* ]]; then
#     echo "Arch Linux detected"
# fi


echo "Looking for the latest tarball..."

tmpdir=$(mktemp -d)
# Download the latest tarball from "https://api.github.com/GustavoWidman/screenium/releases/latest/"
curl -sL $(curl -s https://api.github.com/repos/GustavoWidman/screenium/releases/latest | grep "tarball_url" | awk '{print $2}' | sed 's|[\"\,]*||g') | tar xzf - --directory $tmpdir
folder=$(find $tmpdir -maxdepth 1 -type d -name "GustavoWidman-screenium-*")

# save the current dir and change to the new tmpdir
olddir=$(pwd)
cd $folder

echo "Done. Installing dependencies and building..."
go mod tidy
CGO_ENABLED=0 go build -o screenium -ldflags "-s -w" src/main.go
screenium_path="$folder/screenium"

echo "Binary built at $screenium_path"

cd $olddir

echo "Installing screenium..."

sudo cp $screenium_path /usr/local/bin/screenium
sudo chmod +x /usr/local/bin/screenium

echo "Cleaning up..."
rm -rf $tmpdir

echo "Done!"
echo "Run \"screenium\" to start using it."