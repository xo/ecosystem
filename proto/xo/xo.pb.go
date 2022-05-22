// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.20.1
// source: xo.proto

package xo

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FieldOverride_IndexType int32

const (
	FieldOverride_NONE   FieldOverride_IndexType = 0
	FieldOverride_INDEX  FieldOverride_IndexType = 1
	FieldOverride_UNIQUE FieldOverride_IndexType = 2
)

// Enum value maps for FieldOverride_IndexType.
var (
	FieldOverride_IndexType_name = map[int32]string{
		0: "NONE",
		1: "INDEX",
		2: "UNIQUE",
	}
	FieldOverride_IndexType_value = map[string]int32{
		"NONE":   0,
		"INDEX":  1,
		"UNIQUE": 2,
	}
)

func (x FieldOverride_IndexType) Enum() *FieldOverride_IndexType {
	p := new(FieldOverride_IndexType)
	*p = x
	return p
}

func (x FieldOverride_IndexType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FieldOverride_IndexType) Descriptor() protoreflect.EnumDescriptor {
	return file_xo_proto_enumTypes[0].Descriptor()
}

func (FieldOverride_IndexType) Type() protoreflect.EnumType {
	return &file_xo_proto_enumTypes[0]
}

func (x FieldOverride_IndexType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FieldOverride_IndexType.Descriptor instead.
func (FieldOverride_IndexType) EnumDescriptor() ([]byte, []int) {
	return file_xo_proto_rawDescGZIP(), []int{3, 0}
}

// OneToMany is an entry to create a one-to-many table for.
type OneToMany struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A unique name identifying the one-to-many field.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// A unique suffix identifying the type of the one-to-many field. The
	// specified field must have an associated table.
	TypeSuffix string `protobuf:"bytes,2,opt,name=type_suffix,json=typeSuffix,proto3" json:"type_suffix,omitempty"`
}

func (x *OneToMany) Reset() {
	*x = OneToMany{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OneToMany) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OneToMany) ProtoMessage() {}

func (x *OneToMany) ProtoReflect() protoreflect.Message {
	mi := &file_xo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OneToMany.ProtoReflect.Descriptor instead.
func (*OneToMany) Descriptor() ([]byte, []int) {
	return file_xo_proto_rawDescGZIP(), []int{0}
}

func (x *OneToMany) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *OneToMany) GetTypeSuffix() string {
	if x != nil {
		return x.TypeSuffix
	}
	return ""
}

// Ref is a reference of a field in another type.
type Ref struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A unique suffix identifying the type of the referenced table.
	TypeSuffix string `protobuf:"bytes,1,opt,name=type_suffix,json=typeSuffix,proto3" json:"type_suffix,omitempty"`
	// Name of the field that the current field references.
	FieldName string `protobuf:"bytes,2,opt,name=field_name,json=fieldName,proto3" json:"field_name,omitempty"`
}

func (x *Ref) Reset() {
	*x = Ref{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ref) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ref) ProtoMessage() {}

func (x *Ref) ProtoReflect() protoreflect.Message {
	mi := &file_xo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ref.ProtoReflect.Descriptor instead.
func (*Ref) Descriptor() ([]byte, []int) {
	return file_xo_proto_rawDescGZIP(), []int{1}
}

func (x *Ref) GetTypeSuffix() string {
	if x != nil {
		return x.TypeSuffix
	}
	return ""
}

func (x *Ref) GetFieldName() string {
	if x != nil {
		return x.FieldName
	}
	return ""
}

// MessageOverride is an override of default marshalling behaviour of
// protoc-gen-xo.
type MessageOverride struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Create the table without a default primary key.
	Manual bool `protobuf:"varint,1,opt,name=manual,proto3" json:"manual,omitempty"`
	// Ignore the annotated message and not create a table for it.
	Ignore bool `protobuf:"varint,2,opt,name=ignore,proto3" json:"ignore,omitempty"`
	// Embeds the message as a JSON instead of creating a new table, whenever
	// referenced.
	EmbedAsJson bool `protobuf:"varint,3,opt,name=embed_as_json,json=embedAsJson,proto3" json:"embed_as_json,omitempty"`
	// A list of one-to-many fields to create tables for.
	HasMany []*OneToMany `protobuf:"bytes,4,rep,name=has_many,json=hasMany,proto3" json:"has_many,omitempty"`
}

