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

var HttpExceptionName = reflect.TypeOf(HttpException{}).Name()

type HttpException struct {
	_HttpException
	Elog log.Messages
}

func NewHttpException(parent _HttpException, message log.Message) *HttpException {
	return &HttpException{parent, log.Messages{message}}
}

func (e HttpException) Code() int64 {
	return 3200000
}

func (e HttpException) Name() string {
	return HttpExceptionName
}

func (e HttpException) What() string {
	return "http exception"
}

func (e *HttpException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e HttpException) GetLog() log.Messages {
	return e.Elog
}

func (e HttpException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e HttpException) DetailMessage() string {
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

func (e HttpException) String() string {
	return e.DetailMessage()
}

func (e HttpException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3200000,
		Name: HttpExceptionName,
		What: "http exception",
	}

	return json.Marshal(except)
}

func (e HttpException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*HttpException):
		callback(&e)
		return true
	case func(HttpException):
		callback(e)
		return true
	default:
		return false
	}
}
