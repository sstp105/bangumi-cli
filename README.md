# bangumi-cli

A set of CLI tools for bangumi.tv and plex media streaming related tasks.

## Build & Install

### Windows

Setup environment variables:

```ps1
[Environment]::SetEnvironmentVariable("LOCAL_SERVER_PORT", "8765", "Machine")
[Environment]::SetEnvironmentVariable("BANGUMI_CLIENT_ID", "bangumi APP ID", "Machine")
[Environment]::SetEnvironmentVariable("BANGUMI_CLIENT_SECRET", "bangumi APP Secret", "Machine")
[Environment]::SetEnvironmentVariable("QBITTORRENT_SERVER", "http://localhost:8767", "Machine")
[Environment]::SetEnvironmentVariable("QBITTORRENT_USERNAME", "admin", "Machine")
[Environment]::SetEnvironmentVariable("QBITTORRENT_PASSWORD", "", "Machine")
[Environment]::SetEnvironmentVariable("MIKAN_IDENTITY_COOKIE", "", "Machine")
```

Run the following command to build the binary file:

```sh
go build -o bangumi.exe
```

## Test

```
go test -v ./...
```

## Usage

```sh
# Log in to bangumi.tv and obtain an API access token
bangumi login

# Subscribe to the specified season from Mikan (parse, generate metadata, prepare for pre-download)
bangumi subscribe

# Check if there are new torrents in the subscribed Mikan RSS feeds
bangumi update

# Batch process RSS and send download tasks to the client queue
bangumi download

# Recursively format files to comply with Plex / Jellyfin / Emby media library standards
bangumi format

# Collect organized anime data (batch add entries to bangumi.tv collection, sync Mikan subscription list)
bangumi collect

# Unsubscribe from Mikan
bangumi unsubscribe
```
