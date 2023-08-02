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

**Watcher is used for watching for the files changes** and **logger for viewing logs in real time by given filters**. More comprehence usage examples below.

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

## Notes

- For watching file changes, a custom solution was implemented. However, it has not undergone thorough testing, and in a real-life scenario, it would be better to use a well-tested and widely used library like fsnotify. Another option is to use inotify, but this approach might not be cross-platform compatible.
- The logic related to CLI argument parsing was partially taken from my other home project. It served well there, but in this case, it is overkill and probably not well-suited.
- For small and focused projects, simplicity is often key to achieving the desired outcome without unnecessary overhead or complexity, but I wanted to play with golang a little bit more and that is why some decision were made a little bit more complex than it supposed to be.
- I wanted to implement my own logging with different log levels, but it became complex and not as useful as I initially thought
- Options/arguments passed to **watcher** or **logger** need to be validated. 

## How the watcher works

The watcher operates by recursively scanning all files and directories. Upon initial load, it adds them to a map collection where the file/directory name serves as the key, and additional file information (**FileInfo**) is stored as the corresponding value. Approximately every 100 milliseconds, I retrieve information about files in the directory and poll for file change events in a separate goroutine. By comparing the newly retrieved list of files with the file information stored in the map (which was loaded at application start), I can determine the nature of the changes that have occurred. Subsequently, the worker sends notifications using Go channels, allowing the main application to handle file changes in a separate function.

## How the logger works

The logger relies on the watcher module. It operates independently while also utilizing the **file.Watcher** to receive notifications about modifications made to the log file (**file-filter.log**). During startup, all files are scanned, and only filtered lines are printed. Later, a notification is sent to the logger when a file change occurs, prompting the logger to read the last line of the file and display it if the filter options are satisfied. Filtering is achieved by utilizing the **matchString** function.
