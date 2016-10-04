package logging

import (
	"fmt"
	"time"
	"github.com/troykinsella/crash/util"
	"encoding/json"
)

var jsISODateFormat = "2006-01-02T15:04:05.999999999Z07:00"
// less precision: "2006-01-02T15:04:05Z07:00"

type MessageType uint8

const (
	PLAN MessageType = iota
	SERIAL
	PARALLEL
	ACTION
	CHECK
	INFO
	DEBUG
)

type Logger struct {
	enabled bool
	debug   bool
	json    bool
	sw      *util.StopWatch
}

func NewLogger(enabled bool, debug bool, json bool) (*Logger) {
	return &Logger{
		enabled: enabled,
		debug: debug,
		json: json,
		sw: util.NewStopWatch().Start(),
	}
}

func (l *Logger) Start(t MessageType, msg string, args ...interface{}) {
	l.log(t, -1, true, msg, args)
}

func (l *Logger) Finish(t MessageType, d time.Duration, msg string, args ...interface{}) {
	l.log(t, d, true, msg, args)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.logf(INFO, msg, args)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logf(DEBUG, msg, args)
}

func (l *Logger) logf(t MessageType, str string, args ...interface{}) {
	l.log(t, -1, true, fmt.Sprintf(str, args))
}

func (l *Logger) log(t MessageType,
                     d time.Duration,
                     ok bool,
                     msg string,
                     args ...interface{}) {
	if !l.enabled {
		return
	}

	var typeStr string
	switch t {
	case ACTION:
		if l.json {
			typeStr = "!"
		} else {
			typeStr = "\033[34m!\033[0m"
		}
	case CHECK:
		if l.json {
			typeStr = "✓"
			if !ok {
				typeStr = "✗"
			}
		} else {
			typeStr = "\033[32m✓\033[0m"
			if !ok {
				typeStr = "\033[31m✗\033[0m"
			}
		}
	case PLAN:
		typeStr = "#"
		if !l.json {
			msg = fmt.Sprintf("\033[4m%s\033[24m", msg)
		}
	case SERIAL:
		typeStr = "-"
	case PARALLEL:
		typeStr = "="
	case INFO:
		typeStr = "I"
	case DEBUG:
		if !l.debug {
			return
		}
		typeStr = "D"
	}

	if l.json {
		data := map[string]interface{}{
			"type": typeStr,
			"timestamp": time.Now().UTC().Format(jsISODateFormat),
			"message": msg,
		}

		if l.debug && t == ACTION {
			data["result"] = args[0]
		}

		if t == CHECK {
			data["pass"] = ok
		}

		if d > 0 {
			data["duration"] = l.sw.Since()
		}

		out, _ := json.Marshal(data)
		fmt.Printf("%s\n", out)
	} else {
		nowStr := ""
		if t != CHECK {
			nowStr = fmt.Sprintf(" {%.3fs}", l.sw.Since().Seconds())
		}

		durStr := ""
		if d > 0 {
			durStr = fmt.Sprintf(" (%s)", d.String())
		}

		fmt.Printf("[%s]%s%s %s\n", typeStr, nowStr, durStr, msg)
	}
}

func (l *Logger) Check(ok bool, msg string, data map[string]interface{}) {
	l.log(CHECK, -1, ok, msg, data)
}
