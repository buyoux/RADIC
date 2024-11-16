// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: video.proto

package videoSearch

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type BiliVideo struct {
	Id       string   `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Title    string   `protobuf:"bytes,2,opt,name=Title,proto3" json:"Title,omitempty"`
	PostTime int64    `protobuf:"varint,3,opt,name=PostTime,proto3" json:"PostTime,omitempty"`
	Author   string   `protobuf:"bytes,4,opt,name=Author,proto3" json:"Author,omitempty"`
	View     int32    `protobuf:"varint,5,opt,name=View,proto3" json:"View,omitempty"`
	Like     int32    `protobuf:"varint,6,opt,name=Like,proto3" json:"Like,omitempty"`
	Coin     int32    `protobuf:"varint,7,opt,name=Coin,proto3" json:"Coin,omitempty"`
	Favorite int32    `protobuf:"varint,8,opt,name=Favorite,proto3" json:"Favorite,omitempty"`
	Share    int32    `protobuf:"varint,9,opt,name=Share,proto3" json:"Share,omitempty"`
	Keywords []string `protobuf:"bytes,10,rep,name=Keywords,proto3" json:"Keywords,omitempty"`
}

func (m *BiliVideo) Reset()         { *m = BiliVideo{} }
func (m *BiliVideo) String() string { return proto.CompactTextString(m) }
func (*BiliVideo) ProtoMessage()    {}
func (*BiliVideo) Descriptor() ([]byte, []int) {
	return fileDescriptor_0ad4ea8866efb1e3, []int{0}
}
func (m *BiliVideo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BiliVideo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BiliVideo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BiliVideo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BiliVideo.Merge(m, src)
}
func (m *BiliVideo) XXX_Size() int {
	return m.Size()
}
func (m *BiliVideo) XXX_DiscardUnknown() {
	xxx_messageInfo_BiliVideo.DiscardUnknown(m)
}

var xxx_messageInfo_BiliVideo proto.InternalMessageInfo

func (m *BiliVideo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *BiliVideo) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *BiliVideo) GetPostTime() int64 {
	if m != nil {
		return m.PostTime
	}
	return 0
}

func (m *BiliVideo) GetAuthor() string {
	if m != nil {
		return m.Author
	}
	return ""
}

func (m *BiliVideo) GetView() int32 {
	if m != nil {
		return m.View
	}
	return 0
}

func (m *BiliVideo) GetLike() int32 {
	if m != nil {
		return m.Like
	}
	return 0
}

func (m *BiliVideo) GetCoin() int32 {
	if m != nil {
		return m.Coin
	}
	return 0
}

func (m *BiliVideo) GetFavorite() int32 {
	if m != nil {
		return m.Favorite
	}
	return 0
}

func (m *BiliVideo) GetShare() int32 {
	if m != nil {
		return m.Share
	}
	return 0
}

func (m *BiliVideo) GetKeywords() []string {
	if m != nil {
		return m.Keywords
	}
	return nil
}

func init() {
	proto.RegisterType((*BiliVideo)(nil), "videoSearch.BiliVideo")
}

func init() { proto.RegisterFile("video.proto", fileDescriptor_0ad4ea8866efb1e3) }

var fileDescriptor_0ad4ea8866efb1e3 = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x3c, 0x90, 0xb1, 0x4a, 0xc4, 0x40,
	0x10, 0x86, 0xb3, 0xc9, 0x25, 0x5e, 0x56, 0xb0, 0x58, 0x44, 0x06, 0x8b, 0x25, 0x58, 0xa5, 0xb2,
	0xf1, 0x09, 0x3c, 0x41, 0x38, 0xb4, 0x90, 0xdc, 0x71, 0x7d, 0x74, 0x07, 0x32, 0x78, 0xba, 0xb2,
	0xb7, 0xde, 0xe1, 0x5b, 0xf8, 0x58, 0x96, 0x29, 0x2d, 0x25, 0x79, 0x00, 0x5f, 0x41, 0x66, 0x57,
	0xd2, 0xfd, 0xdf, 0x37, 0xfc, 0xcc, 0x30, 0xf2, 0x78, 0x4f, 0x06, 0xed, 0xe5, 0x9b, 0xb3, 0xde,
	0xaa, 0x08, 0x2b, 0x6c, 0xdd, 0x53, 0x77, 0xf1, 0x2b, 0x64, 0xb9, 0xa0, 0x2d, 0x6d, 0xd8, 0xa9,
	0x13, 0x99, 0x2e, 0x0d, 0x88, 0x4a, 0xd4, 0x65, 0x93, 0x2e, 0x8d, 0x3a, 0x95, 0xf9, 0x9a, 0xfc,
	0x16, 0x21, 0x0d, 0x2a, 0x82, 0x3a, 0x97, 0xf3, 0x07, 0xbb, 0xf3, 0x6b, 0x7a, 0x41, 0xc8, 0x2a,
	0x51, 0x67, 0xcd, 0xc4, 0xea, 0x4c, 0x16, 0xd7, 0xef, 0xbe, 0xb3, 0x0e, 0x66, 0xa1, 0xf2, 0x4f,
	0x4a, 0xc9, 0xd9, 0x86, 0xf0, 0x00, 0x79, 0x25, 0xea, 0xbc, 0x09, 0x99, 0xdd, 0x3d, 0x3d, 0x23,
	0x14, 0xd1, 0x71, 0x66, 0x77, 0x63, 0xe9, 0x15, 0x8e, 0xa2, 0xe3, 0xcc, 0xfb, 0x6e, 0xdb, 0xbd,
	0x75, 0xe4, 0x11, 0xe6, 0xc1, 0x4f, 0xcc, 0x17, 0xae, 0xba, 0xd6, 0x21, 0x94, 0x61, 0x10, 0x81,
	0x1b, 0x77, 0xf8, 0x71, 0xb0, 0xce, 0xec, 0x40, 0x56, 0x59, 0x5d, 0x36, 0x13, 0x2f, 0xe0, 0x6b,
	0xd0, 0xa2, 0x1f, 0xb4, 0xf8, 0x19, 0xb4, 0xf8, 0x1c, 0x75, 0xd2, 0x8f, 0x3a, 0xf9, 0x1e, 0x75,
	0xf2, 0x58, 0x84, 0xff, 0x5c, 0xfd, 0x05, 0x00, 0x00, 0xff, 0xff, 0xb1, 0x38, 0xc2, 0xa0, 0x2e,
	0x01, 0x00, 0x00,
}

func (m *BiliVideo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BiliVideo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BiliVideo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Keywords) > 0 {
		for iNdEx := len(m.Keywords) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Keywords[iNdEx])
			copy(dAtA[i:], m.Keywords[iNdEx])
			i = encodeVarintVideo(dAtA, i, uint64(len(m.Keywords[iNdEx])))
			i--
			dAtA[i] = 0x52
		}
	}
	if m.Share != 0 {
		i = encodeVarintVideo(dAtA, i, uint64(m.Share))
		i--
		dAtA[i] = 0x48
	}
	if m.Favorite != 0 {
		i = encodeVarintVideo(dAtA, i, uint64(m.Favorite))
		i--
		dAtA[i] = 0x40
	}
	if m.Coin != 0 {
		i = encodeVarintVideo(dAtA, i, uint64(m.Coin))
		i--
		dAtA[i] = 0x38
	}
	if m.Like != 0 {
		i = encodeVarintVideo(dAtA, i, uint64(m.Like))
		i--
		dAtA[i] = 0x30
	}
	if m.View != 0 {
		i = encodeVarintVideo(dAtA, i, uint64(m.View))
		i--
		dAtA[i] = 0x28
	}
	if len(m.Author) > 0 {
		i -= len(m.Author)
		copy(dAtA[i:], m.Author)
		i = encodeVarintVideo(dAtA, i, uint64(len(m.Author)))
		i--
		dAtA[i] = 0x22
	}
	if m.PostTime != 0 {
		i = encodeVarintVideo(dAtA, i, uint64(m.PostTime))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintVideo(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintVideo(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintVideo(dAtA []byte, offset int, v uint64) int {
	offset -= sovVideo(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BiliVideo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovVideo(uint64(l))
	}
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovVideo(uint64(l))
	}
	if m.PostTime != 0 {
		n += 1 + sovVideo(uint64(m.PostTime))
	}
	l = len(m.Author)
	if l > 0 {
		n += 1 + l + sovVideo(uint64(l))
	}
	if m.View != 0 {
		n += 1 + sovVideo(uint64(m.View))
	}
	if m.Like != 0 {
		n += 1 + sovVideo(uint64(m.Like))
	}
	if m.Coin != 0 {
		n += 1 + sovVideo(uint64(m.Coin))
	}
	if m.Favorite != 0 {
		n += 1 + sovVideo(uint64(m.Favorite))
	}
	if m.Share != 0 {
		n += 1 + sovVideo(uint64(m.Share))
	}
	if len(m.Keywords) > 0 {
		for _, s := range m.Keywords {
			l = len(s)
			n += 1 + l + sovVideo(uint64(l))
		}
	}
	return n
}

func sovVideo(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozVideo(x uint64) (n int) {
	return sovVideo(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BiliVideo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVideo
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BiliVideo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BiliVideo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthVideo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVideo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthVideo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVideo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PostTime", wireType)
			}
			m.PostTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PostTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Author", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthVideo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVideo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Author = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field View", wireType)
			}
			m.View = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.View |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Like", wireType)
			}
			m.Like = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Like |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Coin", wireType)
			}
			m.Coin = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Coin |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Favorite", wireType)
			}
			m.Favorite = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Favorite |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Share", wireType)
			}
			m.Share = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Share |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Keywords", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthVideo
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVideo
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Keywords = append(m.Keywords, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVideo(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthVideo
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
func skipVideo(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowVideo
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
					return 0, ErrIntOverflowVideo
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowVideo
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
			if length < 0 {
				return 0, ErrInvalidLengthVideo
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupVideo
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthVideo
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthVideo        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowVideo          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupVideo = fmt.Errorf("proto: unexpected end of group")
)
