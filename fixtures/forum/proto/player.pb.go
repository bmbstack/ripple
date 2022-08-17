// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fixtures/forum/proto/player.proto

// @RippleRpc
// @NacosGroup DEFAULT_GROUP
// @NacosCluster playerserver

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

type GetPlayerInfoReq struct {
	// 用户id
	PlayerId             int64    `protobuf:"varint,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPlayerInfoReq) Reset()         { *m = GetPlayerInfoReq{} }
func (m *GetPlayerInfoReq) String() string { return proto.CompactTextString(m) }
func (*GetPlayerInfoReq) ProtoMessage()    {}
func (*GetPlayerInfoReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_ceed7288d3804cbf, []int{0}
}
func (m *GetPlayerInfoReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GetPlayerInfoReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GetPlayerInfoReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GetPlayerInfoReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPlayerInfoReq.Merge(m, src)
}
func (m *GetPlayerInfoReq) XXX_Size() int {
	return m.Size()
}
func (m *GetPlayerInfoReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPlayerInfoReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetPlayerInfoReq proto.InternalMessageInfo

func (m *GetPlayerInfoReq) GetPlayerId() int64 {
	if m != nil {
		return m.PlayerId
	}
	return 0
}

type GetPlayerInfoResponse struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	OpenId               string   `protobuf:"bytes,2,opt,name=open_id,json=openId,proto3" json:"open_id,omitempty"`
	Nickname             string   `protobuf:"bytes,3,opt,name=nickname,proto3" json:"nickname,omitempty"`
	UnionId              string   `protobuf:"bytes,4,opt,name=union_id,json=unionId,proto3" json:"union_id,omitempty"`
	State                int32    `protobuf:"varint,5,opt,name=state,proto3" json:"state,omitempty"`
	Avatar               string   `protobuf:"bytes,6,opt,name=avatar,proto3" json:"avatar,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPlayerInfoResponse) Reset()         { *m = GetPlayerInfoResponse{} }
func (m *GetPlayerInfoResponse) String() string { return proto.CompactTextString(m) }
func (*GetPlayerInfoResponse) ProtoMessage()    {}
func (*GetPlayerInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ceed7288d3804cbf, []int{1}
}
func (m *GetPlayerInfoResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GetPlayerInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GetPlayerInfoResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GetPlayerInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPlayerInfoResponse.Merge(m, src)
}
func (m *GetPlayerInfoResponse) XXX_Size() int {
	return m.Size()
}
func (m *GetPlayerInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPlayerInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetPlayerInfoResponse proto.InternalMessageInfo

func (m *GetPlayerInfoResponse) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *GetPlayerInfoResponse) GetOpenId() string {
	if m != nil {
		return m.OpenId
	}
	return ""
}

func (m *GetPlayerInfoResponse) GetNickname() string {
	if m != nil {
		return m.Nickname
	}
	return ""
}

func (m *GetPlayerInfoResponse) GetUnionId() string {
	if m != nil {
		return m.UnionId
	}
	return ""
}

func (m *GetPlayerInfoResponse) GetState() int32 {
	if m != nil {
		return m.State
	}
	return 0
}

func (m *GetPlayerInfoResponse) GetAvatar() string {
	if m != nil {
		return m.Avatar
	}
	return ""
}

func init() {
	proto.RegisterType((*GetPlayerInfoReq)(nil), "proto.GetPlayerInfoReq")
	proto.RegisterType((*GetPlayerInfoResponse)(nil), "proto.GetPlayerInfoResponse")
}

func init() { proto.RegisterFile("fixtures/forum/proto/player.proto", fileDescriptor_ceed7288d3804cbf) }

var fileDescriptor_ceed7288d3804cbf = []byte{
	// 261 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x41, 0x4a, 0xc4, 0x30,
	0x14, 0x86, 0x27, 0x1d, 0x9b, 0xe9, 0x3c, 0x50, 0x86, 0x87, 0x3a, 0x75, 0x94, 0x52, 0xbb, 0xea,
	0x6a, 0x0a, 0x7a, 0x03, 0x37, 0xda, 0x9d, 0xe4, 0x02, 0x12, 0x4d, 0x0a, 0x41, 0x27, 0xa9, 0x69,
	0x2a, 0x7a, 0x13, 0x2f, 0xe0, 0x5d, 0x5c, 0x7a, 0x04, 0xa9, 0x17, 0x91, 0x49, 0x60, 0xc0, 0xe2,
	0x2a, 0x7c, 0xf9, 0xde, 0x9f, 0xbc, 0xf7, 0xe0, 0xbc, 0x51, 0xaf, 0xae, 0xb7, 0xb2, 0xab, 0x1a,
	0x63, 0xfb, 0x4d, 0xd5, 0x5a, 0xe3, 0x4c, 0xd5, 0x3e, 0xf1, 0x37, 0x69, 0xd7, 0x1e, 0x30, 0xf6,
	0x47, 0x51, 0xc1, 0xe2, 0x5a, 0xba, 0x5b, 0x6f, 0x6a, 0xdd, 0x18, 0x26, 0x9f, 0xf1, 0x14, 0xe6,
	0xa1, 0xf4, 0x4e, 0x89, 0x94, 0xe4, 0xa4, 0x9c, 0xb2, 0x24, 0x5c, 0xd4, 0xa2, 0xf8, 0x20, 0x70,
	0x34, 0x4a, 0x74, 0xad, 0xd1, 0x9d, 0xc4, 0x03, 0x88, 0x76, 0xf5, 0x91, 0x12, 0xb8, 0x84, 0x99,
	0x69, 0xa5, 0xde, 0x3e, 0x12, 0xe5, 0xa4, 0x9c, 0x33, 0xba, 0xc5, 0x5a, 0xe0, 0x0a, 0x12, 0xad,
	0x1e, 0x1e, 0x35, 0xdf, 0xc8, 0x74, 0xea, 0xcd, 0x8e, 0xf1, 0x04, 0x92, 0x5e, 0x2b, 0xe3, 0x53,
	0x7b, 0xde, 0xcd, 0x3c, 0xd7, 0x02, 0x0f, 0x21, 0xee, 0x1c, 0x77, 0x32, 0x8d, 0x73, 0x52, 0xc6,
	0x2c, 0x00, 0x1e, 0x03, 0xe5, 0x2f, 0xdc, 0x71, 0x9b, 0xd2, 0xf0, 0x49, 0xa0, 0x0b, 0x06, 0x34,
	0xf4, 0x88, 0x37, 0xb0, 0xff, 0xa7, 0x61, 0x5c, 0x86, 0x15, 0xac, 0xc7, 0x83, 0xaf, 0xce, 0xfe,
	0x17, 0x61, 0xbe, 0x62, 0x72, 0xb5, 0xf8, 0x1c, 0x32, 0xf2, 0x35, 0x64, 0xe4, 0x7b, 0xc8, 0xc8,
	0xfb, 0x4f, 0x36, 0xb9, 0xa7, 0x3e, 0x70, 0xf9, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x5d, 0x87, 0xe8,
	0x7c, 0x71, 0x01, 0x00, 0x00,
}

func (m *GetPlayerInfoReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetPlayerInfoReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GetPlayerInfoReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.PlayerId != 0 {
		i = encodeVarintPlayer(dAtA, i, uint64(m.PlayerId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *GetPlayerInfoResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetPlayerInfoResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GetPlayerInfoResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Avatar) > 0 {
		i -= len(m.Avatar)
		copy(dAtA[i:], m.Avatar)
		i = encodeVarintPlayer(dAtA, i, uint64(len(m.Avatar)))
		i--
		dAtA[i] = 0x32
	}
	if m.State != 0 {
		i = encodeVarintPlayer(dAtA, i, uint64(m.State))
		i--
		dAtA[i] = 0x28
	}
	if len(m.UnionId) > 0 {
		i -= len(m.UnionId)
		copy(dAtA[i:], m.UnionId)
		i = encodeVarintPlayer(dAtA, i, uint64(len(m.UnionId)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Nickname) > 0 {
		i -= len(m.Nickname)
		copy(dAtA[i:], m.Nickname)
		i = encodeVarintPlayer(dAtA, i, uint64(len(m.Nickname)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.OpenId) > 0 {
		i -= len(m.OpenId)
		copy(dAtA[i:], m.OpenId)
		i = encodeVarintPlayer(dAtA, i, uint64(len(m.OpenId)))
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintPlayer(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintPlayer(dAtA []byte, offset int, v uint64) int {
	offset -= sovPlayer(v)
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
// Gernerated from fixtures/forum/proto/player.proto

const ServiceNameOfPlayer = "PlayerRpc"

//================== interface ===================
type PlayerInterface interface {

	// PlayerInterface can be used for interface verification.
	// GetPlayerInfo is server rpc method as defined
	GetPlayerInfo(ctx context.Context, req *GetPlayerInfoReq, reply *GetPlayerInfoResponse) (err error)
}

//================== server implement demo ===================
//ripple.Default().RegisterRpc("User", &UserRpcDemo{}, "")
type PlayerRpcDemo struct{}

func (this *PlayerRpcDemo) GetPlayerInfo(ctx context.Context, req *GetPlayerInfoReq, reply *GetPlayerInfoResponse) (err error) {
	// TODO: add business logics
	*reply = GetPlayerInfoResponse{}
	return nil
}

//================== client stub ===================
// newXClientForPlayer creates a XClient.
// You can configure this client pool with more options such as etcd registry, serialize type, select algorithm and fail mode.
func newXClientPoolForPlayer() (*client.XClientPool, error) {
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

	d, err := client1.NewNacosDiscovery(ServiceNameOfPlayer, "playerserver", "DEFAULT_GROUP", clientConfig, serverConfig)
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
		failMode = client.Failover
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
	pool := client.NewXClientPool(poolSize, ServiceNameOfPlayer, failMode, selectMode, d, opt)

	return pool, nil
}

// Player is a client wrapped XClient.
type PlayerClient struct {
	XClientPool *client.XClientPool
}

// NewPlayerClient wraps a XClient as PlayerClient.
// You can pass a shared XClient object created by NewXClientForPlayer.
func NewPlayerClient() *PlayerClient {
	pool, err := newXClientPoolForPlayer()
	if err != nil {
		fmt.Println(fmt.Sprintf("Create rpcx client err: playerserver", err.Error()))
		return &PlayerClient{}
	}
	return &PlayerClient{XClientPool: pool}
}

// GetPlayerInfo is client rpc method as defined
func (c *PlayerClient) GetPlayerInfo(ctx context.Context, req *GetPlayerInfoReq) (reply *GetPlayerInfoResponse, err error) {
	reply = &GetPlayerInfoResponse{}
	if c.XClientPool != nil {
		err = c.XClientPool.Get().Call(ctx, "GetPlayerInfo", req, reply)
	}
	return reply, err
}

// ======================================================
func (m *GetPlayerInfoReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PlayerId != 0 {
		n += 1 + sovPlayer(uint64(m.PlayerId))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *GetPlayerInfoResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovPlayer(uint64(m.Id))
	}
	l = len(m.OpenId)
	if l > 0 {
		n += 1 + l + sovPlayer(uint64(l))
	}
	l = len(m.Nickname)
	if l > 0 {
		n += 1 + l + sovPlayer(uint64(l))
	}
	l = len(m.UnionId)
	if l > 0 {
		n += 1 + l + sovPlayer(uint64(l))
	}
	if m.State != 0 {
		n += 1 + sovPlayer(uint64(m.State))
	}
	l = len(m.Avatar)
	if l > 0 {
		n += 1 + l + sovPlayer(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovPlayer(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPlayer(x uint64) (n int) {
	return sovPlayer(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GetPlayerInfoReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPlayer
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
			return fmt.Errorf("proto: GetPlayerInfoReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetPlayerInfoReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PlayerId", wireType)
			}
			m.PlayerId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlayer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PlayerId |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPlayer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPlayer
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
func (m *GetPlayerInfoResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPlayer
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
			return fmt.Errorf("proto: GetPlayerInfoResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetPlayerInfoResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlayer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OpenId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlayer
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
				return ErrInvalidLengthPlayer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPlayer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OpenId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nickname", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlayer
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
				return ErrInvalidLengthPlayer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPlayer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Nickname = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnionId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlayer
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
				return ErrInvalidLengthPlayer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPlayer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UnionId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlayer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Avatar", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlayer
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
				return ErrInvalidLengthPlayer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPlayer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Avatar = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPlayer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPlayer
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
func skipPlayer(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPlayer
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
					return 0, ErrIntOverflowPlayer
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
					return 0, ErrIntOverflowPlayer
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
				return 0, ErrInvalidLengthPlayer
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPlayer
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPlayer
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPlayer        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPlayer          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPlayer = fmt.Errorf("proto: unexpected end of group")
)