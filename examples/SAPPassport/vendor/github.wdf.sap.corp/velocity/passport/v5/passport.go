// (C) 2021 SAP SE or an SAP affiliate company. All rights reserved.
//
// Package passport implements the SAP passport.
package passport

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var (
	eyeCatcher = []byte{0x2A, 0x54, 0x48, 0x2A}
)

// Passport is a data structure for correlation used for request tracing and
// integration monitoring.
//
// The passport data structure consists of the following fields:
//
// 		| Start | Name                    | Type   | Size | Creation | Update |
// 		|-------|-------------------------|--------|------|----------|--------|
// 		| 000   | eye catcher             | []byte |  4   |          |        |
// 		| 004   | version                 | uint8  |  1   |          |        |
// 		| 005   | length                  | uint16 |  2   |          |        |
// 		| 007   | trace flag              | uint16 |  2   | X (opt.) |        |
// 		| 009   | component name          | []byte | 32   | X        |        |
// 		| 041   | service                 | uint16 |  2   |          |        |
// 		| 043   | user ID                 | []byte | 32   | X        |        |
// 		| 075   | action                  | []byte | 40   | X        |        |
// 		| 115   | action type             | uint16 |  2   | X        |        |
// 		| 117   | previous component name | []byte | 32   |          | X      |
// 		| 149   | transaction ID          | []byte | 32   | X        |        |
// 		| 181   | client number           | []byte |  3   |          |        |
// 		| 184   | component type          | uint16 |  2   | X        |        |
// 		| 186   | root context ID         | []byte | 16   | X        |        |
// 		| 202   | connection ID           | []byte | 16   |          | X      |
// 		| 218   | connection counter      | uint32 |  4   |          | X      |
// 		| 222   | variable parts number   | uint16 |  2   |          |        |
// 		| 224   | variable parts offset   | uint16 |  2   |          |        |
// 		| 226   | variable parts          | []byte |  ?   |          |        |
// 		| 226+? | eye catcher             | []byte |  4   |          |        |
//
// Most fields are set during creation of the passport. Only three are updated
// during the passport's lifetime.  The remaining fields are either generated
// automatically, or are currently not used by any passport use case (service &
// client number).
type Passport struct {
	version               uint8
	length                uint16
	traceFlag             uint16
	componentName         string
	service               uint16
	userID                string
	action                string
	actionType            uint16
	previousComponentName string
	transactionID         [16]byte
	clientNumber          []byte
	componentType         uint16
	rootContextID         [16]byte
	connectionID          [16]byte
	connectionCounter     uint32
	variablePartsNumber   uint16
	variablePartsOffset   uint16
	variableParts         []byte
}

func newPassport(traceFlag uint16, componentName string,
	service uint16, userID string, action string, actionType uint16,
	previousComponentName string, transactionID [16]byte, clientNumber []byte,
	componentType uint16, rootContextID [16]byte, connectionID [16]byte,
	connectionCounter uint32, variablePartsNumber uint16, variablePartsOffset uint16,
	variableParts []byte) *Passport {
	var version uint8 = 3
	var length uint16 = lenTotalExclVariablePart + uint16(len(variableParts)) // Do we need to include the variable parts offset?
	return &Passport{
		version, length, traceFlag, componentName,
		service, userID, action, actionType,
		previousComponentName, transactionID, clientNumber,
		componentType, rootContextID, connectionID,
		connectionCounter, variablePartsNumber,
		variablePartsOffset, variableParts,
	}
}

func (p *Passport) copy() *Passport {
	return &Passport{
		p.version, p.length, p.traceFlag, p.componentName,
		p.service, p.userID, p.action, p.actionType,
		p.previousComponentName, p.transactionID, p.clientNumber,
		p.componentType, p.rootContextID, p.connectionID,
		p.connectionCounter, p.variablePartsNumber,
		p.variablePartsOffset, p.variableParts,
	}
}

// Version of the passport.
func (p *Passport) Version() uint8 {
	return p.version
}

// Length of the complete passport data structure (incl. eye catchers).
func (p *Passport) Length() uint16 {
	return p.length
}

func (p *Passport) TraceFlags() uint16 {
	return p.traceFlag
}

// ComponentName returns the name of the component that created the passport.
func (p *Passport) ComponentName() string {
	return p.componentName
}

