package trc

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
	"time"
	"unicode/utf8"
)

func severityFormat(l Level) string {
	return glfSeverityFormats[l.uint64]
}

var Application = ""
var glfSeverityFormats = []string{
	"UNDEF",
	"NONE ",
	"FATAL",
	"ERROR",
	"WARN ",
	"INFO ",
	"DEBUG",
}

var pool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 64))
	},
}

func PrintV2Glf(t *tracer, severity Level, now time.Time, pc uintptr, file string, line int, fn interface{}, format string, v ...interface{}) {
	var function string
	if pc == 0 {
		function = "???"
		file = "???"
		line = 0
	} else {
		f := runtime.FuncForPC(pc)
		if f != nil {
			function = shortFunc(f.Name())
		} else {
			function = "???"
		}
		file = shortFile(file)
	}
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()
	formatHeaderGlf(buf, now, severity)
	tBuf := pool.Get().(*bytes.Buffer)
	tBuf.Reset()
	if fn != nil {
		if format != "" {
			fn.(FprintfFunc)(tBuf, format, v...)
		} else {
			fn.(FprintFunc)(tBuf, v...)
		}
	} else {
		if format != "" {
			fmt.Fprintf(tBuf, format, v...)
		} else {
			fmt.Fprintln(tBuf, v...)
			tBuf.Truncate(tBuf.Len() - 1)
		}
	}
	tbBytes := tBuf.Bytes()
	var last int = 0
	for i, c := range tbBytes {
		if (c-31 >= 93) && (c < utf8.RuneSelf) {
			switch c {
			case '|', '\x1e', '\x1b':
				buf.Write(tbBytes[last:i])
				last = i
				buf.WriteByte('\x1b')
			}
		}
	}
	buf.Write(tbBytes[last:])
	formatFooterGlf(t, buf, function, file, line)
	outLock.Lock()
	buf.WriteTo(out)
	outLock.Unlock()
	pool.Put(tBuf)
	pool.Put(buf)
}

func ttoa(buf *bytes.Buffer, ti time.Time) {
	var b [32]byte
	year, month, day := ti.Date()
	hour, min, sec := ti.Clock()
	nano := ti.Nanosecond() / 1000
	_, off := ti.Zone()
	for i := 3; i >= 0; i-- {
		b[i] = byte(year%10) + '0'
		year /= 10
	}
	b[4] = '-'
	b[5] = byte(month/10) + '0'
	b[6] = byte(month%10) + '0'
	b[7] = '-'
	b[8] = byte(day/10) + '0'
	b[9] = byte(day%10) + '0'
	b[10] = ' '
	b[11] = byte(hour/10) + '0'
	b[12] = byte(hour%10) + '0'
	b[13] = ':'
	b[14] = byte(min/10) + '0'
	b[15] = byte(min%10) + '0'
	b[16] = ':'
	b[17] = byte(sec/10) + '0'
	b[18] = byte(sec%10) + '0'
	b[19] = '.'
	for i := 25; i >= 20; i-- {
		b[i] = byte(nano%10) + '0'
		nano /= 10
	}
	b[26] = '|'
	if off < 0 {
		off *= -1
		b[27] = '-'
	} else {
		b[27] = '+'
	}
	b[28] = byte(off/36000) + '0'
	b[29] = byte(off/3600%10) + '0'
	b[30] = byte(off%3600/600) + '0'
	b[31] = byte(off/60%10) + '0'
	buf.Write(b[:])
}

func formatHeaderGlf(buf *bytes.Buffer, ti time.Time, severity Level) {
	ttoa(buf, ti)
	buf.WriteByte('|')
	buf.WriteString(glfSeverityFormats[severity.uint64])
	buf.WriteByte('|')
}

func formatFooterGlf(t *tracer, buf *bytes.Buffer, function, file string, line int) {
	if flen := len(t.frags); flen > 0 {
		buf.WriteString(" {")
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
		buf.WriteByte('}')
	}
	buf.WriteByte('|')
	buf.WriteString(Application)
	buf.WriteByte('|')
	buf.WriteString(t.topic.Name)
	buf.WriteByte('|')
	itoa(buf, int(GoroutineID()), -1)
	buf.WriteByte('|')
	buf.WriteString(function)
	buf.WriteByte('|')
	buf.WriteString(file)
	buf.WriteByte('(')
	itoa(buf, line, -1)
	buf.WriteString(")\x1e\n")
}
