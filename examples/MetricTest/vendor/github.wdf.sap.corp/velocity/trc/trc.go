// (C) 2016 SAP SE or an SAP affiliate company. All rights reserved.
package trc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
)

var trcTracer Tracer

func init() {
	// register trace topic for this package
	trcTracer = InitTraceTopic("trc", "Tracer of package trc")
}

func newTopicTracer(root *traceTopic, info ...Info) *tracer {
	return &tracer{
		topic: root,
		frags: infoFrags(nil).subFromFrag(serializeInfos(info...)),
	}
}

type Info struct {
	K string
	V interface{}
}

func NewInfo(k string, v interface{}) Info {
	return Info{K: k, V: v}
}

func NewInfoCaller(k string) Info {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	} else {
		file = shortFile(file)
	}
	return NewInfo(k, file+":"+strconv.Itoa(line))
}

type infoFrag []byte
type infoFrags []infoFrag

func (parent infoFrags) subFromFrag(f infoFrag) infoFrags {
	if len(f) <= 0 {
		return parent
	}
	if len(parent) <= 0 {
		return infoFrags{f}
	}
	ret := make(infoFrags, len(parent)+1)
	copy(ret, parent)
	ret[len(parent)] = f
	return ret
}

func (parent infoFrags) subFromFrags(f infoFrags) infoFrags {
	if len(f) <= 0 {
		return parent
	}
	if len(parent) <= 0 {
		return f
	}
	ret := make(infoFrags, len(parent)+len(f))
	copy(ret, parent)
	copy(ret[len(parent):], f)
	return ret
}

// serializeContext converts the given contexts to a json fragment, e.g.: `"k1": "v1", "k2": "v2"`
func serializeInfos(info ...Info) []byte {
	if len(info) == 0 {
		return nil
	}
	const bufSizePerEntry = 25
	buf := bytes.NewBuffer(make([]byte, 0, bufSizePerEntry*len(info)))
	serializeInfosBuf(buf, info...)
	return buf.Bytes()
}

func serializeInfosBuf(buf *bytes.Buffer, info ...Info) {
	enc := json.NewEncoder(buf)
	if len(info) > 0 {
		for i, c := range info {
			// encode key
			switch i {
			case 0:
				_, _ = buf.WriteString(c.K + "=")
			default:
				_, _ = buf.WriteString("," + c.K + "=")
			}
			// encode value
			serializeInfoValue(buf, c, enc)
		}
	}
}

// serializeInfoValue serializes the info value into a buffer. A Json encoder
// can optionally be given to reduce overhead.
func serializeInfoValue(buf *bytes.Buffer, info Info, enc *json.Encoder) {
	// fast-path for trivially serializable types
	switch v := info.V.(type) {
	case int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64, uintptr:
		_, _ = buf.WriteString(fmt.Sprint(v))
		return
	case string:
		b := []byte(v)
		// double quotes and backslashes in strings need to be properly escaped
		if bytes.IndexByte(b, '"') < 0 && bytes.IndexByte(b, '\\') < 0 {
			_ = buf.WriteByte('"')
			_, _ = buf.WriteString(v)
			_ = buf.WriteByte('"')
			return
		}
	}
	// slow-path for everything else
	if enc == nil {
		enc = json.NewEncoder(buf)
	}
	pos := buf.Len()
	err := enc.Encode(info.V)
	if err != nil {
		// enc could have written before the err has been reported
		buf.Truncate(pos)
		_, _ = buf.WriteString("\"<encoding_error>\"")
		// replace encoder (internal error flag cannot be unset)
		*enc = *json.NewEncoder(buf)
	} else {
		// truncate newline introduced by json.Encoder.Encode()
		buf.Truncate(buf.Len() - 1)
	}
}

var out = io.Writer(os.Stderr)
var outMu sync.Mutex
var outLock sync.Locker = &outMu
var outPrint = PrintV2Glf

func SetOutputLock(locker sync.Locker) {
	if locker == nil {
		panic("locker must not be nil")
	}
	outLock = locker
}

// SetOutput should only be used for testing
func SetOutput(newOut io.Writer) (resetFunc func()) {
	old := out
	resetFunc = func() { out = old }
	out = newOut
	return
}

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
// Knows the buffer has capacity.
func itoa(buf *bytes.Buffer, i int, wid int) {
	var u = uint(i)
	if u == 0 && wid <= 1 {
		buf.WriteByte('0')
		return
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	w := wid
	for ; u > 0 || w > 0; u /= 10 {
		bp--
		w--
		b[bp] = byte(u%10) + '0'
	}
	if uint(wid) > uint(len(b)-bp) {
		wid = len(b) - bp
	}
	buf.Write(b[bp : bp+wid])
}

func shortFile(file string) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return short
}

func shortFunc(f string) string {
	for i := len(f) - 1; i >= 0; i-- {
		if f[i] == '.' {
			return f[i+1:]
		}
	}
	return f
}
