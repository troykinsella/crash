Crash
=====

[![Version](https://badge.fury.io/gh/troykinsella%2Fcrash.svg)](https://badge.fury.io/gh/troykinsella%2Fcrash)
[![License](https://img.shields.io/github/license/troykinsella/crash.svg)](https://github.com/troykinsella/crash/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/troykinsella/crash.svg?branch=master)](https://travis-ci.org/troykinsella/crash)

A command-line tool for executing test plans and reporting results. Test plans, `Crashfile`s, defined in YAML,
direct the execution of actions, such as HTTP requests, checks, such as asserting an HTTP
response code of 200, repetitions, and concurrency.

## Documentation

User documentation can be found at https://troykinsella.github.io/crash/

## Usage

### Test Operation

```bash
NAME:
   crash test - Execute a test plan

USAGE:
   crash test [command options] [arguments...]

OPTIONS:
   -j                 Format logging output as JSON
   -q                 Quiet mode; suppress logging
   -v                 Verbose logging; Use -vv or -vvv to increase verbosity
   -s FILE|KEY=VALUE  Variable FILE|KEY=VALUE
   -f FILE            Input test yaml FILE; Defaults to searching for Crashfile.y[a]ml in the current directory
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
