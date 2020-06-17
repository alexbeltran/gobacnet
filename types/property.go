package types

import (
	"bytes"
	"fmt"
	"strings"
)

const defaultSpacing = 4

type PropertyType uint32

const (
	PropAllProperties    PropertyType = 8
	PropDescription      PropertyType = 28
	PropFileSize         PropertyType = 42
	PropFileType         PropertyType = 43
	PropModelName        PropertyType = 70
	PropObjectIdentifier PropertyType = 75
	PropObjectList       PropertyType = 76
	PropObjectName       PropertyType = 77
	PropObjectReference  PropertyType = 78
	PropObjectType       PropertyType = 79
	PropPresentValue     PropertyType = 85
	PropUnits            PropertyType = 117
	PropPriorityArray    PropertyType = 87
)

const (
	DescriptionStr = "Description"
	ObjectNameStr  = "ObjectName"
)

// propertyTypeMapping should be treated as read only.
var propertyTypeMapping = map[string]PropertyType{
	"AllProperties":    PropAllProperties,
	DescriptionStr:     PropDescription,
	"FileSize":         PropFileSize,
	"FileType":         PropFileType,
	"ModelName":        PropModelName,
	"ObjectIdentifier": PropObjectIdentifier,
	"ObjectList":       PropObjectList,
	ObjectNameStr:      PropObjectName,
	"ObjectReference":  PropObjectReference,
	"ObjectType":       PropObjectType,
	"PresentValue":     PropPresentValue,
	"Units":            PropUnits,
	"PriorityArray":    PropPriorityArray,
}

// propertyTypeStrMapping is a human readable printing of the priority
var propertyTypeStrMapping = map[PropertyType]string{
	PropAllProperties:    "All Properties",
	PropDescription:      "Description",
	PropFileSize:         "File Size",
	PropFileType:         "File Type",
	PropModelName:        "Model Name",
	PropObjectIdentifier: "Object Identifier",
	PropObjectList:       "Object List",
	PropObjectName:       "Object Name",
	PropObjectReference:  "Object Reference",
	PropObjectType:       "Object Type",
	PropPresentValue:     "Present Value",
	PropUnits:            "PropUnits",
	PropPriorityArray:    "Priority Array",
}

// listOfKeys should be treated as read only after init
var listOfKeys []string

func init() {
	listOfKeys = make([]string, len(propertyTypeMapping))
	i := 0
	for k := range propertyTypeMapping {
		listOfKeys[i] = k
		i++
	}
}

func Keys() map[string]PropertyType {
	// A copy is made since we do not want outside packages editing our keys by
	// accident
	keys := make(map[string]PropertyType)
	for k, v := range propertyTypeMapping {
		keys[k] = v
	}
	return keys
}

func Get(s string) (PropertyType, error) {
	if v, ok := propertyTypeMapping[s]; ok {
		return v, nil
	}
	err := fmt.Errorf("%s is not a valid property.", s)
	return 0, err
}

// String returns a human readible string of the given property
func String(prop PropertyType) string {
	s, ok := propertyTypeStrMapping[prop]
	if !ok {
		return "Unknown"
	}
	return fmt.Sprintf("%s (%d)", s, prop)
}

// The bool in the map doesn't actually matter since it won't be used.
var deviceProperties = map[PropertyType]bool{
	PropObjectList: true,
}

func IsDeviceProperty(id PropertyType) bool {
	_, ok := deviceProperties[id]
	return ok
}

type Property struct {
	Type       PropertyType
	ArrayIndex uint32
	Data       interface{}
	Priority   NPDUPriority
}

type PropertyData struct {
	InvokeID   uint16
	Object     Object
	ErrorClass uint8
	ErrorCode  uint8
}

type MultiplePropertyData struct {
	Objects    []Object
	ErrorClass uint8
	ErrorCode  uint8
}

// String returns a pretty print of the read multiple property structure
func (rp MultiplePropertyData) String() string {
	buff := bytes.Buffer{}
	spacing := strings.Repeat(" ", defaultSpacing)
	for _, obj := range rp.Objects {
		buff.WriteString(obj.ID.String())
		buff.WriteString("\n")
		for _, prop := range obj.Properties {
			buff.WriteString(spacing)
			buff.WriteString(String(prop.Type))
			buff.WriteString(fmt.Sprintf("[%v]", prop.ArrayIndex))
			buff.WriteString(": ")
			buff.WriteString(fmt.Sprintf("%v", prop.Data))
			buff.WriteString("\n")
		}
		buff.WriteString("\n")
	}
	return buff.String()
}

// PrintAllProperties prints all of the properties within this package. This is only a
// subset of all properties.
func PrintAllProperties() {
	max := func(x map[string]PropertyType) int {
		max := 0
		for k, _ := range x {
			if len(k) > max {
				max = len(k)
			}
		}
		return max
	}(propertyTypeMapping)

	const numOfAdditionalSpaces = 15

	printRow := func(col1, col2 string, maxLen int) {
		spacing := strings.Repeat(" ", maxLen-len(col1)+numOfAdditionalSpaces)
		fmt.Printf("%s%s%s\n", col1, spacing, col2)
	}

	printRow("Key", "Int", max)
	fmt.Println(strings.Repeat("-", max+numOfAdditionalSpaces+6))

	for k, id := range propertyTypeMapping {
		// Spacing
		printRow(k, fmt.Sprintf("%d", id), max)
	}
}
