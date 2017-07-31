package encoding

import (
	bactype "github.com/alexbeltran/gobacnet/types"
)

func (e *Encoder) readMultiPropertyHeader(tagPos uint8, data bactype.ReadPropertyData) uint8 {
	// For each object
	// Tag 0 - Encode Object ID

	// Tag 1 - Opening Tag
	// for each property
	// Tag 0 - Property ID
	// Tag 1 (OPTIONAL) - Array Length
	// endfor
	// Tag 1 - Closing Tag

	// endfor
	return 0
}
