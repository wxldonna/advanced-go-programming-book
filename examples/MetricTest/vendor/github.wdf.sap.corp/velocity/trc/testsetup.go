// (C) 2016 SAP SE or an SAP affiliate company. All rights reserved.
package trc

import (
	"fmt"
	"os"
)

// ApplyTestEnvConfig allies a trace configuration for tests
func ApplyTestEnvConfig() {
	// first configure trace topics
	if val, present := os.LookupEnv("V2N_TRC"); present {
		err := applyConfigString(val, true)
		if err != nil {
			fmt.Println("Error when configuring tracer from env variable:", err)
			os.Exit(1)
		}
	}
}
