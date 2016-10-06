# `http`

Make a request to an HTTP or HTTPS server.

---

## Parameters

Name         | Required | Default      | Description
------------ | -------- | ------------ | -------------
method       | no       | "GET"        | The HTTP request method to use. 
url          | yes      |              | The URL against which a request will be made. Must have an "http://" or "https://" scheme.

## Outputs

Name         | Description
------------ | ------------
body         | 
headers      | A map of HTTP headers returned in the response.
status-code  | The HTTP response status code.

## Examples

Test plan:

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
    - headers.content-type eq 'text/html'
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
[✓] headers.content-type eq 'text/html'
[✗] example failure!
[-] {0.166s} (166.388073ms) 
[#] {0.166s} (166.407687ms) Sooper Site
```

