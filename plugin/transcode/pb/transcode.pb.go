// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.21.12
// source: transcode.proto

package pb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	pb "m7s.live/v5/pb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OverlayConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OverlayStream   string `protobuf:"bytes,1,opt,name=overlay_stream,json=overlayStream,proto3" json:"overlay_stream,omitempty"`       // 叠加流 可为空
	OverlayRegion   string `protobuf:"bytes,2,opt,name=overlay_region,json=overlayRegion,proto3" json:"overlay_region,omitempty"`       // x,y,w,h 可为空,所有区域
	OverlayImage    string `protobuf:"bytes,3,opt,name=overlay_image,json=overlayImage,proto3" json:"overlay_image,omitempty"`          // 图片 base64  可为空 如果图片和视频流都有，则使用图片
	OverlayPosition string `protobuf:"bytes,4,opt,name=overlay_position,json=overlayPosition,proto3" json:"overlay_position,omitempty"` // 位置 x,y
	Text            string `protobuf:"bytes,5,opt,name=text,proto3" json:"text,omitempty"`                                              // 文字
	TimeOffset      int64  `protobuf:"varint,6,opt,name=time_offset,json=timeOffset,proto3" json:"time_offset,omitempty"`               // 时间偏移
	TimeFormat      string `protobuf:"bytes,7,opt,name=time_format,json=timeFormat,proto3" json:"time_format,omitempty"`                // 时间格式
	FontName        string `protobuf:"bytes,8,opt,name=font_name,json=fontName,proto3" json:"font_name,omitempty"`                      // 字体文件名
	FontSize        string `protobuf:"bytes,9,opt,name=font_size,json=fontSize,proto3" json:"font_size,omitempty"`                      // 字体大小
	FontColor       string `protobuf:"bytes,10,opt,name=font_color,json=fontColor,proto3" json:"font_color,omitempty"`                  // r,g,b 颜色
	TextPosition    string `protobuf:"bytes,11,opt,name=text_position,json=textPosition,proto3" json:"text_position,omitempty"`         // x,y 文字在图片上的位置
	LineSpacing     string `protobuf:"bytes,12,opt,name=line_spacing,json=lineSpacing,proto3" json:"line_spacing,omitempty"`            //文字行间距
}

func (x *OverlayConfig) Reset() {
	*x = OverlayConfig{}
	mi := &file_transcode_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OverlayConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OverlayConfig) ProtoMessage() {}

func (x *OverlayConfig) ProtoReflect() protoreflect.Message {
	mi := &file_transcode_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OverlayConfig.ProtoReflect.Descriptor instead.
func (*OverlayConfig) Descriptor() ([]byte, []int) {
	return file_transcode_proto_rawDescGZIP(), []int{0}
}

func (x *OverlayConfig) GetOverlayStream() string {
	if x != nil {
		return x.OverlayStream
	}
	return ""
}

func (x *OverlayConfig) GetOverlayRegion() string {
	if x != nil {
		return x.OverlayRegion
	}
	return ""
}

func (x *OverlayConfig) GetOverlayImage() string {
	if x != nil {
		return x.OverlayImage
	}
	return ""
}

func (x *OverlayConfig) GetOverlayPosition() string {
	if x != nil {
		return x.OverlayPosition
	}
	return ""
}

func (x *OverlayConfig) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *OverlayConfig) GetTimeOffset() int64 {
	if x != nil {
		return x.TimeOffset
	}
	return 0
}

func (x *OverlayConfig) GetTimeFormat() string {
	if x != nil {
		return x.TimeFormat
	}
	return ""
}

func (x *OverlayConfig) GetFontName() string {
	if x != nil {
		return x.FontName
	}
	return ""
}

func (x *OverlayConfig) GetFontSize() string {
	if x != nil {
		return x.FontSize
	}
	return ""
}

func (x *OverlayConfig) GetFontColor() string {
	if x != nil {
		return x.FontColor
	}
	return ""
}

func (x *OverlayConfig) GetTextPosition() string {
	if x != nil {
		return x.TextPosition
	}
	return ""
}

func (x *OverlayConfig) GetLineSpacing() string {
	if x != nil {
		return x.LineSpacing
	}
	return ""
}

type TransRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GlobalOptions  string           `protobuf:"bytes,1,opt,name=global_options,json=globalOptions,proto3" json:"global_options,omitempty"` // 全局参数
	SrcStream      string           `protobuf:"bytes,2,opt,name=src_stream,json=srcStream,proto3" json:"src_stream,omitempty"`             // 原始流
	DstStream      string           `protobuf:"bytes,3,opt,name=dst_stream,json=dstStream,proto3" json:"dst_stream,omitempty"`             // 输出流
	Encodec        string           `protobuf:"bytes,4,opt,name=encodec,proto3" json:"encodec,omitempty"`
	Decodec        string           `protobuf:"bytes,5,opt,name=decodec,proto3" json:"decodec,omitempty"`
	Scale          string           `protobuf:"bytes,6,opt,name=scale,proto3" json:"scale,omitempty"`
	OverlayConfigs []*OverlayConfig `protobuf:"bytes,8,rep,name=overlay_configs,json=overlayConfigs,proto3" json:"overlay_configs,omitempty"`
}

