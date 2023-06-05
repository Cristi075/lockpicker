# Lockpicker

A simple program written in go that checks if users on your linux server have simple passwords.

![Screenshot](/docs/demo.png)

In order to run the checks, the program uses 'expect' (more precisely, Google's 'goexpect' implementation available at [https://github.com/google/goexpect](https://github.com/google/goexpect)).
The "su" command is used in order to attempt to login.

At the moment of writing this, the only passwords that are checked are the following ones:
- password = username
- password = username + 1 digit
- password = username + 2 digits


## Using the tool

### Download a binary
- Go to [https://github.com/Cristi075/lockpicker/releases](https://github.com/Cristi075/lockpicker/releases)
- Download the appropriate binary for your target system (x64 or x86)
- Run the binary


### Compile your own binaries
**Prerequisites**: you must have go 1.20 installed on your system


- Clone this repository
- Run "go build"  

You can also build both the x86 and x64 versions by running:
```
GOARCH=amd64 go build -o lockpicker_x64
GOARCH=386 go build -o lockpicker_x86
```

## Roadmap

This is still a work in progress for me.  
Here is a list of things that I plan to (eventually) add to this script:
- [ ] Improve concurrency. Right now, this is a version that just worked for PoC purposes
- [ ] Add flags to allow the user to change the program's behaviour
- [ ] Allow the user to select a custom password list
- [ ] Allow the user to change how many requests are executed in parallel 
- [ ] Add more options for generating simple passwords (uppercase/lowercase mutations, more digits, common words, etc) 
- [ ] ???

