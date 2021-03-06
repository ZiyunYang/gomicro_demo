// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf.proto

package protobuf

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
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

type DownloadRequest struct {
	Logfile              *Logfile `protobuf:"bytes,1,opt,name=logfile,proto3" json:"logfile,omitempty"`
	Offset               int64    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Whence               int64    `protobuf:"varint,3,opt,name=whence,proto3" json:"whence,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DownloadRequest) Reset()         { *m = DownloadRequest{} }
func (m *DownloadRequest) String() string { return proto.CompactTextString(m) }
func (*DownloadRequest) ProtoMessage()    {}
func (*DownloadRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{0}
}

func (m *DownloadRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DownloadRequest.Unmarshal(m, b)
}
func (m *DownloadRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DownloadRequest.Marshal(b, m, deterministic)
}
func (m *DownloadRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DownloadRequest.Merge(m, src)
}
func (m *DownloadRequest) XXX_Size() int {
	return xxx_messageInfo_DownloadRequest.Size(m)
}
func (m *DownloadRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DownloadRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DownloadRequest proto.InternalMessageInfo

func (m *DownloadRequest) GetLogfile() *Logfile {
	if m != nil {
		return m.Logfile
	}
	return nil
}

func (m *DownloadRequest) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *DownloadRequest) GetWhence() int64 {
	if m != nil {
		return m.Whence
	}
	return 0
}

type DownloadResponse struct {
	Total                int64    `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Times                int64    `protobuf:"varint,2,opt,name=times,proto3" json:"times,omitempty"`
	DataBytes            []byte   `protobuf:"bytes,3,opt,name=dataBytes,proto3" json:"dataBytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DownloadResponse) Reset()         { *m = DownloadResponse{} }
func (m *DownloadResponse) String() string { return proto.CompactTextString(m) }
func (*DownloadResponse) ProtoMessage()    {}
func (*DownloadResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{1}
}

func (m *DownloadResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DownloadResponse.Unmarshal(m, b)
}
func (m *DownloadResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DownloadResponse.Marshal(b, m, deterministic)
}
func (m *DownloadResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DownloadResponse.Merge(m, src)
}
func (m *DownloadResponse) XXX_Size() int {
	return xxx_messageInfo_DownloadResponse.Size(m)
}
func (m *DownloadResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DownloadResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DownloadResponse proto.InternalMessageInfo

func (m *DownloadResponse) GetTotal() int64 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *DownloadResponse) GetTimes() int64 {
	if m != nil {
		return m.Times
	}
	return 0
}

func (m *DownloadResponse) GetDataBytes() []byte {
	if m != nil {
		return m.DataBytes
	}
	return nil
}

type Logfile struct {
	Pod                  string   `protobuf:"bytes,1,opt,name=pod,proto3" json:"pod,omitempty"`
	Container            string   `protobuf:"bytes,2,opt,name=container,proto3" json:"container,omitempty"`
	File                 string   `protobuf:"bytes,3,opt,name=file,proto3" json:"file,omitempty"`
	Topic                string   `protobuf:"bytes,4,opt,name=topic,proto3" json:"topic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Logfile) Reset()         { *m = Logfile{} }
func (m *Logfile) String() string { return proto.CompactTextString(m) }
func (*Logfile) ProtoMessage()    {}
func (*Logfile) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{2}
}

func (m *Logfile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Logfile.Unmarshal(m, b)
}
func (m *Logfile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Logfile.Marshal(b, m, deterministic)
}
func (m *Logfile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Logfile.Merge(m, src)
}
func (m *Logfile) XXX_Size() int {
	return xxx_messageInfo_Logfile.Size(m)
}
func (m *Logfile) XXX_DiscardUnknown() {
	xxx_messageInfo_Logfile.DiscardUnknown(m)
}

var xxx_messageInfo_Logfile proto.InternalMessageInfo

func (m *Logfile) GetPod() string {
	if m != nil {
		return m.Pod
	}
	return ""
}

func (m *Logfile) GetContainer() string {
	if m != nil {
		return m.Container
	}
	return ""
}

