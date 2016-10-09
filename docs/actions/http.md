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
body         | The response body string.
headers      | A map of HTTP headers returned in the response.
status-code  | The HTTP response status code.
raw-body     | The response body bytes.

## Examples

```yaml
# Crashfile
{{ shell "cat docs/actions/http_example.yml" }}
```

Standard output (-vvv):
```
{{ shell "crash test -vvv --nc -f docs/actions/http_example.yml" }}
```

JSON output (-vvv):
```json
{{ shell "crash test -vvv -j -f docs/actions/http_example.yml" }}
```
