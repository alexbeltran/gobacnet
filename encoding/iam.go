package encoding

import (
	"fmt"

	"github.com/alexbeltran/gobacnet/types"
)

func (d *Decoder) IAm(ids []types.ObjectID) error {
	for i := 0; i < len(ids) && d.len() > 0; i++ {
		obj, err := d.AppData()
		// Issue decoding data
		if err != nil {
			return err
		}

		// Check type we receive
		switch t := obj.(type) {
		case types.ObjectID:
			ids[i] = t
		default:
			return fmt.Errorf("Expected type ObjectID, received a %T", t)
		}
	}
	return nil
}
