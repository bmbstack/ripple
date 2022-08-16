// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fixtures/forum/proto/teacher.proto

// @RippleRpc
// @NacosGroup DEFAULT_GROUP
// @NacosCluster ripple

package proto

import (
	context "context"
	errors "errors"
	fmt "fmt"
	ripple "github.com/bmbstack/ripple"
	helper "github.com/bmbstack/ripple/helper"
	constant "github.com/bmbstack/ripple/nacos/nacos-sdk-go/v2/common/constant"
	client1 "github.com/bmbstack/ripple/nacos/rpcxnacos/client"
	proto "github.com/golang/protobuf/proto"
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

type TeachReq struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeachReq) Reset()         { *m = TeachReq{} }
func (m *TeachReq) String() string { return proto.CompactTextString(m) }
func (*TeachReq) ProtoMessage()    {}
func (*TeachReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_eebb822c9cd88554, []int{0}
}
func (m *TeachReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TeachReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TeachReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TeachReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeachReq.Merge(m, src)
}
func (m *TeachReq) XXX_Size() int {
	return m.Size()
}
func (m *TeachReq) XXX_DiscardUnknown() {
	xxx_messageInfo_TeachReq.DiscardUnknown(m)
}

var xxx_messageInfo_TeachReq proto.InternalMessageInfo

func (m *TeachReq) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type TeachReply struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeachReply) Reset()         { *m = TeachReply{} }
func (m *TeachReply) String() string { return proto.CompactTextString(m) }
func (*TeachReply) ProtoMessage()    {}
func (*TeachReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_eebb822c9cd88554, []int{1}
}
func (m *TeachReply) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TeachReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TeachReply.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TeachReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeachReply.Merge(m, src)
}
func (m *TeachReply) XXX_Size() int {
	return m.Size()
}
func (m *TeachReply) XXX_DiscardUnknown() {
	xxx_messageInfo_TeachReply.DiscardUnknown(m)
}

var xxx_messageInfo_TeachReply proto.InternalMessageInfo

func (m *TeachReply) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*TeachReq)(nil), "proto.TeachReq")
	proto.RegisterType((*TeachReply)(nil), "proto.TeachReply")
}

func init() {
	proto.RegisterFile("fixtures/forum/proto/teacher.proto", fileDescriptor_eebb822c9cd88554)
}

var fileDescriptor_eebb822c9cd88554 = []byte{
	// 162 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4a, 0xcb, 0xac, 0x28,
	0x29, 0x2d, 0x4a, 0x2d, 0xd6, 0x4f, 0xcb, 0x2f, 0x2a, 0xcd, 0xd5, 0x2f, 0x28, 0xca, 0x2f, 0xc9,
	0xd7, 0x2f, 0x49, 0x4d, 0x4c, 0xce, 0x48, 0x2d, 0xd2, 0x03, 0xf3, 0x84, 0x58, 0xc1, 0x94, 0x92,
	0x14, 0x17, 0x47, 0x08, 0x48, 0x3c, 0x28, 0xb5, 0x50, 0x88, 0x8f, 0x8b, 0x29, 0x33, 0x45, 0x82,
	0x51, 0x81, 0x51, 0x83, 0x25, 0x88, 0x29, 0x33, 0x45, 0x49, 0x81, 0x8b, 0x0b, 0x2a, 0x57, 0x90,
	0x53, 0x29, 0x24, 0xc4, 0xc5, 0x92, 0x97, 0x98, 0x9b, 0x0a, 0x96, 0xe7, 0x0c, 0x02, 0xb3, 0x8d,
	0x2c, 0xb8, 0xd8, 0x43, 0x20, 0xa6, 0x0a, 0xe9, 0x72, 0xb1, 0x82, 0x99, 0x42, 0xfc, 0x10, 0x0b,
	0xf4, 0x60, 0xc6, 0x4a, 0x09, 0xa2, 0x0a, 0x14, 0xe4, 0x54, 0x2a, 0x31, 0x38, 0x09, 0x9c, 0x78,
	0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x33, 0x1e, 0xcb, 0x31, 0x24,
	0xb1, 0x81, 0x55, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x55, 0xa2, 0x9f, 0x0e, 0xbd, 0x00,
	0x00, 0x00,
}

