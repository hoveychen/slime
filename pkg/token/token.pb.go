// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: github.com/hoveychen/slime/pkg/token/token.proto

package token

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

type AgentToken struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name       string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	ExpireAt   int64    `protobuf:"varint,3,opt,name=expire_at,json=expireAt,proto3" json:"expire_at,omitempty"`
	ScopePaths []string `protobuf:"bytes,4,rep,name=scope_paths,json=scopePaths,proto3" json:"scope_paths,omitempty"`
	Scopes     []string `protobuf:"bytes,5,rep,name=scopes,proto3" json:"scopes,omitempty"`
}

func (x *AgentToken) Reset() {
	*x = AgentToken{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_hoveychen_slime_pkg_token_token_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentToken) ProtoMessage() {}

func (x *AgentToken) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_hoveychen_slime_pkg_token_token_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentToken.ProtoReflect.Descriptor instead.
func (*AgentToken) Descriptor() ([]byte, []int) {
	return file_github_com_hoveychen_slime_pkg_token_token_proto_rawDescGZIP(), []int{0}
}

func (x *AgentToken) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AgentToken) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AgentToken) GetExpireAt() int64 {
	if x != nil {
		return x.ExpireAt
	}
	return 0
}

func (x *AgentToken) GetScopePaths() []string {
	if x != nil {
		return x.ScopePaths
	}
	return nil
}

func (x *AgentToken) GetScopes() []string {
	if x != nil {
		return x.Scopes
	}
	return nil
}

var File_github_com_hoveychen_slime_pkg_token_token_proto protoreflect.FileDescriptor

var file_github_com_hoveychen_slime_pkg_token_token_proto_rawDesc = []byte{
	0x0a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x6f, 0x76,
	0x65, 0x79, 0x63, 0x68, 0x65, 0x6e, 0x2f, 0x73, 0x6c, 0x69, 0x6d, 0x65, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x86, 0x01, 0x0a, 0x0a, 0x41, 0x67,
	0x65, 0x6e, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x08, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x41, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x63, 0x6f,
	0x70, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a,
	0x73, 0x63, 0x6f, 0x70, 0x65, 0x50, 0x61, 0x74, 0x68, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63,
	0x6f, 0x70, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x73, 0x63, 0x6f, 0x70,
	0x65, 0x73, 0x42, 0x26, 0x5a, 0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x68, 0x6f, 0x76, 0x65, 0x79, 0x63, 0x68, 0x65, 0x6e, 0x2f, 0x73, 0x6c, 0x69, 0x6d, 0x65,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_github_com_hoveychen_slime_pkg_token_token_proto_rawDescOnce sync.Once
	file_github_com_hoveychen_slime_pkg_token_token_proto_rawDescData = file_github_com_hoveychen_slime_pkg_token_token_proto_rawDesc
)

func file_github_com_hoveychen_slime_pkg_token_token_proto_rawDescGZIP() []byte {
	file_github_com_hoveychen_slime_pkg_token_token_proto_rawDescOnce.Do(func() {
		file_github_com_hoveychen_slime_pkg_token_token_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_hoveychen_slime_pkg_token_token_proto_rawDescData)
	})
	return file_github_com_hoveychen_slime_pkg_token_token_proto_rawDescData
}

var file_github_com_hoveychen_slime_pkg_token_token_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_hoveychen_slime_pkg_token_token_proto_goTypes = []interface{}{
	(*AgentToken)(nil), // 0: token.AgentToken
}
var file_github_com_hoveychen_slime_pkg_token_token_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_github_com_hoveychen_slime_pkg_token_token_proto_init() }
func file_github_com_hoveychen_slime_pkg_token_token_proto_init() {
	if File_github_com_hoveychen_slime_pkg_token_token_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_hoveychen_slime_pkg_token_token_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgentToken); i {
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
			RawDescriptor: file_github_com_hoveychen_slime_pkg_token_token_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_hoveychen_slime_pkg_token_token_proto_goTypes,
		DependencyIndexes: file_github_com_hoveychen_slime_pkg_token_token_proto_depIdxs,
		MessageInfos:      file_github_com_hoveychen_slime_pkg_token_token_proto_msgTypes,
	}.Build()
	File_github_com_hoveychen_slime_pkg_token_token_proto = out.File
	file_github_com_hoveychen_slime_pkg_token_token_proto_rawDesc = nil
	file_github_com_hoveychen_slime_pkg_token_token_proto_goTypes = nil
	file_github_com_hoveychen_slime_pkg_token_token_proto_depIdxs = nil
}
