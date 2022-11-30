package packets

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrNotFoundProperties = errors.New("malformed packet: properties item")
)

type Properties struct {
	item map[byte]interface{}
}

func (prop *Properties) Decode(t byte, buf []byte, offset int) (int, error) {
	length, next, err := decodeUint8(buf, offset)
	fmt.Printf("读取指定的数据长度的数据 %v \r\n", length)
	if err != nil {
		return 0, err
	}

	if next+int(length) > len(buf) {
		return 0, ErrOffsetBytesOutOfRange
	}

	val_next := next
	for {

		//读取key
		idx, key_next, err := decodeByte(buf, val_next)

		if err != nil {
			return 0, err
		}
		//读取值
		switch idx {
		case PayloadFormatIndicator:
			prop.item[idx], val_next, err = decodeByte(buf, key_next)
		case MessageExpiryInterval:
			prop.item[idx], val_next, err = decodeUint32(buf, key_next)
		case ContentType:
			prop.item[idx], val_next, err = decodeString(buf, key_next)
		case CorrelationData:
			prop.item[idx], val_next, err = decodeBytes(buf, key_next)
		case SubscriptionIdentifier:
			prop.item[idx], val_next, err = decodeBytes(buf, key_next)
		case SessionExpiryInterval:
			prop.item[idx], val_next, err = decodeUint32(buf, key_next)
		case AssignedClientIdentifier:
			prop.item[idx], val_next, err = decodeString(buf, key_next)
		case ServerKeepAlive:
			prop.item[idx], val_next, err = decodeUint16(buf, key_next)
		case AuthenticationMethod:
			prop.item[idx], val_next, err = decodeString(buf, key_next)
		case AuthenticationData:
			prop.item[idx], val_next, err = decodeBytes(buf, key_next)
		case RequestProblemInformation:
			prop.item[idx], val_next, err = decodeByte(buf, key_next)
		case WillDelayInterval:
			prop.item[idx], val_next, err = decodeUint32(buf, key_next)
		case RequestResponseInformation:
			prop.item[idx], val_next, err = decodeByte(buf, key_next)
		case ResponseInformation:
			prop.item[idx], val_next, err = decodeString(buf, key_next)
		case ServerReference:
			prop.item[idx], val_next, err = decodeString(buf, key_next)
		case ReasonString:
			prop.item[idx], val_next, err = decodeString(buf, key_next)
		case ReceiveMaximum:
			prop.item[idx], val_next, err = decodeUint16(buf, key_next)
		case TopicAliasMaximum:
			prop.item[idx], val_next, err = decodeUint16(buf, key_next)
		case TopicAlias:
			prop.item[idx], val_next, err = decodeUint16(buf, key_next)
		case MaximumQoS:
			prop.item[idx], val_next, err = decodeByte(buf, key_next)
		case RetainAvailable:
			prop.item[idx], val_next, err = decodeByte(buf, key_next)
		case UserProperty:
			prop.item[idx], val_next, err = decodeBytes(buf, key_next)
		case MaximumPacketSize:
			prop.item[idx], val_next, err = decodeUint32(buf, key_next)
		case WildcardSubscriptionAvailable:
			prop.item[idx], val_next, err = decodeByte(buf, key_next)
		case SubscriptionIdentifierAvailable:
			prop.item[idx], val_next, err = decodeByte(buf, key_next)
		case SharedSubscriptionAvailable:
			prop.item[idx], val_next, err = decodeByte(buf, key_next)
		default:
			err = ErrNotFoundProperties
			return 0, err
		}
		if err != nil {
			return 0, err
		}
		//fmt.Printf("x: %v\r\n", prop.item[idx])
		if val_next > int(length) {
			break
		}
	}
	// prop 所有的数据
	return next + int(length), nil
}

func (prop *Properties) Encode(t byte, buf *bytes.Buffer) (ret []byte, length int, err error) {
	for idx, val := range prop.item {
		ret = append(ret, idx)
		length += 1
		//读取值
		switch idx {
		case PayloadFormatIndicator:
			ret = append(ret, val.(byte))
		case MessageExpiryInterval:
			ret = append(ret, encodeUint32(val.(uint32))...)
			// prop.item[idx], val_next, err = decodeUint32(buf, key_next)
		case ContentType:
			ret = append(ret, []byte(val.(string))...)
			// case CorrelationData:
			// 	prop.item[idx], val_next, err = decodeBytes(buf, key_next)
			// case SubscriptionIdentifier:
			// 	prop.item[idx], val_next, err = decodeBytes(buf, key_next)
			// case SessionExpiryInterval:
			// 	prop.item[idx], val_next, err = decodeUint32(buf, key_next)
			// case AssignedClientIdentifier:
			// 	prop.item[idx], val_next, err = decodeString(buf, key_next)
			// case ServerKeepAlive:
			// 	prop.item[idx], val_next, err = decodeUint16(buf, key_next)
			// case AuthenticationMethod:
			// 	prop.item[idx], val_next, err = decodeString(buf, key_next)
			// case AuthenticationData:
			// 	prop.item[idx], val_next, err = decodeBytes(buf, key_next)
			// case RequestProblemInformation:
			// 	prop.item[idx], val_next, err = decodeByte(buf, key_next)
			// case WillDelayInterval:
			// 	prop.item[idx], val_next, err = decodeUint32(buf, key_next)
			// case RequestResponseInformation:
			// 	prop.item[idx], val_next, err = decodeByte(buf, key_next)
			// case ResponseInformation:
			// 	prop.item[idx], val_next, err = decodeString(buf, key_next)
			// case ServerReference:
			// 	prop.item[idx], val_next, err = decodeString(buf, key_next)
			// case ReasonString:
			// 	prop.item[idx], val_next, err = decodeString(buf, key_next)
			// case ReceiveMaximum:
			// 	prop.item[idx], val_next, err = decodeUint16(buf, key_next)
			// case TopicAliasMaximum:
			// 	prop.item[idx], val_next, err = decodeUint16(buf, key_next)
			// case TopicAlias:
			// 	prop.item[idx], val_next, err = decodeUint16(buf, key_next)
			// case MaximumQoS:
			// 	prop.item[idx], val_next, err = decodeByte(buf, key_next)
			// case RetainAvailable:
			// 	prop.item[idx], val_next, err = decodeByte(buf, key_next)
			// case UserProperty:
			// 	prop.item[idx], val_next, err = decodeBytes(buf, key_next)
			// case MaximumPacketSize:
			// 	prop.item[idx], val_next, err = decodeUint32(buf, key_next)
			// case WildcardSubscriptionAvailable:
			// 	prop.item[idx], val_next, err = decodeByte(buf, key_next)
			// case SubscriptionIdentifierAvailable:
			// 	prop.item[idx], val_next, err = decodeByte(buf, key_next)
			// case SharedSubscriptionAvailable:
			// 	prop.item[idx], val_next, err = decodeByte(buf, key_next)
			// default:
			// 	err = ErrNotFoundProperties
			// 	return 0, err
		}

	}
	length += len(ret)
	return ret, length, nil
}
