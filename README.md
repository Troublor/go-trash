# Go-trash

A trash files and directories management tool for Linux command line. 

## Features

- Remove files and directories by moving them in to a trash bin
- Retrieve files and directories by moving them out of trash bin to original path
- List all items in trash bin
- Search for items in trash bin
- User isolation: different users have their own trash bin which doesn't interfere with each other. Root users can see all trash items. 

## Features to Come

- Redirect the retrieving target loaction
- Auto clearance of trash bin

## Installation

**The installation of Go-trash requires root authority.**

```
git clone https://github.com/Troublor/trash-go.git
cd trash-go
chmod +x install.sh
./install.sh
```

The executable file of Go-trash is installed in to /usr/local/bin 
The trash bin is located at /etc/gotrash

To uninstall Go-trash, run the uninstall.sh script: 

```
./uninstall.sh
```
