
## Reporting an Issue

When [reporting an issue](https://github.com/troykinsella/crash/issues), please include as much 
detail as you can think of, including:

* `crash` version
* Operating system
* System architecture
* Step-by-step instructions on how to reproduce the issue
* Stack traces and any other relevant output

## Submitting a Pull Request

Pull requests are very welcome! All PRs should be created against the `develop` branch.

Please be sure to include:

* `go fmt` formatted code
* Passing tests
* User documentation (the `docs/` dir)

## Building the `crash` Executable

Set up a [Go workspace](https://golang.org/doc/code.html).

From your `GOPATH`, get the sources:
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

## Runing Tests

From your `GOPATH`:

```sh
go test github.com/troykinsella/crash/...
```