func (m *Logfile) GetFile() string {
	if m != nil {
		return m.File
	}
	return ""
}

func (m *Logfile) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func init() {
	proto.RegisterType((*DownloadRequest)(nil), "DownloadRequest")
	proto.RegisterType((*DownloadResponse)(nil), "DownloadResponse")
	proto.RegisterType((*Logfile)(nil), "Logfile")
}

func init() { proto.RegisterFile("protobuf.proto", fileDescriptor_c77a803fcbc0c059) }

var fileDescriptor_c77a803fcbc0c059 = []byte{
	// 260 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xcb, 0x4f, 0x83, 0x40,
	0x10, 0x87, 0x45, 0x6a, 0x2b, 0xe3, 0x0b, 0x37, 0xc6, 0x10, 0xe3, 0xa1, 0xe1, 0xd4, 0x13, 0x31,
	0xf5, 0xe2, 0xd9, 0xf8, 0xb8, 0x70, 0x5a, 0xaf, 0x5e, 0xb6, 0x30, 0x50, 0x12, 0xba, 0x83, 0xec,
	0x34, 0x8d, 0xff, 0xbd, 0x61, 0x16, 0x43, 0xd2, 0xdb, 0x7c, 0xdf, 0x66, 0x7f, 0xf3, 0x80, 0xeb,
	0xae, 0x27, 0xa6, 0xcd, 0xbe, 0xca, 0xa4, 0x48, 0x11, 0x6e, 0xde, 0xe8, 0x60, 0x5b, 0x32, 0xa5,
	0xc6, 0x9f, 0x3d, 0x3a, 0x56, 0x29, 0x2c, 0x5a, 0xaa, 0xab, 0xa6, 0xc5, 0x24, 0x58, 0x06, 0xab,
	0x8b, 0xf5, 0x79, 0x96, 0x7b, 0xd6, 0xff, 0x0f, 0xea, 0x1e, 0xe6, 0x54, 0x55, 0x0e, 0x39, 0x39,
	0x5d, 0x06, 0xab, 0x50, 0x8f, 0x34, 0xf8, 0xc3, 0x16, 0x6d, 0x81, 0x49, 0xe8, 0xbd, 0xa7, 0xf4,
	0x1b, 0xe2, 0xa9, 0x8d, 0xeb, 0xc8, 0x3a, 0x54, 0x77, 0x70, 0xc6, 0xc4, 0xa6, 0x95, 0x2e, 0xa1,
	0xf6, 0x20, 0xb6, 0xd9, 0xa1, 0x1b, 0x83, 0x3d, 0xa8, 0x47, 0x88, 0x4a, 0xc3, 0xe6, 0xf5, 0x97,
	0xd1, 0x49, 0xf4, 0xa5, 0x9e, 0x44, 0x5a, 0xc0, 0x62, 0x9c, 0x50, 0xc5, 0x10, 0x76, 0x54, 0x4a,
	0x64, 0xa4, 0x87, 0x72, 0xf8, 0x5a, 0x90, 0x65, 0xd3, 0x58, 0xec, 0x25, 0x34, 0xd2, 0x93, 0x50,
	0x0a, 0x66, 0xb2, 0x69, 0x28, 0x0f, 0x52, 0xfb, 0xc1, 0xba, 0xa6, 0x48, 0x66, 0x22, 0x3d, 0xac,
	0xdf, 0x21, 0xca, 0xa9, 0xfe, 0x34, 0xbc, 0xc5, 0x5e, 0xbd, 0xc0, 0x55, 0x4e, 0xf5, 0x47, 0xd3,
	0xe2, 0x17, 0xf7, 0x68, 0x76, 0x2a, 0xce, 0x8e, 0xce, 0xf8, 0x70, 0x9b, 0x1d, 0x6f, 0x9c, 0x9e,
	0x3c, 0x05, 0x9b, 0xb9, 0xdc, 0xfd, 0xf9, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x8f, 0x21, 0x9d, 0xfb,
	0x89, 0x01, 0x00, 0x00,
}
