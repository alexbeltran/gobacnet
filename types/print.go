package types

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alexbeltran/gobacnet/property"
)

const defaultSpacing = 4

const (
	AnalogInput       ObjectType = 0
	AnalogOutput      ObjectType = 1
	AnalogValue       ObjectType = 2
	BinaryInput       ObjectType = 3
	BinaryOutput      ObjectType = 4
	BinaryValue       ObjectType = 5
	DeviceType        ObjectType = 8
	File              ObjectType = 10
	MultiStateInput   ObjectType = 13
	NotificationClass ObjectType = 15
	MultiStateValue   ObjectType = 19
	TrendLog          ObjectType = 20
	CharacterString   ObjectType = 40
)

const (
	AnalogInputStr       = "Analog Input"
	AnalogOutputStr      = "Analog Output"
	AnalogValueStr       = "Analog Value"
	BinaryInputStr       = "Binary Input"
	BinaryOutputStr      = "Binary Output"
	BinaryValueStr       = "Binary Value"
	DeviceTypeStr        = "Device"
	FileStr              = "File"
	NotificationClassStr = "Notification Class"
	MultiStateValueStr   = "Multi-State Value"
	MultiStateInputStr   = "Multi-State Input"
	TrendLogStr          = "Trend Log"
	CharacterStringStr   = "Character String"
)

var objTypeMap = map[ObjectType]string{
	AnalogInput:       AnalogInputStr,
	AnalogOutput:      AnalogOutputStr,
	AnalogValue:       AnalogValueStr,
	BinaryInput:       BinaryInputStr,
	BinaryOutput:      BinaryOutputStr,
	BinaryValue:       BinaryValueStr,
	DeviceType:        DeviceTypeStr,
	File:              FileStr,
	NotificationClass: NotificationClassStr,
	MultiStateValue:   MultiStateValueStr,
	MultiStateInput:   MultiStateInputStr,
	TrendLog:          TrendLogStr,
	CharacterString:   CharacterStringStr,
}

var objStrTypeMap = map[string]ObjectType{
	AnalogInputStr:       AnalogInput,
	AnalogOutputStr:      AnalogOutput,
	AnalogValueStr:       AnalogValue,
	BinaryInputStr:       BinaryInput,
	BinaryOutputStr:      BinaryOutput,
	BinaryValueStr:       BinaryValue,
	DeviceTypeStr:        DeviceType,
	FileStr:              File,
	NotificationClassStr: NotificationClass,
	MultiStateValueStr:   MultiStateValue,
	TrendLogStr:          TrendLog,
	CharacterStringStr:   CharacterString,
}

func GetType(s string) ObjectType {
	t, ok := objStrTypeMap[s]
	if !ok {
		return 0
	}
	return t
}

func (t ObjectType) String() string {
	s, ok := objTypeMap[t]
	if !ok {
		return fmt.Sprintf("Unknown (%d)", t)
	}
	return fmt.Sprintf("%s", s)
}

// String returns a pretty print of the ObjectID structure
func (id ObjectID) String() string {
	return fmt.Sprintf("Instance: %d Type: %s", id.Instance, id.Type.String())
}

// String returns a pretty print of the read multiple property structure
func (rp ReadMultipleProperty) String() string {
	buff := bytes.Buffer{}
	spacing := strings.Repeat(" ", defaultSpacing)
	for _, obj := range rp.Objects {
		buff.WriteString(obj.ID.String())
		buff.WriteString("\n")
		for _, prop := range obj.Properties {
			buff.WriteString(spacing)
			buff.WriteString(property.String(prop.Type))
			buff.WriteString(fmt.Sprintf("[%v]", prop.ArrayIndex))
			buff.WriteString(": ")
			buff.WriteString(fmt.Sprintf("%v", prop.Data))
			buff.WriteString("\n")
		}
		buff.WriteString("\n")
	}
	return buff.String()
}
