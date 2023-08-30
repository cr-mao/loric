// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: game.proto

package pb

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

// 路由号枚举
type Route int32

const (
	Route_Login             Route = 0 // 登录路由号
	Route_LianmentChatEnter Route = 1 //进入联盟聊天室
	Route_LianmengChat      Route = 2 // 发送消息
)

// Enum value maps for Route.
var (
	Route_name = map[int32]string{
		0: "Login",
		1: "LianmentChatEnter",
		2: "LianmengChat",
	}
	Route_value = map[string]int32{
		"Login":             0,
		"LianmentChatEnter": 1,
		"LianmengChat":      2,
	}
)

func (x Route) Enum() *Route {
	p := new(Route)
	*p = x
	return p
}

func (x Route) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Route) Descriptor() protoreflect.EnumDescriptor {
	return file_game_proto_enumTypes[0].Descriptor()
}

func (Route) Type() protoreflect.EnumType {
	return &file_game_proto_enumTypes[0]
}

func (x Route) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Route.Descriptor instead.
func (Route) EnumDescriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{0}
}

type LoginCode_Code int32

const (
	LoginCode_Ok     LoginCode_Code = 0 // 校验成功
	LoginCode_Failed LoginCode_Code = 1 // 校验失败
)

// Enum value maps for LoginCode_Code.
var (
	LoginCode_Code_name = map[int32]string{
		0: "Ok",
		1: "Failed",
	}
	LoginCode_Code_value = map[string]int32{
		"Ok":     0,
		"Failed": 1,
	}
)

func (x LoginCode_Code) Enum() *LoginCode_Code {
	p := new(LoginCode_Code)
	*p = x
	return p
}

func (x LoginCode_Code) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LoginCode_Code) Descriptor() protoreflect.EnumDescriptor {
	return file_game_proto_enumTypes[1].Descriptor()
}

func (LoginCode_Code) Type() protoreflect.EnumType {
	return &file_game_proto_enumTypes[1]
}

func (x LoginCode_Code) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LoginCode_Code.Descriptor instead.
func (LoginCode_Code) EnumDescriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{0, 0}
}

type LianmengEnterCode_Code int32

const (
	LianmengEnterCode_Ok     LianmengEnterCode_Code = 0 // 成功
	LianmengEnterCode_Failed LianmengEnterCode_Code = 1 // 失败
)

// Enum value maps for LianmengEnterCode_Code.
var (
	LianmengEnterCode_Code_name = map[int32]string{
		0: "Ok",
		1: "Failed",
	}
	LianmengEnterCode_Code_value = map[string]int32{
		"Ok":     0,
		"Failed": 1,
	}
)

func (x LianmengEnterCode_Code) Enum() *LianmengEnterCode_Code {
	p := new(LianmengEnterCode_Code)
	*p = x
	return p
}

func (x LianmengEnterCode_Code) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LianmengEnterCode_Code) Descriptor() protoreflect.EnumDescriptor {
	return file_game_proto_enumTypes[2].Descriptor()
}

func (LianmengEnterCode_Code) Type() protoreflect.EnumType {
	return &file_game_proto_enumTypes[2]
}

func (x LianmengEnterCode_Code) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LianmengEnterCode_Code.Descriptor instead.
func (LianmengEnterCode_Code) EnumDescriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{1, 0}
}

type LianmengChatCode_Code int32

const (
	LianmengChatCode_Ok     LianmengChatCode_Code = 0 // 成功
	LianmengChatCode_Failed LianmengChatCode_Code = 1 // 失败
)

// Enum value maps for LianmengChatCode_Code.
var (
	LianmengChatCode_Code_name = map[int32]string{
		0: "Ok",
		1: "Failed",
	}
	LianmengChatCode_Code_value = map[string]int32{
		"Ok":     0,
		"Failed": 1,
	}
)

func (x LianmengChatCode_Code) Enum() *LianmengChatCode_Code {
	p := new(LianmengChatCode_Code)
	*p = x
	return p
}

func (x LianmengChatCode_Code) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LianmengChatCode_Code) Descriptor() protoreflect.EnumDescriptor {
	return file_game_proto_enumTypes[3].Descriptor()
}

func (LianmengChatCode_Code) Type() protoreflect.EnumType {
	return &file_game_proto_enumTypes[3]
}

func (x LianmengChatCode_Code) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LianmengChatCode_Code.Descriptor instead.
func (LianmengChatCode_Code) EnumDescriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{2, 0}
}

type LoginCode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LoginCode) Reset() {
	*x = LoginCode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginCode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginCode) ProtoMessage() {}

func (x *LoginCode) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginCode.ProtoReflect.Descriptor instead.
func (*LoginCode) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{0}
}

