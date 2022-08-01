package trc

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"sync"
)

/* global flag */

var defRegFlag = newFlag()
var onceFlagInit = sync.Once{}

// Must be called before flag.Parse()
func InitFlag() {
	onceFlagInit.Do(func() {
		flag.Var(defRegFlag, "trc", "trace topic configuration")
	})
}

/* flag implementation */

type trcFlag struct {
	v string
}

func newFlag() *trcFlag {
	return &trcFlag{}
}

func (tf *trcFlag) String() string { return tf.v }

func (tf *trcFlag) Set(v string) (err error) {
	err = applyConfigString(v, true)
	tf.v = v
	return
}

func ValidateConfigString(cfg string) error {
	err, _ := parseConfigString(cfg)
	return err
}

// applyConfigString is only to be used by tests or to implmenet
// initial config after initialization of all trace topics.
// Can for example be used to apply config from env variable right
// before flag parsing.
func applyConfigString(v string, abortOnError bool) error {
	err, initOps := parseConfigString(v)
	if err != nil {
		return err
	}
	// do not abort for errors during trace topic registration
	for _, op := range initOps {
		err := op(defReg) // call init on default registry
		if err != nil {
			trcTracer.Error("Unable to register trace topic:", err)
		}
	}
	return nil
}

func parseConfigString(v string) (error, []func(*registry) error) {
	var initOps []func(*registry) error
	initDefs := strings.Split(v, ",")
	var errs []string
	for _, def := range initDefs {
		s := strings.Split(def, ":")
		switch len(s) {
		case 0:
			continue
		case 1:
			switch strings.ToLower(strings.TrimSpace(s[0])) {
			case flagLevelDebug:
				initOps = append(initOps, func(r *registry) error { return r.reconfigGlobal(LevelDebug) })
			case flagLevelInfo:
				initOps = append(initOps, func(r *registry) error { return r.reconfigGlobal(LevelInfo) })
			case flagLevelWarning:
				initOps = append(initOps, func(r *registry) error { return r.reconfigGlobal(LevelWarning) })
			case flagLevelError:
				initOps = append(initOps, func(r *registry) error { return r.reconfigGlobal(LevelError) })
			case flagLevelFatal:
				initOps = append(initOps, func(r *registry) error { return r.reconfigGlobal(LevelFatal) })
			case flagLevelNone:
				initOps = append(initOps, func(r *registry) error { return r.reconfigGlobal(LevelNone) })
			default:
				errs = append(errs, fmt.Sprintf("Illegal global trace level: %q\n", s[0]))
			}
		case 2:
			switch strings.ToLower(strings.TrimSpace(s[1])) {
			case flagLevelDebug:
				initOps = append(initOps, func(r *registry) error { return r.reconfigTopic(strings.TrimSpace(s[0]), LevelDebug) })
			case flagLevelInfo:
				initOps = append(initOps, func(r *registry) error { return r.reconfigTopic(strings.TrimSpace(s[0]), LevelInfo) })
			case flagLevelWarning:
				initOps = append(initOps, func(r *registry) error { return r.reconfigTopic(strings.TrimSpace(s[0]), LevelWarning) })
			case flagLevelError:
				initOps = append(initOps, func(r *registry) error { return r.reconfigTopic(strings.TrimSpace(s[0]), LevelError) })
			case flagLevelFatal:
				initOps = append(initOps, func(r *registry) error { return r.reconfigTopic(strings.TrimSpace(s[0]), LevelFatal) })
			case flagLevelNone:
				initOps = append(initOps, func(r *registry) error { return r.reconfigTopic(strings.TrimSpace(s[0]), LevelNone) })
			default:
				errs = append(errs, fmt.Sprintf("Illegal trace level for topic %q: %q\n", strings.TrimSpace(s[0]), s[1]))
			}
		default:
			errs = append(errs, fmt.Sprintf("Unable to parse trace topic initialization flag: %q\n", def))
		}
	}
	var err error
	if errs != nil {
		err = errors.New(fmt.Sprint(errs))
	}
	return err, initOps
}

const (
	flagLevelDebug   = "debug"
	flagLevelInfo    = "info"
	flagLevelWarning = "warning"
	flagLevelError   = "error"
	flagLevelFatal   = "fatal"
	flagLevelNone    = "none"
)
