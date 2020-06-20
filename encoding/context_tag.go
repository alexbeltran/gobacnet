package encoding

import "github.com/alexbeltran/gobacnet/types"

/* returns the fixed tag type for certain context tagged properties */
func tagTypeInContext(property types.PropertyType, tagNumber uint8) uint8 {
	var tag uint8 = tagNumber
	switch property {
	case types.PROP_ACTUAL_SHED_LEVEL:
	case types.PROP_REQUESTED_SHED_LEVEL:
	case types.PROP_EXPECTED_SHED_LEVEL:
		switch tagNumber {
		case 0:
		case 1:
			tag = tagUint
			break
		case 2:
			tag = tagReal
			break
		default:
			break
		}
		break
	case types.PROP_ACTION:
		switch tagNumber {
		case 0:
		case 1:
			tag = tagObjectID
			break
		case 2:
			tag = tagEnumerated
			break
		case 3:
		case 5:
		case 6:
			tag = tagUint
			break
		case 7:
		case 8:
			tag = tagBool
			break
		case 4: /* propertyValue: abstract syntax */
		default:
			break
		}
		break
	case types.PROP_LIST_OF_GROUP_MEMBERS:
		/* Sequence of ReadAccessSpecification */
		switch tagNumber {
		case 0:
			tag = tagObjectID
			break
		default:
			break
		}
		break
	case types.PROP_EXCEPTION_SCHEDULE:
		switch tagNumber {
		case 1:
			tag = tagObjectID
			break
		case 3:
			tag = tagUint
			break
		case 0: /* calendarEntry: abstract syntax + context */
		case 2: /* list of BACnetTimeValue: abstract syntax */
		default:
			break
		}
		break
	case types.PROP_LOG_DEVICE_OBJECT_PROPERTY:
		switch tagNumber {
		case 0: /* Object ID */
		case 3: /* Device ID */
			tag = tagObjectID
			break
		case 1: /* Property ID */
			tag = tagEnumerated
			break
		case 2: /* Array index */
			tag = tagUint
			break
		default:
			break
		}
		break
	case types.PROP_SUBORDINATE_LIST:
		/* BACnetARRAY[N] of BACnetDeviceObjectReference */
		switch tagNumber {
		case 0: /* Optional Device ID */
		case 1: /* Object ID */
			tag = tagObjectID
			break
		default:
			break
		}
		break

	case types.PROP_RECIPIENT_LIST:
		/* List of BACnetDestination */
		switch tagNumber {
		case 0: /* Device Object ID */
			tag = tagObjectID
			break
		case 1:
			/* 2015.08.22 EKH 135-2012 pg 708
			   todo - Context 1 in Recipient list would be a BACnetAddress, not coded yet...
			   BACnetRecipient::= CHOICE {
			        device  [0] BACnetObjectIdentifier,
			        address  [1] BACnetAddress
			         }
			*/
			break
		default:
			break
		}
		break
	case types.PROP_ACTIVE_COV_SUBSCRIPTIONS:
		/* BACnetCOVSubscription */
		switch tagNumber {
		case 0: /* BACnetRecipientProcess */
		case 1: /* BACnetObjectPropertyReference */
			break
		case 2: /* issueConfirmedNotifications */
			tag = tagBool
			break
		case 3: /* timeRemaining */
			tag = tagUint
			break
		case 4: /* covIncrement */
			tag = tagReal
			break
		default:
			break
		}
		break
	default:
		break
	}

	return tag
}
