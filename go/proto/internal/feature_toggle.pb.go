// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: feature_toggle.proto

package internal

import (
	feature_toggle "github.com/featureguards/featureguards-go/proto/feature_toggle"
	_ "github.com/golang/protobuf/ptypes/timestamp"
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

type EnvironmentVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Version int64  `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *EnvironmentVersion) Reset() {
	*x = EnvironmentVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feature_toggle_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvironmentVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvironmentVersion) ProtoMessage() {}

func (x *EnvironmentVersion) ProtoReflect() protoreflect.Message {
	mi := &file_feature_toggle_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvironmentVersion.ProtoReflect.Descriptor instead.
func (*EnvironmentVersion) Descriptor() ([]byte, []int) {
	return file_feature_toggle_proto_rawDescGZIP(), []int{0}
}

func (x *EnvironmentVersion) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *EnvironmentVersion) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

type EnvironmentFeatureToggles struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string                          `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	StartingVersion int64                           `protobuf:"varint,2,opt,name=starting_version,json=startingVersion,proto3" json:"starting_version,omitempty"`
	EndingVersion   int64                           `protobuf:"varint,3,opt,name=ending_version,json=endingVersion,proto3" json:"ending_version,omitempty"`
	FeaturToggles   []*feature_toggle.FeatureToggle `protobuf:"bytes,4,rep,name=featur_toggles,json=featurToggles,proto3" json:"featur_toggles,omitempty"`
}

func (x *EnvironmentFeatureToggles) Reset() {
	*x = EnvironmentFeatureToggles{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feature_toggle_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvironmentFeatureToggles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvironmentFeatureToggles) ProtoMessage() {}

func (x *EnvironmentFeatureToggles) ProtoReflect() protoreflect.Message {
	mi := &file_feature_toggle_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvironmentFeatureToggles.ProtoReflect.Descriptor instead.
func (*EnvironmentFeatureToggles) Descriptor() ([]byte, []int) {
	return file_feature_toggle_proto_rawDescGZIP(), []int{1}
}

func (x *EnvironmentFeatureToggles) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *EnvironmentFeatureToggles) GetStartingVersion() int64 {
	if x != nil {
		return x.StartingVersion
	}
	return 0
}

func (x *EnvironmentFeatureToggles) GetEndingVersion() int64 {
	if x != nil {
		return x.EndingVersion
	}
	return 0
}

func (x *EnvironmentFeatureToggles) GetFeaturToggles() []*feature_toggle.FeatureToggle {
	if x != nil {
		return x.FeaturToggles
	}
	return nil
}

var File_feature_toggle_proto protoreflect.FileDescriptor

var file_feature_toggle_proto_rawDesc = []byte{
	0x0a, 0x14, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1b, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x5f, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3e,
	0x0a, 0x12, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xc3,
	0x01, 0x0a, 0x19, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x46, 0x65,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x73, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x29, 0x0a, 0x10,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e, 0x64, 0x69, 0x6e,
	0x67, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0d, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x44,
	0x0a, 0x0e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x5f, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x5f, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x2e, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x54,
	0x6f, 0x67, 0x67, 0x6c, 0x65, 0x52, 0x0d, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x54, 0x6f, 0x67,
	0x67, 0x6c, 0x65, 0x73, 0x42, 0x1c, 0x5a, 0x1a, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_feature_toggle_proto_rawDescOnce sync.Once
	file_feature_toggle_proto_rawDescData = file_feature_toggle_proto_rawDesc
)

func file_feature_toggle_proto_rawDescGZIP() []byte {
	file_feature_toggle_proto_rawDescOnce.Do(func() {
		file_feature_toggle_proto_rawDescData = protoimpl.X.CompressGZIP(file_feature_toggle_proto_rawDescData)
	})
	return file_feature_toggle_proto_rawDescData
}

var file_feature_toggle_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_feature_toggle_proto_goTypes = []interface{}{
	(*EnvironmentVersion)(nil),           // 0: internal.EnvironmentVersion
	(*EnvironmentFeatureToggles)(nil),    // 1: internal.EnvironmentFeatureToggles
	(*feature_toggle.FeatureToggle)(nil), // 2: feature_toggle.FeatureToggle
}
var file_feature_toggle_proto_depIdxs = []int32{
	2, // 0: internal.EnvironmentFeatureToggles.featur_toggles:type_name -> feature_toggle.FeatureToggle
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_feature_toggle_proto_init() }
func file_feature_toggle_proto_init() {
	if File_feature_toggle_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_feature_toggle_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvironmentVersion); i {
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
		file_feature_toggle_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvironmentFeatureToggles); i {
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
			RawDescriptor: file_feature_toggle_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_feature_toggle_proto_goTypes,
		DependencyIndexes: file_feature_toggle_proto_depIdxs,
		MessageInfos:      file_feature_toggle_proto_msgTypes,
	}.Build()
	File_feature_toggle_proto = out.File
	file_feature_toggle_proto_rawDesc = nil
	file_feature_toggle_proto_goTypes = nil
	file_feature_toggle_proto_depIdxs = nil
}
