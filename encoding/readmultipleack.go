package encoding

import (
	bactype "github.com/alexbeltran/gobacnet/types"
)

func (e *Encoder) ReadMultiplePropertyAck(invokeID uint8, data bactype.ReadMultipleProperty) error {
	a := bactype.APDU{
		DataType: bactype.ComplexAck,
		Service:  bactype.ServiceConfirmedReadPropMultiple,
		InvokeId: invokeID,
	}
	e.APDU(a)
	err := e.objectsWithData(data.Objects)
	if err != nil {
		return err
	}

	return e.Error()
}

func (e *Encoder) objectsWithData(objects []bactype.Object) error {
	var tag uint8
	for _, obj := range objects {
		tag = 0
		e.contextObjectID(tag, obj.ID.Type, obj.ID.Instance)
		// Tag 1 - Opening Tag
		tag = 1
		e.openingTag(tag)

		e.propertiesWithData(obj.Properties)

		// Tag 1 - Closing Tag
		e.closingTag(tag)
	}
	return nil
}

func (e *Encoder) propertiesWithData(properties []bactype.Property) error {
	var tag uint8
	for _, prop := range properties {
		// Tag 2 - Property ID
		tag = 2
		e.contextEnumerated(tag, prop.Type)

		// Tag 3 (OPTIONAL) - Array Length
		tag++
		if prop.ArrayIndex != ArrayAll {
			e.contextUnsigned(tag, prop.ArrayIndex)
		}

		// Tag 4 Opening Tag
		tag++
		openedTag := tag
		e.openingTag(openedTag)
		e.write(prop.Data)
		e.closingTag(openedTag)

		e.write(prop.Data)
	}
	return e.Error()
}
