// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: plan.proto

/*
	Package plan is a generated protocol buffer package.

	It is generated from these files:
		plan.proto

	It has these top-level messages:
		Block
*/
package plan

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

type Block struct {
	// This is an optional, human-understandable name/label for this Block (i.e. a field-name).
	// In general, a Block's label conforms to the context/protocol its being used in or with (if applicable).
	Label string `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	// Similar to a MIME type, this field self-describes the binary format of Block.content.
	// Anyone handed this Block only needs to look a Block's codec in order to accurately process/deserialize its content.
	// This field is a "multicodec path" -- see: https://github.com/multiformats/multistream
	Codec string `protobuf:"bytes,2,opt,name=codec,proto3" json:"codec,omitempty"`
	// This is alternative numerical repesentation of Block.codec
	// For a list of codes: https://github.com/multiformats/multicodec/blob/master/table.csv
	CodecCode uint32 `protobuf:"varint,3,opt,name=codec_code,json=codecCode,proto3" json:"codec_code,omitempty"`
	// Data of any kind, serialized in conformance with the accompanying codec descriptor (above).
	Content []byte `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	// In addition to it's label and content, a Block can contain nested "sub" blocks.  A Block's sub blocks
	//    can be interpreted or employed any way a client or protocol sees fit.
	Subs []*Block `protobuf:"bytes,5,rep,name=subs" json:"subs,omitempty"`
}

func (m *Block) Reset()                    { *m = Block{} }
func (m *Block) String() string            { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()               {}
func (*Block) Descriptor() ([]byte, []int) { return fileDescriptorPlan, []int{0} }

func (m *Block) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *Block) GetCodec() string {
	if m != nil {
		return m.Codec
	}
	return ""
}

func (m *Block) GetCodecCode() uint32 {
	if m != nil {
		return m.CodecCode
	}
	return 0
}

func (m *Block) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *Block) GetSubs() []*Block {
	if m != nil {
		return m.Subs
	}
	return nil
}

func init() {
	proto.RegisterType((*Block)(nil), "plan.Block")
}
func (m *Block) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Block) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Label) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintPlan(dAtA, i, uint64(len(m.Label)))
		i += copy(dAtA[i:], m.Label)
	}
	if len(m.Codec) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintPlan(dAtA, i, uint64(len(m.Codec)))
		i += copy(dAtA[i:], m.Codec)
	}
	if m.CodecCode != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintPlan(dAtA, i, uint64(m.CodecCode))
	}
	if len(m.Content) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintPlan(dAtA, i, uint64(len(m.Content)))
		i += copy(dAtA[i:], m.Content)
	}
	if len(m.Subs) > 0 {
		for _, msg := range m.Subs {
			dAtA[i] = 0x2a
			i++
			i = encodeVarintPlan(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeVarintPlan(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Block) Size() (n int) {
	var l int
	_ = l
	l = len(m.Label)
	if l > 0 {
		n += 1 + l + sovPlan(uint64(l))
	}
	l = len(m.Codec)
	if l > 0 {
		n += 1 + l + sovPlan(uint64(l))
	}
	if m.CodecCode != 0 {
		n += 1 + sovPlan(uint64(m.CodecCode))
	}
	l = len(m.Content)
	if l > 0 {
		n += 1 + l + sovPlan(uint64(l))
	}
	if len(m.Subs) > 0 {
		for _, e := range m.Subs {
			l = e.Size()
			n += 1 + l + sovPlan(uint64(l))
		}
	}
	return n
}

func sovPlan(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozPlan(x uint64) (n int) {
	return sovPlan(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Block) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPlan
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
			return fmt.Errorf("proto: Block: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Block: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Label", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
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
				return ErrInvalidLengthPlan
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Label = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Codec", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
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
				return ErrInvalidLengthPlan
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Codec = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CodecCode", wireType)
			}
			m.CodecCode = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CodecCode |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Content", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
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
				return ErrInvalidLengthPlan
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Content = append(m.Content[:0], dAtA[iNdEx:postIndex]...)
			if m.Content == nil {
				m.Content = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
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
				return ErrInvalidLengthPlan
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Subs = append(m.Subs, &Block{})
			if err := m.Subs[len(m.Subs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPlan(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPlan
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
func skipPlan(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPlan
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
					return 0, ErrIntOverflowPlan
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
					return 0, ErrIntOverflowPlan
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
				return 0, ErrInvalidLengthPlan
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowPlan
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
				next, err := skipPlan(dAtA[start:])
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
	ErrInvalidLengthPlan = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPlan   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("plan.proto", fileDescriptorPlan) }

var fileDescriptorPlan = []byte{
	// 171 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0xc8, 0x49, 0xcc,
	0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0xb1, 0x95, 0x7a, 0x19, 0xb9, 0x58, 0x9d,
	0x72, 0xf2, 0x93, 0xb3, 0x85, 0x44, 0xb8, 0x58, 0x73, 0x12, 0x93, 0x52, 0x73, 0x24, 0x18, 0x15,
	0x18, 0x35, 0x38, 0x83, 0x20, 0x1c, 0x90, 0x68, 0x72, 0x7e, 0x4a, 0x6a, 0xb2, 0x04, 0x13, 0x44,
	0x14, 0xcc, 0x11, 0x92, 0xe5, 0xe2, 0x02, 0x33, 0xe2, 0x41, 0xa4, 0x04, 0xb3, 0x02, 0xa3, 0x06,
	0x6f, 0x10, 0x27, 0x58, 0xc4, 0x39, 0x3f, 0x25, 0x55, 0x48, 0x82, 0x8b, 0x3d, 0x39, 0x3f, 0xaf,
	0x24, 0x35, 0xaf, 0x44, 0x82, 0x45, 0x81, 0x51, 0x83, 0x27, 0x08, 0xc6, 0x15, 0x92, 0xe7, 0x62,
	0x29, 0x2e, 0x4d, 0x2a, 0x96, 0x60, 0x55, 0x60, 0xd6, 0xe0, 0x36, 0xe2, 0xd6, 0x03, 0xbb, 0x07,
	0x6c, 0x7f, 0x10, 0x58, 0xc2, 0x49, 0xe0, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f,
	0x3c, 0x92, 0x63, 0x9c, 0xf1, 0x58, 0x8e, 0x21, 0x89, 0x0d, 0xec, 0x5c, 0x63, 0x40, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xd9, 0x98, 0xa3, 0xa6, 0xbc, 0x00, 0x00, 0x00,
}