// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.1
// source: mp4.proto

package pb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReqRecordList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StreamPath string `protobuf:"bytes,1,opt,name=streamPath,proto3" json:"streamPath,omitempty"`
	Range      string `protobuf:"bytes,2,opt,name=range,proto3" json:"range,omitempty"`
	Start      string `protobuf:"bytes,3,opt,name=start,proto3" json:"start,omitempty"`
	End        string `protobuf:"bytes,4,opt,name=end,proto3" json:"end,omitempty"`
	Page       uint32 `protobuf:"varint,5,opt,name=page,proto3" json:"page,omitempty"`
	PageSize   uint32 `protobuf:"varint,6,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
}

func (x *ReqRecordList) Reset() {
	*x = ReqRecordList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mp4_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqRecordList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqRecordList) ProtoMessage() {}

func (x *ReqRecordList) ProtoReflect() protoreflect.Message {
	mi := &file_mp4_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqRecordList.ProtoReflect.Descriptor instead.
func (*ReqRecordList) Descriptor() ([]byte, []int) {
	return file_mp4_proto_rawDescGZIP(), []int{0}
}

func (x *ReqRecordList) GetStreamPath() string {
	if x != nil {
		return x.StreamPath
	}
	return ""
}

func (x *ReqRecordList) GetRange() string {
	if x != nil {
		return x.Range
	}
	return ""
}

func (x *ReqRecordList) GetStart() string {
	if x != nil {
		return x.Start
	}
	return ""
}

func (x *ReqRecordList) GetEnd() string {
	if x != nil {
		return x.End
	}
	return ""
}

func (x *ReqRecordList) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ReqRecordList) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type RecordFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	FilePath   string                 `protobuf:"bytes,2,opt,name=filePath,proto3" json:"filePath,omitempty"`
	StreamPath string                 `protobuf:"bytes,3,opt,name=streamPath,proto3" json:"streamPath,omitempty"`
	StartTime  *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=startTime,proto3" json:"startTime,omitempty"`
	EndTime    *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=endTime,proto3" json:"endTime,omitempty"`
}

func (x *RecordFile) Reset() {
	*x = RecordFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mp4_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordFile) ProtoMessage() {}

func (x *RecordFile) ProtoReflect() protoreflect.Message {
	mi := &file_mp4_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordFile.ProtoReflect.Descriptor instead.
func (*RecordFile) Descriptor() ([]byte, []int) {
	return file_mp4_proto_rawDescGZIP(), []int{1}
}

func (x *RecordFile) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RecordFile) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

func (x *RecordFile) GetStreamPath() string {
	if x != nil {
		return x.StreamPath
	}
	return ""
}

func (x *RecordFile) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *RecordFile) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

type ResponseList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code     int32         `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message  string        `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Total    uint32        `protobuf:"varint,3,opt,name=total,proto3" json:"total,omitempty"`
	Page     uint32        `protobuf:"varint,4,opt,name=page,proto3" json:"page,omitempty"`
	PageSize uint32        `protobuf:"varint,5,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
	Data     []*RecordFile `protobuf:"bytes,6,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *ResponseList) Reset() {
	*x = ResponseList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mp4_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseList) ProtoMessage() {}

func (x *ResponseList) ProtoReflect() protoreflect.Message {
	mi := &file_mp4_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseList.ProtoReflect.Descriptor instead.
func (*ResponseList) Descriptor() ([]byte, []int) {
	return file_mp4_proto_rawDescGZIP(), []int{2}
}

func (x *ResponseList) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ResponseList) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ResponseList) GetTotal() uint32 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *ResponseList) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ResponseList) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ResponseList) GetData() []*RecordFile {
	if x != nil {
		return x.Data
	}
	return nil
}

type Catalog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StreamPath string                 `protobuf:"bytes,1,opt,name=streamPath,proto3" json:"streamPath,omitempty"`
	Count      uint32                 `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	StartTime  *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=startTime,proto3" json:"startTime,omitempty"`
	EndTime    *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=endTime,proto3" json:"endTime,omitempty"`
}

func (x *Catalog) Reset() {
	*x = Catalog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mp4_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Catalog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Catalog) ProtoMessage() {}

func (x *Catalog) ProtoReflect() protoreflect.Message {
	mi := &file_mp4_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Catalog.ProtoReflect.Descriptor instead.
func (*Catalog) Descriptor() ([]byte, []int) {
	return file_mp4_proto_rawDescGZIP(), []int{3}
}

func (x *Catalog) GetStreamPath() string {
	if x != nil {
		return x.StreamPath
	}
	return ""
}

func (x *Catalog) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *Catalog) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *Catalog) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

