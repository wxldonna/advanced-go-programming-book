package main

import (
	"context"
	"flag"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.wdf.sap.corp/velocity/trc"
)

const (
	trcInfoOperation    = "operation"
	trcInfoDuration     = "duration"
	trcInfoStartedUTC   = "started-at"
	trcInfoSpanID       = "span-id"
	trcInfoParentSpanID = "parent-span-id"
	trcInfoTagPrefix    = "tag:"
)

func init() {
	trc.Application = "MyApp"
}
func gettracerinfo() []trc.Info {
	infos := make([]trc.Info, 0)
	infos = append(infos,
		trc.NewInfo("key1", "value1"),
		trc.NewInfo("key2", "value2"),
	)
	return infos
}

var tracer = trc.InitTraceTopic("topic", "executable package")
var trcFlag = flag.String("trc", "", "e.g. -trc=debug,main:warning")

func main() {
	flag.Parse()
	if err := trc.ReconfigFromString(*trcFlag); err != nil {
		log.Fatal(err)
	}
	tracer = tracer.Sub(gettracerinfo()...)
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
