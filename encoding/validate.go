package encoding

import (
	"fmt"
)

func isValidObjectType(idType uint16) error {
	if idType > MaxObject {
		return fmt.Errorf("Object types is %d which must be less then %d", idType, MaxObject)
	}
	return nil
}

func isValidPropertyType(propType uint32) error {
	if propType < MaxPropertyID {
		return fmt.Errorf("Object types is %d which must be less then %d", propType, MaxPropertyID)
	}
	return nil
}