// Service is currently (2021-06) not used by any passport use case.
func (p *Passport) Service() uint16 {
	return p.service
}

func (p *Passport) UserID() string {
	return p.userID
}

func (p *Passport) Action() string {
	return p.action
}

func (p *Passport) ActionType() uint16 {
	return p.actionType
}

// PreviousComponentName refers to the component that forwarded the passport to the current component.
func (p *Passport) PreviousComponentName() string {
	return p.previousComponentName
}

// TransactionID refers to the technical transaction. Multiple TransactionIDs
// can exist for one RootContextID.
func (p *Passport) TransactionID() [16]byte {
	return p.transactionID
}

// TransactionIDString returns the transaction id as UUID string.
func (p *Passport) TransactionIDString() string {
	obj, _ := uuid.FromBytes(p.transactionID[:])
	return obj.String()
}

func (p *Passport) transactionIDHexString() string {
	return strings.ToUpper(hex.EncodeToString(p.transactionID[:]))
}

// ClientNumber is currently (2021-06) not used by any passport use case.
func (p *Passport) ClientNumber() []byte {
	return p.clientNumber
}

func (p *Passport) ComponentType() uint16 {
	return p.componentType
}

// RootContextID refers to the initial context within a complex scenario.
func (p *Passport) RootContextID() [16]byte {
	return p.rootContextID
}

// RootContextIDString returns the root context id as UUID string.
func (p *Passport) RootContextIDString() string {
	obj, _ := uuid.FromBytes(p.rootContextID[:])
	return obj.String()
}

// ConnectionID identifies an outgoing connection.
func (p *Passport) ConnectionID() [16]byte {
	return p.connectionID
}

// ConnectionIDString returns the connection id as UUID string.
func (p *Passport) ConnectionIDString() string {
	obj, _ := uuid.FromBytes(p.connectionID[:])
	return obj.String()
}

// ConnectionCounter together with ConnectionID uniquely identifies an outgoing
// request. (Starts at 1.)
func (p *Passport) ConnectionCounter() uint32 {
	return p.connectionCounter
}

// VariablePartsNumber specifies the number of elements contained in the
// VariableParts data structure.
func (p *Passport) VariablePartsNumber() uint16 {
	return p.variablePartsNumber
}

// VariablePartsOffset specifies the offset to elements contained in the
// VariableParts data structure.
func (p *Passport) VariablePartsOffset() uint16 {
	return p.variablePartsOffset
}

// VariableParts contains additional data that is piggy-backed on the passport.
// It's use is optional.
func (p *Passport) VariableParts() []byte {
	return p.variableParts
}

func (p *Passport) IsTraceLevelLow() bool {
	return p.traceFlag == traceLevelLow
}

func (p *Passport) IsTraceLevelMedium() bool {
	return p.traceFlag == traceLevelMedium
}

func (p *Passport) IsTraceLevelHigh() bool {
	return p.traceFlag == traceLevelHigh
}

func (p *Passport) IsTraceLevelSet() bool {
	return p.IsTraceLevelLow() || p.IsTraceLevelMedium() || p.IsTraceLevelHigh()
}

// Component layers are defined here
// https://sap.sharepoint.com/teams/IntelligentEnterpriseSuite-SupportFunctions/Shared%20Documents/04_Detect%20to%20correct/Integration%20Monitoring/Specifications/SAP_Passport.pdf#page=13
func (p *Passport) IsComponentLayerRuntime() bool {
	return p.componentType >= 1 && p.componentType <= 10
}

// Component layers are defined here
// https://sap.sharepoint.com/teams/IntelligentEnterpriseSuite-SupportFunctions/Shared%20Documents/04_Detect%20to%20correct/Integration%20Monitoring/Specifications/SAP_Passport.pdf#page=13
func (p *Passport) IsComponentLayerFramework() bool {
	return p.componentType >= 101 && p.componentType <= 1000
}

// Component layers are defined here
// https://sap.sharepoint.com/teams/IntelligentEnterpriseSuite-SupportFunctions/Shared%20Documents/04_Detect%20to%20correct/Integration%20Monitoring/Specifications/SAP_Passport.pdf#page=13
func (p *Passport) IsComponentLayerApplication() bool {
	return (p.componentType >= 21 && p.componentType <= 50) ||
		(p.componentType >= 1001 && p.componentType <= 32000)
}

