package trc

import (
	"math"
	"strconv"
	"sync/atomic"
)

func newTraceTopic(name, descr string) *traceTopic {
	return &traceTopic{
		topicLevels: newTopicLevels(LevelUndefined, LevelUndefined),
		Name:        name,
		Description: descr,
	}
}

// traceTopic represents a trace, identified by a unique ID, which is responsible for a part of a system.
// At runtime, trace levels can be dynamically configured per trace topic.
type traceTopic struct {
	topicLevels
	Name        string
	Description string
}

// topicLevels encodes the global and local level currently configured for a trace topic.
// Must be accessed via set*() and active() functions only! (All those functions use atomic read/writes.)
// set*() functions must be synchronized by the caller.
type topicLevels struct {
	uint64 // uint16 (global) + uint16 (local) + uint32 (effective)
}

// setters are not safe for parallel use

const (
	globalMask    = uint64(math.MaxUint16) << 48
	localMask     = uint64(math.MaxUint16) << 32
	effectiveMask = uint64(math.MaxUint32)
)

func newTopicLevels(global, local Level) topicLevels {
	var tl topicLevels
	tl.set(global, local)
	return tl
}

func (tl *topicLevels) Stringer() string {
	global, local := tl.levels()
	return "|global:" + global.String() + "(" + strconv.FormatUint(uint64(global.uint64), 2) + ")|local:" + local.String() + "(" + strconv.FormatUint(uint64(local.uint64), 2) + ")|"
}

func (tl *topicLevels) set(global, local Level) {
	effective := local
	if effective.uint64 == LevelUndefined.uint64 {
		effective = global
	}
	atomic.StoreUint64(&(tl.uint64), (global.uint64<<48)|(local.uint64<<32)|effective.uint64)
}

func (tl *topicLevels) setGlobal(global Level) {
	// non-atomic load and store is ok
	// writing operations are assumed to be synchronized on a higher level (see registry)
	curr := atomic.LoadUint64(&(tl.uint64))
	effective := (curr & localMask) >> 32
	if effective == LevelUndefined.uint64 {
		effective = global.uint64
	}
	val := (global.uint64 << 48) | (curr & localMask) | effective
	atomic.StoreUint64(&(tl.uint64), val)
}

func (tl *topicLevels) setLocal(local Level) {
	// non-atomic load and store is ok
	// writing operations are assumed to be synchronized on a higher level (see registry)
	curr := atomic.LoadUint64(&(tl.uint64))
	effective := local.uint64
	if effective == LevelUndefined.uint64 {
		effective = curr >> 48
	}
	val := (curr & globalMask) | local.uint64<<32 | effective
	atomic.StoreUint64(&(tl.uint64), val)
}

func (tl *topicLevels) local() Level {
	curr := atomic.LoadUint64(&(tl.uint64))
	return Level{uint64: (curr & localMask) >> 32}
}

func (tl *topicLevels) global() Level {
	curr := atomic.LoadUint64(&(tl.uint64))
	return Level{uint64: curr >> 48}
}

func (tl *topicLevels) levels() (global Level, local Level) {
	curr := atomic.LoadUint64(&(tl.uint64))
	return Level{uint64: (curr >> 48)}, Level{uint64: (curr & localMask) >> 32}
}

func (tl *topicLevels) active(level Level) bool {
	curr := atomic.LoadUint64(&(tl.uint64))
	return uint32(curr) >= uint32(level.uint64)
}
