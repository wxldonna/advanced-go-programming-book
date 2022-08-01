// (C) 2021 SAP SE or an SAP affiliate company. All rights reserved.

package passport

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.wdf.sap.corp/velocity/trc"
)

type passportKey int

const (
	contextKey passportKey = iota
	HeaderKey  string      = "SAP-PASSPORT"
)

// FromHexString deserializes the passport from a hex string.
func FromHexString(hexString string) (*Passport, error) {
	byteString, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	return Deserialize(byteString)
}

// ToHexString serializes the passport to a hex string.
func ToHexString(pp *Passport) string {
	return strings.ToUpper(hex.EncodeToString(pp.Serialize()))
}

// FromHTTPHeader extracts the passport from the header, deserializes it an returns it.
func FromHTTPHeader(header *http.Header) (*Passport, error) {
	ppHeader := header.Get(HeaderKey)
	if len(ppHeader) == 0 {
		return nil, nil
	}
	return FromHexString(ppHeader)
}

// ToHTTPHeader adds the passport to the http header.
func ToHTTPHeader(header *http.Header, pp *Passport) {
	ser := pp.Serialize()
	byteString := hex.EncodeToString(ser)
	header.Set(HeaderKey, byteString)
}

// DeleteFromHTTPHeader deletes passport from http header if it exists.
func DeleteFromHTTPHeader(header *http.Header) {
	header.Del(HeaderKey)
}

// FromHTTPRequest extracts the passport from the request, deserializes it an returns it.
func FromHTTPRequest(req *http.Request) (*Passport, error) {
	// consuming error here, to be backwards compatible
	pp, err := FromHTTPHeader(&req.Header)
	return pp, err
}

// ToHTTPRequest adds the passport to the http request's header.
func ToHTTPRequest(req *http.Request, pp *Passport) {
	ToHTTPHeader(&req.Header, pp)
}

// DeleteFromHTTPRequest deletes passport from http request header if it exists.
func DeleteFromHTTPRequest(req *http.Request) {
	DeleteFromHTTPHeader(&req.Header)
}

// ToContext creates a new context with a passport pointer included.
func ToContext(ctx context.Context, pp *Passport) context.Context {
	if pp != nil {
		return context.WithValue(ctx, contextKey, *pp)
	}
	return ctx
}

// FromContext extracts the passport from the context.
func FromContext(ctx context.Context) *Passport {
	pp := ctx.Value(contextKey)
	if pp != nil {
		val, ok := pp.(Passport)
		if ok {
			return &val
		}
	}
	return nil
}

// Converts passport into trc info fields
func ToTrcInfos(pp *Passport) []trc.Info {
	var infos []trc.Info
	transactionIDString := pp.TransactionIDString()
	rootContextIDString := pp.RootContextIDString()
	connectionIDString := pp.ConnectionIDString()
	infos = append(infos, trc.Info{K: "sap_passport_transaction_id", V: transactionIDString})
	infos = append(infos, trc.Info{K: "sap_passport_previous_component_name", V: pp.PreviousComponentName()})
	infos = append(infos, trc.Info{K: "sap_passport_component_name", V: pp.ComponentName()})
	infos = append(infos, trc.Info{K: "sap_passport_root_context_id", V: rootContextIDString})
	infos = append(infos, trc.Info{K: "sap_passport_connection_id", V: connectionIDString})
	infos = append(infos, trc.Info{K: "sap_passport_connection_counter", V: pp.ConnectionCounter()})
	return infos
}

// AttachTrcInfosToContext converts the passport to trc info objects and adds it to the context.
func AttachTrcInfosToContext(ctx context.Context, pp *Passport) context.Context {
	infos := ToTrcInfos(pp)
	return trc.AttachTrcInfo(ctx, infos...)
}

// AttachTrcInfosToContext converts the passport to trc info objects and adds it to the request context.
func AttachTrcInfosToRequest(r *http.Request, pp *Passport) *http.Request {
	return r.WithContext(AttachTrcInfosToContext(r.Context(), pp))
}

