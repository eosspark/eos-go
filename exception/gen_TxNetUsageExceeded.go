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

var TxNetUsageExceededName = reflect.TypeOf(TxNetUsageExceeded{}).Name()

type TxNetUsageExceeded struct {
	_ResourceExhaustedException
	Elog log.Messages
}

func NewTxNetUsageExceeded(parent _ResourceExhaustedException, message log.Message) *TxNetUsageExceeded {
	return &TxNetUsageExceeded{parent, log.Messages{message}}
}

func (e TxNetUsageExceeded) Code() int64 {
	return 3080002
}

func (e TxNetUsageExceeded) Name() string {
	return TxNetUsageExceededName
}

func (e TxNetUsageExceeded) What() string {
	return "Transaction exceeded the current network usage limit imposed on the transaction"
}

func (e *TxNetUsageExceeded) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e TxNetUsageExceeded) GetLog() log.Messages {
	return e.Elog
}

func (e TxNetUsageExceeded) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e TxNetUsageExceeded) DetailMessage() string {
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

func (e TxNetUsageExceeded) String() string {
	return e.DetailMessage()
}

func (e TxNetUsageExceeded) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3080002,
		Name: TxNetUsageExceededName,
		What: "Transaction exceeded the current network usage limit imposed on the transaction",
	}

	return json.Marshal(except)
}

func (e TxNetUsageExceeded) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*TxNetUsageExceeded):
		callback(&e)
		return true
	case func(TxNetUsageExceeded):
		callback(e)
		return true
	default:
		return false
	}
}
