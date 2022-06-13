// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fixtures/forum/proto/student.proto

// @RippleRpc
// @NacosGroup DEFAULT_GROUP
// @NacosCluster orderserver

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

type HelloReq struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloReq) Reset()         { *m = HelloReq{} }
func (m *HelloReq) String() string { return proto.CompactTextString(m) }
func (*HelloReq) ProtoMessage()    {}
func (*HelloReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_7a63a0c0de64b8ab, []int{0}
}
func (m *HelloReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HelloReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HelloReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HelloReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloReq.Merge(m, src)
}
func (m *HelloReq) XXX_Size() int {
	return m.Size()
}
func (m *HelloReq) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloReq.DiscardUnknown(m)
}

var xxx_messageInfo_HelloReq proto.InternalMessageInfo

func (m *HelloReq) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type HelloReply struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloReply) Reset()         { *m = HelloReply{} }
func (m *HelloReply) String() string { return proto.CompactTextString(m) }
func (*HelloReply) ProtoMessage()    {}
func (*HelloReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_7a63a0c0de64b8ab, []int{1}
}
func (m *HelloReply) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HelloReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HelloReply.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HelloReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloReply.Merge(m, src)
}
func (m *HelloReply) XXX_Size() int {
	return m.Size()
}
func (m *HelloReply) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloReply.DiscardUnknown(m)
}

var xxx_messageInfo_HelloReply proto.InternalMessageInfo

func (m *HelloReply) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*HelloReq)(nil), "proto.HelloReq")
	proto.RegisterType((*HelloReply)(nil), "proto.HelloReply")
}

func init() {
	proto.RegisterFile("fixtures/forum/proto/student.proto", fileDescriptor_7a63a0c0de64b8ab)
}

var fileDescriptor_7a63a0c0de64b8ab = []byte{
	// 164 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4a, 0xcb, 0xac, 0x28,
	0x29, 0x2d, 0x4a, 0x2d, 0xd6, 0x4f, 0xcb, 0x2f, 0x2a, 0xcd, 0xd5, 0x2f, 0x28, 0xca, 0x2f, 0xc9,
	0xd7, 0x2f, 0x2e, 0x29, 0x4d, 0x49, 0xcd, 0x2b, 0xd1, 0x03, 0xf3, 0x84, 0x58, 0xc1, 0x94, 0x92,
	0x14, 0x17, 0x87, 0x47, 0x6a, 0x4e, 0x4e, 0x7e, 0x50, 0x6a, 0xa1, 0x10, 0x1f, 0x17, 0x53, 0x66,
	0x8a, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x4b, 0x10, 0x53, 0x66, 0x8a, 0x92, 0x02, 0x17, 0x17, 0x54,
	0xae, 0x20, 0xa7, 0x52, 0x48, 0x88, 0x8b, 0x25, 0x2f, 0x31, 0x37, 0x15, 0x2c, 0xcf, 0x19, 0x04,
	0x66, 0x1b, 0x59, 0x70, 0xb1, 0x07, 0x43, 0x4c, 0x15, 0xd2, 0xe5, 0x62, 0x05, 0x2b, 0x16, 0xe2,
	0x87, 0x58, 0xa0, 0x07, 0x33, 0x56, 0x4a, 0x10, 0x55, 0xa0, 0x20, 0xa7, 0x52, 0x89, 0xc1, 0x49,
	0xe0, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf1, 0x58,
	0x8e, 0x21, 0x89, 0x0d, 0xac, 0xca, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x43, 0x92, 0xdb, 0x5e,
	0xbd, 0x00, 0x00, 0x00,
}

func (m *HelloReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HelloReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HelloReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Id != 0 {
		i = encodeVarintStudent(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *HelloReply) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HelloReply) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HelloReply) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
		i = encodeVarintStudent(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintStudent(dAtA []byte, offset int, v uint64) int {
	offset -= sovStudent(v)
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
// Gernerated from fixtures/forum/proto/student.proto

const ServiceNameOfStudent = "StudentRpc"

//================== interface ===================
type StudentInterface interface {

	// StudentInterface can be used for interface verification.
	// Hello is server rpc method as defined
	Hello(ctx context.Context, req *HelloReq, reply *HelloReply) (err error)
}

//================== server implement demo ===================
//ripple.Default().RegisterRpc("User", &UserRpcDemo{}, "")
type StudentRpcDemo struct{}

func (this *StudentRpcDemo) Hello(ctx context.Context, req *HelloReq, reply *HelloReply) (err error) {
	// TODO: add business logics
	*reply = HelloReply{}
	return nil
}

//================== client stub ===================
// newXClientForStudent creates a XClient.
// You can configure this client with more options such as etcd registry, serialize type, select algorithm and fail mode.
func newXClientForStudent() (client.XClient, error) {
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

	d, err := client1.NewNacosDiscovery(ServiceNameOfStudent, "orderserver", "DEFAULT_GROUP", clientConfig, serverConfig)
	if err != nil {
		return nil, err
	}

	opt := client.DefaultOption
	opt.SerializeType = protocol.ProtoBuffer

	xclient := client.NewXClient(ServiceNameOfStudent, client.Failtry, client.RoundRobin, d, opt)

	return xclient, nil
}

// Student is a client wrapped XClient.
type StudentClient struct {
	xclient client.XClient
}

// NewStudentClient wraps a XClient as StudentClient.
// You can pass a shared XClient object created by NewXClientForStudent.
func NewStudentClient() *StudentClient {
	xc, err := newXClientForStudent()
	if err != nil {
		fmt.Println(fmt.Sprintf("Create rpcx client err: orderserver", err.Error()))
		return &StudentClient{}
	}
	return &StudentClient{xclient: xc}
}

// Hello is client rpc method as defined
func (c *StudentClient) Hello(ctx context.Context, req *HelloReq) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	err = c.xclient.Call(ctx, "Hello", req, reply)
	return reply, err
}

// ======================================================
func (m *HelloReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovStudent(uint64(m.Id))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *HelloReply) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovStudent(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovStudent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozStudent(x uint64) (n int) {
	return sovStudent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *HelloReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStudent
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
			return fmt.Errorf("proto: HelloReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HelloReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStudent
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
			skippy, err := skipStudent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthStudent
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
func (m *HelloReply) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStudent
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
			return fmt.Errorf("proto: HelloReply: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HelloReply: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStudent
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
				return ErrInvalidLengthStudent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStudent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStudent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthStudent
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
func skipStudent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStudent
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
					return 0, ErrIntOverflowStudent
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
					return 0, ErrIntOverflowStudent
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
				return 0, ErrInvalidLengthStudent
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupStudent
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthStudent
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthStudent        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStudent          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupStudent = fmt.Errorf("proto: unexpected end of group")
)
