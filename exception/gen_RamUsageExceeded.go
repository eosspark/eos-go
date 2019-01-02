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

var RamUsageExceededName = reflect.TypeOf(RamUsageExceeded{}).Name()

type RamUsageExceeded struct {
	_ResourceExhaustedException
	Elog log.Messages
}

func NewRamUsageExceeded(parent _ResourceExhaustedException, message log.Message) *RamUsageExceeded {
	return &RamUsageExceeded{parent, log.Messages{message}}
}

func (e RamUsageExceeded) Code() int64 {
	return 3080001
}

func (e RamUsageExceeded) Name() string {
	return RamUsageExceededName
}

func (e RamUsageExceeded) What() string {
	return "Account using more than allotted RAM usage"
}

func (e *RamUsageExceeded) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e RamUsageExceeded) GetLog() log.Messages {
	return e.Elog
}

func (e RamUsageExceeded) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e RamUsageExceeded) DetailMessage() string {
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

func (e RamUsageExceeded) String() string {
	return e.DetailMessage()
}

func (e RamUsageExceeded) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3080001,
		Name: RamUsageExceededName,
		What: "Account using more than allotted RAM usage",
	}

	return json.Marshal(except)
}

func (e RamUsageExceeded) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*RamUsageExceeded):
		callback(&e)
		return true
	case func(RamUsageExceeded):
		callback(e)
		return true
	default:
		return false
	}
}