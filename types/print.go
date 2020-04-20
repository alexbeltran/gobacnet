package types

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alexbeltran/gobacnet/property"
)

const defaultSpacing = 4

var unitsTypeMap = map[int]string{
	//--Area
	0: "square-meters",
	1: "square-feet",

	//--currency
	105: "currencyl",
	106: "currency2",
	107: "currency3",
	108: "currency4",
	109: "currency5",
	110: "currency6",
	111: "currency7",
	112: "currency8",
	113: "currency9",
	114: "currency10",

	//--Electrical
	2:   "milliamperes",
	3:   "amperes",
	4:   "ohms",
	122: "kilohms",
	123: "megohms",
	5:   "volts",
	124: "millivolts",
	6:   "kilovolts",
	7:   "megavolts",
	8:   "volt-amperes",
	9:   "kilovolt-amperes",
	10:  "megavolt-amperes",
	11:  "volt-amperes-reactive",
	12:  "kilovolt-amperes-reactive",
	13:  "megavolt-amperes-reactive",
	14:  "degrees-phase",
	15:  "power-factor",

	//--Energy
	16:  "joules",
	17:  "kilojoules",
	125: "kilojoules-per-kilogram",
	126: "megajoules",
	18:  "watt-hours",
	19:  "kilowatt-hours",
	20:  "btus",
	21:  "therms",
	22:  "ton-hours",

	//--Enthalpy
	23: "joules-per-kilogram-dry-air",
	24: "btus-per-pound-dry-air",

	//--Entropy
	127: "joules-per-degree-Kelvin",
	128: "joules-per-kilogram-degree-Kelvin",

	//--Frequency
	25:  "cycles-per-hour",
	26:  "cycles-per-minute",
	27:  "hertz",
	129: "kilohertz",
	130: "megahertz",
	131: "per-hour",

	//--Humidity
	28: "grams-of-water-per-kilogram-dry-air",
	29: "percent-relative-humidity",

	//--Length
	30: "millimeters",
	31: "meters",
	32: "inches",
	33: "feet",

	//--Light
	34: "watts-per-square-foot",
	35: "Watts-per-square-meter",
	36: "lumens",
	37: "luxes",
	38: "foot-candles",

	//--Mass
	39: "kilograms",
	40: "pounds-mass",
	41: "tons",

	//--Mass Flow
	42: "kilograms-per-second",
	43: "kilograms -per-minute",
	44: "kilograms-per-hour",
	45: "pounds-mass-per-minute",
	46: "pounds-mass-per-hour",

	//--Power
	132: "milliwatts",
	47:  "watts",
	48:  "kilowatts",
	49:  "megawatts",
	50:  "btus-per-hour",
	51:  "horsepower",
	52:  "tons-refrigeration",

	//--Pressure
	53:  "pascals",
	133: "hectopascals",
	54:  "kilopascals",
	134: "millibars",
	55:  "bars",
	56:  "pounds-force-per-square-inch",
	57:  "centimeters-of-water",
	58:  "inches-of-water",
	59:  "millimeters-of-mercury",
	60:  "centimeters-of-mercury",
	61:  "inches-of-mercury",

	//--Temperature
	62: "degrees-Celsius",
	63: "degrees-Kelvin",
	64: "degrees-Fahrenheit",
	65: "degree-days-Celsius",
	66: "degree-days-Fahrenheit",

	//--Time
	67: "years",
	68: "months",
	69: "weeks",
	70: "days",
	71: "hours",
	72: "minutes",
	73: "seconds",

	//--Velocity
	74: "meters-per-second",
	75: "kilometers-per-hour",
	76: "feet-per-second",
	77: "feet-per-minute",
	78: "miles-per-hour",

	//--volume
	79: "cubic-feet",
	80: "cubic-meters",
	81: "imperial-gallons",
	82: "liters",
	83: "us-gallons",

	//--Volumetric Flow
	84:  "cubic-feet-per-minute",
	85:  "cubic-meters-per-second",
	135: "cubic-meters-per-hour",
	86:  "imperial-gallons-per-minute",
	87:  "liters-per-second",
	88:  "liters-per-minute",
	136: "liters-per-hour",
	89:  "us-gallons-per-minute",

	//--Other
	90:  "degrees-angular",
	91:  "degrees-Celsius-per-hour",
	92:  "degrees-Celsius-per-minute",
	93:  "degrees-Fahrenheit-per-hour",
	94:  "degrees-Fahrenheit-per-minute",
	137: "kilowatt-hours-per-square-meter",
	138: "kilowatt-hours-per-square-foot",
	139: "megajoules-per-square-meter",
	140: "megajoules-per-square-foot",
	95:  "no-units",
	96:  "parts-per-million",
	97:  "parts-per-billion",
	98:  "percent",
	99:  "percent-per-second",
	100: "per-minute",
	101: "per-second",
	102: "psi-per-degree-Fahrenheit",
	103: "radians",
	104: "revolutions-per-minute",
	141: "watts-per-square-meter-degree-kelvin",
}

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