func (x *MessageOverride) Reset() {
	*x = MessageOverride{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageOverride) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageOverride) ProtoMessage() {}

func (x *MessageOverride) ProtoReflect() protoreflect.Message {
	mi := &file_xo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageOverride.ProtoReflect.Descriptor instead.
func (*MessageOverride) Descriptor() ([]byte, []int) {
	return file_xo_proto_rawDescGZIP(), []int{2}
}

func (x *MessageOverride) GetManual() bool {
	if x != nil {
		return x.Manual
	}
	return false
}

func (x *MessageOverride) GetIgnore() bool {
	if x != nil {
		return x.Ignore
	}
	return false
}

func (x *MessageOverride) GetEmbedAsJson() bool {
	if x != nil {
		return x.EmbedAsJson
	}
	return false
}

func (x *MessageOverride) GetHasMany() []*OneToMany {
	if x != nil {
		return x.HasMany
	}
	return nil
}

// FieldOverride is an override of default marshalling behaviour of
// protoc-gen-xo.
type FieldOverride struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Marks the index type for the field.
	Index FieldOverride_IndexType `protobuf:"varint,1,opt,name=index,proto3,enum=xo.options.FieldOverride_IndexType" json:"index,omitempty"`
	// Ignore the annotated field and not create a column and/or associated
	// tables for it.
	Ignore bool `protobuf:"varint,2,opt,name=ignore,proto3" json:"ignore,omitempty"`
	// Embeds the field as a JSON instead of creating a new table.
	EmbedAsJson bool `protobuf:"varint,3,opt,name=embed_as_json,json=embedAsJson,proto3" json:"embed_as_json,omitempty"`
	// SQL Expression for the default value for the annotated field.
	DefaultValue string `protobuf:"bytes,4,opt,name=default_value,json=defaultValue,proto3" json:"default_value,omitempty"`
	// The field referenced by the overridden field.
	Ref *Ref `protobuf:"bytes,5,opt,name=ref,proto3" json:"ref,omitempty"`
	// Mark the annotated field as nullable.
	Nullable bool `protobuf:"varint,6,opt,name=nullable,proto3" json:"nullable,omitempty"`
}

func (x *FieldOverride) Reset() {
	*x = FieldOverride{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldOverride) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldOverride) ProtoMessage() {}

func (x *FieldOverride) ProtoReflect() protoreflect.Message {
	mi := &file_xo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldOverride.ProtoReflect.Descriptor instead.
func (*FieldOverride) Descriptor() ([]byte, []int) {
	return file_xo_proto_rawDescGZIP(), []int{3}
}

func (x *FieldOverride) GetIndex() FieldOverride_IndexType {
	if x != nil {
		return x.Index
	}
	return FieldOverride_NONE
}

func (x *FieldOverride) GetIgnore() bool {
	if x != nil {
		return x.Ignore
	}
	return false
}

func (x *FieldOverride) GetEmbedAsJson() bool {
	if x != nil {
		return x.EmbedAsJson
	}
	return false
}

func (x *FieldOverride) GetDefaultValue() string {
	if x != nil {
		return x.DefaultValue
	}
	return ""
}

func (x *FieldOverride) GetRef() *Ref {
	if x != nil {
		return x.Ref
	}
	return nil
}

func (x *FieldOverride) GetNullable() bool {
	if x != nil {
		return x.Nullable
	}
	return false
}

// FileOverride is an override of default marshalling behaviour of
// protoc-gen-xo.
type FileOverride struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Skip the package prefix for all messages within the file.
	SkipPrefix bool `protobuf:"varint,1,opt,name=skip_prefix,json=skipPrefix,proto3" json:"skip_prefix,omitempty"`
}

func (x *FileOverride) Reset() {
	*x = FileOverride{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileOverride) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileOverride) ProtoMessage() {}