// Component layers are defined here
// https://sap.sharepoint.com/teams/IntelligentEnterpriseSuite-SupportFunctions/Shared%20Documents/04_Detect%20to%20correct/Integration%20Monitoring/Specifications/SAP_Passport.pdf#page=13
func (p *Passport) IsComponentLayerUndefined() bool {
	return !p.IsComponentLayerRuntime() && !p.IsComponentLayerFramework() && !p.IsComponentLayerApplication()
}

func parseTraceFlags(traceFlag uint16) []string {
	var flags []string
	if traceFlag&TraceFlagAbapCondens1 > 0 {
		flags = append(flags, traceFlagStringAbapCondens1)
	}
	if traceFlag&TraceFlagAbapCondens2 > 0 {
		flags = append(flags, traceFlagStringAbapCondens2)
	}
	if traceFlag&TraceFlagDsrSat > 0 {
		flags = append(flags, traceFlagStringDsrSat)
	}
	if traceFlag&TraceFlagSql > 0 {
		flags = append(flags, traceFlagStringSql)
	}
	if traceFlag&TraceFlagBuffer > 0 {
		flags = append(flags, traceFlagStringBuffer)
	}
	if traceFlag&TraceFlagEnqueue > 0 {
		flags = append(flags, traceFlagStringEnqueue)
	}
	if traceFlag&TraceFlagRfc > 0 {
		flags = append(flags, traceFlagStringRfc)
	}
	if traceFlag&TraceFlagAuthTrace > 0 {
		flags = append(flags, traceFlagStringAuthTrace)
	}
	if traceFlag&TraceFlagCFunction > 0 {
		flags = append(flags, traceFlagStringCFunction)
	}
	if traceFlag&TraceFlagUserTrace > 0 {
		flags = append(flags, traceFlagStringUserTrace)
	}
	if traceFlag&TraceFlagDsrAbapTrace > 0 {
		flags = append(flags, traceFlagStringDsrAbapTrace)
	}
	return flags
}

// WithPreviousComponentName creates a copy of the passport with the given
// previousComponentName. previousComponentName must be at most 32 bytes long.
//
// Use this method to create a passport for all outbound connections.
func (p *Passport) WithPreviousComponentName(previousComponentName string) (*Passport, error) {
	if len(previousComponentName) > lenComponentName {
		return nil, fmt.Errorf("previousComponentName length exceeds %d bytes", lenComponentName)
	}
	pp2 := p.copy()
	pp2.previousComponentName = previousComponentName
	return pp2, nil
}

// WithConnectionIDNew creates a copy of the passport with a new generated
// connection ID.
//
// Use this method to create a passport for each outbound connections.
func (p *Passport) WithConnectionIDNew() *Passport {
	pp2 := p.copy()
	pp2.connectionID = uuid.New()
	pp2.connectionCounter = defaultConnectionCounter
	return pp2
}

// WithConnectionID creates a copy of the passport with the specified
// connection ID.
//
// Use this method to create a passport for each outbound connections.
func (p *Passport) WithConnectionID(connectionID [16]byte) (*Passport, error) {
	pp2 := p.copy()
	pp2.connectionID = connectionID
	pp2.connectionCounter = defaultConnectionCounter
	return pp2, nil
}

// IncrementConnectionCounter increments the connection counter of the passport.
//
// Use this method to create a passport for each outbound request on one
// connection.
func (p *Passport) IncrementConnectionCounter() uint32 {
	p.connectionCounter += 1
	return p.connectionCounter
}

// WithoutVariablePart removes a copy of the passport without the variable part. In some scenarios the
// propagation of the variable part can be a security risk.
func (p *Passport) WithoutVariablePart() *Passport {
	pp2 := p.copy()
	pp2.variableParts = []byte{}
	pp2.variablePartsNumber = 0
	pp2.variablePartsOffset = 0
	return pp2
}

func trimByteString(arraySlice []byte) []byte {
	endOfString := len(arraySlice)

	// location of zero character
	if idx := bytes.IndexByte(arraySlice, 0); idx >= 0 {
		endOfString = idx
	}

	// remove all trailing spaces
	for ; endOfString > 0 && arraySlice[endOfString-1] == 32; endOfString-- {
	}
	return arraySlice[:endOfString]
}

