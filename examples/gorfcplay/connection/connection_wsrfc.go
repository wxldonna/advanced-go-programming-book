// © 2020-2021 SAP SE or an SAP affiliate company. All rights reserved.
package connection

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.wdf.sap.corp/velocity/axino/src/core/axino"
	"github.wdf.sap.corp/velocity/axino/src/vsystem/cm"
	"github.wdf.sap.corp/velocity/gorfc/gorfc"
	"github.wdf.sap.corp/velocity/trc"
)

//nolint:gochecknoglobals
var wsrfcTracer = trc.InitTraceTopic("wsrfc", "wsrfc connection")

func NewWSRFCConnParameters(
	ctx context.Context, info *cm.Info, certs []cm.Certificate,
) (*gorfc.ConnectionParameter, error) {
	tracer := wsrfcTracer.SubFromContext(ctx)

	crypto, err := newSAPCrypto(ctx)
	if err != nil {
		return nil, err
	}

	total := len(certs)
	if total > 0 {
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
		Client: info.ContentData.Client,
		Lang:   info.ContentData.Language,
		Sysid:  info.ContentData.SysID,

		// NOTE(daniel) WebSocketRFC goes through UCON, which is configured such in cloud systems
		// that an Alias User must be used instead of a backend user. In on-premise systems, you
		// can configure UCON to allow using a backend user, but since we can't distinguish between
		// the systems, we require the Alias User is used there as well!
		Alias_User:                   info.ContentData.User,
		Passwd:                       info.ContentData.Password,
		Wshost:                       info.ContentData.WsHost,
		Wsport:                       fmt.Sprintf("%d", info.ContentData.WsPort),
		Use_TLS:                      "1",
		Tls_Client_PSE:               crypto.GetPSEFile(),
		Tls_Trust_All:                axino.GetTLSTrustAllFromEnv(),
		Trace:                        axino.GetRFCTraceLevel(),
		Tls_Client_Certificate_Logon: 1,

		// NOTE(daniel) These options aren't supported in cloud systems (due to missing
		// authorization for RFC_METADATA_GET) As we can't distinguish cloud and on-premise, we
		// need to disable these until the authorization issue has been fixed :(
		//		Use_Repository_Roundtrip_Optimization: "1",
		//		Serialization_Format:                  "CB_SERIALIZATION",
		//		Compression_Type:                      "LAN",
	}

	return param, nil
}

/*
func newWSRFCConnection(ctx context.Context, info *cm.Info, certs []cm.Certificate) (Connection, error) {
	c, err := newRFCConnection(info)
	if err != nil {
		return nil, err
	}

	rfc := c.(*abapRFCConnection)

	rfc.param, err = NewWSRFCConnParameters(ctx, info, certs)
	if err != nil {
		return nil, err
	}

	return rfc, nil
}
*/
//nolint:gochecknoglobals
var globalCrypto gorfc.SAPCrypto

func newSAPCrypto(ctx context.Context) (gorfc.SAPCrypto, error) {
	tracer := wsrfcTracer.SubFromContext(ctx)

	if globalCrypto != nil {
		tracer.Info("reusing existing crypto PSE")

		return globalCrypto, nil
	}

	var (
		src  *os.File
		err  error
		path string
	)

	if os.Getenv("AXINO_HOME") != "" {
		path = os.ExpandEnv("${AXINO_HOME}/root.pse")

		src, err = os.Open(path)
		if err != nil {
			src = nil
		}
	}

	// NOTE(daniel) if either the environment variable doesn't exist, or we can't
	// open the root.pse file, fallback to generating an empty PSE file.
	if src == nil {
		tracer.Info("generating empty crypto PSE")

		// TODO(daniel) should we generate the DName based on the user/tenant?
		globalCrypto, err = gorfc.NewSAPCrypto(
			gorfc.WithDName("CN=axino,OU=datahub,O=SAP,C=DE"),
			gorfc.WithTemporaryPSE,
		)
		if err != nil {
			return nil, err
		}

		return globalCrypto, nil
	}

	tracer.Infof("using crypto PSE from %q", path)

	// NOTE(daniel) copy pse file containing all root certificates to a new
	// temporary file, and initialize SAPCrypto to use this.
	defer src.Close()

	tgt, err := ioutil.TempFile("", "my.*.pse")
	if err != nil {
		return nil, err
	}
	defer tgt.Close()

	_, err = io.Copy(tgt, src)
	if err != nil {
		return nil, err
	}

	globalCrypto, err = gorfc.NewSAPCrypto(
		gorfc.WithPSE(tgt.Name()),
	)
	if err != nil {
		return nil, err
	}

	return globalCrypto, nil
}

// © 2020-2021 SAP SE or an SAP affiliate company. All rights reserved.
