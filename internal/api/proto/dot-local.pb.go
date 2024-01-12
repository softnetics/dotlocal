// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: proto/dot-local.proto

package api

import (
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

type CreateMappingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host       *string `protobuf:"bytes,1,req,name=host" json:"host,omitempty"`
	PathPrefix *string `protobuf:"bytes,2,req,name=path_prefix,json=pathPrefix" json:"path_prefix,omitempty"`
	Target     *string `protobuf:"bytes,3,req,name=target" json:"target,omitempty"`
}

func (x *CreateMappingRequest) Reset() {
	*x = CreateMappingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_dot_local_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateMappingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMappingRequest) ProtoMessage() {}

func (x *CreateMappingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_dot_local_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMappingRequest.ProtoReflect.Descriptor instead.
func (*CreateMappingRequest) Descriptor() ([]byte, []int) {
	return file_proto_dot_local_proto_rawDescGZIP(), []int{0}
}

func (x *CreateMappingRequest) GetHost() string {
	if x != nil && x.Host != nil {
		return *x.Host
	}
	return ""
}

func (x *CreateMappingRequest) GetPathPrefix() string {
	if x != nil && x.PathPrefix != nil {
		return *x.PathPrefix
	}
	return ""
}

func (x *CreateMappingRequest) GetTarget() string {
	if x != nil && x.Target != nil {
		return *x.Target
	}
	return ""
}

type MappingKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host       *string `protobuf:"bytes,1,req,name=host" json:"host,omitempty"`
	PathPrefix *string `protobuf:"bytes,2,req,name=path_prefix,json=pathPrefix" json:"path_prefix,omitempty"`
}

func (x *MappingKey) Reset() {
	*x = MappingKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_dot_local_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MappingKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MappingKey) ProtoMessage() {}

func (x *MappingKey) ProtoReflect() protoreflect.Message {
	mi := &file_proto_dot_local_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MappingKey.ProtoReflect.Descriptor instead.
func (*MappingKey) Descriptor() ([]byte, []int) {
	return file_proto_dot_local_proto_rawDescGZIP(), []int{1}
}

func (x *MappingKey) GetHost() string {
	if x != nil && x.Host != nil {
		return *x.Host
	}
	return ""
}

func (x *MappingKey) GetPathPrefix() string {
	if x != nil && x.PathPrefix != nil {
		return *x.PathPrefix
	}
	return ""
}

type Mapping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host       *string                `protobuf:"bytes,1,req,name=host" json:"host,omitempty"`
	PathPrefix *string                `protobuf:"bytes,2,req,name=path_prefix,json=pathPrefix" json:"path_prefix,omitempty"`
	Id         *string                `protobuf:"bytes,3,req,name=id" json:"id,omitempty"`
	Target     *string                `protobuf:"bytes,4,req,name=target" json:"target,omitempty"`
	ExpiresAt  *timestamppb.Timestamp `protobuf:"bytes,5,req,name=expires_at,json=expiresAt" json:"expires_at,omitempty"`
}

func (x *Mapping) Reset() {
	*x = Mapping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_dot_local_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Mapping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Mapping) ProtoMessage() {}

func (x *Mapping) ProtoReflect() protoreflect.Message {
	mi := &file_proto_dot_local_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Mapping.ProtoReflect.Descriptor instead.
func (*Mapping) Descriptor() ([]byte, []int) {
	return file_proto_dot_local_proto_rawDescGZIP(), []int{2}
}

func (x *Mapping) GetHost() string {
	if x != nil && x.Host != nil {
		return *x.Host
	}
	return ""
}

func (x *Mapping) GetPathPrefix() string {
	if x != nil && x.PathPrefix != nil {
		return *x.PathPrefix
	}
	return ""
}

func (x *Mapping) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

func (x *Mapping) GetTarget() string {
	if x != nil && x.Target != nil {
		return *x.Target
	}
	return ""
}

func (x *Mapping) GetExpiresAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ExpiresAt
	}
	return nil
}

type ListMappingsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mappings []*Mapping `protobuf:"bytes,1,rep,name=mappings" json:"mappings,omitempty"`
}

func (x *ListMappingsResponse) Reset() {
	*x = ListMappingsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_dot_local_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListMappingsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMappingsResponse) ProtoMessage() {}