type ResponseCatalog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32      `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string     `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    []*Catalog `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *ResponseCatalog) Reset() {
	*x = ResponseCatalog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mp4_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseCatalog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseCatalog) ProtoMessage() {}

func (x *ResponseCatalog) ProtoReflect() protoreflect.Message {
	mi := &file_mp4_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseCatalog.ProtoReflect.Descriptor instead.
func (*ResponseCatalog) Descriptor() ([]byte, []int) {
	return file_mp4_proto_rawDescGZIP(), []int{4}
}

func (x *ResponseCatalog) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ResponseCatalog) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ResponseCatalog) GetData() []*Catalog {
	if x != nil {
		return x.Data
	}
	return nil
}

type ReqRecordDelete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StreamPath string                 `protobuf:"bytes,1,opt,name=streamPath,proto3" json:"streamPath,omitempty"`
	Ids        []uint32               `protobuf:"varint,2,rep,packed,name=ids,proto3" json:"ids,omitempty"`
	StartTime  *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=startTime,proto3" json:"startTime,omitempty"`
	EndTime    *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=endTime,proto3" json:"endTime,omitempty"`
}

func (x *ReqRecordDelete) Reset() {
	*x = ReqRecordDelete{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mp4_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqRecordDelete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqRecordDelete) ProtoMessage() {}

func (x *ReqRecordDelete) ProtoReflect() protoreflect.Message {
	mi := &file_mp4_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqRecordDelete.ProtoReflect.Descriptor instead.
func (*ReqRecordDelete) Descriptor() ([]byte, []int) {
	return file_mp4_proto_rawDescGZIP(), []int{5}
}

func (x *ReqRecordDelete) GetStreamPath() string {
	if x != nil {
		return x.StreamPath
	}
	return ""
}

func (x *ReqRecordDelete) GetIds() []uint32 {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *ReqRecordDelete) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *ReqRecordDelete) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

type ResponseDelete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ResponseDelete) Reset() {
	*x = ResponseDelete{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mp4_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseDelete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseDelete) ProtoMessage() {}

func (x *ResponseDelete) ProtoReflect() protoreflect.Message {
	mi := &file_mp4_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseDelete.ProtoReflect.Descriptor instead.
func (*ResponseDelete) Descriptor() ([]byte, []int) {
	return file_mp4_proto_rawDescGZIP(), []int{6}
}

func (x *ResponseDelete) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ResponseDelete) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_mp4_proto protoreflect.FileDescriptor

var file_mp4_proto_rawDesc = []byte{
	0x0a, 0x09, 0x6d, 0x70, 0x34, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x6d, 0x70, 0x34,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9d, 0x01, 0x0a,
	0x0d, 0x52, 0x65, 0x71, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1e,
	0x0a, 0x0a, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x50, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x50, 0x61, 0x74, 0x68, 0x12, 0x14,
	0x0a, 0x05, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72,
	0x61, 0x6e, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0xc8, 0x01, 0x0a,
	0x0a, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x50, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x50, 0x61, 0x74, 0x68, 0x12, 0x38, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x34, 0x0a, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07,
	0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0xa7, 0x01, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x23, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x70, 0x34,
	0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x22, 0xaf, 0x01, 0x0a, 0x07, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x12, 0x1e, 0x0a,
	0x0a, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x50, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x50, 0x61, 0x74, 0x68, 0x12, 0x14, 0x0a,
	0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x38, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x34, 0x0a,
	0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54,
	0x69, 0x6d, 0x65, 0x22, 0x61, 0x0a, 0x0f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43,
	0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x20, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x70, 0x34, 0x2e, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xb3, 0x01, 0x0a, 0x0f, 0x52, 0x65, 0x71, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x50, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x50, 0x61, 0x74, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x38, 0x0a, 0x09,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x34, 0x0a, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x3e, 0x0a, 0x0e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x8f, 0x02, 0x0a,
	0x03, 0x61, 0x70, 0x69, 0x12, 0x54, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x2e, 0x6d,
	0x70, 0x34, 0x2e, 0x52, 0x65, 0x71, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x4c, 0x69, 0x73, 0x74,
	0x1a, 0x11, 0x2e, 0x6d, 0x70, 0x34, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4c,
	0x69, 0x73, 0x74, 0x22, 0x25, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1f, 0x12, 0x1d, 0x2f, 0x6d, 0x70,
	0x34, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6c, 0x69, 0x73, 0x74, 0x2f, 0x7b, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x50, 0x61, 0x74, 0x68, 0x3d, 0x2a, 0x2a, 0x7d, 0x12, 0x51, 0x0a, 0x07, 0x43, 0x61,
	0x74, 0x61, 0x6c, 0x6f, 0x67, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x14, 0x2e,
	0x6d, 0x70, 0x34, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x61, 0x74, 0x61,
	0x6c, 0x6f, 0x67, 0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x12, 0x10, 0x2f, 0x6d, 0x70,
	0x34, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x12, 0x5f, 0x0a,
	0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x14, 0x2e, 0x6d, 0x70, 0x34, 0x2e, 0x52, 0x65,
	0x71, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x1a, 0x13, 0x2e,
	0x6d, 0x70, 0x34, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x22, 0x2a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x24, 0x22, 0x1f, 0x2f, 0x6d, 0x70, 0x34,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x2f, 0x7b, 0x73, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x50, 0x61, 0x74, 0x68, 0x3d, 0x2a, 0x2a, 0x7d, 0x3a, 0x01, 0x2a, 0x42, 0x1c,
	0x5a, 0x1a, 0x6d, 0x37, 0x73, 0x2e, 0x6c, 0x69, 0x76, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x2f, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x6d, 0x70, 0x34, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mp4_proto_rawDescOnce sync.Once
	file_mp4_proto_rawDescData = file_mp4_proto_rawDesc
)

func file_mp4_proto_rawDescGZIP() []byte {
	file_mp4_proto_rawDescOnce.Do(func() {
		file_mp4_proto_rawDescData = protoimpl.X.CompressGZIP(file_mp4_proto_rawDescData)
	})
	return file_mp4_proto_rawDescData
}

var file_mp4_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_mp4_proto_goTypes = []interface{}{
	(*ReqRecordList)(nil),         // 0: mp4.ReqRecordList
	(*RecordFile)(nil),            // 1: mp4.RecordFile
	(*ResponseList)(nil),          // 2: mp4.ResponseList
	(*Catalog)(nil),               // 3: mp4.Catalog
	(*ResponseCatalog)(nil),       // 4: mp4.ResponseCatalog
	(*ReqRecordDelete)(nil),       // 5: mp4.ReqRecordDelete
	(*ResponseDelete)(nil),        // 6: mp4.ResponseDelete
	(*timestamppb.Timestamp)(nil), // 7: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 8: google.protobuf.Empty
}
var file_mp4_proto_depIdxs = []int32{
	7,  // 0: mp4.RecordFile.startTime:type_name -> google.protobuf.Timestamp
	7,  // 1: mp4.RecordFile.endTime:type_name -> google.protobuf.Timestamp
	1,  // 2: mp4.ResponseList.data:type_name -> mp4.RecordFile
	7,  // 3: mp4.Catalog.startTime:type_name -> google.protobuf.Timestamp
	7,  // 4: mp4.Catalog.endTime:type_name -> google.protobuf.Timestamp
	3,  // 5: mp4.ResponseCatalog.data:type_name -> mp4.Catalog
	7,  // 6: mp4.ReqRecordDelete.startTime:type_name -> google.protobuf.Timestamp
	7,  // 7: mp4.ReqRecordDelete.endTime:type_name -> google.protobuf.Timestamp
	0,  // 8: mp4.api.List:input_type -> mp4.ReqRecordList
	8,  // 9: mp4.api.Catalog:input_type -> google.protobuf.Empty
	5,  // 10: mp4.api.Delete:input_type -> mp4.ReqRecordDelete
	2,  // 11: mp4.api.List:output_type -> mp4.ResponseList
	4,  // 12: mp4.api.Catalog:output_type -> mp4.ResponseCatalog
	6,  // 13: mp4.api.Delete:output_type -> mp4.ResponseDelete
	11, // [11:14] is the sub-list for method output_type
	8,  // [8:11] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_mp4_proto_init() }
func file_mp4_proto_init() {
	if File_mp4_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mp4_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqRecordList); i {
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
		file_mp4_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordFile); i {
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
		file_mp4_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseList); i {
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
		file_mp4_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Catalog); i {
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
		file_mp4_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseCatalog); i {
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
		file_mp4_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqRecordDelete); i {
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
		file_mp4_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseDelete); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_mp4_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mp4_proto_goTypes,
		DependencyIndexes: file_mp4_proto_depIdxs,
		MessageInfos:      file_mp4_proto_msgTypes,
	}.Build()
	File_mp4_proto = out.File
	file_mp4_proto_rawDesc = nil
	file_mp4_proto_goTypes = nil
	file_mp4_proto_depIdxs = nil
}
