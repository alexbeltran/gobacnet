package gobacnet

import (
	"context"
	"fmt"
	"time"

	"github.com/alexbeltran/gobacnet/encoding"
	bactype "github.com/alexbeltran/gobacnet/types"
)

func (c *Client) WriteProperty(dest bactype.Device, wp bactype.PropertyData) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	id, err := c.tsm.ID(ctx)
	if err != nil {
		return fmt.Errorf("unable to get an transaction id: %v", err)
	}
	defer c.tsm.Put(id)

	enc := encoding.NewEncoder()
	enc.NPDU(bactype.NPDU{
		Version:               bactype.ProtocolVersion,
		Destination:           &dest.Addr,
		Source:                c.dataLink.GetMyAddress(),
		IsNetworkLayerMessage: false,
		ExpectingReply:        true,
		Priority:              bactype.Normal,
		HopCount:              bactype.DefaultHopCount,
	})

	enc.WriteProperty(uint8(id), wp)
	if enc.Error() != nil {
		return enc.Error()
	}

	// the value filled doesn't matter. it just needs to be non nil
	err = fmt.Errorf("go")
	for count := 0; err != nil && count < 2; count++ {
		var b []byte
		var raw interface{}
		_, err = c.Send(dest.Addr, enc.Bytes())
		if err != nil {
			continue
		}

		raw, err = c.tsm.Receive(id, time.Duration(5)*time.Second)
		if err != nil {
			continue
		}
		switch v := raw.(type) {
		case error:
			return err
		case []byte:
			b = v
		default:
			return fmt.Errorf("received unknown datatype %T", raw)
		}

		dec := encoding.NewDecoder(b)

		var apdu bactype.APDU
		if err = dec.APDU(&apdu); err != nil {
			continue
		}
		if apdu.Error.Class != 0 || apdu.Error.Code != 0 {
			err = fmt.Errorf("received error, class: %d, code: %d", apdu.Error.Class, apdu.Error.Code)
			continue
		}

		return err
	}
	return err
}
