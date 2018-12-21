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

var FcExceptionName = reflect.TypeOf(FcException{}).Name()

type FcException struct {
	Exception
	Elog log.Messages
}

func NewFcException(parent Exception, message log.Message) *FcException {
	return &FcException{parent, log.Messages{message}}
}

func (e FcException) Code() int64 {
	return UnspecifiedExceptionCode
}

func (e FcException) Name() string {
	return FcExceptionName
}

func (e FcException) What() string {
	return "unspecified"
}

func (e *FcException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e FcException) GetLog() log.Messages {
	return e.Elog
}

func (e FcException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e FcException) DetailMessage() string {
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

func (e FcException) String() string {
	return e.DetailMessage()
}

func (e FcException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: UnspecifiedExceptionCode,
		Name: FcExceptionName,
		What: "unspecified",
	}

	return json.Marshal(except)
}

func (e FcException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*FcException):
		callback(&e)
		return true
	case func(FcException):
		callback(e)
		return true
	default:
		return false
	}
}
