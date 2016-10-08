package logging

import (
	"fmt"
	"time"
	"github.com/troykinsella/crash/util"
	"encoding/json"
)

var jsISODateFormat = "2006-01-02T15:04:05.999999999Z07:00"
// less precision: "2006-01-02T15:04:05Z07:00"

type Level uint8
type MessageType uint8

const (
	PLAN MessageType = iota
	SERIAL
	PARALLEL
	ACTION
	CHECK
	INFO
	DEBUG
	TRACE
)

const (
	L_OFF Level = iota
	L_DEFAULT
	L_INFO
	L_DEBUG
	L_TRACE
)

type Logger struct {
	level    Level
	json     bool
	colorize bool
	sw       *util.StopWatch
}

func NewLogger(level Level,
               colorize bool,
               json bool) (*Logger) {
	return &Logger{
		level: level,
		colorize: colorize,
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

func (l *Logger) Trace(msg string, args ...interface{}) {
	l.logf(TRACE, msg, args)
}

func (l *Logger) logf(t MessageType, str string, args ...interface{}) {
	l.log(t, -1, true, fmt.Sprintf(str, args))
}

func (l *Logger) enabled(level Level) bool {
	return l.level >= level
}

func (l *Logger) levelForMessage(t MessageType) Level {
	switch t {
	case PLAN, ACTION: return L_DEFAULT
	case SERIAL, PARALLEL, INFO: return L_INFO
	case CHECK, DEBUG: return L_DEBUG
	case TRACE: return L_TRACE
	default: return L_OFF
	}
}

func (l *Logger) log(t MessageType,
                     d time.Duration,
                     ok bool,
                     msg string,
                     args ...interface{}) {
	if l.level == L_OFF {
		return
	}
	if !l.enabled(l.levelForMessage(t)) {
		return
	}

	if l.json {
		l.logJSON(t, d, ok, msg, args...)
	} else {
		l.logHuman(t, d, ok, msg, args...)
	}
}

func (l *Logger) logJSON(t MessageType,
						 d time.Duration,
						 ok bool,
						 msg string,
						 args ...interface{}) {
	var typeStr string
	switch t {
	case ACTION:
		typeStr = "action"
	case CHECK:
		typeStr = "check"
	case PLAN:
		typeStr = "plan"
	case SERIAL:
		typeStr = "serial"
	case PARALLEL:
		typeStr = "parallel"
	case INFO:
		typeStr = "info"
	case DEBUG:
		typeStr = "debug"
	}

	data := map[string]interface{}{
		"type": typeStr,
		"timestamp": time.Now().UTC().Format(jsISODateFormat),
		"message": msg,
	}

	if l.enabled(L_DEBUG) && t == ACTION {
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
}

func (l *Logger) logHuman(t MessageType,
						  d time.Duration,
						  ok bool,
						  msg string,
						  args ...interface{}) {
	var typeStr string
	switch t {
	case ACTION:
		if l.colorize {
			typeStr = "\033[34m!\033[0m"
		} else {
			typeStr = "!"
		}
	case CHECK:
		if l.colorize {
			if ok {
				typeStr = "\033[32m✓\033[0m"
			} else {
				typeStr = "\033[31m✗\033[0m"
			}
		} else {
			if ok {
				typeStr = "✓"
			} else {
				typeStr = "✗"
			}
		}
	case PLAN:
		if l.colorize {
			msg = fmt.Sprintf("\033[4m%s\033[24m", msg)
		}
		typeStr = "#"
	case SERIAL:
		typeStr = "-"
	case PARALLEL:
		typeStr = "="
	case INFO:
		typeStr = "I"
	case DEBUG:
		typeStr = "D"
	}

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

func (l *Logger) Check(ok bool, msg string, data map[string]interface{}) {
	l.log(CHECK, -1, ok, msg, data)
}
