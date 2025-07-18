# bangumi-cli

A set of CLI tools for bangumi.tv and plex media streaming related tasks.

## Environment Setup

|                         | Description                                                                                                  |
| ----------------------- | ------------------------------------------------------------------------------------------------------------ |
| `LOCAL_SERVER_PORT`     | The local server port number for the application and is used for listening for bangumi API callback.         |
| `BANGUMI_CLIENT_ID`     | Bangumi app's Client ID used. https://bangumi.tv/dev/app                                                     |
| `BANGUMI_CLIENT_SECRET` | Bangumi app's Client Secret. https://bangumi.tv/dev/app                                                      |
| `QBITTORRENT_SERVER`    | qBittorrent Web UI server, can be configured in qBittorrent client settings (e.g., `http://localhost:8767`). |
| `QBITTORRENT_USERNAME`  | Username for logging into the qBittorrent Web UI (default is `admin`).                                       |
| `QBITTORRENT_PASSWORD`  | Password for the qBittorrent Web UI, can be configured in qBittorrent client settings.                       |
| `MIKAN_IDENTITY_COOKIE` | Mikan identity cookie `.AspNetCore.Identity.Application`                                                     |

To setup environment variables:

```sh
# Windows: Open powershell as administrator then run the following commands under `bangumi-cli` directory
cd scripts
.\setup.ps1

# MacOS: Run the following commands under bangumi-cli directory
chmod +x setup.sh
./setup.sh
source ~/.zshrc
```

## Build

```sh
# Windows
go build -o bangumi.exe

# MacOS
go build -o bangumi
```

## Test

```
go test -v ./...
```

## Usage

```sh
bangumi login -h

bangumi subscribe -h

bangumi update -h

bangumi download -h

bangumi format -h

bangumi collect -h

bangumi unsubscribe -h
```
