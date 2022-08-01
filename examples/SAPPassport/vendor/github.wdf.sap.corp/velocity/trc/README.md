# trc
Golang tracer

### Example
```go
package main

import (
    "context"
    "flag"
    "log"

    "github.com/davecgh/go-spew"
    "github.wdf.sap.corp/velocity/trc"
)

func init() {
    trc.Application = "MyApp"
}

var tracer = trc.InitTraceTopic("main", "executable package")
var trcFlag = flag.String("trc", "", "e.g. -trc=debug,main:warning")

func main() {
    flag.Parse()
    if err := trc.ReconfigFromString(*trcFlag); err != nil {
        log.Fatal(err)
    }

    tracer.Info("app started")

    ctx := context.Background()
    for worker := 0; worker < 10; worker++ {
        tracer.Debugf("starting worker %d", worker)
        
        workerCtx := trc.AttachTrcInfo(ctx, trc.NewInfo("worker", worker))
        go work(workerCtx)
    }
}

func work(ctx context.Context) {
    sub := tracer.SubFromContext(ctx)

    sub.DebugFnf(spew.Fprintf, "context: %#+v", ctx)

    if sub.IsWarning() {
        sub.Warning("this", "is", "expensive")
    }
}
```

### Remarks on usage
#### InitTraceTopic
The value of `trc.Application` is printed into the "Tracer" column in GLF. If a Go library initializes a topic, it SHOULD prefix its components with a library specific name seprarated by a dot. Example:
```
var libTracer = trc.InitTraceTopic("v2auth.srvauth", "Server authenticator")
```
This allows to correlate the output of a reused library to a specific application.

#### Trace functions
The call of a function with varadic interfaces as arguments is quite expensive. The more arguments the more has the runtime to convert to interfaces. Therefore, it is adviced to use the level check functions in performance sensitive sections of code:
```go
if tracer.IsDebug() {
    tracer.Debug(largeArray...)
}
```

#### Custom message formatting
The `*Fn` functions can be used to format the arguments in a user defined way. You can pass in functions with `fmt.Fprintf`-like signature. A prominent example that provides those functions is `go-spew` for general object serialization. But you can also use closures:
```
b := []byte{123}
tracer.DebugFn(func(out io.Writer, args ...interface{}) {
    out.Write([]byte(hex.EncodeToString(b)))
})
```
Due to late evaluation of closures, this is quite performant, if you don't pass in any arguments `args`.

#### Custom output (default stdout)
If you want to change the output of the trace entries, you can specify your own `io.Writer` by calling `trc.SetOutput`. For example:
```
trc.SetOutput(ioutil.Discard)
```
Futhermore, if you need to synchronize the tracer with other writers to the same `io.Writer`, you can set a custom lock that fulfills the `sync.Locker` interface with the `trc.SetOutputLock` function:
```
var myLock sync.Mutex
trc.SetOutputLock(&myLock)
```

#### Command line flag
For convenience, you can call `trc.InitFlag()` somewhere before `flag.Parse()` to register a common `-trc` flag, which gets automatically parsed and configures the topics.

### Benchmark to /dev/null
```
Benchmark_trc/Check/Inactive/Serial-16          200000000                6.23 ns/op            0 B/op          0 allocs/op
Benchmark_trc/Check/Active/Serial-16            200000000                6.26 ns/op            0 B/op          0 allocs/op
Benchmark_trc/Trace/Inactive/Serial-16          10000000               129 ns/op              32 B/op          2 allocs/op
Benchmark_trc/Trace/Active/Serial-16              500000              2424 ns/op              32 B/op          2 allocs/op
Benchmark_trc/Check/Inactive/Parallel(200)-16   2000000000               1.23 ns/op            0 B/op          0 allocs/op
Benchmark_trc/Check/Active/Parallel(200)-16     2000000000               1.23 ns/op            0 B/op          0 allocs/op
Benchmark_trc/Trace/Inactive/Parallel(200)-16   20000000                64.7 ns/op            32 B/op          2 allocs/op
Benchmark_trc/Trace/Active/Parallel(200)-16      3000000               578 ns/op              33 B/op          2 allocs/op

Benchmark_vasglf/Inactive/Serial-16             10000000               230 ns/op              32 B/op          2 allocs/op
Benchmark_vasglf/Active/Serial-16                 500000              3786 ns/op             192 B/op          4 allocs/op
Benchmark_vasglf/Inactive/Parallel(200)-16      10000000               153 ns/op              32 B/op          2 allocs/op
Benchmark_vasglf/Active/Parallel(200)-16          300000              4797 ns/op             194 B/op          4 allocs/op
```