func (m *TeachReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TeachReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TeachReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Id != 0 {
		i = encodeVarintTeacher(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TeachReply) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TeachReply) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TeachReply) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
		i = encodeVarintTeacher(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTeacher(dAtA []byte, offset int, v uint64) int {
	offset -= sovTeacher(v)
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
// Gernerated from fixtures/forum/proto/teacher.proto

const ServiceNameOfTeacher = "TeacherRpc"

//================== interface ===================
type TeacherInterface interface {

	// TeacherInterface can be used for interface verification.
	// Teach is server rpc method as defined
	Teach(ctx context.Context, req *TeachReq, reply *TeachReply) (err error)
}

//================== server implement demo ===================
//ripple.Default().RegisterRpc("User", &UserRpcDemo{}, "")
type TeacherRpcDemo struct{}

func (this *TeacherRpcDemo) Teach(ctx context.Context, req *TeachReq, reply *TeachReply) (err error) {
	// TODO: add business logics
	*reply = TeachReply{}
	return nil
}

//================== client stub ===================
// newXClientForTeacher creates a XClient.
// You can configure this client pool with more options such as etcd registry, serialize type, select algorithm and fail mode.
func newXClientPoolForTeacher() (*client.XClientPool, error) {
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

	d, err := client1.NewNacosDiscovery(ServiceNameOfTeacher, "ripple", "DEFAULT_GROUP", clientConfig, serverConfig)
	if err != nil {
		return nil, err
	}

	opt := client.DefaultOption
	opt.SerializeType = protocol.ProtoBuffer

	var failMode client.FailMode
	switch config.Nacos.FailMode {
	case "failover":
		failMode = client.Failover
	case "failfast":
		failMode = client.Failfast
	case "failbackup":
		failMode = client.Failbackup
	default:
		failMode = client.Failtry
	}

	var selectMode client.SelectMode
	switch config.Nacos.SelectMode {
	case "randomSelect":
		selectMode = client.RandomSelect
	case "weightedRoundRobin":
		selectMode = client.WeightedRoundRobin
	case "weightedICMP":
		selectMode = client.WeightedICMP
	case "consistentHash":
		selectMode = client.ConsistentHash
	case "closest":
		selectMode = client.Closest
	case "selectByUser":
		selectMode = client.SelectByUser
	default:
		selectMode = client.RoundRobin
	}
	poolSize := config.Nacos.ClientPoolSize
	if poolSize == 0 {
		poolSize = 10
	}
	pool := client.NewXClientPool(poolSize, ServiceNameOfTeacher, failMode, selectMode, d, opt)

	return pool, nil
}

// Teacher is a client wrapped XClient.
type TeacherClient struct {
	XClientPool *client.XClientPool
}

// NewTeacherClient wraps a XClient as TeacherClient.
// You can pass a shared XClient object created by NewXClientForTeacher.
func NewTeacherClient() *TeacherClient {
	pool, err := newXClientPoolForTeacher()
	if err != nil {
		fmt.Println(fmt.Sprintf("Create rpcx client err: ripple", err.Error()))
		return &TeacherClient{}
	}
	return &TeacherClient{XClientPool: pool}
}

// Teach is client rpc method as defined
func (c *TeacherClient) Teach(ctx context.Context, req *TeachReq) (reply *TeachReply, err error) {
	reply = &TeachReply{}
	err = c.XClientPool.Get().Call(ctx, "Teach", req, reply)
	return reply, err
}

// ======================================================
func (m *TeachReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovTeacher(uint64(m.Id))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TeachReply) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovTeacher(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTeacher(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTeacher(x uint64) (n int) {
	return sovTeacher(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TeachReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTeacher
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
			return fmt.Errorf("proto: TeachReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TeachReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTeacher
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
			skippy, err := skipTeacher(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTeacher
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
func (m *TeachReply) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTeacher
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
			return fmt.Errorf("proto: TeachReply: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TeachReply: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTeacher
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
				return ErrInvalidLengthTeacher
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTeacher
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTeacher(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTeacher
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
func skipTeacher(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTeacher
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
					return 0, ErrIntOverflowTeacher
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
					return 0, ErrIntOverflowTeacher
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
				return 0, ErrInvalidLengthTeacher
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTeacher
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTeacher
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTeacher        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTeacher          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTeacher = fmt.Errorf("proto: unexpected end of group")
)
