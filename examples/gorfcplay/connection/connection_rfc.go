package connection

import (
	"context"
	"fmt"
	"strconv"

	"github.wdf.sap.corp/velocity/axino/src/vsystem/cm"
	"github.wdf.sap.corp/velocity/gorfc/gorfc"
	"github.wdf.sap.corp/velocity/trc"
)

//export GONOSUMDB=*
//nolint:gochecknoglobals
/*
var rfcTracer = trc.InitTraceTopic("rfc", "rfc connection")

type rfcEndpoint struct {
	name           string
	graphManager   string
	graphVersion   string
	graphRoundtrip string
	functionExists string
}

type rfcParams map[string]interface{}

// NOTE(daniel) These are all the functions we need from the GoRFC library.
type GoRFCConnection interface {
	GetFunctionDescription(string) (gorfc.FunctionDescription, error)
	RemoveFunctionDescription(string) error
	Call(string, interface{}) (map[string]interface{}, error)
	Close() error
}

type GoRFCFacade interface {
	ConnectionFromParams(gorfc.ConnectionParameter) (GoRFCConnection, error)
}

type goRFCFacade struct{}

func (g *goRFCFacade) ConnectionFromParams(param gorfc.ConnectionParameter) (GoRFCConnection, error) {
	return gorfc.ConnectionFromParams(param)
}

func newGoRFCFacade() GoRFCFacade {
	return &goRFCFacade{}
}

type abapRFCConnection struct {
	conn  GoRFCConnection
	gorfc GoRFCFacade

	// NOTE(daniel) The NWRFCLIB doesn't support concurrent calls on the same connection, so all
	// gorfc-calls need to be protected by this mutex!
	mutex *sync.Mutex

	endpoint rfcEndpoint
	info     *cm.Info
	param    *gorfc.ConnectionParameter
}
*/
var rfcTracer = trc.InitTraceTopic("rfc", "rfc connection")

func NewRFCConnParameters(
	info *cm.Info, certs []cm.Certificate,
) (*gorfc.ConnectionParameter, error) {

	total := len(certs)
	if total > 0 {
		ctx := context.TODO()
		tracer := wsrfcTracer.SubFromContext(ctx)
		crypto, err := newSAPCrypto(ctx)
		if err != nil {
			return nil, err
		}
		tracer.Infof("loading %d user certificates", total)

		for i, cert := range certs {
			if err := crypto.TrustServer(cert.CData); err != nil {
				tracer.Errorf("- (%d/%d) failed importing certificate '%s': %v", i+1, total, cert.Filename, err)
			} else {
				tracer.Debugf("- (%d/%d) imported certificate '%s'", i+1, total, cert.Filename)
			}
		}
	}

	param := &gorfc.ConnectionParameter{
		Ashost:    info.ContentData.AsHost,
		Mshost:    info.ContentData.MsHost,
		Msserv:    strconv.Itoa(info.ContentData.MsServ),
		Group:     info.ContentData.Group,
		Client:    info.ContentData.Client,
		Lang:      info.ContentData.Language,
		Sysid:     info.ContentData.SysID,
		Sysnr:     info.ContentData.SysNr,
		Gwhost:    info.ContentData.GwHost,
		Gwserv:    info.ContentData.GwServ,
		User:      info.ContentData.User,
		Passwd:    info.ContentData.Password,
		Saprouter: info.ContentData.Router,
		//	Tls_Client_PSE:               crypto.GetPSEFile(),
		//	Tls_Trust_All:                axino.GetTLSTrustAllFromEnv(),
		//	Tls_Client_Certificate_Logon: 1,
		//getsso2: 1,

		//	Trace:     axino.GetRFCTraceLevel(),
	}

	if info.Gateway.Subaccount != "" {
		param.Connectivity_Proxy_Host = info.Gateway.Host
		param.Connectivity_Proxy_Port = fmt.Sprintf("%d", info.Gateway.Port)
		param.Connectivity_Subaccount = info.Gateway.Subaccount
		param.Connectivity_Location_Id = info.Gateway.LocationID
		param.Connectivity_Proxy_Authentication = info.Gateway.Authentication
	}

	return param, nil
}

