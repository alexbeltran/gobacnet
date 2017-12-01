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

func (c *Client) objectList(dev *bactype.Device) error {
	dev.Objects = make(map[uint32]bactype.Object)

	l, err := c.objectListLen(*dev)
	if err != nil {
		return err
	}

	// Scan size is broken
	scanSize := int(dev.MaxApdu) / readPropRequestSize
	i := 0
	for i = 0; i < l/scanSize; i++ {
		start := i*scanSize + 1
		end := (i + 1) * scanSize
		log.Printf("%d -> %d", start, end)

		objs, err := c.objectsRange(*dev, start, end)
		if err != nil {
			return err
		}

		for _, o := range objs {
			dev.Objects[o.ID.Instance] = o
		}
	}
	start := i*scanSize + 1
	end := l
	if start <= end {
		log.Printf("%d -> %d", start, end)
		objs, err := c.objectsRange(*dev, start, end)
		if err != nil {
			return err
		}
		for _, o := range objs {
			dev.Objects[o.ID.Instance] = o
		}
	}
	return nil
}

func (c *Client) objectInformation(dev *bactype.Device) error {
	rpm := bactype.ReadMultipleProperty{
		Objects: []bactype.Object{},
	}

	// Often times the map will re arrange the order it spits out
	// so we need to keep track since the response will be in the
	// same order we issue the commands.
	keys := make([]uint32, len(dev.Objects))
	counter := 0
	for i, o := range dev.Objects {
		keys[counter] = i
		counter++
		rpm.Objects = append(rpm.Objects, bactype.Object{
			ID: o.ID,
			Properties: []bactype.Property{
				bactype.Property{
					Type:       property.ObjectName,
					ArrayIndex: bactype.ArrayAll,
				},
				bactype.Property{
					Type:       property.Description,
					ArrayIndex: bactype.ArrayAll,
				},
			},
		})

	}
	resp, err := c.ReadMultiProperty(*dev, rpm)
	if err != nil {
		return err
	}
	var name, description string
	var ok bool
	for i, r := range resp.Objects {
		name, ok = r.Properties[0].Data.(string)
		if !ok {
			return fmt.Errorf("Incorrect data returned")
		}
		description, ok = r.Properties[1].Data.(string)
		if !ok {
			return fmt.Errorf("Incorrect data returned")
		}
		obj := dev.Objects[keys[i]]
		obj.Name = name
		obj.Description = description
		dev.Objects[keys[i]] = obj
	}

	return nil
}

// Objects scans for all objects within the device. It will also gather
// additional information from the object such as the name and description of
// the objects.
func (c *Client) Objects(dev bactype.Device) (bactype.Device, error) {
	err := c.objectList(&dev)
	if err != nil {
		return dev, nil
	}
	log.Println(dev)
	err = c.objectInformation(&dev)
	log.Println(dev)
	return dev, err
}
