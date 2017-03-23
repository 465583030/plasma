// Code generated by protoc-gen-go.
// source: stream.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	stream.proto

It has these top-level messages:
	Request
	EventType
	Payload
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type Request struct {
	Events []*EventType `protobuf:"bytes,1,rep,name=Events" json:"Events,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto1.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetEvents() []*EventType {
	if m != nil {
		return m.Events
	}
	return nil
}

type EventType struct {
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
}

func (m *EventType) Reset()                    { *m = EventType{} }
func (m *EventType) String() string            { return proto1.CompactTextString(m) }
func (*EventType) ProtoMessage()               {}
func (*EventType) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *EventType) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type Payload struct {
	EventType *EventType `protobuf:"bytes,1,opt,name=eventType" json:"eventType,omitempty"`
	Data      string     `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
}

func (m *Payload) Reset()                    { *m = Payload{} }
func (m *Payload) String() string            { return proto1.CompactTextString(m) }
func (*Payload) ProtoMessage()               {}
func (*Payload) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Payload) GetEventType() *EventType {
	if m != nil {
		return m.EventType
	}
	return nil
}

func (m *Payload) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto1.RegisterType((*Request)(nil), "proto.Request")
	proto1.RegisterType((*EventType)(nil), "proto.EventType")
	proto1.RegisterType((*Payload)(nil), "proto.Payload")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for StreamService service

type StreamServiceClient interface {
	Events(ctx context.Context, in *Request, opts ...grpc.CallOption) (StreamService_EventsClient, error)
}

type streamServiceClient struct {
	cc *grpc.ClientConn
}

func NewStreamServiceClient(cc *grpc.ClientConn) StreamServiceClient {
	return &streamServiceClient{cc}
}

func (c *streamServiceClient) Events(ctx context.Context, in *Request, opts ...grpc.CallOption) (StreamService_EventsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_StreamService_serviceDesc.Streams[0], c.cc, "/proto.StreamService/Events", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamServiceEventsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StreamService_EventsClient interface {
	Recv() (*Payload, error)
	grpc.ClientStream
}

type streamServiceEventsClient struct {
	grpc.ClientStream
}

func (x *streamServiceEventsClient) Recv() (*Payload, error) {
	m := new(Payload)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for StreamService service

type StreamServiceServer interface {
	Events(*Request, StreamService_EventsServer) error
}

func RegisterStreamServiceServer(s *grpc.Server, srv StreamServiceServer) {
	s.RegisterService(&_StreamService_serviceDesc, srv)
}

func _StreamService_Events_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamServiceServer).Events(m, &streamServiceEventsServer{stream})
}

type StreamService_EventsServer interface {
	Send(*Payload) error
	grpc.ServerStream
}

type streamServiceEventsServer struct {
	grpc.ServerStream
}

func (x *streamServiceEventsServer) Send(m *Payload) error {
	return x.ServerStream.SendMsg(m)
}

var _StreamService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.StreamService",
	HandlerType: (*StreamServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Events",
			Handler:       _StreamService_Events_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "stream.proto",
}

func init() { proto1.RegisterFile("stream.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 180 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x2e, 0x29, 0x4a,
	0x4d, 0xcc, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0xc6, 0x5c, 0xec,
	0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x1a, 0x5c, 0x6c, 0xae, 0x65, 0xa9, 0x79, 0x25,
	0xc5, 0x12, 0x8c, 0x0a, 0xcc, 0x1a, 0xdc, 0x46, 0x02, 0x10, 0x95, 0x7a, 0x60, 0xc1, 0x90, 0xca,
	0x82, 0xd4, 0x20, 0xa8, 0xbc, 0x92, 0x3c, 0x17, 0x27, 0x5c, 0x50, 0x48, 0x88, 0x8b, 0xa5, 0xa4,
	0xb2, 0x20, 0x55, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xcc, 0x56, 0xf2, 0xe5, 0x62, 0x0f,
	0x48, 0xac, 0xcc, 0xc9, 0x4f, 0x4c, 0x11, 0xd2, 0xe3, 0xe2, 0x4c, 0x85, 0xa9, 0x05, 0xab, 0xc1,
	0x66, 0x30, 0x42, 0x09, 0xc8, 0xb8, 0x94, 0xc4, 0x92, 0x44, 0x09, 0x26, 0x88, 0x71, 0x20, 0xb6,
	0x91, 0x2d, 0x17, 0x6f, 0x30, 0xd8, 0xed, 0xc1, 0xa9, 0x45, 0x65, 0x99, 0xc9, 0xa9, 0x42, 0x3a,
	0x30, 0xa7, 0x0a, 0xf1, 0x41, 0xcd, 0x82, 0x7a, 0x42, 0x0a, 0xc6, 0x87, 0x5a, 0xaf, 0xc4, 0x60,
	0xc0, 0x98, 0xc4, 0x06, 0x16, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x84, 0xd2, 0xdd, 0x34,
	0x01, 0x01, 0x00, 0x00,
}