/*
func newRFCConnection(info *cm.Info) (Connection, error) {
	param, err := NewRFCConnParameters(info, nil)
	if err != nil {
		return nil, err
	}

	// NOTE(daniel) Repository Roundtrip Optimization requires a recent kernel. Since this doesn't
	// add much benefit for our very simple function modules anyway, turn it off.
	param.Use_Repository_Roundtrip_Optimization = "0"

	// NOTE(daniel) CB_SERIALIZATION in combination with the 722 ABAP kernel triggers a bug that
	// causes data to be skipped in the RFC libraries. As we anyway do not benefit from this due
	// to the structure of our RFC parameters, disable this.

	// param.Serialization_Format = "CB_SERIALIZATION"
	param.Compression_Type = "LAN"

	rfc := &abapRFCConnection{
		gorfc: newGoRFCFacade(),
		info:  info,
		mutex: &sync.Mutex{},
		param: param,
	}

	return rfc, nil
}

func (c *abapRFCConnection) GetSystemInfo(ctx context.Context) (*SystemInfoResponse, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	endpoint, err := c.getEndpoint(ctx)
	if err != nil {
		return nil, err
	}

	// NOTE(daniel) If the importing parameter `IV_VERSION` exists in the
	// ABAP system, fill it with the DH Version. Otherwise skip it.
	funDes, err := c.conn.GetFunctionDescription(endpoint.graphVersion)
	if err != nil {
		return nil, err
	}

	param := rfcParams{}

	for _, v := range funDes.Parameters {
		if v.Name == "IV_VERSION" {
			param[v.Name] = "v1" //axino.GetDHversion()

			break
		}
	}

	r, err := c.conn.Call(endpoint.graphVersion, param)
	if err != nil {
		return nil, fmt.Errorf("RFC call failed: %w", err)
	}

	res := SystemInfoResponse{
		Engine:   r["EV_VERSION"].(string),
		Messages: mapMsgsFromAbap(r["ET_MSG"].([]interface{})),
	}

	return &res, nil
}

// TODO(daniel) do we have anything to cleanup? Should we make a last ABAP call to gracefully cleanup the
// ABAP session?
func (c *abapRFCConnection) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil {
		conn := c.conn
		c.conn = nil

		return conn.Close()
	}

	return nil
}

func (c *abapRFCConnection) getEndpoint(ctx context.Context) (rfcEndpoint, error) {
	if _, err := c.getConnection(ctx); err != nil {
		return rfcEndpoint{}, fmt.Errorf("ape[RFC] open connection: %w", err)
	}

	if c.endpoint == (rfcEndpoint{}) {
		c.endpoint = c.probe(ctx)
	}

	if c.endpoint == (rfcEndpoint{}) {
		return rfcEndpoint{}, ErrAPENotDetected
	}

	return c.endpoint, nil
}

//nolint:unparam
func (c *abapRFCConnection) getConnection(ctx context.Context) (GoRFCConnection, error) {
	tracer := rfcTracer.SubFromContext(ctx)

	if c.conn != nil {
		tracer.Debug("ape[RFC] reuse connection")

		return c.conn, nil
	}

	for c.conn == nil {
		// NOTE(daniel) this opens the RFC connection
		conn, err := c.gorfc.ConnectionFromParams(*c.param)
		if err != nil {
			tracer.Debugf("ape[RFC] open connection: %v", err)

			if c.param.Serialization_Format != "" {
				tracer.Info("ape[RFC] open connection: switch off cb_serialization")

				c.param.Serialization_Format = ""

				continue
			}

			if c.param.Use_Repository_Roundtrip_Optimization != "0" {
				tracer.Info("ape[RFC] open connection: switch off roundtrip optimization")

				c.param.Use_Repository_Roundtrip_Optimization = "0"

				continue
			}

			//nolint:stylecheck,golint
			return nil, fmt.Errorf("%w. Please refer to SAP Note 2849542 for more information.", err)
		}

		c.conn = conn
	}

	return c.conn, nil
}
*/
