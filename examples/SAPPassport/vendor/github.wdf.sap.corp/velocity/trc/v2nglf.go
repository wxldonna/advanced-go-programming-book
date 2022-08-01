package trc

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"time"
)

func PrintV2nGlf(t *tracer, severity Level, now time.Time, pc uintptr, file string, line int, fn interface{}, format string, v ...interface{}) {
	if pc == 0 {
		file = "???"
		line = 0
	} else {
		file = shortFile(file)
	}
	buf := pool.Get().(*bytes.Buffer)
	defer pool.Put(buf)
	buf.Reset()
	formatHeaderV2nGlf(t, buf, now, severity, file, line)
	if fn != nil {
		if format != "" {
			fn.(FprintfFunc)(buf, format, v...)
		} else {
			fn.(FprintFunc)(buf, v...)
		}
	} else {
		if format != "" {
			fmt.Fprintf(buf, format, v...)
		} else {
			fmt.Fprintln(buf, v...)
			buf.Truncate(buf.Len() - 1)
		}
	}
	formatFooterV2nGlf(t, buf)
	outLock.Lock()
	defer outLock.Unlock()
	buf.WriteTo(out)
}

var (
	locationLen  int64
	componentLen int64
)

const (
	empty                 = "                                "
	maxLocationLen  int64 = 32
	maxComponentLen int64 = 14
)

func writePadding(b *bytes.Buffer, p *int64, m int64, l int) {
	l64 := int64(l)
	var c int64
	for {
		c = atomic.LoadInt64(p)
		if c >= l64 {
			break
		}
		if l64 > m || atomic.CompareAndSwapInt64(p, c, l64) {
			return
		}
	}
	b.WriteString(empty[:c-l64])
}

func formatHeaderV2nGlf(t *tracer, buf *bytes.Buffer, ti time.Time, severity Level, file string, line int) {
	buf.Grow(128)
	ttoa(buf, ti)
	buf.WriteByte('|')
	buf.WriteString(t.topic.Name)
	writePadding(buf, &componentLen, maxComponentLen, len(t.topic.Name))
	buf.WriteByte('|')
	buf.WriteString(file)
	buf.WriteByte(':')
	l := buf.Len()
	itoa(buf, line, -1)
	l = buf.Len() - l
	writePadding(buf, &locationLen, maxLocationLen, len(file)+1+l)
	buf.WriteByte('|')
	buf.WriteString(glfSeverityFormats[severity.uint64])
	buf.WriteByte('|')
	if flen := len(t.frags); flen > 0 {
		buf.WriteString("{")
		if flen == 1 {
			buf.Write(t.frags[0])
		} else {
			end := flen - 1
			for i := 0; i < end; i++ {
				buf.Write(t.frags[i])
				buf.WriteByte(',')
			}
			buf.Write(t.frags[end])
		}
		buf.WriteString("} ")
	}
}

func formatFooterV2nGlf(t *tracer, buf *bytes.Buffer) {
	buf.WriteString("\x1e\n")
}
