# Crash

Crash is a command-line tool for executing test plans and reporting results, written in Go.
It works on Darwin, Linux, and Windows systems, with 64-bit architecture.

---

## Installation

### Binary Distribution

Crash is distributed as a single OS-/Architecture-specific executable binary.

Head over to [Releases](https://github.com/troykinsella/crash/releases/) and download the appropriate 
binary for your system. Then, move the binary to a convenient location for execution.

Or, run:
```bash
VERSION={{ shell "curl -s https://api.github.com/repos/troykinsella/crash/releases | jq -r '.[0].tag_name' | sed 's/v//'" }}
OS=linux # or darwin, or windows
curl -SsL -o /usr/local/bin/crash https://github.com/troykinsella/crash/releases/download/v${VERSION}/crash_${OS}_amd64
chmod +x /usr/local/bin/crash
```

## Getting Started

### Create a `Crashfile`

```yaml
{{ shell "cat docs/actions/http_example.yml" }}
```
This `Crashfile` has a single plan which runs an action that performs an HTTP request.
After a response is received, three assertions are performed against the response data, 
specifically, `status-code`, `body`, and `headers`, which, among others, are created by `crash`'s
`http` action.

### Run `crash test`

In the same directory as your `Crashfile`, run `crash test`. For clarity's sake, let's
add the `-vv` (add verbosity) option to get a little more logging output.

```bash
crash test -vv
```
... and examine the output:

```
{{ shell "crash test -vv --nc -f docs/actions/http_example.yml" }}
```
