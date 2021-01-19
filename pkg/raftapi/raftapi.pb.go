// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: raftapi.proto

package raftapi

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

//
//Basic Types
type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftapi_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_raftapi_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_raftapi_proto_rawDescGZIP(), []int{0}
}

type Bool struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status bool `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Bool) Reset() {
	*x = Bool{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftapi_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bool) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bool) ProtoMessage() {}

func (x *Bool) ProtoReflect() protoreflect.Message {
	mi := &file_raftapi_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bool.ProtoReflect.Descriptor instead.
func (*Bool) Descriptor() ([]byte, []int) {
	return file_raftapi_proto_rawDescGZIP(), []int{1}
}

func (x *Bool) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value int64 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Value) Reset() {
	*x = Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftapi_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_raftapi_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_raftapi_proto_rawDescGZIP(), []int{2}
}

func (x *Value) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type LogEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term  uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Value int64  `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *LogEntry) Reset() {
	*x = LogEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftapi_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogEntry) ProtoMessage() {}

func (x *LogEntry) ProtoReflect() protoreflect.Message {
	mi := &file_raftapi_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogEntry.ProtoReflect.Descriptor instead.
func (*LogEntry) Descriptor() ([]byte, []int) {
	return file_raftapi_proto_rawDescGZIP(), []int{3}
}

func (x *LogEntry) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *LogEntry) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type WhoAmI struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Port uint32 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	Role string `protobuf:"bytes,3,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *WhoAmI) Reset() {
	*x = WhoAmI{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftapi_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WhoAmI) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WhoAmI) ProtoMessage() {}

func (x *WhoAmI) ProtoReflect() protoreflect.Message {
	mi := &file_raftapi_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WhoAmI.ProtoReflect.Descriptor instead.
func (*WhoAmI) Descriptor() ([]byte, []int) {
	return file_raftapi_proto_rawDescGZIP(), []int{4}
}

func (x *WhoAmI) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *WhoAmI) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *WhoAmI) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

//
//Raft Protocol Messages
type RequestVoteMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term        uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Candidate   string `protobuf:"bytes,2,opt,name=candidate,proto3" json:"candidate,omitempty"`
	LogSize     uint64 `protobuf:"varint,3,opt,name=logSize,proto3" json:"logSize,omitempty"`
	LastLogTerm uint32 `protobuf:"varint,4,opt,name=lastLogTerm,proto3" json:"lastLogTerm,omitempty"`
}

func (x *RequestVoteMessage) Reset() {
	*x = RequestVoteMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftapi_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestVoteMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestVoteMessage) ProtoMessage() {}

func (x *RequestVoteMessage) ProtoReflect() protoreflect.Message {
	mi := &file_raftapi_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestVoteMessage.ProtoReflect.Descriptor instead.
func (*RequestVoteMessage) Descriptor() ([]byte, []int) {
	return file_raftapi_proto_rawDescGZIP(), []int{5}
}

func (x *RequestVoteMessage) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *RequestVoteMessage) GetCandidate() string {
	if x != nil {
		return x.Candidate
	}
	return ""
}

func (x *RequestVoteMessage) GetLogSize() uint64 {
	if x != nil {
		return x.LogSize
	}
	return 0
}

func (x *RequestVoteMessage) GetLastLogTerm() uint32 {
	if x != nil {
		return x.LastLogTerm
	}
	return 0
}

type AppendEntryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term        uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Leader      string `protobuf:"bytes,2,opt,name=leader,proto3" json:"leader,omitempty"`
	PrevLogId   int64  `protobuf:"varint,3,opt,name=prevLogId,proto3" json:"prevLogId,omitempty"`
	PrevLogTerm uint64 `protobuf:"varint,4,opt,name=prevLogTerm,proto3" json:"prevLogTerm,omitempty"`
	// Types that are assignable to LogEntry:
	//	*AppendEntryRequest_Entry
	LogEntry isAppendEntryRequest_LogEntry `protobuf_oneof:"logEntry"`
}

func (x *AppendEntryRequest) Reset() {
	*x = AppendEntryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftapi_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntryRequest) ProtoMessage() {}

func (x *AppendEntryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_raftapi_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntryRequest.ProtoReflect.Descriptor instead.
func (*AppendEntryRequest) Descriptor() ([]byte, []int) {
	return file_raftapi_proto_rawDescGZIP(), []int{6}
}

func (x *AppendEntryRequest) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntryRequest) GetLeader() string {
	if x != nil {
		return x.Leader
	}
	return ""
}

func (x *AppendEntryRequest) GetPrevLogId() int64 {
	if x != nil {
		return x.PrevLogId
	}
	return 0
}

func (x *AppendEntryRequest) GetPrevLogTerm() uint64 {
	if x != nil {
		return x.PrevLogTerm
	}
	return 0
}

func (m *AppendEntryRequest) GetLogEntry() isAppendEntryRequest_LogEntry {
	if m != nil {
		return m.LogEntry
	}
	return nil
}

func (x *AppendEntryRequest) GetEntry() *LogEntry {
	if x, ok := x.GetLogEntry().(*AppendEntryRequest_Entry); ok {
		return x.Entry
	}
	return nil
}

type isAppendEntryRequest_LogEntry interface {
	isAppendEntryRequest_LogEntry()
}

type AppendEntryRequest_Entry struct {
	Entry *LogEntry `protobuf:"bytes,5,opt,name=entry,proto3,oneof"`
}

func (*AppendEntryRequest_Entry) isAppendEntryRequest_LogEntry() {}

type AppendEntryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term    uint64 `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`
	Success bool   `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *AppendEntryResponse) Reset() {
	*x = AppendEntryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftapi_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntryResponse) ProtoMessage() {}

