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

var ChainIdTypeExceptionName = reflect.TypeOf(ChainIdTypeException{}).Name()

type ChainIdTypeException struct {
	_ChainException
	Elog log.Messages
}

func NewChainIdTypeException(parent _ChainException, message log.Message) *ChainIdTypeException {
	return &ChainIdTypeException{parent, log.Messages{message}}
}

func (e ChainIdTypeException) Code() int64 {
	return 3010012
}

func (e ChainIdTypeException) Name() string {
	return ChainIdTypeExceptionName
}

func (e ChainIdTypeException) What() string {
	return "Invalid chain ID"
}

func (e *ChainIdTypeException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e ChainIdTypeException) GetLog() log.Messages {
	return e.Elog
}

func (e ChainIdTypeException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e ChainIdTypeException) DetailMessage() string {
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

func (e ChainIdTypeException) String() string {
	return e.DetailMessage()
}

func (e ChainIdTypeException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3010012,
		Name: ChainIdTypeExceptionName,
		What: "Invalid chain ID",
	}

	return json.Marshal(except)
}

func (e ChainIdTypeException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*ChainIdTypeException):
		callback(&e)
		return true
	case func(ChainIdTypeException):
		callback(e)
		return true
	default:
		return false
	}
}