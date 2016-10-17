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
)

const (
	L_OFF Level = iota
	L_DEFAULT
	L_DETAIL
	L_INFO
	L_DEBUG
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

func (l *Logger) Start(t MessageType, msg string) {
	l.log(t, -1, nil, nil, msg, nil)
}

func (l *Logger) Finish(t MessageType, ok bool, d time.Duration, msg string, data interface{}) {
	l.log(t, d, util.BoolFor(ok), nil, msg, data)
}

func (l *Logger) Error(t MessageType, d time.Duration, err error, msg string) {
	l.log(t, d, util.False, err, msg, nil)
}

func (l *Logger) Info(msg string) {
	l.log(INFO, -1, nil, nil, msg, nil)
}

func (l *Logger) Debug(msg string) {
	l.log(DEBUG, -1, nil, nil, msg, nil)
}

func (l *Logger) Check(ok bool, msg string, data map[string]interface{}) {
	l.log(CHECK, -1, util.BoolFor(ok), nil, msg, data)
}

func (l *Logger) enabled(level Level) bool {
	return l.level >= level
}

func (l *Logger) levelForMessage(t MessageType, ok *util.Bool) Level {
	start := ok == nil

	switch t {
	case PLAN, ACTION:
		if start {
			return L_DETAIL
		}
		return L_DEFAULT
	case CHECK:
		if ok.Value() {
			return L_DETAIL
		}
		return L_DEFAULT
	case SERIAL, PARALLEL, INFO:
		return L_INFO
	case DEBUG:
		return L_DEBUG
	default:
		return L_OFF
	}
}

func (l *Logger) log(t MessageType,
					 d time.Duration,
					 ok *util.Bool,
					 err error,
					 msg string,
					 data interface{}) {
	if l.level == L_OFF {
		return
	}
	if !l.enabled(l.levelForMessage(t, ok)) {
		return
	}

	if l.json {
		l.logJSON(t, d, ok, err, msg, data)
	} else {
		l.logHuman(t, d, ok, err, msg)
	}
}

func (l *Logger) logJSON(t MessageType,
						 d time.Duration,
						 ok *util.Bool,
						 err error,
						 msg string,
						 data interface{}) {
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

	m := map[string]interface{}{
		"type": typeStr,
		"timestamp": time.Now().UTC().Format(jsISODateFormat),
		"message": msg,
	}

	if l.enabled(L_DEBUG) && t == ACTION {
		m["result"] = data
	}

	if t == CHECK && ok != nil {
		m["pass"] = ok.Value()
	}

	if err != nil {
		m["error"] = err.Error()
	}

	if d > 0 {
		m["duration"] = l.sw.Since()
	}

	out, _ := json.Marshal(m)
	fmt.Printf("%s\n", out)
}

func (l *Logger) logHuman(t MessageType,
						  d time.Duration,
						  ok *util.Bool,
						  err error,
						  msg string) {
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
			typeStr = "\033[35m?\033[0m"
		} else {
			typeStr = "?"
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

	okStr := "."
	if ok != nil {
		if l.colorize {
			if ok.Value() {
				okStr = "\033[32m✓\033[0m"
			} else {
				okStr = "\033[31m✗\033[0m"
			}
		} else {
			if ok.Value() {
				okStr = "✓"
			} else {
				okStr = "✗"
			}
		}
	}

	nowStr := ""
	if t != CHECK {
		nowStr = fmt.Sprintf(" %.1fs", l.sw.Since().Seconds())
	}

	durStr := ""
	if d > 0 {
		durStr = fmt.Sprintf(" (%s)", d.String())
	}

	if err != nil {
		msg = fmt.Sprintf("%s: %s", msg, err.Error())
	}

	fmt.Printf("%s %s%s%s %s\n", typeStr, okStr, nowStr, durStr, msg)
}