// DeleteFromContext removes the passport from the context, but not the key.
func DeleteFromContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, nil)
}

type passportInternal struct {
	TraceFlag             uint16 `json:"traceFlag"`
	ComponentName         string `json:"componentName"`
	Service               uint16 `json:"Service"`
	UserID                string `json:"UserID"`
	Action                string `json:"action"`
	ActionType            uint16 `json:"actionType"`
	PreviousComponentName string `json:"previousComponentName"`
	TransactionID         string `json:"transactionID"`
	ClientNumber          string `json:"clientNumber"`
	ComponentType         uint16 `json:"componentType"`
	RootContextID         string `json:"rootContextID"`
	ConnectionID          string `json:"connectionID"`
	ConnectionCounter     uint32 `json:"connectionCounter"`
	VariablePartsNumber   uint16 `json:"variablePartsNumber"`
	VariablePartsOffset   uint16 `json:"variablePartsOffset"`
	VariableParts         string `json:"variableParts"`
}

// ToJSON converts a passport to JSON.
func ToJSON(pp *Passport) ([]byte, error) {
	return json.Marshal(passportInternal{
		TraceFlag:             pp.traceFlag,
		ComponentName:         pp.componentName,
		Service:               pp.service,
		UserID:                pp.UserID(),
		Action:                pp.action,
		ActionType:            pp.actionType,
		PreviousComponentName: pp.previousComponentName,
		TransactionID:         pp.TransactionIDString(),
		ClientNumber:          hex.EncodeToString(pp.ClientNumber()),
		ComponentType:         pp.componentType,
		RootContextID:         pp.RootContextIDString(),
		ConnectionID:          pp.ConnectionIDString(),
		ConnectionCounter:     pp.connectionCounter,
		VariablePartsNumber:   pp.variablePartsNumber,
		VariablePartsOffset:   pp.VariablePartsOffset(),
		VariableParts:         hex.EncodeToString(pp.VariableParts()),
	})
}

// FromJSON converts a json byte string to a passport
func FromJSON(jsonString []byte) (*Passport, error) {
	var ppInternal passportInternal
	err := json.Unmarshal(jsonString, &ppInternal)
	if err != nil {
		return nil, err
	}
	builder := NewBuilder()
	builder.SetTraceFlags(ppInternal.TraceFlag)
	err = builder.SetComponentName(ppInternal.ComponentName)
	if err != nil {
		return nil, err
	}
	// service is currently not used by anybody
	builder.setService(ppInternal.Service)
	err = builder.SetUserID(ppInternal.UserID)
	if err != nil {
		return nil, err
	}
	err = builder.SetAction(ppInternal.Action)
	if err != nil {
		return nil, err
	}
	builder.SetActionType(ppInternal.ActionType)
	err = builder.setPreviousComponentName(ppInternal.PreviousComponentName)
	if err != nil {
		return nil, err
	}
	err = builder.SetTransactionIDString(ppInternal.TransactionID)
	if err != nil {
		return nil, err
	}
	clientNumber, err := hex.DecodeString(ppInternal.ClientNumber)
	if err != nil {
		return nil, err
	}
	err = builder.setClientNumber(clientNumber)
	if err != nil {
		return nil, err
	}
	builder.SetComponentType(ppInternal.ComponentType)
	err = builder.SetRootContextIDString(ppInternal.RootContextID)
	if err != nil {
		return nil, err
	}
	err = builder.setConnectionIDString(ppInternal.ConnectionID)
	if err != nil {
		return nil, err
	}
	builder.setConnectionCounter(ppInternal.ConnectionCounter)
	builder.setVariablePartsNumber(ppInternal.VariablePartsNumber)
	builder.setVariablePartsOffset(ppInternal.VariablePartsOffset)
	variableParts, err := hex.DecodeString(ppInternal.VariableParts)
	if err != nil {
		return nil, err
	}
	builder.setVariableParts(variableParts)
	pp, err := builder.Create()
	return pp, err
}
