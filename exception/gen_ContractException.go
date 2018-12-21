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

var ContractExceptionName = reflect.TypeOf(ContractException{}).Name()

type ContractException struct {
	_ContractException
	Elog log.Messages
}

func NewContractException(parent _ContractException, message log.Message) *ContractException {
	return &ContractException{parent, log.Messages{message}}
}

func (e ContractException) Code() int64 {
	return 3160000
}

func (e ContractException) Name() string {
	return ContractExceptionName
}

func (e ContractException) What() string {
	return "Contract exception"
}

func (e *ContractException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e ContractException) GetLog() log.Messages {
	return e.Elog
}

func (e ContractException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e ContractException) DetailMessage() string {
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

func (e ContractException) String() string {
	return e.DetailMessage()
}

func (e ContractException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3160000,
		Name: ContractExceptionName,
		What: "Contract exception",
	}

	return json.Marshal(except)
}

func (e ContractException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*ContractException):
		callback(&e)
		return true
	case func(ContractException):
		callback(e)
		return true
	default:
		return false
	}
}
