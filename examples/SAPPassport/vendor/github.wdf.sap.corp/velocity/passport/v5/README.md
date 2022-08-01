[//]: # ( (C) 2021 SAP SE or an SAP affiliate company. All rights reserved.)

# SAP passport in Go

Here you can find a go implementation of the SAP Passport by the Data Intelligence team. For _true_ reference passport implementations please look [here](https://github.wdf.sap.corp/xdsr).

The SAP passport is a correlation mechanism that can be used for integration monitoring (across components) and request tracing (across services of one component). For more information on the SAP Passport in general please look at the [specification from 2020](https://sap.sharepoint.com/teams/IntelligentEnterpriseSuite-SupportFunctions/Shared%20Documents/04_Detect%20to%20correct/Integration%20Monitoring/Specifications/SAP_Passport.pdf) and the [earlier documentation](https://wiki.wdf.sap.corp/wiki/display/Introscope/Passport+Handling) (includes a description of trace flags).

## Installation

Import the passport library in your go program.

```go
import (
    "github.wdf.sap.corp/velocity/passport/v5"
    "github.wdf.sap.corp/velocity/passport/v5/middleware"
)
```

And download the package.

```bash
go get github.wdf.sap.corp/velocity/passport/v5/...
```

## Usage

The [specification from 2020](https://sap.sharepoint.com/teams/IntelligentEnterpriseSuite-SupportFunctions/Shared%20Documents/04_Detect%20to%20correct/Integration%20Monitoring/Specifications/SAP_Passport.pdf) documents how Cloud services should use the passport. The passport needs to be logged both when received with a request from another Cloud Service (inbound path) and when transferred with a request to another Cloud Service (outbound path). For the typical kubernetes based micro-service architecture, these two paths are usually handled by different micro services. We provide implementation patterns for both below.

We provide a go package documentation with many examples here:
- [**passport documentation**](https://github.wdf.sap.corp/pages/velocity/passport/docs/pkg/github.wdf.sap.corp/velocity/passport/v5/index.html)
- [**passport middleware documentation**](https://github.wdf.sap.corp/pages/velocity/passport/docs/pkg/github.wdf.sap.corp/velocity/passport/v5/middleware/index.html) (A [gorilla/mux
middleware](https://github.com/gorilla/mux))

## Inbound path

On the inbound path to a Cloud Service (e.g., Data Intelligence), requests arrive either with or without passports. For requests arriving without passport, a passport has to be created. For requests arriving with a passport, the passport has to first be logged and then updated with the information of the Cloud Service. In both cases, the passport has to be forwarded with the request through the Cloud Service.

### Using a middleware

We can use [gorilla/mux
middleware](https://github.com/gorilla/mux) to simplify the inbound process.

```go
// middleware/example_middleware_test.go

// (C) 2021 SAP SE or an SAP affiliate company. All rights reserved.
package middleware_test

import (
	"context"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.wdf.sap.corp/velocity/passport/v5"
	"github.wdf.sap.corp/velocity/passport/v5/middleware"
	"github.wdf.sap.corp/velocity/trc"
)

var tracer = trc.InitTraceTopic("example_tracer", "example description")

func Example_inboundMiddleware() {
	router := mux.NewRouter()
	router.Use(
		createInboundMiddleware(),
	)
	// router.HandleFunc("/", handler)
	server := httptest.NewServer(router)
	server.Close()
}

func createInboundMiddleware() mux.MiddlewareFunc {
	tenantId := uint64(0) // add tenant id here
	componentName := passport.MakeComponentName(passport.ComponentTypeDIP, tenantId)

	createPassport := func(_ context.Context) (*passport.Passport, error) {
		builder := passport.DummyPassportBuilder() // replace by real builder
		builder.SetComponentName(componentName)
		return builder.Create()
	}

	updatePassport := func(_ context.Context) (string, error) {
		return componentName, nil
	}

	return middleware.Chain(
		middleware.FromRequestHeader(),
		middleware.TrcLogEntry(tracer),
		middleware.UpdatePreviousComponentName(updatePassport),
		middleware.Create(createPassport),
		middleware.UpdateConnectionID(),
	)
}

```

You can see a full example is currently implemented in vsystem. You find the code [here](https://github.wdf.sap.corp/velocity/vsystem/blob/master/src/util/http/passport/main.go).

For more details on the individual middleware please look at the [**passport middleware documentation**](https://github.wdf.sap.corp/pages/velocity/passport/docs/pkg/github.wdf.sap.corp/velocity/passport/v5/middleware/index.html).
### Using the passport handler

Alternatively we may handle the passport through the handlers. This is a good option, if you don't receive requests via HTTP or do not want to use the [gorilla/mux
API routers](https://github.com/gorilla/mux).

```go
// example_inbound_test.go

// (C) 2021 SAP SE or an SAP affiliate company. All rights reserved.
package passport_test

import (
	"net/http"

	"github.wdf.sap.corp/velocity/passport/v5"
	"github.wdf.sap.corp/velocity/trc"
)

var inTracer = trc.InitTraceTopic("inbound_tracer", "description")

func Example_inboundPassport() {
	var pp *passport.Passport
	tenantId := uint64(0) // add tenant id here
	componentName := passport.MakeComponentName(passport.ComponentTypeDIP, tenantId)

	req := setupExampleRequest()

	pp, err := passport.FromHTTPHeader(&req.Header)
	if err != nil {
		inTracer.Debugf("cannot parse passport: %w")
		return
	}
	if pp != nil {
		subTracer := inTracer.Sub(passport.ToTrcInfos(pp)...)
		subTracer.Info("sap-passport")
		pp = updatePassport(pp, componentName)
	} else {
		pp = createPassport(componentName)
	}
}

func updatePassport(pp *passport.Passport, name string) *passport.Passport {
	pp, err := pp.WithPreviousComponentName(name)
	if err != nil {
		inTracer.Debugf("cannot set previous component name: %w", err)
		return nil
	}
	return pp
}

func createPassport(name string) *passport.Passport {
	builder := passport.NewBuilder()
	err := builder.SetComponentName(name)
	if err != nil {
		inTracer.Debugf("cannot set component name: %w", err)
		return nil
	}
	// You must set all other fields using the builder
	// We omit this here for brevity
	pp, err := builder.Create()
	if err != nil {
		inTracer.Debugf("cannot create passport: %w", err)
		return nil
	}
	return pp
}

// Only use this for testing
func setupExampleRequest() *http.Request {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	ppTest, _ := passport.DummyPassportBuilder().Create()
	passport.ToHTTPHeader(&req.Header, ppTest)
	return req
}

```

For details on the passport library, look at the [**passport documentation**](https://github.wdf.sap.corp/pages/velocity/passport/docs/pkg/github.wdf.sap.corp/velocity/passport/v5/index.html).

### Outbound path

On the outbound path of the Cloud Service, the passport has to be updated before it is forwarded to another Cloud Service. After the passport has been passed with the request, the passport needs to be logged with the status of the response.

```go
// example_outbound_test.go

// (C) 2021 SAP SE or an SAP affiliate company. All rights reserved.
package passport_test

import (
	"net/http"
	"os"

	"github.wdf.sap.corp/velocity/passport/v5"
	"github.wdf.sap.corp/velocity/trc"
)

var outTracer = trc.InitTraceTopic("outbound_tracer", "description")

func Example_outboundPassport() {
	trc.SetOutput(os.Stdout) // ignore this

	var pp *passport.Passport
	var err error

	setupExampleEnvironmentVariable()

	// collect the passport
	// e.g., by reading it from an environment variable
	ppAsHexString := os.Getenv("PASSPORT")
	if len(ppAsHexString) == 0 {
		outTracer.Debug("found no passport environment variable")
		return
	}

	pp, err = passport.FromHexString(ppAsHexString)
	if err != nil {
		outTracer.Debugf("cannot parse passport from environment variable: %w")
		return
	}

	// we advise to not trust excessively long variable parts by default
	// a length of 2000 characters should suffice for all scenarios
	if len(pp.VariableParts()) > 2000 {
		pp = pp.WithoutVariablePart()
	}

	client := &http.Client{}
	numRequests := 1

	// create a new connection id for every connection
	pp = pp.WithConnectionIDNew()

	// for every request
	for i := 0; i < numRequests; i++ {
		// increment the connection counter
		pp.IncrementConnectionCounter()

		// send request with passport
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		passport.ToHTTPHeader(&req.Header, pp)
		resp, _ := client.Do(req)

		// record the response status by using either the INFO or ERROR log level
		ppInfos := passport.ToTrcInfos(pp)
		subTracer := outTracer.Sub(ppInfos...)
		if resp.StatusCode == 200 { // replace by your success condition
			subTracer.Info("sap-passport")
		} else {
			subTracer.Error("sap-passport")
		}
	}
}

// Only use this for testing
func setupExampleEnvironmentVariable() {
	pp, _ := passport.DummyPassportBuilder().Create()
	hexString := passport.ToHexString(pp)
	os.Setenv("PASSPORT", hexString)
}

```

For details on the passport library, look at the [**passport documentation**](https://github.wdf.sap.corp/pages/velocity/passport/docs/pkg/github.wdf.sap.corp/velocity/passport/v5/index.html).

Code examples were generated directly from the files using the embedme tool: `npx embedme README.md`

## Passport command-line tool

We provide a small command-line tool to convert passports between formats.

### Installation

```bash
$> go install github.wdf.sap.corp/velocity/passport/v5/passportcli
```

### Usage

```
$> passportcli
Usage: passportcli [OPTIONS] <passport>
  -from string
        input format: json or hex (default "hex")
  -quiet
        no output is produced
  -to string
        output format: json or hex (default "json")
```

### Example

```bash
PASSPORT="2a54482a03010d890a5341505f4532455f54415f506c7567496e20202020202020202020202020202000005341505f4532455f54415f5573657220202020202020202020202020202020205341505f4532455f54415f526571756573742020202020202020202020202020202020202020202000055341505f4532455f54415f506c7567496e2020202020202020202020202020203334303238363037433939323145443539463937304530353641363943323832202020000734028607c9921ed59f970e056a69a2820000000000000000000000000000000000000000000100e22a54482a0100270100020003000200010400085800020002040008300002000302000b000000002a54482a"
passportcli $PASSPORT | jq
{
  "traceFlag": 2697,
  "componentName": "SAP_E2E_TA_PlugIn",
  "Service": 0,
  "UserID": "5341505f4532455f54415f55736572",
  "action": "SAP_E2E_TA_Request",
  "actionType": 5,
  "previousComponentName": "SAP_E2E_TA_PlugIn",
  "transactionID": "34028607-c992-1ed5-9f97-0e056a69c282",
  "clientNumber": "202020",
  "componentType": 7,
  "rootContextID": "34028607-c992-1ed5-9f97-0e056a69a282",
  "connectionID": "00000000-0000-0000-0000-000000000000",
  "connectionCounter": 0,
  "variablePartsNumber": 1,
  "variablePartsOffset": 226,
  "variableParts": "2a54482a0100270100020003000200010400085800020002040008300002000302000b00000000"
}
```


[//]: # ( (C) 2021 SAP SE or an SAP affiliate company. All rights reserved.)
