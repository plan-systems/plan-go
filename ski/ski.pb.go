// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ski.proto

/*
	Package ski is a generated protocol buffer package.

	It is generated from these files:
		ski.proto

	It has these top-level messages:
		KeyEntry
		KeyList
*/
package ski

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type KeyEntry struct {
	KeyInfo     uint32 `protobuf:"varint,1,opt,name=key_info,json=keyInfo,proto3" json:"key_info,omitempty"`
	TimeCreated int64  `protobuf:"varint,2,opt,name=time_created,json=timeCreated,proto3" json:"time_created,omitempty"`
	PrivKey     []byte `protobuf:"bytes,3,opt,name=priv_key,json=privKey,proto3" json:"priv_key,omitempty"`
	PubKey      []byte `protobuf:"bytes,4,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty"`
}

func (m *KeyEntry) Reset()                    { *m = KeyEntry{} }
func (m *KeyEntry) String() string            { return proto.CompactTextString(m) }
func (*KeyEntry) ProtoMessage()               {}
func (*KeyEntry) Descriptor() ([]byte, []int) { return fileDescriptorSki, []int{0} }

func (m *KeyEntry) GetKeyInfo() uint32 {
	if m != nil {
		return m.KeyInfo
	}
	return 0
}

func (m *KeyEntry) GetTimeCreated() int64 {
	if m != nil {
		return m.TimeCreated
	}
	return 0
}

func (m *KeyEntry) GetPrivKey() []byte {
	if m != nil {
		return m.PrivKey
	}
	return nil
}

func (m *KeyEntry) GetPubKey() []byte {
	if m != nil {
		return m.PubKey
	}
	return nil
}

