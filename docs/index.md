# Crash

Crash is a command-line tool for executing test plans and reporting results, written in Go.
It works on Darwin, Linux, and Windows systems, with 64-bit architecture.

---

## Installation

### Binary Distribution

Crash is distributed as a single OS-/Architecture-specific executable binary.

Head over to [Releases](https://github.com/troykinsella/crash/releases/) and download the appropriate 
binary for your system. Then, move the binary to a convenient location:

```bash
$ sudo mv ~/Downloads/crash_[OS]_[Arch] /usr/local/bin/crash
$ sudo chmod +x /usr/local/bin/crash
```

### Building from Source

Set up a [Go workspace](https://golang.org/doc/code.html).

Get the sources:
```sh
$ go get github.com/troykinsella/crash
```

Install dependencies:
```sh
$ go get -d -v ./...
```

Build the binary:
```sh
$ go build -o crash -v github.com/troykinsella/crash/cmd
```

## Getting Started