type LianmengEnterCode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LianmengEnterCode) Reset() {
	*x = LianmengEnterCode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LianmengEnterCode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LianmengEnterCode) ProtoMessage() {}

func (x *LianmengEnterCode) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LianmengEnterCode.ProtoReflect.Descriptor instead.
func (*LianmengEnterCode) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{1}
}

type LianmengChatCode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LianmengChatCode) Reset() {
	*x = LianmengChatCode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LianmengChatCode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LianmengChatCode) ProtoMessage() {}

func (x *LianmengChatCode) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LianmengChatCode.ProtoReflect.Descriptor instead.
func (*LianmengChatCode) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{2}
}

// 登录请求
type LoginReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=Token,proto3" json:"Token,omitempty"` // token
}

func (x *LoginReq) Reset() {
	*x = LoginReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginReq) ProtoMessage() {}

func (x *LoginReq) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginReq.ProtoReflect.Descriptor instead.
func (*LoginReq) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{3}
}

func (x *LoginReq) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

// 登录响应
type LoginRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code LoginCode_Code `protobuf:"varint,1,opt,name=Code,proto3,enum=pb.LoginCode_Code" json:"Code,omitempty"` // 返回码
}

func (x *LoginRes) Reset() {
	*x = LoginRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRes) ProtoMessage() {}

func (x *LoginRes) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRes.ProtoReflect.Descriptor instead.
func (*LoginRes) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{4}
}

func (x *LoginRes) GetCode() LoginCode_Code {
	if x != nil {
		return x.Code
	}
	return LoginCode_Ok
}

type LianmengEnterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code LianmengEnterCode_Code `protobuf:"varint,1,opt,name=Code,proto3,enum=pb.LianmengEnterCode_Code" json:"Code,omitempty"` //返回码
}

func (x *LianmengEnterResponse) Reset() {
	*x = LianmengEnterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LianmengEnterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LianmengEnterResponse) ProtoMessage() {}

func (x *LianmengEnterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LianmengEnterResponse.ProtoReflect.Descriptor instead.
func (*LianmengEnterResponse) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{5}
}

func (x *LianmengEnterResponse) GetCode() LianmengEnterCode_Code {
	if x != nil {
		return x.Code
	}
	return LianmengEnterCode_Ok
}

type LianmengChatMsgReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=Msg,proto3" json:"Msg,omitempty"`
}

func (x *LianmengChatMsgReq) Reset() {
	*x = LianmengChatMsgReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LianmengChatMsgReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LianmengChatMsgReq) ProtoMessage() {}

func (x *LianmengChatMsgReq) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LianmengChatMsgReq.ProtoReflect.Descriptor instead.
func (*LianmengChatMsgReq) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{6}
}

func (x *LianmengChatMsgReq) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type LianmengChatSendMsgRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code     LianmengChatCode_Code `protobuf:"varint,1,opt,name=Code,proto3,enum=pb.LianmengChatCode_Code" json:"Code,omitempty"`
	Msg      string                `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`           // 广播消息内容
	UserName string                `protobuf:"bytes,3,opt,name=UserName,proto3" json:"UserName,omitempty"` //用户名
}

func (x *LianmengChatSendMsgRes) Reset() {
	*x = LianmengChatSendMsgRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LianmengChatSendMsgRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LianmengChatSendMsgRes) ProtoMessage() {}

