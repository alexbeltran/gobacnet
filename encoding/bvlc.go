package encoding

import (
	bactype "github.com/alexbeltran/gobacnet/types"
)

// Bacnet Virtual Layer Control

func (e *Encoder) BVLC(b bactype.BVLC) error {
	// Set packet type
	e.write(b.Type)
	e.write(b.Function)
	e.write(b.Length)
	e.write(b.Data)
	return e.Error()
}

func (d *Decoder) BVLC(b *bactype.BVLC) error {
	d.decode(&b.Type)
	d.decode(&b.Function)
	d.decode(&b.Length)
	d.decode(&b.Data)
	return d.Error()
}
