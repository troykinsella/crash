# Introduction

Running `crash test` produces output that details the passing and failing elements of
the test plan execution.

## Logging Levels

Output verbosity is controlled by the `-v` option, passed to `crash test`. Messages
are printed with further verbosity by supplying `-vv`, and further still with `-vvv`.
Additionally, passing `-q` (quiet) turns off logging output altogether.

The verbosity options correlate to logging level as follows:

Option | Logging Level
------ | -------------
`-q`   | `OFF`
(Neither `-q` nor `-v`) | `DEFAULT`
`-v`   | `DETAIL`
`-vv`  | `INFO`
`-vvv` | `DEBUG`

## Output for Humans

Crash events are printed according to this format:
```
<message-type> <success> [elapsed-time] [step-duration] <message>
```

The presence of `elapsed-time` and `step-duration` columns are optional, and are 
predictable based on the `message-type`.

### First Column: Message Type

Symbol | When       | Level     | Description 
------ | ---------- | --------- | ----------- 
#      | Start      | `DETAIL`  | Plan execution started
#      | Finish     | `DEFAULT` | Plan execution finished
-      | Start      | `INFO`    | Serial execution started
-      | Finish     | `INFO`    | Serial execution finished
=      | Start      | `INFO`    | Parallel execution started
=      | Finish     | `INFO`    | Parallel execution finished
!      | Start      | `DETAIL`  | Action execution started
!      | Finish     | `DEFAULT` | Action execution finished
?      | Occurrence | `DETAIL` when pass, `DEFAULT` when fail | Check, a.k.a. assertion
I      | Occurrence | `INFO`    | Info log message
D      | Occurrence | `DEBUG`   | Debug log message

### Second Column: Success

Symbol | Description
------ | -----------
.      | Success is not relevant to the message
✓      | Pass
✗      | Failure

### Elapsed Time

Reports the time elapsed in seconds since the beginning of the `crash test` execution. Not present for checks
(when the message type symbol is `?`).

### Step Duration

Shows the time elapsed upon completion of a step, since the beginning of the step, 
as a string in the form "3m0.5s". Durations less than one second use a smaller unit 
(milli-, micro-, or nanoseconds).

### Message

The meat and potatoes of the logged event.

## Output for Machines

TODO
