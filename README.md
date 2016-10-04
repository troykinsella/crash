Crash
=====

[![Version](https://badge.fury.io/gh/troykinsella%2Fcrash.svg)](https://badge.fury.io/gh/troykinsella%2Fcrash)
[![License](https://img.shields.io/github/license/troykinsella/crash.svg)](https://github.com/troykinsella/crash/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/troykinsella/crash.svg?branch=master)](https://travis-ci.org/troykinsella/crash)

A command-line tool for executing test plans and reporting results. Test plans, defined in YAML,
direct the execution of actions, such as HTTP requests, checks, such as asserting an HTTP
response code of 200, repetitions, and concurrency.

## Features

* Is wicked
* Human readable test output
* JSON test event output
* Query DSL for extracting data and performing assertions

## Usage

### Test Operation

```bash
> crash test -h
NAME:
   crash test - Execute a test plan

USAGE:
   crash test [command options] [arguments...]

OPTIONS:
   -d            Debug mode; increase logging verbosity
   -j            Format logging output as JSON
   -q            Quiet mode; suppress logging
   -s key=value  Variable key pair key=value
   -f FILE       Input test yaml FILE
   -v FILE       Variables yaml FILE
```

Example execution:

```bash
> crash test -f mytest.yml -s base_url=http://example/foo
```

## Test File

The test YAML file includes these properties:

* `plans` - Required. A list of plans which define steps to execute.
* `vars` - Optional. Key-value pairs available to the test runtime.

### Plans

```yaml
---
plans:
- plan: Sooper Site
  steps:
  - run:
      name: Home Page
      type: http
      params:
        method: get
        url: http://example.com
    check:
    - status-code in 200, 299 // http status ${status-code} is 2xx
    - body contains '<!doctype html>' // has html5 doctype declaration
    - headers.Content-Type eq 'text/html'
    - body contains 'something crazy' // example failure!
```

Running the above plan yields the following output:

```
[#] {0.000s} Sooper Site
[-] {0.000s} serial...
[!] {0.000s} Home Page
[I] {0.000s} GET http://example.com
[I] {0.166s} GET http://example.com -> 200
[!] {0.166s} (166.220661ms) Home Page
[✓] http status 200 is 2xx
[✓] has html5 doctype declaration
[✓] headers.Content-Type eq 'text/html'
[✗] example failure!
[-] {0.166s} (166.388073ms) 
[#] {0.166s} (166.407687ms) Sooper Site
```

## Roadmap

* Way more fucking documentation. WAY more.
* HTTPS support
* Configurable input sources and output targets
* Shell script execution

## License

MIT © Troy Kinsella
