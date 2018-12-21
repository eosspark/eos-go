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

var ResourceExhaustedExceptionName = reflect.TypeOf(ResourceExhaustedException{}).Name()

type ResourceExhaustedException struct {
	_ResourceExhaustedException
	Elog log.Messages
}

func NewResourceExhaustedException(parent _ResourceExhaustedException, message log.Message) *ResourceExhaustedException {
	return &ResourceExhaustedException{parent, log.Messages{message}}
}

func (e ResourceExhaustedException) Code() int64 {
	return 3080000
}

func (e ResourceExhaustedException) Name() string {
	return ResourceExhaustedExceptionName
}

func (e ResourceExhaustedException) What() string {
	return "Resource exhausted exception"
}

func (e *ResourceExhaustedException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e ResourceExhaustedException) GetLog() log.Messages {
	return e.Elog
}

func (e ResourceExhaustedException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e ResourceExhaustedException) DetailMessage() string {
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

func (e ResourceExhaustedException) String() string {
	return e.DetailMessage()
}

func (e ResourceExhaustedException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3080000,
		Name: ResourceExhaustedExceptionName,
		What: "Resource exhausted exception",
	}

	return json.Marshal(except)
}

func (e ResourceExhaustedException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*ResourceExhaustedException):
		callback(&e)
		return true
	case func(ResourceExhaustedException):
		callback(e)
		return true
	default:
		return false
	}
}