func convertSliceToArray(input []byte) (output [16]byte) {
	copy(output[:], input)
	return output
}

// Deserialize passport from a byte slice. The process fails if inconsistencies
// are found (length, eye catchers, ...).
func Deserialize(arr []byte) (*Passport, error) {
	if len(arr) < lenTotalExclVariablePart {
		return nil, fmt.Errorf("invalid passport: provided byte array is too small for a passport: %d", len(arr))
	}
	buf := bytes.NewBuffer(arr)
	if !bytes.Equal(eyeCatcher, buf.Next(lenEyeCatcher)) {
		return nil, fmt.Errorf("invalid passport: byte array did not contain eyeCatcher at beginning")
	}
	var version uint8 = buf.Next(1)[0]
	if version != 3 {
		return nil, fmt.Errorf("invalid passport: only passports of version 3 are supported")
	}
	builder := NewBuilder()
	// var length uint16 = binary.BigEndian.Uint16()
	var length uint16 = binary.BigEndian.Uint16(buf.Next(lenLength)) // only used for validation
	buf.Next(lenTraceFlags)
	builder.SetTraceFlags(traceLevelLow)
	//builder.SetTraceFlags(uint16(binary.LittleEndian.Uint16(buf.Next(lenTraceFlags))))
	err := builder.SetComponentName(string(trimByteString(buf.Next(lenComponentName))))
	if err != nil {
		return nil, err
	}
	builder.setService(binary.BigEndian.Uint16(buf.Next(lenService)))
	err = builder.SetUserID(string(trimByteString(buf.Next(lenUser))))
	if err != nil {
		return nil, err
	}
	err = builder.SetAction(string(trimByteString(buf.Next(lenAction))))
	if err != nil {
		return nil, err
	}
	builder.SetActionType(uint16(binary.BigEndian.Uint16(buf.Next(lenActionType))))
	err = builder.setPreviousComponentName(string(trimByteString(buf.Next(lenComponentName))))
	if err != nil {
		return nil, err
	}
	var transactionID []byte
	temp := buf.Next(lenTransactionID)
	transactionID, err = hex.DecodeString(string(temp))
	if err != nil {
		// if the transaction id is not a hex string, we take the first 16 bytes
		transactionID = temp[:16]
	}
	builder.SetTransactionID(convertSliceToArray(transactionID))
	err = builder.setClientNumber(buf.Next(lenClientNumber))
	if err != nil {
		return nil, err
	}
	builder.SetComponentType(uint16(binary.BigEndian.Uint16(buf.Next(lenComponentType))))
	builder.SetRootContextID(convertSliceToArray(buf.Next(lenRootContextID)))
	builder.setConnectionID(convertSliceToArray(buf.Next(lenConnectionID)))
	builder.setConnectionCounter(binary.BigEndian.Uint32(buf.Next(lenConnectionCounter)))
	builder.setVariablePartsNumber(binary.BigEndian.Uint16(buf.Next(lenVariablePartsNumber)))
	builder.setVariablePartsOffset(binary.BigEndian.Uint16(buf.Next(lenVariablePartsOffset)))
	builder.setVariableParts(buf.Next(len(arr) - lenTotalExclVariablePart))
	pp, err := builder.Create()
	if err != nil {
		return nil, fmt.Errorf("invalid passport: could not parse passport: %w", err)
	}
	if !bytes.Equal(eyeCatcher, buf.Next(lenEyeCatcher)) {
		return nil, fmt.Errorf("invalid passport: byte array did not contain eyeCatcher at end")
	}
	if int(length) != len(arr) {
		return nil, fmt.Errorf("invalid passport: passport length %d does not equal byte array length %d", int(length), len(arr))
	}
	return pp, nil
}

func copyIntoFixedLengthByteArray(buf *bytes.Buffer, field []byte, length int) {
	temp := make([]byte, length)
	copy(temp, field)
	_ = binary.Write(buf, binary.BigEndian, temp)
}

