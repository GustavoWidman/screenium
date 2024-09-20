# screenium

## Description

screenium is a modern remake of GNU "screen" in golang.

## Features

- Create, attach, detach, kill, list, and help commands
- Supports multiple shells
- Supports multiple sessions

## Installation

Currently there is a [install.sh](install.sh) script that will download the latest pre-compiled binary and place it in `/usr/local/bin/screenium`. There also is a [build.sh](build.sh) script that will download the latest copy of the repo, build the binary from source and place it in `/usr/local/bin/screenium`.

Please note that to be able to build screnium from source, you will need to have the latest version of the Go compiler installed on your system. If you're not sure how to do that, please refer to the [official Go documentation](https://go.dev/doc/install).

You can also run these one-liners to install and build screenium (they just download the script and run it for you):

### Install

```bash
curl -s https://links.r3dlust.com/screenium-install | bash
```

### Build

```bash
curl -s https://links.r3dlust.com/screenium-build | bash
```

## Author

[GustavoWidman](https://github.com/GustavoWidman)

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
