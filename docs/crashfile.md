# Introduction

Crashfiles, defined in YAML, tell `crash` what to do.  

---

## The Crashfile

A note on examples: When you see `# ...` it denotes that the example is just a snippet
of a larger document, which, by itself, is not valid. The context in which the snippet
must live is made clear by the accompanied documentation.

### Plans

Defined at the root of the test plan document, `plans` is a list of plan objects, 
and every test plan document must define it. Plans are executed in the order defined, 
serially. Must have at least one entry.

### Plan

A plan object defines a name and the list of steps that `crash` will execute in a test run. 

Properties | Required | Description
---------- | -------- | -----------
plan       | yes      | The name of the plan which appears in test output.
steps      | yes      | A list of steps to execute. This step list is really an implicit [serial](#serial) step, which executes the steps sequentially, serially.

Example:
```yaml
---
plans:
- plan: Plan A
  steps:
  # ...
```

### Step

A step can be one of the types listed in the [Plan Steps](#plan-steps) section.
Any kind of step object may have the following properties:

Properties | Required | Description
---------- | -------- | -----------
check      | no       | A list of assertions to perform after the execution of the step is complete.
timeout    | no       | A time duration after which the step is aborted and is considered to have failed. A duration string is a sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "1.5s" or "2s300ms". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h". Zero or negative timeouts are not permitted.
with       | no       | A structure that controls the step execution repetition.

### Vars

Defined at the root of the test plan document, `vars` defines constants that are available
at execution time.

```yaml
---
vars:
  key: value
  foo: bar
# ...
```

### With

The `with` directive controls repetition of the enclosing step.

The value of `with` is an object having one of the following properties:

Properties | Description
---------- | -----------
list       | A list literal (in yaml). The step will be executed once for every element in the list.
item       | An expression that evaluates to an iterable value.

The `with` object also accepts an optional `as` property. The value overrides the default
variable name of `i`, and stores the value of the current element being iterated.

Examples:
```
#  with: { list: [10, 20, 30] }
#  with: { item: 1000, as: i }
#  with: { item: user_ids, as: user_id }
```

## Plan Steps

### Parallel

Execute a list of steps in parallel. The parallel step, itself, completes when all of the nested steps
have completed.

Example:
```yaml
# ...
- parallel:
  - # step 1 ...
  - # step 2 ...
```

### Run

Run an action. Available actions can be browsed in the "Action Reference" from the main menu.

Actions are represented by an object having a `run` property that has an object value. The object
value has the following properties:

Properties | Required | Description
---------- | -------- | -----------
name       | yes      | The name of the action which appears in test output.
type       | yes      | Dictates which action is selected for execution.
params     | yes      | An object defining action-specific key-value pairs that are passed into the action.

Example:
```yaml
# ...
- run:
    name: Home page sample
    type: http
    params:
      url: http://fooland.org
```

### Serial

Execute a list of steps one after the other. The serial step, itself, completes when the last nested step completes.

Example:
```yaml
# ...
- serial:
  - # step 1 ...
  - # step 2 ...
```

## Example

This example intends to show all basic `crash` features.

```yaml
---
vars:
  base_url: http://example.com

plans:
- plan: Sooper Site
  steps:
  - parallel:
    - run:
        name: home page
        type: http
        params:
          url: $base_url
      check:
      - status-code in 200, 299 // http status $status-code is 4xx
      - body contains '<!doctype html>' // has doctype declaration
      - headers.content-type eq 'text/html'
    - run:
        name: about page
        type: http
        params:
          url: $base_url/about.html
      check:
      - status-code eq 200
      - body contains 'Copyright'
```
