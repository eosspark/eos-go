// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/eosspark/eos-go/log"
)

// template type Exception(PARENT,CODE,WHAT)

var DuplicateAbiTypeDefExceptionName = reflect.TypeOf(DuplicateAbiTypeDefException{}).Name()

type DuplicateAbiTypeDefException struct {
	_AbiException
	Elog log.Messages
}

func NewDuplicateAbiTypeDefException(parent _AbiException, message log.Message) *DuplicateAbiTypeDefException {
	return &DuplicateAbiTypeDefException{parent, log.Messages{message}}
}

func (e DuplicateAbiTypeDefException) Code() int64 {
	return 3150005
}

func (e DuplicateAbiTypeDefException) Name() string {
	return DuplicateAbiTypeDefExceptionName
}

func (e DuplicateAbiTypeDefException) What() string {
	return "Duplicate type definition in the ABI"
}

func (e *DuplicateAbiTypeDefException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e DuplicateAbiTypeDefException) GetLog() log.Messages {
	return e.Elog
}

func (e DuplicateAbiTypeDefException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e DuplicateAbiTypeDefException) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteString(" ")
	buffer.WriteString(e.Name())
	buffer.WriteString(": ")
	buffer.WriteString(e.What())
	buffer.WriteString("\n")
	for _, l := range e.Elog {
		buffer.WriteString("[")
		buffer.WriteString(l.GetMessage())
		buffer.WriteString("]")
		buffer.WriteString("\n")
		buffer.WriteString(l.GetContext().String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (e DuplicateAbiTypeDefException) String() string {
	return e.DetailMessage()
}

func (e DuplicateAbiTypeDefException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3150005,
		Name: DuplicateAbiTypeDefExceptionName,
		What: "Duplicate type definition in the ABI",
	}

	return json.Marshal(except)
}

func (e DuplicateAbiTypeDefException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*DuplicateAbiTypeDefException):
		callback(&e)
		return true
	case func(DuplicateAbiTypeDefException):
		callback(e)
		return true
	default:
		return false
	}
}