func (x *FileOverride) ProtoReflect() protoreflect.Message {
	mi := &file_xo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileOverride.ProtoReflect.Descriptor instead.
func (*FileOverride) Descriptor() ([]byte, []int) {
	return file_xo_proto_rawDescGZIP(), []int{4}
}

func (x *FileOverride) GetSkipPrefix() bool {
	if x != nil {
		return x.SkipPrefix
	}
	return false
}

var file_xo_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*MessageOverride)(nil),
		Field:         1147,
		Name:          "xo.options.msg_overrides",
		Tag:           "bytes,1147,opt,name=msg_overrides",
		Filename:      "xo.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*FieldOverride)(nil),
		Field:         1147,
		Name:          "xo.options.field_overrides",
		Tag:           "bytes,1147,opt,name=field_overrides",
		Filename:      "xo.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*FileOverride)(nil),
		Field:         1147,
		Name:          "xo.options.file_overrides",
		Tag:           "bytes,1147,opt,name=file_overrides",
		Filename:      "xo.proto",
	},
}

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional xo.options.MessageOverride msg_overrides = 1147;
	E_MsgOverrides = &file_xo_proto_extTypes[0]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional xo.options.FieldOverride field_overrides = 1147;
	E_FieldOverrides = &file_xo_proto_extTypes[1]
)

// Extension fields to descriptorpb.FileOptions.
var (
	// optional xo.options.FileOverride file_overrides = 1147;
	E_FileOverrides = &file_xo_proto_extTypes[2]
)

var File_xo_proto protoreflect.FileDescriptor

var file_xo_proto_rawDesc = []byte{
	0x0a, 0x08, 0x78, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x78, 0x6f, 0x2e, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x40, 0x0a, 0x09, 0x4f, 0x6e, 0x65, 0x54,
	0x6f, 0x4d, 0x61, 0x6e, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x79, 0x70,
	0x65, 0x5f, 0x73, 0x75, 0x66, 0x66, 0x69, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x74, 0x79, 0x70, 0x65, 0x53, 0x75, 0x66, 0x66, 0x69, 0x78, 0x22, 0x45, 0x0a, 0x03, 0x52, 0x65,
	0x66, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x73, 0x75, 0x66, 0x66, 0x69, 0x78,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x79, 0x70, 0x65, 0x53, 0x75, 0x66, 0x66,
	0x69, 0x78, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x4e, 0x61, 0x6d,
	0x65, 0x22, 0x97, 0x01, 0x0a, 0x0f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x76, 0x65,
	0x72, 0x72, 0x69, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x61, 0x6e, 0x75, 0x61, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x6d, 0x61, 0x6e, 0x75, 0x61, 0x6c, 0x12, 0x16, 0x0a,
	0x06, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69,
	0x67, 0x6e, 0x6f, 0x72, 0x65, 0x12, 0x22, 0x0a, 0x0d, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x5f, 0x61,
	0x73, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x65, 0x6d,
	0x62, 0x65, 0x64, 0x41, 0x73, 0x4a, 0x73, 0x6f, 0x6e, 0x12, 0x30, 0x0a, 0x08, 0x68, 0x61, 0x73,
	0x5f, 0x6d, 0x61, 0x6e, 0x79, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x78, 0x6f,
	0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4f, 0x6e, 0x65, 0x54, 0x6f, 0x4d, 0x61,
	0x6e, 0x79, 0x52, 0x07, 0x68, 0x61, 0x73, 0x4d, 0x61, 0x6e, 0x79, 0x22, 0x98, 0x02, 0x0a, 0x0d,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x12, 0x39, 0x0a,
	0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x23, 0x2e, 0x78,
	0x6f, 0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f,
	0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x67, 0x6e, 0x6f,
	0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65,
	0x12, 0x22, 0x0a, 0x0d, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x5f, 0x61, 0x73, 0x5f, 0x6a, 0x73, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x41, 0x73,
	0x4a, 0x73, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x65, 0x66,
	0x61, 0x75, 0x6c, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x0a, 0x03, 0x72, 0x65, 0x66,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x78, 0x6f, 0x2e, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x52, 0x65, 0x66, 0x52, 0x03, 0x72, 0x65, 0x66, 0x12, 0x1a, 0x0a, 0x08,
	0x6e, 0x75, 0x6c, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08,
	0x6e, 0x75, 0x6c, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x2c, 0x0a, 0x09, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12,
	0x09, 0x0a, 0x05, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x55, 0x4e,
	0x49, 0x51, 0x55, 0x45, 0x10, 0x02, 0x22, 0x2f, 0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x76,
	0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x6b, 0x69, 0x70, 0x5f, 0x70,
	0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x73, 0x6b, 0x69,
	0x70, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x3a, 0x62, 0x0a, 0x0d, 0x6d, 0x73, 0x67, 0x5f, 0x6f,
	0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x73, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xfb, 0x08, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1b, 0x2e, 0x78, 0x6f, 0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x52, 0x0c, 0x6d,
	0x73, 0x67, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x73, 0x3a, 0x62, 0x0a, 0x0f, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x73, 0x12, 0x1d,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xfb, 0x08,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x78, 0x6f, 0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x52,
	0x0e, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x73, 0x3a,
	0x5e, 0x0a, 0x0e, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65,
	0x73, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0xfb, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x78, 0x6f, 0x2e, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65,
	0x52, 0x0d, 0x66, 0x69, 0x6c, 0x65, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x73, 0x42,
	0x22, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x78, 0x6f,
	0x2f, 0x65, 0x63, 0x6f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x78, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_xo_proto_rawDescOnce sync.Once
	file_xo_proto_rawDescData = file_xo_proto_rawDesc
)

