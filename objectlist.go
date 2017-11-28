package gobacnet

import (
	"fmt"
	"log"

	"github.com/alexbeltran/gobacnet/property"
	bactype "github.com/alexbeltran/gobacnet/types"
)

func (c *Client) objectListLen(dev bactype.Device) (int, error) {
	rp := bactype.ReadPropertyData{
		Object: bactype.Object{
			ID: dev.ID,
			Properties: []bactype.Property{
				bactype.Property{
					Type:       property.ObjectList,
					ArrayIndex: 0,
				},
			},
		},
	}

	resp, err := c.ReadProperty(dev, rp)
	if err != nil {
		return 0, err
	}
	data, ok := resp.Object.Properties[0].Data.(uint32)
	if !ok {
		return 0, fmt.Errorf("Unable to get object length")
	}
	return int(data), nil
}

func (c *Client) objectsRange(dev bactype.Device, start, end int) ([]bactype.Object, error) {
	rpm := bactype.ReadMultipleProperty{
		Objects: []bactype.Object{
			bactype.Object{
				ID: dev.ID,
			},
		},
	}

	for i := start; i <= end; i++ {
		rpm.Objects[0].Properties = append(rpm.Objects[0].Properties, bactype.Property{
			Type:       property.ObjectList,
			ArrayIndex: uint32(i),
		})
	}
	resp, err := c.ReadMultiProperty(dev, rpm)
	if err != nil {
		return nil, err
	}
	if len(resp.Objects) == 0 {
		return nil, fmt.Errorf("No data was returned")
	}

	objs := make([]bactype.Object, len(resp.Objects[0].Properties))

	for i, prop := range resp.Objects[0].Properties {
		id, ok := prop.Data.(bactype.ObjectID)
		if !ok {
			return nil, fmt.Errorf("Expected type Object ID, got %T", prop.Data)
		}
		objs[i].ID = id
	}

	return objs, nil
}

const readPropRequestSize = 16

func (c *Client) Objects(dev bactype.Device) (bactype.Device, error) {
	dev.Objects = make(map[uint32]bactype.Object)

	l, err := c.objectListLen(dev)
	if err != nil {
		return dev, err
	}

	// Scan size is broken
	scanSize := int(dev.MaxApdu) / readPropRequestSize
	i := 0
	for i = 0; i < l/scanSize; i++ {
		start := i*scanSize + 1
		end := (i + 1) * scanSize
		log.Printf("%d -> %d", start, end)

		objs, err := c.objectsRange(dev, start, end)
		if err != nil {
			return dev, err
		}

		for _, o := range objs {
			dev.Objects[o.ID.Instance] = o
		}
	}
	start := i*scanSize + 1
	end := l
	if start <= end {
		log.Printf("%d -> %d", start, end)
		objs, err := c.objectsRange(dev, start, end)
		if err != nil {
			return dev, err
		}
		for _, o := range objs {
			dev.Objects[o.ID.Instance] = o
		}
	}
	return dev, nil
}