func (x *TransRequest) Reset() {
	*x = TransRequest{}
	mi := &file_transcode_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransRequest) ProtoMessage() {}

func (x *TransRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transcode_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransRequest.ProtoReflect.Descriptor instead.
func (*TransRequest) Descriptor() ([]byte, []int) {
	return file_transcode_proto_rawDescGZIP(), []int{1}
}

func (x *TransRequest) GetGlobalOptions() string {
	if x != nil {
		return x.GlobalOptions
	}
	return ""
}

func (x *TransRequest) GetSrcStream() string {
	if x != nil {
		return x.SrcStream
	}
	return ""
}

func (x *TransRequest) GetDstStream() string {
	if x != nil {
		return x.DstStream
	}
	return ""
}

func (x *TransRequest) GetEncodec() string {
	if x != nil {
		return x.Encodec
	}
	return ""
}

func (x *TransRequest) GetDecodec() string {
	if x != nil {
		return x.Decodec
	}
	return ""
}

func (x *TransRequest) GetScale() string {
	if x != nil {
		return x.Scale
	}
	return ""
}

func (x *TransRequest) GetOverlayConfigs() []*OverlayConfig {
	if x != nil {
		return x.OverlayConfigs
	}
	return nil
}

type TransTwin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrcStream string `protobuf:"bytes,1,opt,name=src_stream,proto3" json:"src_stream,omitempty"` // 原始流
	DstStream string `protobuf:"bytes,2,opt,name=dst_stream,proto3" json:"dst_stream,omitempty"` // 输出流
}

func (x *TransTwin) Reset() {
	*x = TransTwin{}
	mi := &file_transcode_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransTwin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransTwin) ProtoMessage() {}

