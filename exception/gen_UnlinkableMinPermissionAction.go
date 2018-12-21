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

var UnlinkableMinPermissionActionName = reflect.TypeOf(UnlinkableMinPermissionAction{}).Name()

type UnlinkableMinPermissionAction struct {
	_AuthorizationException
	Elog log.Messages
}

func NewUnlinkableMinPermissionAction(parent _AuthorizationException, message log.Message) *UnlinkableMinPermissionAction {
	return &UnlinkableMinPermissionAction{parent, log.Messages{message}}
}

func (e UnlinkableMinPermissionAction) Code() int64 {
	return 3090008
}

func (e UnlinkableMinPermissionAction) Name() string {
	return UnlinkableMinPermissionActionName
}

func (e UnlinkableMinPermissionAction) What() string {
	return "The action is not allowed to be linked with minimum permission"
}

func (e *UnlinkableMinPermissionAction) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e UnlinkableMinPermissionAction) GetLog() log.Messages {
	return e.Elog
}

func (e UnlinkableMinPermissionAction) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e UnlinkableMinPermissionAction) DetailMessage() string {
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

func (e UnlinkableMinPermissionAction) String() string {
	return e.DetailMessage()
}

func (e UnlinkableMinPermissionAction) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3090008,
		Name: UnlinkableMinPermissionActionName,
		What: "The action is not allowed to be linked with minimum permission",
	}

	return json.Marshal(except)
}

func (e UnlinkableMinPermissionAction) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*UnlinkableMinPermissionAction):
		callback(&e)
		return true
	case func(UnlinkableMinPermissionAction):
		callback(e)
		return true
	default:
		return false
	}
}
