package grid

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *HandlerID) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 uint8
		zb0001, err = dc.ReadUint8()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = HandlerID(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z HandlerID) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteUint8(uint8(z))
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z HandlerID) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendUint8(o, uint8(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *HandlerID) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 uint8
		zb0001, bts, err = msgp.ReadUint8Bytes(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = HandlerID(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z HandlerID) Msgsize() (s int) {
	s = msgp.Uint8Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Op) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 uint8
		zb0001, err = dc.ReadUint8()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = Op(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Op) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteUint8(uint8(z))
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Op) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendUint8(o, uint8(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Op) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 uint8
		zb0001, bts, err = msgp.ReadUint8Bytes(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		(*z) = Op(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Op) Msgsize() (s int) {
	s = msgp.Uint8Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *connectReq) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "ID":
			err = dc.ReadExactBytes((z.ID)[:])
			if err != nil {
				err = msgp.WrapError(err, "ID")
				return
			}
		case "Host":
			z.Host, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Host")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *connectReq) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "ID"
	err = en.Append(0x82, 0xa2, 0x49, 0x44)
	if err != nil {
		return
	}
	err = en.WriteBytes((z.ID)[:])
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	// write "Host"
	err = en.Append(0xa4, 0x48, 0x6f, 0x73, 0x74)
	if err != nil {
		return
	}
	err = en.WriteString(z.Host)
	if err != nil {
		err = msgp.WrapError(err, "Host")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *connectReq) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "ID"
	o = append(o, 0x82, 0xa2, 0x49, 0x44)
	o = msgp.AppendBytes(o, (z.ID)[:])
	// string "Host"
	o = append(o, 0xa4, 0x48, 0x6f, 0x73, 0x74)
	o = msgp.AppendString(o, z.Host)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *connectReq) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "ID":
			bts, err = msgp.ReadExactBytes(bts, (z.ID)[:])
			if err != nil {
				err = msgp.WrapError(err, "ID")
				return
			}
		case "Host":
			z.Host, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Host")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *connectReq) Msgsize() (s int) {
	s = 1 + 3 + msgp.ArrayHeaderSize + (16 * (msgp.ByteSize)) + 5 + msgp.StringPrefixSize + len(z.Host)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *connectResp) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "ID":
			err = dc.ReadExactBytes((z.ID)[:])
			if err != nil {
				err = msgp.WrapError(err, "ID")
				return
			}
		case "Accepted":
			z.Accepted, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "Accepted")
				return
			}
		case "RejectedReason":
			z.RejectedReason, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "RejectedReason")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *connectResp) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "ID"
	err = en.Append(0x83, 0xa2, 0x49, 0x44)
	if err != nil {
		return
	}
	err = en.WriteBytes((z.ID)[:])
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	// write "Accepted"
	err = en.Append(0xa8, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64)
	if err != nil {
		return
	}
	err = en.WriteBool(z.Accepted)
	if err != nil {
		err = msgp.WrapError(err, "Accepted")
		return
	}
	// write "RejectedReason"
	err = en.Append(0xae, 0x52, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x65, 0x64, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteString(z.RejectedReason)
	if err != nil {
		err = msgp.WrapError(err, "RejectedReason")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *connectResp) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "ID"
	o = append(o, 0x83, 0xa2, 0x49, 0x44)
	o = msgp.AppendBytes(o, (z.ID)[:])
	// string "Accepted"
	o = append(o, 0xa8, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64)
	o = msgp.AppendBool(o, z.Accepted)
	// string "RejectedReason"
	o = append(o, 0xae, 0x52, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x65, 0x64, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.RejectedReason)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *connectResp) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "ID":
			bts, err = msgp.ReadExactBytes(bts, (z.ID)[:])
			if err != nil {
				err = msgp.WrapError(err, "ID")
				return
			}
		case "Accepted":
			z.Accepted, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Accepted")
				return
			}
		case "RejectedReason":
			z.RejectedReason, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "RejectedReason")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *connectResp) Msgsize() (s int) {
	s = 1 + 3 + msgp.ArrayHeaderSize + (16 * (msgp.ByteSize)) + 9 + msgp.BoolSize + 15 + msgp.StringPrefixSize + len(z.RejectedReason)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *message) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 6 {
		err = msgp.ArrayError{Wanted: 6, Got: zb0001}
		return
	}
	z.MuxID, err = dc.ReadUint64()
	if err != nil {
		err = msgp.WrapError(err, "MuxID")
		return
	}
	z.Seq, err = dc.ReadUint32()
	if err != nil {
		err = msgp.WrapError(err, "Seq")
		return
	}
	{
		var zb0002 uint8
		zb0002, err = dc.ReadUint8()
		if err != nil {
			err = msgp.WrapError(err, "Handler")
			return
		}
		z.Handler = HandlerID(zb0002)
	}
	{
		var zb0003 uint8
		zb0003, err = dc.ReadUint8()
		if err != nil {
			err = msgp.WrapError(err, "Op")
			return
		}
		z.Op = Op(zb0003)
	}
	z.Flags, err = dc.ReadUint8()
	if err != nil {
		err = msgp.WrapError(err, "Flags")
		return
	}
	z.Payload, err = dc.ReadBytes(z.Payload)
	if err != nil {
		err = msgp.WrapError(err, "Payload")
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *message) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 6
	err = en.Append(0x96)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.MuxID)
	if err != nil {
		err = msgp.WrapError(err, "MuxID")
		return
	}
	err = en.WriteUint32(z.Seq)
	if err != nil {
		err = msgp.WrapError(err, "Seq")
		return
	}
	err = en.WriteUint8(uint8(z.Handler))
	if err != nil {
		err = msgp.WrapError(err, "Handler")
		return
	}
	err = en.WriteUint8(uint8(z.Op))
	if err != nil {
		err = msgp.WrapError(err, "Op")
		return
	}
	err = en.WriteUint8(z.Flags)
	if err != nil {
		err = msgp.WrapError(err, "Flags")
		return
	}
	err = en.WriteBytes(z.Payload)
	if err != nil {
		err = msgp.WrapError(err, "Payload")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *message) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 6
	o = append(o, 0x96)
	o = msgp.AppendUint64(o, z.MuxID)
	o = msgp.AppendUint32(o, z.Seq)
	o = msgp.AppendUint8(o, uint8(z.Handler))
	o = msgp.AppendUint8(o, uint8(z.Op))
	o = msgp.AppendUint8(o, z.Flags)
	o = msgp.AppendBytes(o, z.Payload)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *message) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 6 {
		err = msgp.ArrayError{Wanted: 6, Got: zb0001}
		return
	}
	z.MuxID, bts, err = msgp.ReadUint64Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "MuxID")
		return
	}
	z.Seq, bts, err = msgp.ReadUint32Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Seq")
		return
	}
	{
		var zb0002 uint8
		zb0002, bts, err = msgp.ReadUint8Bytes(bts)
		if err != nil {
			err = msgp.WrapError(err, "Handler")
			return
		}
		z.Handler = HandlerID(zb0002)
	}
	{
		var zb0003 uint8
		zb0003, bts, err = msgp.ReadUint8Bytes(bts)
		if err != nil {
			err = msgp.WrapError(err, "Op")
			return
		}
		z.Op = Op(zb0003)
	}
	z.Flags, bts, err = msgp.ReadUint8Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Flags")
		return
	}
	z.Payload, bts, err = msgp.ReadBytesBytes(bts, z.Payload)
	if err != nil {
		err = msgp.WrapError(err, "Payload")
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *message) Msgsize() (s int) {
	s = 1 + msgp.Uint64Size + msgp.Uint32Size + msgp.Uint8Size + msgp.Uint8Size + msgp.Uint8Size + msgp.BytesPrefixSize + len(z.Payload)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *muxConnectError) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Error":
			z.Error, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Error")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z muxConnectError) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Error"
	err = en.Append(0x81, 0xa5, 0x45, 0x72, 0x72, 0x6f, 0x72)
	if err != nil {
		return
	}
	err = en.WriteString(z.Error)
	if err != nil {
		err = msgp.WrapError(err, "Error")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z muxConnectError) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Error"
	o = append(o, 0x81, 0xa5, 0x45, 0x72, 0x72, 0x6f, 0x72)
	o = msgp.AppendString(o, z.Error)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *muxConnectError) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Error":
			z.Error, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Error")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z muxConnectError) Msgsize() (s int) {
	s = 1 + 6 + msgp.StringPrefixSize + len(z.Error)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *pongMsg) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "f":
			z.NotFound, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "NotFound")
				return
			}
		case "e":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					err = msgp.WrapError(err, "Err")
					return
				}
				z.Err = nil
			} else {
				if z.Err == nil {
					z.Err = new(string)
				}
				*z.Err, err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "Err")
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *pongMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "f"
	err = en.Append(0x82, 0xa1, 0x66)
	if err != nil {
		return
	}
	err = en.WriteBool(z.NotFound)
	if err != nil {
		err = msgp.WrapError(err, "NotFound")
		return
	}
	// write "e"
	err = en.Append(0xa1, 0x65)
	if err != nil {
		return
	}
	if z.Err == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = en.WriteString(*z.Err)
		if err != nil {
			err = msgp.WrapError(err, "Err")
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *pongMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "f"
	o = append(o, 0x82, 0xa1, 0x66)
	o = msgp.AppendBool(o, z.NotFound)
	// string "e"
	o = append(o, 0xa1, 0x65)
	if z.Err == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendString(o, *z.Err)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *pongMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "f":
			z.NotFound, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "NotFound")
				return
			}
		case "e":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Err = nil
			} else {
				if z.Err == nil {
					z.Err = new(string)
				}
				*z.Err, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Err")
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *pongMsg) Msgsize() (s int) {
	s = 1 + 2 + msgp.BoolSize + 2
	if z.Err == nil {
		s += msgp.NilSize
	} else {
		s += msgp.StringPrefixSize + len(*z.Err)
	}
	return
}
