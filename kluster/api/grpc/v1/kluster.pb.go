// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v5.29.3
// source: api/grpc/v1/kluster.proto

package v1

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

type CreateNamespaceRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Labels        map[string]string      `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Annotations   map[string]string      `protobuf:"bytes,3,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateNamespaceRequest) Reset() {
	*x = CreateNamespaceRequest{}
	mi := &file_api_grpc_v1_kluster_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateNamespaceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateNamespaceRequest) ProtoMessage() {}

func (x *CreateNamespaceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_v1_kluster_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateNamespaceRequest.ProtoReflect.Descriptor instead.
func (*CreateNamespaceRequest) Descriptor() ([]byte, []int) {
	return file_api_grpc_v1_kluster_proto_rawDescGZIP(), []int{0}
}

func (x *CreateNamespaceRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateNamespaceRequest) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *CreateNamespaceRequest) GetAnnotations() map[string]string {
	if x != nil {
		return x.Annotations
	}
	return nil
}

type CreateNamespaceResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Pong          string                 `protobuf:"bytes,1,opt,name=pong,proto3" json:"pong,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateNamespaceResponse) Reset() {
	*x = CreateNamespaceResponse{}
	mi := &file_api_grpc_v1_kluster_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateNamespaceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateNamespaceResponse) ProtoMessage() {}

func (x *CreateNamespaceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_v1_kluster_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateNamespaceResponse.ProtoReflect.Descriptor instead.
func (*CreateNamespaceResponse) Descriptor() ([]byte, []int) {
	return file_api_grpc_v1_kluster_proto_rawDescGZIP(), []int{1}
}

func (x *CreateNamespaceResponse) GetPong() string {
	if x != nil {
		return x.Pong
	}
	return ""
}

type SyncOpenTelemetryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ClusterName   string                 `protobuf:"bytes,1,opt,name=cluster_name,json=clusterName,proto3" json:"cluster_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncOpenTelemetryRequest) Reset() {
	*x = SyncOpenTelemetryRequest{}
	mi := &file_api_grpc_v1_kluster_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncOpenTelemetryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncOpenTelemetryRequest) ProtoMessage() {}