func (x *AppendEntryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_raftapi_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntryResponse.ProtoReflect.Descriptor instead.
func (*AppendEntryResponse) Descriptor() ([]byte, []int) {
	return file_raftapi_proto_rawDescGZIP(), []int{7}
}

func (x *AppendEntryResponse) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntryResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_raftapi_proto protoreflect.FileDescriptor

var file_raftapi_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x72, 0x61, 0x66, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x72, 0x61, 0x66, 0x74, 0x61, 0x70, 0x69, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x1e, 0x0a, 0x04, 0x42, 0x6f, 0x6f, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x22, 0x1d, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x34, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x44, 0x0a, 0x06, 0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0x82, 0x01, 0x0a,
	0x12, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x61, 0x6e, 0x64, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x61, 0x6e, 0x64,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6c, 0x6f, 0x67, 0x53, 0x69, 0x7a, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x6c, 0x6f, 0x67, 0x53, 0x69, 0x7a, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72,
	0x6d, 0x22, 0xb7, 0x01, 0x0a, 0x12, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x16, 0x0a, 0x06,
	0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67,
	0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72,
	0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x70, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67,
	0x54, 0x65, 0x72, 0x6d, 0x12, 0x29, 0x0a, 0x05, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x6f,
	0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x48, 0x00, 0x52, 0x05, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x42,
	0x0a, 0x0a, 0x08, 0x6c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x22, 0x43, 0x0a, 0x13, 0x41,
	0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x32, 0xac, 0x02, 0x0a, 0x0b, 0x52, 0x61, 0x66, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x29, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x0e, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x61,
	0x70, 0x69, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0f, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x61,
	0x70, 0x69, 0x2e, 0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49, 0x22, 0x00, 0x12, 0x2b, 0x0a, 0x08, 0x53,
	0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x12, 0x0e, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x61, 0x70,
	0x69, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x61, 0x70,
	0x69, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x22, 0x00, 0x12, 0x2e, 0x0a, 0x0b, 0x41, 0x70, 0x70, 0x65,
	0x6e, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x0e, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x61, 0x70,
	0x69, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x0d, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x61, 0x70,
	0x69, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x22, 0x00, 0x12, 0x49, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x1b, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x61, 0x70,
	0x69, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x1a, 0x1b, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x0b, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x1b, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x70, 0x70,
	0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1c, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x0b, 0x5a, 0x09, 0x2e, 0x3b, 0x72, 0x61, 0x66, 0x74, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_raftapi_proto_rawDescOnce sync.Once
	file_raftapi_proto_rawDescData = file_raftapi_proto_rawDesc
)

func file_raftapi_proto_rawDescGZIP() []byte {
	file_raftapi_proto_rawDescOnce.Do(func() {
		file_raftapi_proto_rawDescData = protoimpl.X.CompressGZIP(file_raftapi_proto_rawDescData)
	})
	return file_raftapi_proto_rawDescData
}

var file_raftapi_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_raftapi_proto_goTypes = []interface{}{
	(*Empty)(nil),               // 0: raftapi.Empty
	(*Bool)(nil),                // 1: raftapi.Bool
	(*Value)(nil),               // 2: raftapi.Value
	(*LogEntry)(nil),            // 3: raftapi.LogEntry
	(*WhoAmI)(nil),              // 4: raftapi.WhoAmI
	(*RequestVoteMessage)(nil),  // 5: raftapi.RequestVoteMessage
	(*AppendEntryRequest)(nil),  // 6: raftapi.AppendEntryRequest
	(*AppendEntryResponse)(nil), // 7: raftapi.AppendEntryResponse
}
var file_raftapi_proto_depIdxs = []int32{
	3, // 0: raftapi.AppendEntryRequest.entry:type_name -> raftapi.LogEntry
	0, // 1: raftapi.RaftService.Ping:input_type -> raftapi.Empty
	0, // 2: raftapi.RaftService.Shutdown:input_type -> raftapi.Empty
	2, // 3: raftapi.RaftService.AppendValue:input_type -> raftapi.Value
	5, // 4: raftapi.RaftService.RequestVote:input_type -> raftapi.RequestVoteMessage
	6, // 5: raftapi.RaftService.AppendEntry:input_type -> raftapi.AppendEntryRequest
	4, // 6: raftapi.RaftService.Ping:output_type -> raftapi.WhoAmI
	1, // 7: raftapi.RaftService.Shutdown:output_type -> raftapi.Bool
	1, // 8: raftapi.RaftService.AppendValue:output_type -> raftapi.Bool
	5, // 9: raftapi.RaftService.RequestVote:output_type -> raftapi.RequestVoteMessage
	7, // 10: raftapi.RaftService.AppendEntry:output_type -> raftapi.AppendEntryResponse
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_raftapi_proto_init() }
func file_raftapi_proto_init() {
	if File_raftapi_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_raftapi_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftapi_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Bool); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftapi_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Value); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftapi_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogEntry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftapi_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WhoAmI); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftapi_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestVoteMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftapi_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendEntryRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftapi_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendEntryResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_raftapi_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*AppendEntryRequest_Entry)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_raftapi_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_raftapi_proto_goTypes,
		DependencyIndexes: file_raftapi_proto_depIdxs,
		MessageInfos:      file_raftapi_proto_msgTypes,
	}.Build()
	File_raftapi_proto = out.File
	file_raftapi_proto_rawDesc = nil
	file_raftapi_proto_goTypes = nil
	file_raftapi_proto_depIdxs = nil
}