func file_xo_proto_rawDescGZIP() []byte {
	file_xo_proto_rawDescOnce.Do(func() {
		file_xo_proto_rawDescData = protoimpl.X.CompressGZIP(file_xo_proto_rawDescData)
	})
	return file_xo_proto_rawDescData
}

var file_xo_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_xo_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_xo_proto_goTypes = []interface{}{
	(FieldOverride_IndexType)(0),        // 0: xo.options.FieldOverride.IndexType
	(*OneToMany)(nil),                   // 1: xo.options.OneToMany
	(*Ref)(nil),                         // 2: xo.options.Ref
	(*MessageOverride)(nil),             // 3: xo.options.MessageOverride
	(*FieldOverride)(nil),               // 4: xo.options.FieldOverride
	(*FileOverride)(nil),                // 5: xo.options.FileOverride
	(*descriptorpb.MessageOptions)(nil), // 6: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),   // 7: google.protobuf.FieldOptions
	(*descriptorpb.FileOptions)(nil),    // 8: google.protobuf.FileOptions
}
var file_xo_proto_depIdxs = []int32{
	1, // 0: xo.options.MessageOverride.has_many:type_name -> xo.options.OneToMany
	0, // 1: xo.options.FieldOverride.index:type_name -> xo.options.FieldOverride.IndexType
	2, // 2: xo.options.FieldOverride.ref:type_name -> xo.options.Ref
	6, // 3: xo.options.msg_overrides:extendee -> google.protobuf.MessageOptions
	7, // 4: xo.options.field_overrides:extendee -> google.protobuf.FieldOptions
	8, // 5: xo.options.file_overrides:extendee -> google.protobuf.FileOptions
	3, // 6: xo.options.msg_overrides:type_name -> xo.options.MessageOverride
	4, // 7: xo.options.field_overrides:type_name -> xo.options.FieldOverride
	5, // 8: xo.options.file_overrides:type_name -> xo.options.FileOverride
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	6, // [6:9] is the sub-list for extension type_name
	3, // [3:6] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_xo_proto_init() }
func file_xo_proto_init() {
	if File_xo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_xo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OneToMany); i {
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
		file_xo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ref); i {
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
		file_xo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageOverride); i {
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
		file_xo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldOverride); i {
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
		file_xo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileOverride); i {
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
			RawDescriptor: file_xo_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 3,
			NumServices:   0,
		},
		GoTypes:           file_xo_proto_goTypes,
		DependencyIndexes: file_xo_proto_depIdxs,
		EnumInfos:         file_xo_proto_enumTypes,
		MessageInfos:      file_xo_proto_msgTypes,
		ExtensionInfos:    file_xo_proto_extTypes,
	}.Build()
	File_xo_proto = out.File
	file_xo_proto_rawDesc = nil
	file_xo_proto_goTypes = nil
	file_xo_proto_depIdxs = nil
}
