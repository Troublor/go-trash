# Go-trash

A trash files and directories management tool for Linux command line. 

## Features

- Remove files and directories by moving them in to a trash bin
- Retrieve files and directories by moving them out of trash bin to original path
- Redirect the files or directories to un-remove to the retrieving target location
- List all items in trash bin
- Search for items in trash bin
- User isolation: different users have their own trash bin which doesn't interfere with each other. Root users can see all trash items. 

## Features to Come

- Auto clearance of trash bin
- Command completion for trash items

## Tests

To execute go tests:

```bash
make test
```

## Installation

**The installation of Go-trash requires root authentication.**

```
git clone https://github.com/Troublor/go-trash.git
cd go-trash
make install 
```
This commands will install `gotrash` at `~/bin`
If you want to install `gotrash` in another place, do `export CMD_PATH=<another folder>` before `make install` 

## GOTRASH_PATH

GOTRASH_PATH is the base directory to hold data of `gotrash`. The default `GOTRASH_PATH` is `/etc/gotrash` for `root` 
users, and `~/.gotrash` for other users. 

You can do `export GOTRASH_PATH=<an empty folder>` before `make install` to specifiy `GOTRASH_PATH`. This will make all 
users' `GOTRASH_PATH` the same as what you specified.

## Uninstall

To uninstall gotrash, just simply delete the `GOTRASH_PATH` folder and the `gotrash` command file. 
** Note that this will permanently delete any trash in the trash bin.**
