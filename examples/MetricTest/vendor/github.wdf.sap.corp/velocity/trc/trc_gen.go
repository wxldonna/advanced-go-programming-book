package trc

/*
 *
 *
 * DO NOT MODIFY
 *
 * generated via 'go run gen/main.go | gofmt > trc_gen.go'
 *
 *
 */

import (
	"context"
	"io"
	"runtime"
	"sync/atomic"
	"time"
)

/*
 * Tracer interface
 */

type Tracer interface {
	Sub(...Info) Tracer
	SubFromContext(context.Context) Tracer

	IsFatal() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	FatalFn(fn FprintFunc, args ...interface{})
	FatalFnf(fn FprintfFunc, format string, args ...interface{})

	IsError() bool
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	ErrorFn(fn FprintFunc, args ...interface{})
	ErrorFnf(fn FprintfFunc, format string, args ...interface{})

	IsWarning() bool
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	WarningFn(fn FprintFunc, args ...interface{})
	WarningFnf(fn FprintfFunc, format string, args ...interface{})

	IsInfo() bool
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	InfoFn(fn FprintFunc, args ...interface{})
	InfoFnf(fn FprintfFunc, format string, args ...interface{})

	IsDebug() bool
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	DebugFn(fn FprintFunc, args ...interface{})
	DebugFnf(fn FprintfFunc, format string, args ...interface{})
}

type FprintFunc func(out io.Writer, args ...interface{}) (int, error)
type FprintfFunc func(out io.Writer, format string, args ...interface{}) (int, error)

/*
 * Levels
 */

// Level (severity) of a log entry
type Level struct{ uint64 } // do not expose actual numeric value

var (
	LevelUndefined = Level{lvlUndefined}
	LevelNone      = Level{lvlNone}
	LevelFatal     = Level{lvlFatal}
	LevelError     = Level{lvlError}
	LevelWarning   = Level{lvlWarning}
	LevelInfo      = Level{lvlInfo}
	LevelDebug     = Level{lvlDebug}
)

const (
	lvlUndefined = iota
	lvlNone
	lvlFatal
	lvlError
	lvlWarning
	lvlInfo
	lvlDebug
)

func (l Level) String() string {
	switch l.uint64 {
	case lvlUndefined:
		return "UNDEFINED"
	case lvlNone:
		return "NONE"
	case lvlFatal:
		return "FATAL"
	case lvlError:
		return "ERROR"
	case lvlWarning:
		return "WARNING"
	case lvlInfo:
		return "INFO"
	case lvlDebug:
		return "DEBUG"
	default:
		return "!ILLEGAL!"
	}
}

/*
 * Tracer implementations
 */

type tracer struct {
	topic *traceTopic
	frags infoFrags
	ctx   context.Context
}

func (t *tracer) IsFatal() bool {
	curr := atomic.LoadUint64(&(t.topic.topicLevels.uint64))
	return uint32(curr) >= uint32(lvlFatal)
}
func (t *tracer) IsError() bool {
	curr := atomic.LoadUint64(&(t.topic.topicLevels.uint64))
	return uint32(curr) >= uint32(lvlError)
}
func (t *tracer) IsWarning() bool {
	curr := atomic.LoadUint64(&(t.topic.topicLevels.uint64))
	return uint32(curr) >= uint32(lvlWarning)
}
func (t *tracer) IsInfo() bool {
	curr := atomic.LoadUint64(&(t.topic.topicLevels.uint64))
	return uint32(curr) >= uint32(lvlInfo)
}
func (t *tracer) IsDebug() bool {
	curr := atomic.LoadUint64(&(t.topic.topicLevels.uint64))
	return uint32(curr) >= uint32(lvlDebug)
}

func (t *tracer) Fatal(args ...interface{}) {
	if t.IsFatal() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelFatal, now, pc, file, line, nil, "", args...)
	}
}

