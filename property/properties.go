package property

import "fmt"

const (
	FileSize         uint32 = 42
	FileType         uint32 = 43
	ModelName        uint32 = 70
	ObjectIdentifier uint32 = 75
	ObjectList       uint32 = 76
	ObjectName       uint32 = 77
	ObjectReference  uint32 = 78
	ObjectType       uint32 = 79
	PresentValue     uint32 = 86
	Units            uint32 = 117
)

const (
	ObjectNameStr = "ObjectName"
)

// enumMapping should be treated as read only.
var enumMapping = map[string]uint32{
	"FileSize":         FileSize,
	"FileType":         FileType,
	"ModelName":        ModelName,
	"ObjectIdentifier": ObjectIdentifier,
	"ObjectList":       ObjectList,
	ObjectNameStr:      ObjectName,
	"ObjectReference":  ObjectReference,
	"ObjectType":       ObjectType,
	"PresentValue":     PresentValue,
	"Units":            Units,
}

// listOfKeys should be treated as read only after init
var listOfKeys []string

func init() {
	listOfKeys := make([]string, len(enumMapping))
	i := 0
	for k := range enumMapping {
		listOfKeys[i] = k
		i++
	}
}

func Keys() []string {
	// A copy is made since we do not want outside packages editing our keys by
	// accident
	var keys []string
	copy(keys, listOfKeys)
	return keys
}

func Get(s string) (uint32, error) {
	if v, ok := enumMapping[s]; ok {
		return v, nil
	}
	err := fmt.Errorf("%s is not a valid property.", s)
	return 0, err
}