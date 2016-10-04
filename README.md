Crash
=====

[![Version](https://badge.fury.io/gh/troykinsella%2Fcrash.svg)](https://badge.fury.io/gh/troykinsella%2Fcrash)
[![License](https://img.shields.io/github/license/troykinsella/crash.svg)](https://github.com/troykinsella/crash/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/troykinsella/crash.svg?branch=master)](https://travis-ci.org/troykinsella/crash)

A command-line tool for executing test plans and reporting results.

## Test Plan

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
```

Output:

```
[#] {0.000s} Sooper Site
[-] {0.000s} serial...
[!] {0.000s} Home Page
[I] {0.000s} GET http://example.com
[I] {0.138s} GET http://example.com -> 200
[!] {0.139s} (138.548468ms) Home Page
[✓] http status 200 is 2xx
[✓] has html5 doctype declaration
[✓] headers.Content-Type eq 'text/html'
[-] {0.139s} (138.664306ms) 
[#] {0.139s} (138.676668ms) Sooper Site
```

## TODO

Way more fucking documentation. WAY more.