func (t *tracer) Fatalf(format string, args ...interface{}) {
	if t.IsFatal() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelFatal, now, pc, file, line, nil, format, args...)
	}
}

func (t *tracer) FatalFn(fn FprintFunc, args ...interface{}) {
	if t.IsFatal() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelFatal, now, pc, file, line, fn, "", args...)
	}
}

func (t *tracer) FatalFnf(fn FprintfFunc, format string, args ...interface{}) {
	if t.IsFatal() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelFatal, now, pc, file, line, fn, format, args...)
	}
}

func (t *tracer) Error(args ...interface{}) {
	if t.IsError() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelError, now, pc, file, line, nil, "", args...)
	}
}

func (t *tracer) Errorf(format string, args ...interface{}) {
	if t.IsError() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelError, now, pc, file, line, nil, format, args...)
	}
}

func (t *tracer) ErrorFn(fn FprintFunc, args ...interface{}) {
	if t.IsError() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelError, now, pc, file, line, fn, "", args...)
	}
}

func (t *tracer) ErrorFnf(fn FprintfFunc, format string, args ...interface{}) {
	if t.IsError() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelError, now, pc, file, line, fn, format, args...)
	}
}

func (t *tracer) Warning(args ...interface{}) {
	if t.IsWarning() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelWarning, now, pc, file, line, nil, "", args...)
	}
}

func (t *tracer) Warningf(format string, args ...interface{}) {
	if t.IsWarning() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelWarning, now, pc, file, line, nil, format, args...)
	}
}

func (t *tracer) WarningFn(fn FprintFunc, args ...interface{}) {
	if t.IsWarning() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelWarning, now, pc, file, line, fn, "", args...)
	}
}

func (t *tracer) WarningFnf(fn FprintfFunc, format string, args ...interface{}) {
	if t.IsWarning() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelWarning, now, pc, file, line, fn, format, args...)
	}
}

func (t *tracer) Info(args ...interface{}) {
	if t.IsInfo() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelInfo, now, pc, file, line, nil, "", args...)
	}
}

func (t *tracer) Infof(format string, args ...interface{}) {
	if t.IsInfo() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelInfo, now, pc, file, line, nil, format, args...)
	}
}

func (t *tracer) InfoFn(fn FprintFunc, args ...interface{}) {
	if t.IsInfo() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelInfo, now, pc, file, line, fn, "", args...)
	}
}

func (t *tracer) InfoFnf(fn FprintfFunc, format string, args ...interface{}) {
	if t.IsInfo() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelInfo, now, pc, file, line, fn, format, args...)
	}
}

func (t *tracer) Debug(args ...interface{}) {
	if t.IsDebug() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelDebug, now, pc, file, line, nil, "", args...)
	}
}

func (t *tracer) Debugf(format string, args ...interface{}) {
	if t.IsDebug() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelDebug, now, pc, file, line, nil, format, args...)
	}
}

func (t *tracer) DebugFn(fn FprintFunc, args ...interface{}) {
	if t.IsDebug() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelDebug, now, pc, file, line, fn, "", args...)
	}
}

func (t *tracer) DebugFnf(fn FprintfFunc, format string, args ...interface{}) {
	if t.IsDebug() {
		now := time.Now()
		t.activate()
		pc, file, line, _ := runtime.Caller(1)
		outPrint(t, LevelDebug, now, pc, file, line, fn, format, args...)
	}
}

func (t *tracer) Sub(info ...Info) Tracer {
	return &tracer{
		topic: t.topic,
		frags: t.frags.subFromFrag(serializeInfos(info...)),
	}
}

func (t *tracer) SubFromContext(ctx context.Context) Tracer {
	return &tracer{
		topic: t.topic,
		frags: t.frags,
		ctx:   ctx,
	}
}

func (t *tracer) activate() {
	if t.ctx != nil {
		t.frags = t.frags.subFromFrags(infoFragmentsFromContext(t.ctx))
		t.ctx = nil
	}
}