// Serialize a passport to a byte slice.
func (p *Passport) Serialize() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, p.Length()))
	_ = binary.Write(buf, binary.BigEndian, eyeCatcher)
	_ = binary.Write(buf, binary.BigEndian, p.Version())
	_ = binary.Write(buf, binary.BigEndian, p.Length())
	_ = binary.Write(buf, binary.LittleEndian, p.TraceFlags())
	copyIntoFixedLengthByteArray(buf, []byte(p.ComponentName()), lenComponentName)
	_ = binary.Write(buf, binary.BigEndian, p.Service())
	copyIntoFixedLengthByteArray(buf, []byte(p.UserID()), lenUser)
	copyIntoFixedLengthByteArray(buf, []byte(p.Action()), lenAction)
	_ = binary.Write(buf, binary.BigEndian, p.ActionType())
	copyIntoFixedLengthByteArray(buf, []byte(p.PreviousComponentName()), lenComponentName)
	// The transacationID was originally a 32-byte field.
	// Therefore, we encode it additionally as hex string.
	copyIntoFixedLengthByteArray(buf, []byte(p.transactionIDHexString()), lenTransactionID)
	copyIntoFixedLengthByteArray(buf, p.ClientNumber(), lenClientNumber)
	_ = binary.Write(buf, binary.BigEndian, p.ComponentType())
	_ = binary.Write(buf, binary.BigEndian, p.RootContextID())
	_ = binary.Write(buf, binary.BigEndian, p.ConnectionID())
	_ = binary.Write(buf, binary.BigEndian, p.ConnectionCounter())
	_ = binary.Write(buf, binary.BigEndian, p.VariablePartsNumber())
	_ = binary.Write(buf, binary.BigEndian, p.VariablePartsOffset())
	_ = binary.Write(buf, binary.BigEndian, p.VariableParts())
	_ = binary.Write(buf, binary.BigEndian, eyeCatcher)
	return buf.Bytes()
}

func parseTraceLevel(p *Passport) string {
	if p.IsTraceLevelLow() {
		return "low"
	}
	if p.IsTraceLevelMedium() {
		return "medium"
	}
	if p.IsTraceLevelHigh() {
		return "high"
	}
	return "not set"
}

// String returns a human-readable string. For a json object consider the
// LogString method.
func (p *Passport) String() string {
	var items []string
	sep := "    "
	items = append(items, sep+fmt.Sprintf("Version: %d", p.Version()))
	items = append(items, sep+fmt.Sprintf("Length: %d", p.Length()))
	items = append(items, sep+fmt.Sprintf("TraceFlag: 0x%04X", p.TraceFlags()))
	items = append(items, sep+sep+fmt.Sprintf("TraceLevel: %v", parseTraceLevel(p)))
	// This is necessary for the examples.
	if p.traceFlag > 0 {
		items = append(items, sep+sep+fmt.Sprintf("TraceTags: %s",
			strings.Join(parseTraceFlags(p.TraceFlags()), ",")))
	} else {
		items = append(items, sep+sep+"TraceTags:")
	}
	items = append(items, sep+fmt.Sprintf("ComponentName: %q", p.ComponentName()))
	items = append(items, sep+fmt.Sprintf("Service: %d", p.Service()))
	items = append(items, sep+fmt.Sprintf("User: %q", p.UserID()))
	items = append(items, sep+fmt.Sprintf("Action: %q", p.Action()))
	items = append(items, sep+fmt.Sprintf("ActionType: %d", p.ActionType()))
	items = append(items, sep+fmt.Sprintf("PreviousComponentName: %q", p.PreviousComponentName()))
	items = append(items, sep+fmt.Sprintf("TransactionID: %q", p.TransactionIDString()))
	items = append(items, sep+fmt.Sprintf("ClientNumber: %q", hex.EncodeToString(p.ClientNumber())))
	items = append(items, sep+fmt.Sprintf("ComponentType: %d", p.ComponentType()))
	items = append(items, sep+fmt.Sprintf("RootContextID: %q", p.RootContextIDString()))
	items = append(items, sep+fmt.Sprintf("ConnectionID: %q", p.ConnectionIDString()))
	items = append(items, sep+fmt.Sprintf("ConnectionCounter: %d", p.ConnectionCounter()))
	items = append(items, sep+fmt.Sprintf("number of variable parts: %d", p.VariablePartsNumber()))
	items = append(items, sep+fmt.Sprintf("offset of variable parts: %d", p.VariablePartsOffset()))
	return fmt.Sprintf("SAP-Passport\n%s", strings.Join(items, "\n"))
}
