Crash
=====

[![Version](https://badge.fury.io/gh/troykinsella%2Fcrash.svg)](https://badge.fury.io/gh/troykinsella%2Fcrash)
[![License](https://img.shields.io/github/license/troykinsella/crash.svg)](https://github.com/troykinsella/crash/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/troykinsella/crash.svg?branch=master)](https://travis-ci.org/troykinsella/crash)

A command-line tool for executing test plans and reporting results. Test plans, defined in YAML,
direct the execution of actions, such as HTTP requests, checks, such as asserting an HTTP
response code of 200, repetitions, and concurrency.

## Features

* Human readable test output
* JSON test event output
* Query DSL for extracting data and performing assertions

## Documentation

User documentation can be found at https://troykinsella.github.io/crash/

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
> crash test -f mytest.yml -s base_url=http://example.com/foo
```


## Roadmap

* Way more fucking documentation. WAY more.
* HTTPS support
* Configurable input sources and output targets
* Shell script execution

## License

MIT © Troy Kinsella
