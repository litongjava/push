# push

A simple Go tool to execute commands from `push.txt` in the current directory.

## Features

* Reads commands line by line from `push.txt`
* Executes each command sequentially
* Prints execution logs in English
* Skips empty lines
* Continues even if a command fails

## Requirements

* Go 1.18 or later

## Installation

Build the binary:

```bash
go build -o push
```

## Usage

1. Create a `push.txt` file in the current directory:

```txt
git push origin main
git push gitee main
```

2. Run the tool:

```bash
./push
```

## Example Output

```txt
Starting push execution...
Executing line 1: git push origin main
Command on line 1 executed successfully
Executing line 2: git push gitee main
Command on line 2 executed successfully
Push execution finished.
```

## Notes

* Make sure all commands in `push.txt` are valid shell commands
* Each line is executed independently
* Errors are logged but do not stop execution