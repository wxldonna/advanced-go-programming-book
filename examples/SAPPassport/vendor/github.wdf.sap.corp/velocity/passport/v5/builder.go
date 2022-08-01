// (C) 2021 SAP SE or an SAP affiliate company. All rights reserved.

package passport

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Builder struct {
	passport           *Passport
	componentNameIsSet bool
	userIsSet          bool
	actionIsSet        bool
	actionTypeIsSet    bool
	transactionIdIsSet bool
	componentTypeIsSet bool
	rootContextIdIsSet bool
}

// NewBuilder creates a new builder for the passport.
func NewBuilder() *Builder {
	b := &Builder{
		passport: &Passport{},
	}
	b.setDefaultValues()
	return b
}

var (
	defaultClientNumber  = []byte{0, 0, 0}
	defaultConnectionID  = [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	defaultVariableParts = []byte{}
)

func (b *Builder) setDefaultValues() {
	// previous component name is not set when creating a new passport
	b.passport.previousComponentName = defaultPreviousComponentName

	// service and client number are not used by any use case in the specs
	b.passport.service = defaultService
	b.passport.clientNumber = defaultClientNumber

	b.passport.traceFlag = traceLevelLow
	b.passport.connectionID = defaultConnectionID
	b.passport.connectionCounter = uint32(defaultConnectionCounter)
	b.passport.variablePartsNumber = defaultVariablePartsNumber
	b.passport.variablePartsOffset = defaultVariablePartsOffset
	b.passport.variableParts = defaultVariableParts
}

// GenerateTransactionID sets the transactionID field to a generated UUID.
//
// The UUID uses only 16 byte of the 32 byte long field.
func (b *Builder) GenerateTransactionID() {
	b.passport.transactionID = uuid.New()
	b.transactionIdIsSet = true
}

// GenerateRootContextID sets the rootContextID field to a generated UUID.
func (b *Builder) GenerateRootContextID() {
	b.passport.rootContextID = uuid.New()
	b.rootContextIdIsSet = true
}

// SetTraceFlags sets multiple trace flags.
//
// You can only set either trace flags or trace levels. Setting a traceFlag
// overrides previously set trace levels.
func (b *Builder) SetTraceFlags(traceFlags ...uint16) {
	var traceFlag uint16 = 0
	for _, tf := range traceFlags {
		traceFlag = traceFlag | tf
	}
	b.passport.traceFlag = traceFlag
}

// SetTraceLevelLow sets the trace level to low.
//
// You can only set either trace flags or trace levels. Setting a trace
// level overrides previously set trace flag.
func (b *Builder) SetTraceLevelLow() {
	b.passport.traceFlag = traceLevelLow
}

// SetTraceLevelMedium sets the trace level to medium.
//
// You can only set either trace flags or trace levels. Setting a trace
// level overrides previously set trace flag.
func (b *Builder) SetTraceLevelMedium() {
	b.passport.traceFlag = traceLevelMedium
}

// SetTraceLevelHigh sets the trace level to high.
//
// You can only set either trace flags or trace levels. Setting a trace level
// overrides previously set trace flag.
func (b *Builder) SetTraceLevelHigh() {
	b.passport.traceFlag = traceLevelHigh
}

// SetComponentName sets the component name.
//
// The function MakeComponentName can be used to generate the byte string.
func (b *Builder) SetComponentName(componentName string) error {
	if err := guardStringLength(componentName, "componentName", lenComponentName); err != nil {
		return err
	}
	b.passport.componentName = componentName
	b.componentNameIsSet = true
	return nil
}

// setService sets the service field. This field is not used by current
// (2021-06) passport use cases.
func (b *Builder) setService(service uint16) {
	b.passport.service = service
}

// SetUserID takes string with a length of 32 or fewer characters.
func (b *Builder) SetUserID(userID string) error {
	if err := guardStringLength(userID, "user", lenUser); err != nil {
		return err
	}
	b.passport.userID = userID
	b.userIsSet = true
	return nil
}

// SetAction takes string with a length of 32 or fewer characters.
func (b *Builder) SetAction(action string) error {
	if err := guardStringLength(action, "action", lenAction); err != nil {
		return err
	}
	b.passport.action = action
	b.actionIsSet = true
	return nil
}

// SetActionType sets the type of the action. You can use the predefined action
// constants to set the action type.
func (b *Builder) SetActionType(actionType uint16) {
	b.passport.actionType = actionType
	b.actionTypeIsSet = true
}

// SetPreviousComponentName takes string with a length of 16 or fewer characters.
func (b *Builder) setPreviousComponentName(previousComponentName string) error {
	if err := guardStringLength(previousComponentName, "previousComponentName", lenComponentName); err != nil {
		return err
	}
	b.passport.previousComponentName = previousComponentName
	return nil
}

// SetTransactionID sets the transaction id. This
// field can also be generated.
func (b *Builder) SetTransactionID(transactionID [16]byte) {
	b.passport.transactionID = transactionID
	b.transactionIdIsSet = true
}

// SetTransactionIDString is a convenience function to set transaction ids from a uuid string.
func (b *Builder) SetTransactionIDString(transactionID string) error {
	obj, err := uuid.Parse(transactionID)
	if err != nil {
		return fmt.Errorf("cannot parse root context id %w", err)
	}
	b.SetTransactionID(obj)
	return nil
}

// SetClientNumber takes byte slices with a length of 3 or fewer bytes.
// ClientNumber is by default set to []byte{0, 0, 0}. It is currently (2021-06)
// not used by any passport use case.
func (b *Builder) setClientNumber(clientNumber []byte) error {
	if err := guardByteLength(clientNumber, "clientNumber", lenClientNumber); err != nil {
		return err
	}
	field := make([]byte, lenClientNumber)
	copy(field, clientNumber)
	b.passport.clientNumber = clientNumber
	return nil
}

// SetComponentType sets the type of the component. You can use the predefined
// action constants to set the component type. A complete list can be found
// here:
// https://github.wdf.sap.corp/xdsr/epp-component-type/blob/master/sap.crun.nameservicedemo-CloudService.csv
func (b *Builder) SetComponentType(componentType uint16) {
	b.passport.componentType = componentType
	b.componentTypeIsSet = true
}

// SetRootContextID sets the root context id. This field can also be generated.
func (b *Builder) SetRootContextID(rootContextID [16]byte) {
	b.passport.rootContextID = rootContextID
	b.rootContextIdIsSet = true
}

// setRootContextIDString is a convenience function to set transaction ids from a UUID string.
func (b *Builder) SetRootContextIDString(rootContextID string) error {
	obj, err := uuid.Parse(rootContextID)
	if err != nil {
		return fmt.Errorf("cannot parse root context id %w", err)
	}
	b.SetRootContextID(obj)
	return nil
}

// setConnectionID sets the connection id. This field can also be generated.
func (b *Builder) setConnectionID(connID [16]byte) {
	b.passport.connectionID = connID
}

// setConnectionIDString is a convenience function to set transaction ids from a UUID string.
func (b *Builder) setConnectionIDString(connID string) error {
	obj, err := uuid.Parse(connID)
	if err != nil {
		return fmt.Errorf("cannot parse connection id %w", err)
	}
	b.setConnectionID(obj)
	return nil
}

func (b *Builder) setConnectionCounter(connCounter uint32) {
	b.passport.connectionCounter = connCounter
}

// setVariablePartsNumber sets the number of variable parts that you want to include.
//
// This should not be part of the builder but should be a manipulation of the passport itself.
func (b *Builder) setVariablePartsNumber(number uint16) {
	b.passport.variablePartsNumber = number
}

// setVariablePartsOffset sets an offset for the variable parts.
//
// This should not be part of the builder but should be a manipulation of the passport itself.
func (b *Builder) setVariablePartsOffset(offset uint16) {
	b.passport.variablePartsOffset = offset
}

// setVariablePartsOffset sets a byte slice containing the variable parts.
//
// This should not be part of the builder but should be a manipulation of the passport itself.
func (b *Builder) setVariableParts(varParts []byte) {
	b.passport.variableParts = varParts
}

func (b *Builder) collectMissingFields() []string {
	var fieldsNotSet []string
	if !b.componentNameIsSet {
		fieldsNotSet = append(fieldsNotSet, "componentName")
	}
	if !b.userIsSet {
		fieldsNotSet = append(fieldsNotSet, "user")
	}
	if !b.actionIsSet {
		fieldsNotSet = append(fieldsNotSet, "action")
	}
	if !b.actionTypeIsSet {
		fieldsNotSet = append(fieldsNotSet, "actionType")
	}
	if !b.transactionIdIsSet {
		fieldsNotSet = append(fieldsNotSet, "transactionID")
	}
	if !b.componentTypeIsSet {
		fieldsNotSet = append(fieldsNotSet, "componentType")
	}
	if !b.rootContextIdIsSet {
		fieldsNotSet = append(fieldsNotSet, "rootContextID")
	}
	return fieldsNotSet
}

// Create builds a new passport, unless a mandatory field was not set.
func (b *Builder) Create() (*Passport, error) {
	missingFields := b.collectMissingFields()
	if len(missingFields) > 0 {
		return nil, fmt.Errorf("%d fields were not set: %s", len(missingFields), strings.Join(missingFields, ","))
	}
	return newPassport(
		b.passport.traceFlag, b.passport.componentName, b.passport.service,
		b.passport.userID, b.passport.action, b.passport.actionType,
		b.passport.previousComponentName, b.passport.transactionID,
		b.passport.clientNumber, b.passport.componentType,
		b.passport.rootContextID, b.passport.connectionID,
		b.passport.connectionCounter, b.passport.variablePartsNumber,
		b.passport.variablePartsOffset, b.passport.variableParts,
	), nil
}

// Copy creates a copy of a partially initialized builder.
func (b *Builder) Copy() *Builder {
	return &Builder{
		passport:           b.passport.copy(),
		componentNameIsSet: b.componentNameIsSet,
		userIsSet:          b.userIsSet,
		actionIsSet:        b.actionIsSet,
		actionTypeIsSet:    b.actionTypeIsSet,
		transactionIdIsSet: b.transactionIdIsSet,
		componentTypeIsSet: b.componentTypeIsSet,
		rootContextIdIsSet: b.rootContextIdIsSet,
	}
}

// DummyPassportBuilder returns a dummy passport builder for testing
func DummyPassportBuilder() *Builder {
	// TestDummyPassportBuilder checks this
	builder := NewBuilder()
	builder.SetTraceLevelMedium()
	_ = builder.SetComponentName("SomeComponent")
	_ = builder.SetUserID("SomeUser")
	_ = builder.SetAction("Request")
	builder.SetActionType(ActionTypeHTTP)
	builder.SetComponentType(ComponentTypeUndefined)
	_ = builder.SetTransactionIDString("9bb0f2d2-b5d8-4d58-9bf9-908cd40cdd74")
	_ = builder.SetRootContextIDString("5c66c21c-7282-4a9e-b146-9bcbfe9ffeb7")
	_ = builder.setConnectionIDString("d5c0b704-aea9-4c09-b227-855b7a3431ee")
	return builder
}

func guardByteLength(something []byte, name string, expectedLength int) error {
	if len(something) > expectedLength {
		return fmt.Errorf("%s length exceeds %d byte", name, expectedLength)
	}
	return nil
}

func guardStringLength(something string, name string, expectedLength int) error {
	if len(something) > expectedLength {
		return fmt.Errorf("%s length exceeds %d byte", name, expectedLength)
	}
	return nil
}
