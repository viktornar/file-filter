# file-filter

Here's a simple application to demonstrate Golang skills. The application has two main responsibilities:

- Watch for file changes in the hot path and move them to the backup path whenever any change occurs.
- Filter and display application logs.

Both parts work independently of each other. I wanted to have two separate processes: one for backup and another for log filtering and viewing. However, most of the solutions I've come across are over-engineered :(

Primary development was done on Linux, but solution should be cross platform and it should work from windows as well.

## How to build

Project do not depend on any external dependency. The only requirements is to have golang 1.20 and upper. For build just run

```bash
go build
```

or you can just run:

```bash
go run main.go
```

## How to use

Application is quite easy to use. You can just type:

```bash
go run main.go
```

You should see usage help:

```bash
file-filter version 0.0.1

Usage:
  file-filter <logger|watcher> <arguments>
```

Watcher is used for watching for the files changes and logger for viewing logs by given filters. More comprehence usage examples below.

For watcher:

```bash
go run main.go watcher ./demo/hot ./demo/backup info
```

where:

```bash
Usage:
  watcher <hotPath> <backupPath> <logLevel>
```

For logger:

```bash
go run main.go logger 2023/08/01 test
```

where:

```bash
Usage:
  logger <dateFilter> <nameFilter>
```

## What is wrong

- For watching file changes, a custom solution was implemented. However, it has not undergone thorough testing, and in a real-life scenario, it would be better to use a well-tested and widely used library like fsnotify. Another option is to use inotify, but this approach might not be cross-platform compatible.
- The logic related to CLI argument parsing was partially taken from my other home project. It served well there, but in this case, it is overkill and probably not well-suited.
- For small and focused projects, simplicity is often key to achieving the desired outcome without unnecessary overhead or complexity, but I wanted to play with golang a little bit more and that is why some decision were made a little bit more complex than it supposed to be.
- I wanted to implement my own logging with different log levels, but it became complex and not as useful as I initially thought

## How watcher is working

The watcher recursively scans all the files and directories, and on initial load, it adds them to a map collection. The file/directory name serves as the key, and additional file information (FileInfo) is stored as the value. Every 100 milliseconds, I retrieve information about files in the directory and poll for file change events in a separate goroutine. By comparing the new retrieved list of files with the file information stored in the map (which was loaded at application start), I can determine what kind of changes occurred. The worker then sends notifications using Go channels, and the main application is able to handle file changes in a separate function
