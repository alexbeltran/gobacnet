package encoding

import (
	"fmt"

	"github.com/alexbeltran/gobacnet/property"
	"github.com/alexbeltran/gobacnet/types"
	bactype "github.com/alexbeltran/gobacnet/types"
)

func isValidObjectType(idType types.ObjectType) error {
	if idType > bactype.MaxObject {
		return fmt.Errorf("Object types is %d which must be less then %d", idType, bactype.MaxObject)
	}
	return nil
}

func isValidPropertyType(propType property.PropertyID) error {
	if propType > MaxPropertyID {
		return fmt.Errorf("Object types is %d which must be less then %d", propType, MaxPropertyID)
	}
	return nil
}