// A serialized KeyList is named as ski.KeyListProtobufCodec
type KeyList struct {
	// Any protocol-dependent label (or not), typically employed as a name for this Block instance (i.e. a field-name).
	Label string `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	// key_entry_codec is a "multicodec path" -- aka multistream, see https://github.com/multiformats/multistream
	// It expresses expectations on how the key entries on this KeyList can be used or what they conform to.
	KeysCodec string      `protobuf:"bytes,2,opt,name=keys_codec,json=keysCodec,proto3" json:"keys_codec,omitempty"`
	Keys      []*KeyEntry `protobuf:"bytes,3,rep,name=keys" json:"keys,omitempty"`
}

func (m *KeyList) Reset()                    { *m = KeyList{} }
func (m *KeyList) String() string            { return proto.CompactTextString(m) }
func (*KeyList) ProtoMessage()               {}
func (*KeyList) Descriptor() ([]byte, []int) { return fileDescriptorSki, []int{1} }

func (m *KeyList) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *KeyList) GetKeysCodec() string {
	if m != nil {
		return m.KeysCodec
	}
	return ""
}

func (m *KeyList) GetKeys() []*KeyEntry {
	if m != nil {
		return m.Keys
	}
	return nil
}

func init() {
	proto.RegisterType((*KeyEntry)(nil), "ski.KeyEntry")
	proto.RegisterType((*KeyList)(nil), "ski.KeyList")
}
func (m *KeyEntry) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *KeyEntry) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.KeyInfo != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintSki(dAtA, i, uint64(m.KeyInfo))
	}
	if m.TimeCreated != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintSki(dAtA, i, uint64(m.TimeCreated))
	}
	if len(m.PrivKey) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintSki(dAtA, i, uint64(len(m.PrivKey)))
		i += copy(dAtA[i:], m.PrivKey)
	}
	if len(m.PubKey) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintSki(dAtA, i, uint64(len(m.PubKey)))
		i += copy(dAtA[i:], m.PubKey)
	}
	return i, nil
}

func (m *KeyList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *KeyList) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Label) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintSki(dAtA, i, uint64(len(m.Label)))
		i += copy(dAtA[i:], m.Label)
	}
	if len(m.KeysCodec) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintSki(dAtA, i, uint64(len(m.KeysCodec)))
		i += copy(dAtA[i:], m.KeysCodec)
	}
	if len(m.Keys) > 0 {
		for _, msg := range m.Keys {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintSki(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeVarintSki(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *KeyEntry) Size() (n int) {
	var l int
	_ = l
	if m.KeyInfo != 0 {
		n += 1 + sovSki(uint64(m.KeyInfo))
	}
	if m.TimeCreated != 0 {
		n += 1 + sovSki(uint64(m.TimeCreated))
	}
	l = len(m.PrivKey)
	if l > 0 {
		n += 1 + l + sovSki(uint64(l))
	}
	l = len(m.PubKey)
	if l > 0 {
		n += 1 + l + sovSki(uint64(l))
	}
	return n
}

func (m *KeyList) Size() (n int) {
	var l int
	_ = l
	l = len(m.Label)
	if l > 0 {
		n += 1 + l + sovSki(uint64(l))
	}
	l = len(m.KeysCodec)
	if l > 0 {
		n += 1 + l + sovSki(uint64(l))
	}
	if len(m.Keys) > 0 {
		for _, e := range m.Keys {
			l = e.Size()
			n += 1 + l + sovSki(uint64(l))
		}
	}
	return n
}

func sovSki(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozSki(x uint64) (n int) {
	return sovSki(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *KeyEntry) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSki
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: KeyEntry: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: KeyEntry: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field KeyInfo", wireType)
			}
			m.KeyInfo = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSki
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.KeyInfo |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TimeCreated", wireType)
			}
			m.TimeCreated = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSki
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TimeCreated |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PrivKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSki
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSki
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PrivKey = append(m.PrivKey[:0], dAtA[iNdEx:postIndex]...)
			if m.PrivKey == nil {
				m.PrivKey = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSki
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSki
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubKey = append(m.PubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.PubKey == nil {
				m.PubKey = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSki(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSki
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *KeyList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSki
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: KeyList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: KeyList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Label", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSki
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSki
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Label = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field KeysCodec", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSki
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSki
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.KeysCodec = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Keys", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSki
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSki
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Keys = append(m.Keys, &KeyEntry{})
			if err := m.Keys[len(m.Keys)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSki(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSki
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSki(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSki
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSki
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSki
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthSki
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowSki
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipSki(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthSki = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSki   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("ski.proto", fileDescriptorSki) }

var fileDescriptorSki = []byte{
	// 241 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x8f, 0x41, 0x4a, 0xc3, 0x40,
	0x14, 0x86, 0x1d, 0xa7, 0x36, 0xcd, 0x6b, 0x0b, 0x32, 0x08, 0xc6, 0x85, 0x21, 0xed, 0x2a, 0xab,
	0x2e, 0xf4, 0x06, 0x16, 0x17, 0x12, 0x57, 0x73, 0x81, 0x90, 0xa4, 0xaf, 0x30, 0x4c, 0xcd, 0x84,
	0x99, 0xa9, 0xf0, 0xc0, 0x83, 0x78, 0x24, 0x97, 0x1e, 0x41, 0xe2, 0x45, 0x64, 0x26, 0xb8, 0xfc,
	0xbf, 0x8f, 0xc7, 0xe3, 0x83, 0xd4, 0x69, 0xb5, 0x1b, 0xac, 0xf1, 0x46, 0x70, 0xa7, 0xd5, 0xf6,
	0x03, 0x16, 0x15, 0xd2, 0x73, 0xef, 0x2d, 0x89, 0x3b, 0x58, 0x68, 0xa4, 0x5a, 0xf5, 0x47, 0x93,
	0xb1, 0x82, 0x95, 0x6b, 0x99, 0x68, 0xa4, 0x97, 0xfe, 0x68, 0xc4, 0x06, 0x56, 0x5e, 0xbd, 0x61,
	0xdd, 0x59, 0x6c, 0x3c, 0x1e, 0xb2, 0xcb, 0x82, 0x95, 0x5c, 0x2e, 0x03, 0xdb, 0x4f, 0x28, 0x5c,
	0x0f, 0x56, 0xbd, 0xd7, 0x1a, 0x29, 0xe3, 0x05, 0x2b, 0x57, 0x32, 0x09, 0xbb, 0x42, 0x12, 0xb7,
	0x90, 0x0c, 0xe7, 0x36, 0x9a, 0x59, 0x34, 0xf3, 0xe1, 0xdc, 0x56, 0x48, 0xdb, 0x06, 0x92, 0x0a,
	0xe9, 0x55, 0x39, 0x2f, 0x6e, 0xe0, 0xea, 0xd4, 0xb4, 0x78, 0x8a, 0x9f, 0x53, 0x39, 0x0d, 0x71,
	0x0f, 0xa0, 0x91, 0x5c, 0xdd, 0x99, 0x03, 0x76, 0xf1, 0x6b, 0x2a, 0xd3, 0x40, 0xf6, 0x01, 0x88,
	0x0d, 0xcc, 0xc2, 0xc8, 0x78, 0xc1, 0xcb, 0xe5, 0xc3, 0x7a, 0x17, 0xe2, 0xfe, 0x73, 0x64, 0x54,
	0x4f, 0xd7, 0x5f, 0x63, 0xce, 0xbe, 0xc7, 0x9c, 0xfd, 0x8c, 0x39, 0xfb, 0xfc, 0xcd, 0x2f, 0xda,
	0x79, 0xcc, 0x7f, 0xfc, 0x0b, 0x00, 0x00, 0xff, 0xff, 0xe4, 0x07, 0x6c, 0x6d, 0x0b, 0x01, 0x00,
	0x00,
}
