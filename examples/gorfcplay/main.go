package main

import (
	"fmt"
	"time"

	"gorfcplay/connection"

	"github.wdf.sap.corp/velocity/axino/src/vsystem/cm"
	"github.wdf.sap.corp/velocity/gorfc/gorfc"
)

func abapSystem() *gorfc.ConnectionParameter {

	info := &cm.Info{
		ContentData: cm.Data{
			AsHost:   "ldcists.devint.net.sap",
			Client:   "800",
			GwHost:   "ldcists.devint.net.sap",
			Protocol: "RFC",
			User:     "WANGXIAO24",
			Password: "Wxldonna112358",
			SysNr:    "11",
			SysID:    "STS",
		},
	}

	/*
		wsinfo := &cm.Info{
			ContentData: cm.Data{
				WsHost:   "cc3-715-api.wdf.sap.corp",
				WsPort:   443,
				Client:   "715",
				Protocol: "WebSocket RFC",
				User:     "XL_CC3_715",
				Password: "Wxl@112358",
				SysNr:    "00",
				SysID:    "CC3",
			},
		}
	*/
	//ctx := context.TODO()
	//certs, _ := cm.GetCertificates(ctx)

	param, _ := connection.NewRFCConnParameters(info, nil)
	return param
	//ctx := context.TODO()
	//certs, _ := cm.GetCertificates(ctx)
	//wsrfcParamters, _ := connection.NewWSRFCConnParameters(ctx, wsinfo, certs)
	//fmt.Printf("wsrfc parameters is %v \n", wsrfcParamters)
	return param
}

func main() {
	c, _ := gorfc.ConnectionFromParams(*abapSystem())
	params := map[string]interface{}{
		"IMPORTSTRUCT": map[string]interface{}{
			"RFCFLOAT": 1.23456789,
			"RFCCHAR1": "A",
			"RFCCHAR2": "BC",
			"RFCCHAR4": "ÄBC",
			"RFCINT1":  0xfe,
			"RFCINT2":  0x7ffe,
			"RFCINT4":  999999999,
			"RFCHEX3":  []byte{255, 254, 253},
			"RFCTIME":  time.Now(),
			"RFCDATE":  time.Now(),
			"RFCDATA1": "HELLÖ SÄP",
			"RFCDATA2": "DATA222",
		},
	}
	r, _ := c.Call("STFC_STRUCTURE", params)
	//	ticketid := c.RfcGetPartnerSSOTicket()
	fmt.Printf("result is %v  \n", r["ECHOSTRUCT"])
	//fmt.Println("ticketid is %v", ticketid)
	c.Close()
}