func (x *LianmengChatSendMsgRes) ProtoReflect() protoreflect.Message {
	mi := &file_game_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LianmengChatSendMsgRes.ProtoReflect.Descriptor instead.
func (*LianmengChatSendMsgRes) Descriptor() ([]byte, []int) {
	return file_game_proto_rawDescGZIP(), []int{7}
}

func (x *LianmengChatSendMsgRes) GetCode() LianmengChatCode_Code {
	if x != nil {
		return x.Code
	}
	return LianmengChatCode_Ok
}

func (x *LianmengChatSendMsgRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *LianmengChatSendMsgRes) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

var File_game_proto protoreflect.FileDescriptor

var file_game_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62,
	0x22, 0x27, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x1a, 0x0a,
	0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x6b, 0x10, 0x00, 0x12, 0x0a, 0x0a,
	0x06, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x10, 0x01, 0x22, 0x2f, 0x0a, 0x11, 0x4c, 0x69, 0x61,
	0x6e, 0x6d, 0x65, 0x6e, 0x67, 0x45, 0x6e, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x1a,
	0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x6b, 0x10, 0x00, 0x12, 0x0a,
	0x0a, 0x06, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x10, 0x01, 0x22, 0x2e, 0x0a, 0x10, 0x4c, 0x69,
	0x61, 0x6e, 0x6d, 0x65, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x1a,
	0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x6b, 0x10, 0x00, 0x12, 0x0a,
	0x0a, 0x06, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x10, 0x01, 0x22, 0x20, 0x0a, 0x08, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x32, 0x0a, 0x08,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x43, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65,
	0x22, 0x47, 0x0a, 0x15, 0x4c, 0x69, 0x61, 0x6e, 0x6d, 0x65, 0x6e, 0x67, 0x45, 0x6e, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x43, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x70, 0x62, 0x2e, 0x4c, 0x69, 0x61,
	0x6e, 0x6d, 0x65, 0x6e, 0x67, 0x45, 0x6e, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x2e, 0x43,
	0x6f, 0x64, 0x65, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x26, 0x0a, 0x12, 0x4c, 0x69, 0x61,
	0x6e, 0x6d, 0x65, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x12,
	0x10, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4d, 0x73,
	0x67, 0x22, 0x75, 0x0a, 0x16, 0x4c, 0x69, 0x61, 0x6e, 0x6d, 0x65, 0x6e, 0x67, 0x43, 0x68, 0x61,
	0x74, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x73, 0x12, 0x2d, 0x0a, 0x04, 0x43,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x70, 0x62, 0x2e, 0x4c,
	0x69, 0x61, 0x6e, 0x6d, 0x65, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x2e,
	0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x4d, 0x73,
	0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4d, 0x73, 0x67, 0x12, 0x1a, 0x0a, 0x08,
	0x55, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x55, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x2a, 0x3b, 0x0a, 0x05, 0x52, 0x6f, 0x75, 0x74,
	0x65, 0x12, 0x09, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11,
	0x4c, 0x69, 0x61, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x74, 0x45, 0x6e, 0x74, 0x65,
	0x72, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x4c, 0x69, 0x61, 0x6e, 0x6d, 0x65, 0x6e, 0x67, 0x43,
	0x68, 0x61, 0x74, 0x10, 0x02, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_game_proto_rawDescOnce sync.Once
	file_game_proto_rawDescData = file_game_proto_rawDesc
)

func file_game_proto_rawDescGZIP() []byte {
	file_game_proto_rawDescOnce.Do(func() {
		file_game_proto_rawDescData = protoimpl.X.CompressGZIP(file_game_proto_rawDescData)
	})
	return file_game_proto_rawDescData
}

var file_game_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_game_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_game_proto_goTypes = []interface{}{
	(Route)(0),                     // 0: pb.Route
	(LoginCode_Code)(0),            // 1: pb.LoginCode.Code
	(LianmengEnterCode_Code)(0),    // 2: pb.LianmengEnterCode.Code
	(LianmengChatCode_Code)(0),     // 3: pb.LianmengChatCode.Code
	(*LoginCode)(nil),              // 4: pb.LoginCode
	(*LianmengEnterCode)(nil),      // 5: pb.LianmengEnterCode
	(*LianmengChatCode)(nil),       // 6: pb.LianmengChatCode
	(*LoginReq)(nil),               // 7: pb.LoginReq
	(*LoginRes)(nil),               // 8: pb.LoginRes
	(*LianmengEnterResponse)(nil),  // 9: pb.LianmengEnterResponse
	(*LianmengChatMsgReq)(nil),     // 10: pb.LianmengChatMsgReq
	(*LianmengChatSendMsgRes)(nil), // 11: pb.LianmengChatSendMsgRes
}
var file_game_proto_depIdxs = []int32{
	1, // 0: pb.LoginRes.Code:type_name -> pb.LoginCode.Code
	2, // 1: pb.LianmengEnterResponse.Code:type_name -> pb.LianmengEnterCode.Code
	3, // 2: pb.LianmengChatSendMsgRes.Code:type_name -> pb.LianmengChatCode.Code
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_game_proto_init() }
func file_game_proto_init() {
	if File_game_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_game_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginCode); i {
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
		file_game_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LianmengEnterCode); i {
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
		file_game_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LianmengChatCode); i {
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
		file_game_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginReq); i {
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
		file_game_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginRes); i {
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
		file_game_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LianmengEnterResponse); i {
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
		file_game_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LianmengChatMsgReq); i {
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
		file_game_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LianmengChatSendMsgRes); i {
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
			RawDescriptor: file_game_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_game_proto_goTypes,
		DependencyIndexes: file_game_proto_depIdxs,
		EnumInfos:         file_game_proto_enumTypes,
		MessageInfos:      file_game_proto_msgTypes,
	}.Build()
	File_game_proto = out.File
	file_game_proto_rawDesc = nil
	file_game_proto_goTypes = nil
	file_game_proto_depIdxs = nil
}