func (x *SyncOpenTelemetryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_v1_kluster_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncOpenTelemetryRequest.ProtoReflect.Descriptor instead.
func (*SyncOpenTelemetryRequest) Descriptor() ([]byte, []int) {
	return file_api_grpc_v1_kluster_proto_rawDescGZIP(), []int{2}
}

func (x *SyncOpenTelemetryRequest) GetClusterName() string {
	if x != nil {
		return x.ClusterName
	}
	return ""
}

type SyncOpenTelemetryResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Pong          string                 `protobuf:"bytes,1,opt,name=pong,proto3" json:"pong,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncOpenTelemetryResponse) Reset() {
	*x = SyncOpenTelemetryResponse{}
	mi := &file_api_grpc_v1_kluster_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncOpenTelemetryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncOpenTelemetryResponse) ProtoMessage() {}

func (x *SyncOpenTelemetryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_v1_kluster_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncOpenTelemetryResponse.ProtoReflect.Descriptor instead.
func (*SyncOpenTelemetryResponse) Descriptor() ([]byte, []int) {
	return file_api_grpc_v1_kluster_proto_rawDescGZIP(), []int{3}
}

func (x *SyncOpenTelemetryResponse) GetPong() string {
	if x != nil {
		return x.Pong
	}
	return ""
}

var File_api_grpc_v1_kluster_proto protoreflect.FileDescriptor

var file_api_grpc_v1_kluster_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x6b, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x75, 0x6c, 0x74,
	0x61, 0x72, 0x79, 0x2e, 0x6b, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x22, 0xd4,
	0x02, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x4d, 0x0a,
	0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x35, 0x2e,
	0x75, 0x6c, 0x74, 0x61, 0x72, 0x79, 0x2e, 0x6b, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x5c, 0x0a, 0x0b,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x3a, 0x2e, 0x75, 0x6c, 0x74, 0x61, 0x72, 0x79, 0x2e, 0x6b, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3e, 0x0a, 0x10, 0x41, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x2d, 0x0a, 0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x70, 0x6f, 0x6e, 0x67, 0x22, 0x3d, 0x0a, 0x18, 0x53, 0x79, 0x6e, 0x63, 0x4f, 0x70, 0x65, 0x6e,
	0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0x2f, 0x0a, 0x19, 0x53, 0x79, 0x6e, 0x63, 0x4f, 0x70, 0x65, 0x6e, 0x54,
	0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x70, 0x6f, 0x6e, 0x67, 0x32, 0xe7, 0x01, 0x0a, 0x07, 0x4b, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x12, 0x6a, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x12, 0x29, 0x2e, 0x75, 0x6c, 0x74, 0x61, 0x72, 0x79, 0x2e, 0x6b, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a,
	0x2e, 0x75, 0x6c, 0x74, 0x61, 0x72, 0x79, 0x2e, 0x6b, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x70, 0x0a, 0x11,
	0x53, 0x79, 0x6e, 0x63, 0x4f, 0x70, 0x65, 0x6e, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72,
	0x79, 0x12, 0x2b, 0x2e, 0x75, 0x6c, 0x74, 0x61, 0x72, 0x79, 0x2e, 0x6b, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x4f, 0x70, 0x65, 0x6e, 0x54, 0x65,
	0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c,
	0x2e, 0x75, 0x6c, 0x74, 0x61, 0x72, 0x79, 0x2e, 0x6b, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x4f, 0x70, 0x65, 0x6e, 0x54, 0x65, 0x6c, 0x65, 0x6d,
	0x65, 0x74, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x63,
	0x0a, 0x15, 0x63, 0x6f, 0x2e, 0x75, 0x6c, 0x74, 0x61, 0x72, 0x79, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x4b, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x75, 0x6c, 0x74, 0x61, 0x72, 0x79, 0x2f, 0x70, 0x6f, 0x6c, 0x79, 0x6b,
	0x75, 0x62, 0x65, 0x2f, 0x6b, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0xa2, 0x02, 0x09, 0x4b, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_grpc_v1_kluster_proto_rawDescOnce sync.Once
	file_api_grpc_v1_kluster_proto_rawDescData = file_api_grpc_v1_kluster_proto_rawDesc
)

func file_api_grpc_v1_kluster_proto_rawDescGZIP() []byte {
	file_api_grpc_v1_kluster_proto_rawDescOnce.Do(func() {
		file_api_grpc_v1_kluster_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_grpc_v1_kluster_proto_rawDescData)
	})
	return file_api_grpc_v1_kluster_proto_rawDescData
}

var file_api_grpc_v1_kluster_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_grpc_v1_kluster_proto_goTypes = []any{
	(*CreateNamespaceRequest)(nil),    // 0: ultary.kluster.v1.CreateNamespaceRequest
	(*CreateNamespaceResponse)(nil),   // 1: ultary.kluster.v1.CreateNamespaceResponse
	(*SyncOpenTelemetryRequest)(nil),  // 2: ultary.kluster.v1.SyncOpenTelemetryRequest
	(*SyncOpenTelemetryResponse)(nil), // 3: ultary.kluster.v1.SyncOpenTelemetryResponse
	nil,                               // 4: ultary.kluster.v1.CreateNamespaceRequest.LabelsEntry
	nil,                               // 5: ultary.kluster.v1.CreateNamespaceRequest.AnnotationsEntry
}
var file_api_grpc_v1_kluster_proto_depIdxs = []int32{
	4, // 0: ultary.kluster.v1.CreateNamespaceRequest.labels:type_name -> ultary.kluster.v1.CreateNamespaceRequest.LabelsEntry
	5, // 1: ultary.kluster.v1.CreateNamespaceRequest.annotations:type_name -> ultary.kluster.v1.CreateNamespaceRequest.AnnotationsEntry
	0, // 2: ultary.kluster.v1.Kluster.CreateNamespace:input_type -> ultary.kluster.v1.CreateNamespaceRequest
	2, // 3: ultary.kluster.v1.Kluster.SyncOpenTelemetry:input_type -> ultary.kluster.v1.SyncOpenTelemetryRequest
	1, // 4: ultary.kluster.v1.Kluster.CreateNamespace:output_type -> ultary.kluster.v1.CreateNamespaceResponse
	3, // 5: ultary.kluster.v1.Kluster.SyncOpenTelemetry:output_type -> ultary.kluster.v1.SyncOpenTelemetryResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_grpc_v1_kluster_proto_init() }
func file_api_grpc_v1_kluster_proto_init() {
	if File_api_grpc_v1_kluster_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_grpc_v1_kluster_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_grpc_v1_kluster_proto_goTypes,
		DependencyIndexes: file_api_grpc_v1_kluster_proto_depIdxs,
		MessageInfos:      file_api_grpc_v1_kluster_proto_msgTypes,
	}.Build()
	File_api_grpc_v1_kluster_proto = out.File
	file_api_grpc_v1_kluster_proto_rawDesc = nil
	file_api_grpc_v1_kluster_proto_goTypes = nil
	file_api_grpc_v1_kluster_proto_depIdxs = nil
}
