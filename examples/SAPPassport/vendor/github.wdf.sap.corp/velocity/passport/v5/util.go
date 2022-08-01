// (C) 2021 SAP SE or an SAP affiliate company. All rights reserved.

package passport

import (
	"bytes"
	"fmt"
)

type differenceElement struct {
	field  string
	value1 interface{}
	value2 interface{}
}

func regularCompare(mylist []differenceElement, name string, value1 int, value2 int) []differenceElement {
	if value1 != value2 {
		return append(mylist, differenceElement{name, value1, value2})
	}
	return mylist
}

func arrayCompare(mylist []differenceElement, name string, value1 [16]byte, value2 [16]byte) []differenceElement {
	if value1 != value2 {
		return append(mylist, differenceElement{name, value1, value2})
	}
	return mylist
}

func byteCompare(mylist []differenceElement, name string, value1 []byte, value2 []byte) []differenceElement {
	if !bytes.Equal(value1, value2) {
		return append(mylist, differenceElement{name, value1, value2})
	}
	return mylist
}

func stringCompare(mylist []differenceElement, name string, value1 string, value2 string) []differenceElement {
	if value1 != value2 {
		return append(mylist, differenceElement{name, value1, value2})
	}
	return mylist
}

// Differences identifies the number of differences between two
// passports. Additionally, it returns a string that shows the differences
// between the two passports.
func Differences(pp1 *Passport, pp2 *Passport) (int, string) {
	differentFields := []differenceElement{}
	differentFields = regularCompare(differentFields, "Version", int(pp1.Version()), int(pp2.Version()))
	differentFields = regularCompare(differentFields, "Length", int(pp1.Length()), int(pp2.Length()))
	differentFields = regularCompare(differentFields, "TraceFlag", int(pp1.TraceFlags()), int(pp2.TraceFlags()))
	differentFields = stringCompare(differentFields, "ComponentName", pp1.ComponentName(), pp2.ComponentName())
	differentFields = regularCompare(differentFields, "Service", int(pp1.Service()), int(pp2.Service()))
	differentFields = stringCompare(differentFields, "User", pp1.UserID(), pp2.UserID())
	differentFields = stringCompare(differentFields, "Action", pp1.Action(), pp2.Action())
	differentFields = regularCompare(differentFields, "ActionType", int(pp1.ActionType()), int(pp2.ActionType()))
	differentFields = stringCompare(differentFields, "PreviousComponentName", pp1.PreviousComponentName(), pp2.PreviousComponentName())
	differentFields = arrayCompare(differentFields, "TransactionID", pp1.TransactionID(), pp2.TransactionID())
	differentFields = byteCompare(differentFields, "ClientNumber", pp1.ClientNumber(), pp2.ClientNumber())
	differentFields = regularCompare(differentFields, "ComponentType", int(pp1.ComponentType()), int(pp2.ComponentType()))
	differentFields = arrayCompare(differentFields, "RootContextID", pp1.RootContextID(), pp2.RootContextID())
	differentFields = arrayCompare(differentFields, "ConnectionID", pp1.ConnectionID(), pp2.ConnectionID())
	differentFields = regularCompare(differentFields, "ConnectionCounter", int(pp1.ConnectionCounter()), int(pp2.ConnectionCounter()))
	differentFields = regularCompare(differentFields, "VariablePartsNumber", int(pp1.VariablePartsNumber()), int(pp2.VariablePartsNumber()))
	differentFields = regularCompare(differentFields, "VariablePartsOffset", int(pp1.VariablePartsOffset()), int(pp2.VariablePartsOffset()))
	differentFields = byteCompare(differentFields, "VarPart", pp1.VariableParts(), pp2.VariableParts())
	differences := ""
	for _, el := range differentFields {
		differences += fmt.Sprintf("Field %s: '%v' vs '%v'\n", el.field, el.value1, el.value2)
	}
	return len(differentFields), differences
}

// Equal returns true if two passports are identical and false otherwise.
func Equal(pp1 *Passport, pp2 *Passport) bool {
	if pp1.Version() != pp2.Version() {
		return false
	}
	if pp1.Length() != pp2.Length() {
		return false
	}
	if pp1.TraceFlags() != pp2.TraceFlags() {
		return false
	}
	if pp1.ComponentName() != pp2.ComponentName() {
		return false
	}
	if pp1.Service() != pp2.Service() {
		return false
	}
	if pp1.UserID() != pp2.UserID() {
		return false
	}
	if pp1.Action() != pp2.Action() {
		return false
	}
	if pp1.ActionType() != pp2.ActionType() {
		return false
	}
	if pp1.PreviousComponentName() != pp2.PreviousComponentName() {
		return false
	}
	if pp1.TransactionID() != pp2.TransactionID() {
		return false
	}
	if !bytes.Equal(pp1.ClientNumber(), pp2.ClientNumber()) {
		return false
	}
	if pp1.ComponentType() != pp2.ComponentType() {
		return false
	}
	if pp1.RootContextID() != pp2.RootContextID() {
		return false
	}
	if pp1.ConnectionID() != pp2.ConnectionID() {
		return false
	}
	if pp1.ConnectionCounter() != pp2.ConnectionCounter() {
		return false
	}
	if pp1.VariablePartsNumber() != pp2.VariablePartsNumber() {
		return false
	}
	if pp1.VariablePartsOffset() != pp2.VariablePartsOffset() {
		return false
	}
	if pp1.ConnectionCounter() != pp2.ConnectionCounter() {
		return false
	}
	if !bytes.Equal(pp1.VariableParts(), pp2.VariableParts()) {
		return false
	}
	return true
}

// MakeComponentName creates a correct component name.
// The serviceId is the same as the component type and can be found here:
// https://github.wdf.sap.corp/xdsr/epp-component-type/blob/master/sap.crun.nameservicedemo-CloudService.csv
// The tenantId is the the CLD (Cloud Landscape Directory) tenant ID. See here:
// https://wiki.wdf.sap.corp/wiki/pages/viewpage.action?spaceKey=CLMAM&title=CLD+Tenant+ID
func MakeComponentName(serviceID uint16, tenantID uint64) string {
	return fmt.Sprintf("%04d@%018d", serviceID, tenantID)
}

func collectErrors(errors []error) []string {
	var errorStrings []string
	for _, err := range errors {
		if err != nil {
			errorStrings = append(errorStrings, err.Error())
		}
	}
	return errorStrings
}