func (x *ListMappingsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_dot_local_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMappingsResponse.ProtoReflect.Descriptor instead.
func (*ListMappingsResponse) Descriptor() ([]byte, []int) {
	return file_proto_dot_local_proto_rawDescGZIP(), []int{3}
}

func (x *ListMappingsResponse) GetMappings() []*Mapping {
	if x != nil {
		return x.Mappings
	}
	return nil
}

var File_proto_dot_local_proto protoreflect.FileDescriptor

var file_proto_dot_local_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x6f, 0x74, 0x2d, 0x6c, 0x6f, 0x63, 0x61,
	0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x63, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d,
	0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73,
	0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78,
	0x18, 0x02, 0x20, 0x02, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x61, 0x74, 0x68, 0x50, 0x72, 0x65, 0x66,
	0x69, 0x78, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x03, 0x20, 0x02,
	0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x22, 0x41, 0x0a, 0x0a, 0x4d, 0x61,
	0x70, 0x70, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74,
	0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b,
	0x70, 0x61, 0x74, 0x68, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x02, 0x20, 0x02, 0x28,
	0x09, 0x52, 0x0a, 0x70, 0x61, 0x74, 0x68, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x22, 0xa1, 0x01,
	0x0a, 0x07, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73,
	0x74, 0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x1f, 0x0a,
	0x0b, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x02, 0x20, 0x02,
	0x28, 0x09, 0x52, 0x0a, 0x70, 0x61, 0x74, 0x68, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x02, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x04, 0x20, 0x02, 0x28, 0x09, 0x52, 0x06,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65,
	0x73, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x02, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x41,
	0x74, 0x22, 0x3c, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x08, 0x6d, 0x61, 0x70,
	0x70, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x4d, 0x61,
	0x70, 0x70, 0x69, 0x6e, 0x67, 0x52, 0x08, 0x6d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x73, 0x32,
	0xb7, 0x01, 0x0a, 0x08, 0x44, 0x6f, 0x74, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x12, 0x32, 0x0a, 0x0d,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x12, 0x15, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x08, 0x2e, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x22, 0x00,
	0x12, 0x36, 0x0a, 0x0d, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e,
	0x67, 0x12, 0x0b, 0x2e, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74,
	0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x15, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x24, 0x5a, 0x22, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x66, 0x74, 0x6e, 0x65, 0x74, 0x69,
	0x63, 0x73, 0x2f, 0x64, 0x6f, 0x74, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69,
}

var (
	file_proto_dot_local_proto_rawDescOnce sync.Once
	file_proto_dot_local_proto_rawDescData = file_proto_dot_local_proto_rawDesc
)

func file_proto_dot_local_proto_rawDescGZIP() []byte {
	file_proto_dot_local_proto_rawDescOnce.Do(func() {
		file_proto_dot_local_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_dot_local_proto_rawDescData)
	})
	return file_proto_dot_local_proto_rawDescData
}

var file_proto_dot_local_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_dot_local_proto_goTypes = []interface{}{
	(*CreateMappingRequest)(nil),  // 0: CreateMappingRequest
	(*MappingKey)(nil),            // 1: MappingKey
	(*Mapping)(nil),               // 2: Mapping
	(*ListMappingsResponse)(nil),  // 3: ListMappingsResponse
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 5: google.protobuf.Empty
}
var file_proto_dot_local_proto_depIdxs = []int32{
	4, // 0: Mapping.expires_at:type_name -> google.protobuf.Timestamp
	2, // 1: ListMappingsResponse.mappings:type_name -> Mapping
	0, // 2: DotLocal.CreateMapping:input_type -> CreateMappingRequest
	1, // 3: DotLocal.RemoveMapping:input_type -> MappingKey
	5, // 4: DotLocal.ListMappings:input_type -> google.protobuf.Empty
	2, // 5: DotLocal.CreateMapping:output_type -> Mapping
	5, // 6: DotLocal.RemoveMapping:output_type -> google.protobuf.Empty
	3, // 7: DotLocal.ListMappings:output_type -> ListMappingsResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_dot_local_proto_init() }
func file_proto_dot_local_proto_init() {
	if File_proto_dot_local_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_dot_local_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateMappingRequest); i {
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
		file_proto_dot_local_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MappingKey); i {
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
		file_proto_dot_local_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Mapping); i {
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
		file_proto_dot_local_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListMappingsResponse); i {
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
			RawDescriptor: file_proto_dot_local_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_dot_local_proto_goTypes,
		DependencyIndexes: file_proto_dot_local_proto_depIdxs,
		MessageInfos:      file_proto_dot_local_proto_msgTypes,
	}.Build()
	File_proto_dot_local_proto = out.File
	file_proto_dot_local_proto_rawDesc = nil
	file_proto_dot_local_proto_goTypes = nil
	file_proto_dot_local_proto_depIdxs = nil
}
