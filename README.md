# bangumi-cli

A set of CLI tools for managing seasonal animation collections and Plex Media Streaming related tasks (Work in Progress).

- Rename media files to match Plex media server format.
- Generate folders from bangumi user vollection with metadata config file (from TMDB, bangumi.tv).
- Automate seasonal animation batch download, group, and renaming.

Why?

The goal is to automate the Plex Home Media Streaming workflow that involves:

1. Fetching and organizing seasonal animation from a variety of sources (e.g., acgsecrets.hk) and sync with bangumi.tv.
2. Marked subjects as watching or want to watch.
3. Creating a folder for each subject.
4. Manually searching for the subject keyword on torrent distribution sites (e.g., share.dmhy.org, acgrip.art ...).
5. Filtering false positives, fansub group, languages...
6. Downloading each episode (bulk downloads normally require waiting 1-2 weeks after the current season ends, depends on each fansub group schedule).
7. Renaming file names to match Plex Media formats.

## Build

To build the CLI tool for your platform, use the following commands:

```sh
# windows
go build -o bangumi-cli.exe

# linux
go build -o bangumi-cli
```

## Install

Move the binary to a directory included in your system's environment variable path to make it executable. For example, you can move it to the Go bin path.

```sh
# windows
mv .\bangumi-cli.exe C:\Users\<usrname>\go\bin

# linux
mv bangumi-cli /usr/local/go/bin
```

## General Usage

### Generate Folders from bangumi.tv Collections

The command will fetch user collections from bangumi.tv and generate a folder for each subject. It will also generate a metadata file for each subject that is useful for the subsequent renaming process (e.g., year, episode details, season, etc.).

```sh
bangumi-cli list -h

# generate folders 在看 动画 from user sstp105's collections
bangumi-cli list -c -s 3 -t 2 -u "sstp105"
```

### Rename Files to Match Plex TVShows Format

To stream media content on Plex Media Server, the files must follow Plex TV Show naming conventions.

The metadata are parsed by:

1. The directory name if metadata file does not exist (e.g. manually created the directory). Only title and season can be parsed for this option.
2. From the metadata file if folders are generated using the `bangumi list`.

For the option 1, some subjects cannot be parsed from the title alone. For example, **夏目友人帐 柒**, where the season is 7 for the series **夏目友人帐**, but this cannot be parsed from the name alone.

It is recommended to rename the folders in one of the following formats before rename:

- **Title 第 N 季**, where N can be a Chinese number or Arabic numeral. For example, **第 2 季**, **第 12 季**, **第三季**, **第十一季**.
- **Title 第 N 期**, where N can be a Chinese number or Arabic numeral. For example, **第 2 期**, **第 12 期**, **第三期**, **第十一期**.
- **Title S{NN}**, where {NN} represents a two-digit Arabic numeral. For example, **S01**, **S12**.

```sh
bangumi-cli -h

banggumi-cli rename
```
