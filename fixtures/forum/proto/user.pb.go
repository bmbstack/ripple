// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto/user.proto

package proto

import (
	context "context"
	errors "errors"
	fmt "fmt"
	ripple "github.com/bmbstack/ripple"
	helper "github.com/bmbstack/ripple/helper"
	proto "github.com/golang/protobuf/proto"
	constant "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	client1 "github.com/rpcxio/rpcx-nacos/client"
	client "github.com/smallnest/rpcx/client"
	protocol "github.com/smallnest/rpcx/protocol"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GetInfoReq struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetInfoReq) Reset()         { *m = GetInfoReq{} }
func (m *GetInfoReq) String() string { return proto.CompactTextString(m) }
func (*GetInfoReq) ProtoMessage()    {}
func (*GetInfoReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_d570e3e37e5899c5, []int{0}
}
func (m *GetInfoReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GetInfoReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GetInfoReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GetInfoReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetInfoReq.Merge(m, src)
}
func (m *GetInfoReq) XXX_Size() int {
	return m.Size()
}
func (m *GetInfoReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetInfoReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetInfoReq proto.InternalMessageInfo

func (m *GetInfoReq) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type GetInfoReply struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetInfoReply) Reset()         { *m = GetInfoReply{} }
func (m *GetInfoReply) String() string { return proto.CompactTextString(m) }
func (*GetInfoReply) ProtoMessage()    {}
func (*GetInfoReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_d570e3e37e5899c5, []int{1}
}
func (m *GetInfoReply) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GetInfoReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GetInfoReply.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GetInfoReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetInfoReply.Merge(m, src)
}
func (m *GetInfoReply) XXX_Size() int {
	return m.Size()
}
func (m *GetInfoReply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetInfoReply.DiscardUnknown(m)
}

var xxx_messageInfo_GetInfoReply proto.InternalMessageInfo

func (m *GetInfoReply) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*GetInfoReq)(nil), "proto.GetInfoReq")
	proto.RegisterType((*GetInfoReply)(nil), "proto.GetInfoReply")
}

func init() { proto.RegisterFile("proto/user.proto", fileDescriptor_d570e3e37e5899c5) }

var fileDescriptor_d570e3e37e5899c5 = []byte{
	// 149 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x28, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x2f, 0x2d, 0x4e, 0x2d, 0xd2, 0x03, 0x33, 0x85, 0x58, 0xc1, 0x94, 0x92, 0x0c, 0x17,
	0x97, 0x7b, 0x6a, 0x89, 0x67, 0x5e, 0x5a, 0x7e, 0x50, 0x6a, 0xa1, 0x10, 0x1f, 0x17, 0x53, 0x66,
	0x8a, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x4b, 0x10, 0x53, 0x66, 0x8a, 0x92, 0x12, 0x17, 0x0f, 0x5c,
	0xb6, 0x20, 0xa7, 0x52, 0x48, 0x88, 0x8b, 0x25, 0x2f, 0x31, 0x37, 0x15, 0xac, 0x82, 0x33, 0x08,
	0xcc, 0x36, 0xb2, 0xe6, 0x62, 0x09, 0x2d, 0x4e, 0x2d, 0x12, 0x32, 0xe6, 0x62, 0x87, 0xaa, 0x15,
	0x12, 0x84, 0xd8, 0xa1, 0x87, 0x30, 0x59, 0x4a, 0x18, 0x5d, 0xa8, 0x20, 0xa7, 0x52, 0x89, 0xc1,
	0x49, 0xe0, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf1,
	0x58, 0x8e, 0x21, 0x89, 0x0d, 0xac, 0xce, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x8b, 0xcb, 0x46,
	0xc9, 0xb2, 0x00, 0x00, 0x00,
}

func (m *GetInfoReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetInfoReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GetInfoReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Id != 0 {
		i = encodeVarintUser(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *GetInfoReply) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetInfoReply) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GetInfoReply) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintUser(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintUser(dAtA []byte, offset int, v uint64) int {
	offset -= sovUser(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

// This following code was generated by ripple
// Gernerated from proto/user.proto

const ServiceNameOfUser = "UserRpc"

//================== interface ===================
type UserInterface interface {

	// UserInterface can be used for interface verification.
	// GetInfo is server rpc method as defined
	GetInfo(ctx context.Context, req *GetInfoReq, reply *GetInfoReply) (err error)
}

//================== server implement demo ===================
//ripple.Default().RegisterRpc("User", &UserRpcDemo{}, "")
type UserRpcDemo struct{}

// GetInfo is server rpc method as defined
func (s *UserRpcDemo) GetInfo(ctx context.Context, req *GetInfoReq, reply *GetInfoReply) (err error) {
	// TODO: add business logics
	*reply = GetInfoReply{}
	return nil
}

//================== client stub ===================
// newXClientForUser creates a XClient.
// You can configure this client with more options such as etcd registry, serialize type, select algorithm and fail mode.
func newXClientForUser(group, cluster string) (client.XClient, error) {
	config := ripple.GetBaseConfig()
	if helper.IsEmpty(config.Nacos) {
		return nil, errors.New("yaml nacos config is null")
	}
	clientConfig := constant.ClientConfig{
		TimeoutMs:            10 * 1000,
		ListenInterval:       30 * 1000,
		BeatInterval:         5 * 1000,
		NamespaceId:          config.Nacos.NamespaceId,
		CacheDir:             config.Nacos.CacheDir,
		LogDir:               config.Nacos.LogDir,
		UpdateThreadNum:      20,
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: true,
	}

	serverConfig := []constant.ServerConfig{{
		IpAddr: config.Nacos.Host,
		Port:   config.Nacos.Port,
	}}

	d, err := client1.NewNacosDiscovery(ServiceNameOfUser, cluster, group, clientConfig, serverConfig)
	if err != nil {
		return nil, err
	}

	opt := client.DefaultOption
	opt.SerializeType = protocol.ProtoBuffer

	xclient := client.NewXClient(ServiceNameOfUser, client.Failtry, client.RoundRobin, d, opt)

	return xclient, nil
}

// User is a client wrapped XClient.
type UserClient struct {
	xclient client.XClient
}

// NewUserClient wraps a XClient as UserClient.
// You can pass a shared XClient object created by NewXClientForUser.
func NewUserClient(group, cluster string) *UserClient {
	xc, err := newXClientForUser(group, cluster)
	if err != nil {
		fmt.Println(fmt.Sprintf("Create rpcx client err: %!s(MISSING)", err.Error()))
		return &UserClient{}
	}
	return &UserClient{xclient: xc}
}

// GetInfo is client rpc method as defined
func (c *UserClient) GetInfo(ctx context.Context, req *GetInfoReq) (reply *GetInfoReply, err error) {
	reply = &GetInfoReply{}
	err = c.xclient.Call(ctx, "GetInfo", req, reply)
	return reply, err
}

// ======================================================
func (m *GetInfoReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovUser(uint64(m.Id))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *GetInfoReply) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovUser(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovUser(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozUser(x uint64) (n int) {
	return sovUser(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GetInfoReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowUser
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
			return fmt.Errorf("proto: GetInfoReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetInfoReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipUser(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthUser
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GetInfoReply) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowUser
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
			return fmt.Errorf("proto: GetInfoReply: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetInfoReply: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUser
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
				return ErrInvalidLengthUser
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUser
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipUser(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthUser
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipUser(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowUser
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
					return 0, ErrIntOverflowUser
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
					return 0, ErrIntOverflowUser
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
				return 0, ErrInvalidLengthUser
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupUser
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthUser
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthUser        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowUser          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupUser = fmt.Errorf("proto: unexpected end of group")
)
