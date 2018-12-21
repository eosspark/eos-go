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

var InvalidArgExceptionName = reflect.TypeOf(InvalidArgException{}).Name()

type InvalidArgException struct {
	Exception
	Elog log.Messages
}

func NewInvalidArgException(parent Exception, message log.Message) *InvalidArgException {
	return &InvalidArgException{parent, log.Messages{message}}
}

func (e InvalidArgException) Code() int64 {
	return InvalidArgExceptionCode
}

func (e InvalidArgException) Name() string {
	return InvalidArgExceptionName
}

func (e InvalidArgException) What() string {
	return "Invalid Argument"
}

func (e *InvalidArgException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e InvalidArgException) GetLog() log.Messages {
	return e.Elog
}

func (e InvalidArgException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e InvalidArgException) DetailMessage() string {
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

func (e InvalidArgException) String() string {
	return e.DetailMessage()
}

func (e InvalidArgException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: InvalidArgExceptionCode,
		Name: InvalidArgExceptionName,
		What: "Invalid Argument",
	}

	return json.Marshal(except)
}

func (e InvalidArgException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*InvalidArgException):
		callback(&e)
		return true
	case func(InvalidArgException):
		callback(e)
		return true
	default:
		return false
	}
}
