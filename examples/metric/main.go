package main

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/nonrecording"
	"go.opentelemetry.io/otel/metric/unit"
)

func main() {
	// In a library or program this would be provided by otel.GetMeterProvider().
	meterProvider := nonrecording.NewNoopMeterProvider()

	workDuration, err := meterProvider.Meter("go.opentelemetry.io/otel/metric#SyncExample").SyncInt64().Histogram(
		"workDuration",
		instrument.WithUnit(unit.Milliseconds))
	if err != nil {
		fmt.Println("Failed to register instrument")
		panic(err)
	}

	startTime := time.Now()
	ctx := context.Background()
	// Do work
	// ...
	workDuration.Record(ctx, time.Since(startTime).Milliseconds())

}