func (x *TransTwin) ProtoReflect() protoreflect.Message {
	mi := &file_transcode_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransTwin.ProtoReflect.Descriptor instead.
func (*TransTwin) Descriptor() ([]byte, []int) {
	return file_transcode_proto_rawDescGZIP(), []int{2}
}

func (x *TransTwin) GetSrcStream() string {
	if x != nil {
		return x.SrcStream
	}
	return ""
}

func (x *TransTwin) GetDstStream() string {
	if x != nil {
		return x.DstStream
	}
	return ""
}

type TransListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32        `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string       `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    []*TransTwin `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *TransListResponse) Reset() {
	*x = TransListResponse{}
	mi := &file_transcode_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransListResponse) ProtoMessage() {}

func (x *TransListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_transcode_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransListResponse.ProtoReflect.Descriptor instead.
func (*TransListResponse) Descriptor() ([]byte, []int) {
	return file_transcode_proto_rawDescGZIP(), []int{3}
}

func (x *TransListResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *TransListResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *TransListResponse) GetData() []*TransTwin {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_transcode_proto protoreflect.FileDescriptor

var file_transcode_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa4, 0x03, 0x0a, 0x0d, 0x4f, 0x76, 0x65, 0x72, 0x6c, 0x61,
	0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x25, 0x0a, 0x0e, 0x6f, 0x76, 0x65, 0x72, 0x6c,
	0x61, 0x79, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x25,
	0x0a, 0x0e, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x79, 0x5f, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x79, 0x52,
	0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x79,
	0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6f, 0x76,
	0x65, 0x72, 0x6c, 0x61, 0x79, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x6f, 0x76,
	0x65, 0x72, 0x6c, 0x61, 0x79, 0x5f, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x79, 0x50, 0x6f, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x69, 0x6d,
	0x65, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a,
	0x74, 0x69, 0x6d, 0x65, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x69,
	0x6d, 0x65, 0x5f, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x74, 0x69, 0x6d, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66,
	0x6f, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x66, 0x6f, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x6f, 0x6e, 0x74,
	0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x6f, 0x6e,
	0x74, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x6f, 0x6e, 0x74, 0x5f, 0x63, 0x6f,
	0x6c, 0x6f, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x6f, 0x6e, 0x74, 0x43,
	0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x6f, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x65, 0x78,
	0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x6c, 0x69, 0x6e,
	0x65, 0x5f, 0x73, 0x70, 0x61, 0x63, 0x69, 0x6e, 0x67, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x6c, 0x69, 0x6e, 0x65, 0x53, 0x70, 0x61, 0x63, 0x69, 0x6e, 0x67, 0x22, 0x80, 0x02, 0x0a,
	0x0c, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a,
	0x0e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x72, 0x63, 0x5f, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x72, 0x63, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x73, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x73, 0x74, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x12, 0x18, 0x0a, 0x07,
	0x64, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64,
	0x65, 0x63, 0x6f, 0x64, 0x65, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x12, 0x41, 0x0a, 0x0f,
	0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x79, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x18,
	0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x6f, 0x64,
	0x65, 0x2e, 0x4f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52,
	0x0e, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x22,
	0x4b, 0x0a, 0x09, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x54, 0x77, 0x69, 0x6e, 0x12, 0x1e, 0x0a, 0x0a,
	0x73, 0x72, 0x63, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x73, 0x72, 0x63, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x1e, 0x0a, 0x0a,
	0x64, 0x73, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x64, 0x73, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x22, 0x6b, 0x0a, 0x11,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x28, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x54,
	0x77, 0x69, 0x6e, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0xed, 0x02, 0x0a, 0x03, 0x61, 0x70,
	0x69, 0x12, 0x5c, 0x0a, 0x06, 0x6c, 0x61, 0x75, 0x6e, 0x63, 0x68, 0x12, 0x17, 0x2e, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x53, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x20, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x3a, 0x01, 0x2a, 0x22, 0x15, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x63, 0x6f, 0x64, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6c, 0x61, 0x75, 0x6e, 0x63, 0x68, 0x12,
	0x57, 0x0a, 0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x14, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x63, 0x6f, 0x64, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x54, 0x77, 0x69, 0x6e, 0x1a, 0x17,
	0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x3a,
	0x01, 0x2a, 0x22, 0x14, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x59, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1c, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x63, 0x6f, 0x64, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x12, 0x13,
	0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6c,
	0x69, 0x73, 0x74, 0x12, 0x54, 0x0a, 0x05, 0x65, 0x78, 0x69, 0x73, 0x74, 0x12, 0x14, 0x2e, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x54, 0x77,
	0x69, 0x6e, 0x1a, 0x17, 0x2e, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x53, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1c, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x16, 0x12, 0x14, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x65, 0x78, 0x69, 0x73, 0x74, 0x42, 0x25, 0x5a, 0x23, 0x6d, 0x37, 0x73,
	0x2e, 0x6c, 0x69, 0x76, 0x65, 0x2f, 0x6d, 0x37, 0x73, 0x2f, 0x76, 0x35, 0x2f, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x6f, 0x64, 0x65, 0x2f, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_transcode_proto_rawDescOnce sync.Once
	file_transcode_proto_rawDescData = file_transcode_proto_rawDesc
)

func file_transcode_proto_rawDescGZIP() []byte {
	file_transcode_proto_rawDescOnce.Do(func() {
		file_transcode_proto_rawDescData = protoimpl.X.CompressGZIP(file_transcode_proto_rawDescData)
	})
	return file_transcode_proto_rawDescData
}

var file_transcode_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_transcode_proto_goTypes = []any{
	(*OverlayConfig)(nil),      // 0: transcode.OverlayConfig
	(*TransRequest)(nil),       // 1: transcode.TransRequest
	(*TransTwin)(nil),          // 2: transcode.TransTwin
	(*TransListResponse)(nil),  // 3: transcode.TransListResponse
	(*emptypb.Empty)(nil),      // 4: google.protobuf.Empty
	(*pb.SuccessResponse)(nil), // 5: global.SuccessResponse
}
var file_transcode_proto_depIdxs = []int32{
	0, // 0: transcode.TransRequest.overlay_configs:type_name -> transcode.OverlayConfig
	2, // 1: transcode.TransListResponse.data:type_name -> transcode.TransTwin
	1, // 2: transcode.api.launch:input_type -> transcode.TransRequest
	2, // 3: transcode.api.close:input_type -> transcode.TransTwin
	4, // 4: transcode.api.list:input_type -> google.protobuf.Empty
	2, // 5: transcode.api.exist:input_type -> transcode.TransTwin
	5, // 6: transcode.api.launch:output_type -> global.SuccessResponse
	5, // 7: transcode.api.close:output_type -> global.SuccessResponse
	3, // 8: transcode.api.list:output_type -> transcode.TransListResponse
	5, // 9: transcode.api.exist:output_type -> global.SuccessResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_transcode_proto_init() }
func file_transcode_proto_init() {
	if File_transcode_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_transcode_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_transcode_proto_goTypes,
		DependencyIndexes: file_transcode_proto_depIdxs,
		MessageInfos:      file_transcode_proto_msgTypes,
	}.Build()
	File_transcode_proto = out.File
	file_transcode_proto_rawDesc = nil
	file_transcode_proto_goTypes = nil
	file_transcode_proto_depIdxs = nil
}
