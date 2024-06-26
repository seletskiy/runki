// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.3
// source: message.proto

package messages

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

type AddOrUpdateRequest_AddMode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NotetypeId int64 `protobuf:"varint,1,opt,name=notetype_id,json=notetypeId,proto3" json:"notetype_id,omitempty"`
	DeckId     int64 `protobuf:"varint,2,opt,name=deck_id,json=deckId,proto3" json:"deck_id,omitempty"`
}

func (x *AddOrUpdateRequest_AddMode) Reset() {
	*x = AddOrUpdateRequest_AddMode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddOrUpdateRequest_AddMode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddOrUpdateRequest_AddMode) ProtoMessage() {}

func (x *AddOrUpdateRequest_AddMode) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddOrUpdateRequest_AddMode.ProtoReflect.Descriptor instead.
func (*AddOrUpdateRequest_AddMode) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{0}
}

func (x *AddOrUpdateRequest_AddMode) GetNotetypeId() int64 {
	if x != nil {
		return x.NotetypeId
	}
	return 0
}

func (x *AddOrUpdateRequest_AddMode) GetDeckId() int64 {
	if x != nil {
		return x.DeckId
	}
	return 0
}

type AddOrUpdateRequest_EditMode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NoteId int64 `protobuf:"varint,1,opt,name=note_id,json=noteId,proto3" json:"note_id,omitempty"`
}

func (x *AddOrUpdateRequest_EditMode) Reset() {
	*x = AddOrUpdateRequest_EditMode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddOrUpdateRequest_EditMode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddOrUpdateRequest_EditMode) ProtoMessage() {}

func (x *AddOrUpdateRequest_EditMode) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddOrUpdateRequest_EditMode.ProtoReflect.Descriptor instead.
func (*AddOrUpdateRequest_EditMode) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{1}
}

func (x *AddOrUpdateRequest_EditMode) GetNoteId() int64 {
	if x != nil {
		return x.NoteId
	}
	return 0
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fields []string `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty"`
	Tags   string   `protobuf:"bytes,2,opt,name=tags,proto3" json:"tags,omitempty"`
	// Types that are assignable to Mode:
	//
	//	*Message_Add
	//	*Message_Edit
	Mode isMessage_Mode `protobuf_oneof:"mode"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{2}
}

func (x *Message) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *Message) GetTags() string {
	if x != nil {
		return x.Tags
	}
	return ""
}

func (m *Message) GetMode() isMessage_Mode {
	if m != nil {
		return m.Mode
	}
	return nil
}

func (x *Message) GetAdd() *AddOrUpdateRequest_AddMode {
	if x, ok := x.GetMode().(*Message_Add); ok {
		return x.Add
	}
	return nil
}

func (x *Message) GetEdit() *AddOrUpdateRequest_EditMode {
	if x, ok := x.GetMode().(*Message_Edit); ok {
		return x.Edit
	}
	return nil
}

type isMessage_Mode interface {
	isMessage_Mode()
}

type Message_Add struct {
	Add *AddOrUpdateRequest_AddMode `protobuf:"bytes,3,opt,name=add,proto3,oneof"`
}

type Message_Edit struct {
	Edit *AddOrUpdateRequest_EditMode `protobuf:"bytes,4,opt,name=edit,proto3,oneof"`
}

func (*Message_Add) isMessage_Mode() {}

func (*Message_Edit) isMessage_Mode() {}

var File_message_proto protoreflect.FileDescriptor

var file_message_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x56, 0x0a, 0x1a, 0x41, 0x64, 0x64, 0x4f, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x41, 0x64, 0x64, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x1f, 0x0a,
	0x0b, 0x6e, 0x6f, 0x74, 0x65, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x6e, 0x6f, 0x74, 0x65, 0x74, 0x79, 0x70, 0x65, 0x49, 0x64, 0x12, 0x17,
	0x0a, 0x07, 0x64, 0x65, 0x63, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x64, 0x65, 0x63, 0x6b, 0x49, 0x64, 0x22, 0x36, 0x0a, 0x1b, 0x41, 0x64, 0x64, 0x4f, 0x72,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x45, 0x64,
	0x69, 0x74, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x74, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6e, 0x6f, 0x74, 0x65, 0x49, 0x64, 0x22,
	0xa2, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x2f, 0x0a, 0x03, 0x61, 0x64, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x41, 0x64, 0x64, 0x4f, 0x72, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x41, 0x64, 0x64, 0x4d, 0x6f, 0x64,
	0x65, 0x48, 0x00, 0x52, 0x03, 0x61, 0x64, 0x64, 0x12, 0x32, 0x0a, 0x04, 0x65, 0x64, 0x69, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x41, 0x64, 0x64, 0x4f, 0x72, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x45, 0x64, 0x69, 0x74,
	0x4d, 0x6f, 0x64, 0x65, 0x48, 0x00, 0x52, 0x04, 0x65, 0x64, 0x69, 0x74, 0x42, 0x06, 0x0a, 0x04,
	0x6d, 0x6f, 0x64, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_proto_rawDescOnce sync.Once
	file_message_proto_rawDescData = file_message_proto_rawDesc
)

func file_message_proto_rawDescGZIP() []byte {
	file_message_proto_rawDescOnce.Do(func() {
		file_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_proto_rawDescData)
	})
	return file_message_proto_rawDescData
}

var file_message_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_message_proto_goTypes = []interface{}{
	(*AddOrUpdateRequest_AddMode)(nil),  // 0: AddOrUpdateRequest_AddMode
	(*AddOrUpdateRequest_EditMode)(nil), // 1: AddOrUpdateRequest_EditMode
	(*Message)(nil),                     // 2: Message
}
var file_message_proto_depIdxs = []int32{
	0, // 0: Message.add:type_name -> AddOrUpdateRequest_AddMode
	1, // 1: Message.edit:type_name -> AddOrUpdateRequest_EditMode
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_message_proto_init() }
func file_message_proto_init() {
	if File_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddOrUpdateRequest_AddMode); i {
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
		file_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddOrUpdateRequest_EditMode); i {
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
		file_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
	file_message_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Message_Add)(nil),
		(*Message_Edit)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_proto_goTypes,
		DependencyIndexes: file_message_proto_depIdxs,
		MessageInfos:      file_message_proto_msgTypes,
	}.Build()
	File_message_proto = out.File
	file_message_proto_rawDesc = nil
	file_message_proto_goTypes = nil
	file_message_proto_depIdxs = nil
}
