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

var EofExceptionName = reflect.TypeOf(EofException{}).Name()

type EofException struct {
	Exception
	Elog log.Messages
}

func NewEofException(parent Exception, message log.Message) *EofException {
	return &EofException{parent, log.Messages{message}}
}

func (e EofException) Code() int64 {
	return EofExceptionCode
}

func (e EofException) Name() string {
	return EofExceptionName
}

func (e EofException) What() string {
	return "End Of File"
}

func (e *EofException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e EofException) GetLog() log.Messages {
	return e.Elog
}

func (e EofException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e EofException) DetailMessage() string {
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

func (e EofException) String() string {
	return e.DetailMessage()
}

func (e EofException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: EofExceptionCode,
		Name: EofExceptionName,
		What: "End Of File",
	}

	return json.Marshal(except)
}

func (e EofException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*EofException):
		callback(&e)
		return true
	case func(EofException):
		callback(e)
		return true
	default:
		return false
	}
}
