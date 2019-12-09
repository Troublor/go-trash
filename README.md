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

## Installation

**The installation of Go-trash requires root authentication.**

```
git clone https://github.com/Troublor/go-trash.git
cd go-trash
chmod +x install.sh
./install.sh
```

The executable file of Go-trash is installed in to /usr/local/bin 
The trash bin is located at /etc/gotrash

To uninstall Go-trash, run the uninstall.sh script: 

```
./uninstall.sh
```
