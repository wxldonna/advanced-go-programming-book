package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jarcoal/httpmock"
)

func main() {

	httpmock.Activate()
	httpmock.Deactivate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "=~http://api.*", httpmock.NewStringResponder(200, `{"name":"sap.wangxiaoliang"}`))

	req, err := http.NewRequestWithContext(context.TODO(), "GET", "http://api.abc.com/articles1", nil)
	if err != nil {
		log.Printf("make request failed with error %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	log.Printf("response : %+v", string(data))

}
