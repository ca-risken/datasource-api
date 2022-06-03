// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: osint/entities.proto

package osint

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

// Status
type Status int32

const (
	Status_UNKNOWN     Status = 0
	Status_OK          Status = 1
	Status_CONFIGURED  Status = 2
	Status_IN_PROGRESS Status = 3
	Status_ERROR       Status = 4
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0: "UNKNOWN",
		1: "OK",
		2: "CONFIGURED",
		3: "IN_PROGRESS",
		4: "ERROR",
	}
	Status_value = map[string]int32{
		"UNKNOWN":     0,
		"OK":          1,
		"CONFIGURED":  2,
		"IN_PROGRESS": 3,
		"ERROR":       4,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_osint_entities_proto_enumTypes[0].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_osint_entities_proto_enumTypes[0]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_osint_entities_proto_rawDescGZIP(), []int{0}
}

type Osint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OsintId      uint32 `protobuf:"varint,1,opt,name=osint_id,json=osintId,proto3" json:"osint_id,omitempty"`
	ProjectId    uint32 `protobuf:"varint,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	ResourceType string `protobuf:"bytes,3,opt,name=resource_type,json=resourceType,proto3" json:"resource_type,omitempty"`
	ResourceName string `protobuf:"bytes,4,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	CreatedAt    int64  `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt    int64  `protobuf:"varint,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Osint) Reset() {
	*x = Osint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_osint_entities_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Osint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Osint) ProtoMessage() {}

func (x *Osint) ProtoReflect() protoreflect.Message {
	mi := &file_osint_entities_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Osint.ProtoReflect.Descriptor instead.
func (*Osint) Descriptor() ([]byte, []int) {
	return file_osint_entities_proto_rawDescGZIP(), []int{0}
}

func (x *Osint) GetOsintId() uint32 {
	if x != nil {
		return x.OsintId
	}
	return 0
}

func (x *Osint) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *Osint) GetResourceType() string {
	if x != nil {
		return x.ResourceType
	}
	return ""
}

func (x *Osint) GetResourceName() string {
	if x != nil {
		return x.ResourceName
	}
	return ""
}

func (x *Osint) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Osint) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

type OsintForUpsert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OsintId      uint32 `protobuf:"varint,1,opt,name=osint_id,json=osintId,proto3" json:"osint_id,omitempty"`
	ProjectId    uint32 `protobuf:"varint,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	ResourceType string `protobuf:"bytes,3,opt,name=resource_type,json=resourceType,proto3" json:"resource_type,omitempty"`
	ResourceName string `protobuf:"bytes,4,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
}

func (x *OsintForUpsert) Reset() {
	*x = OsintForUpsert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_osint_entities_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OsintForUpsert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OsintForUpsert) ProtoMessage() {}

func (x *OsintForUpsert) ProtoReflect() protoreflect.Message {
	mi := &file_osint_entities_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OsintForUpsert.ProtoReflect.Descriptor instead.
func (*OsintForUpsert) Descriptor() ([]byte, []int) {
	return file_osint_entities_proto_rawDescGZIP(), []int{1}
}

func (x *OsintForUpsert) GetOsintId() uint32 {
	if x != nil {
		return x.OsintId
	}
	return 0
}

func (x *OsintForUpsert) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *OsintForUpsert) GetResourceType() string {
	if x != nil {
		return x.ResourceType
	}
	return ""
}

func (x *OsintForUpsert) GetResourceName() string {
	if x != nil {
		return x.ResourceName
	}
	return ""
}

type OsintDataSource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OsintDataSourceId uint32  `protobuf:"varint,1,opt,name=osint_data_source_id,json=osintDataSourceId,proto3" json:"osint_data_source_id,omitempty"`
	Name              string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description       string  `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	MaxScore          float32 `protobuf:"fixed32,4,opt,name=max_score,json=maxScore,proto3" json:"max_score,omitempty"`
	CreatedAt         int64   `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt         int64   `protobuf:"varint,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *OsintDataSource) Reset() {
	*x = OsintDataSource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_osint_entities_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OsintDataSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OsintDataSource) ProtoMessage() {}

func (x *OsintDataSource) ProtoReflect() protoreflect.Message {
	mi := &file_osint_entities_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OsintDataSource.ProtoReflect.Descriptor instead.
func (*OsintDataSource) Descriptor() ([]byte, []int) {
	return file_osint_entities_proto_rawDescGZIP(), []int{2}
}

func (x *OsintDataSource) GetOsintDataSourceId() uint32 {
	if x != nil {
		return x.OsintDataSourceId
	}
	return 0
}

func (x *OsintDataSource) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *OsintDataSource) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *OsintDataSource) GetMaxScore() float32 {
	if x != nil {
		return x.MaxScore
	}
	return 0
}

func (x *OsintDataSource) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *OsintDataSource) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

type OsintDataSourceForUpsert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OsintDataSourceId uint32  `protobuf:"varint,1,opt,name=osint_data_source_id,json=osintDataSourceId,proto3" json:"osint_data_source_id,omitempty"`
	Name              string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description       string  `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	MaxScore          float32 `protobuf:"fixed32,4,opt,name=max_score,json=maxScore,proto3" json:"max_score,omitempty"`
}

func (x *OsintDataSourceForUpsert) Reset() {
	*x = OsintDataSourceForUpsert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_osint_entities_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OsintDataSourceForUpsert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OsintDataSourceForUpsert) ProtoMessage() {}

func (x *OsintDataSourceForUpsert) ProtoReflect() protoreflect.Message {
	mi := &file_osint_entities_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OsintDataSourceForUpsert.ProtoReflect.Descriptor instead.
func (*OsintDataSourceForUpsert) Descriptor() ([]byte, []int) {
	return file_osint_entities_proto_rawDescGZIP(), []int{3}
}

func (x *OsintDataSourceForUpsert) GetOsintDataSourceId() uint32 {
	if x != nil {
		return x.OsintDataSourceId
	}
	return 0
}

func (x *OsintDataSourceForUpsert) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *OsintDataSourceForUpsert) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *OsintDataSourceForUpsert) GetMaxScore() float32 {
	if x != nil {
		return x.MaxScore
	}
	return 0
}

type RelOsintDataSource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RelOsintDataSourceId uint32 `protobuf:"varint,1,opt,name=rel_osint_data_source_id,json=relOsintDataSourceId,proto3" json:"rel_osint_data_source_id,omitempty"`
	OsintDataSourceId    uint32 `protobuf:"varint,2,opt,name=osint_data_source_id,json=osintDataSourceId,proto3" json:"osint_data_source_id,omitempty"`
	OsintId              uint32 `protobuf:"varint,3,opt,name=osint_id,json=osintId,proto3" json:"osint_id,omitempty"`
	ProjectId            uint32 `protobuf:"varint,4,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Status               Status `protobuf:"varint,5,opt,name=status,proto3,enum=osint.osint.Status" json:"status,omitempty"`
	StatusDetail         string `protobuf:"bytes,6,opt,name=status_detail,json=statusDetail,proto3" json:"status_detail,omitempty"`
	ScanAt               int64  `protobuf:"varint,7,opt,name=scan_at,json=scanAt,proto3" json:"scan_at,omitempty"`
	CreatedAt            int64  `protobuf:"varint,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt            int64  `protobuf:"varint,9,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *RelOsintDataSource) Reset() {
	*x = RelOsintDataSource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_osint_entities_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RelOsintDataSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelOsintDataSource) ProtoMessage() {}

func (x *RelOsintDataSource) ProtoReflect() protoreflect.Message {
	mi := &file_osint_entities_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelOsintDataSource.ProtoReflect.Descriptor instead.
func (*RelOsintDataSource) Descriptor() ([]byte, []int) {
	return file_osint_entities_proto_rawDescGZIP(), []int{4}
}

func (x *RelOsintDataSource) GetRelOsintDataSourceId() uint32 {
	if x != nil {
		return x.RelOsintDataSourceId
	}
	return 0
}

func (x *RelOsintDataSource) GetOsintDataSourceId() uint32 {
	if x != nil {
		return x.OsintDataSourceId
	}
	return 0
}

func (x *RelOsintDataSource) GetOsintId() uint32 {
	if x != nil {
		return x.OsintId
	}
	return 0
}

func (x *RelOsintDataSource) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *RelOsintDataSource) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_UNKNOWN
}

func (x *RelOsintDataSource) GetStatusDetail() string {
	if x != nil {
		return x.StatusDetail
	}
	return ""
}

func (x *RelOsintDataSource) GetScanAt() int64 {
	if x != nil {
		return x.ScanAt
	}
	return 0
}

func (x *RelOsintDataSource) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *RelOsintDataSource) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

type RelOsintDataSourceForUpsert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RelOsintDataSourceId uint32 `protobuf:"varint,1,opt,name=rel_osint_data_source_id,json=relOsintDataSourceId,proto3" json:"rel_osint_data_source_id,omitempty"`
	OsintDataSourceId    uint32 `protobuf:"varint,2,opt,name=osint_data_source_id,json=osintDataSourceId,proto3" json:"osint_data_source_id,omitempty"`
	OsintId              uint32 `protobuf:"varint,3,opt,name=osint_id,json=osintId,proto3" json:"osint_id,omitempty"`
	ProjectId            uint32 `protobuf:"varint,4,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Status               Status `protobuf:"varint,5,opt,name=status,proto3,enum=osint.osint.Status" json:"status,omitempty"`
	StatusDetail         string `protobuf:"bytes,6,opt,name=status_detail,json=statusDetail,proto3" json:"status_detail,omitempty"`
	ScanAt               int64  `protobuf:"varint,7,opt,name=scan_at,json=scanAt,proto3" json:"scan_at,omitempty"`
}

func (x *RelOsintDataSourceForUpsert) Reset() {
	*x = RelOsintDataSourceForUpsert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_osint_entities_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RelOsintDataSourceForUpsert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelOsintDataSourceForUpsert) ProtoMessage() {}

func (x *RelOsintDataSourceForUpsert) ProtoReflect() protoreflect.Message {
	mi := &file_osint_entities_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelOsintDataSourceForUpsert.ProtoReflect.Descriptor instead.
func (*RelOsintDataSourceForUpsert) Descriptor() ([]byte, []int) {
	return file_osint_entities_proto_rawDescGZIP(), []int{5}
}

func (x *RelOsintDataSourceForUpsert) GetRelOsintDataSourceId() uint32 {
	if x != nil {
		return x.RelOsintDataSourceId
	}
	return 0
}

func (x *RelOsintDataSourceForUpsert) GetOsintDataSourceId() uint32 {
	if x != nil {
		return x.OsintDataSourceId
	}
	return 0
}

func (x *RelOsintDataSourceForUpsert) GetOsintId() uint32 {
	if x != nil {
		return x.OsintId
	}
	return 0
}

func (x *RelOsintDataSourceForUpsert) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *RelOsintDataSourceForUpsert) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_UNKNOWN
}

func (x *RelOsintDataSourceForUpsert) GetStatusDetail() string {
	if x != nil {
		return x.StatusDetail
	}
	return ""
}

func (x *RelOsintDataSourceForUpsert) GetScanAt() int64 {
	if x != nil {
		return x.ScanAt
	}
	return 0
}

type OsintDetectWord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OsintDetectWordId    uint32 `protobuf:"varint,1,opt,name=osint_detect_word_id,json=osintDetectWordId,proto3" json:"osint_detect_word_id,omitempty"`
	RelOsintDataSourceId uint32 `protobuf:"varint,2,opt,name=rel_osint_data_source_id,json=relOsintDataSourceId,proto3" json:"rel_osint_data_source_id,omitempty"`
	Word                 string `protobuf:"bytes,3,opt,name=word,proto3" json:"word,omitempty"`
	ProjectId            uint32 `protobuf:"varint,4,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	CreatedAt            int64  `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt            int64  `protobuf:"varint,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *OsintDetectWord) Reset() {
	*x = OsintDetectWord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_osint_entities_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OsintDetectWord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OsintDetectWord) ProtoMessage() {}

func (x *OsintDetectWord) ProtoReflect() protoreflect.Message {
	mi := &file_osint_entities_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OsintDetectWord.ProtoReflect.Descriptor instead.
func (*OsintDetectWord) Descriptor() ([]byte, []int) {
	return file_osint_entities_proto_rawDescGZIP(), []int{6}
}

func (x *OsintDetectWord) GetOsintDetectWordId() uint32 {
	if x != nil {
		return x.OsintDetectWordId
	}
	return 0
}

func (x *OsintDetectWord) GetRelOsintDataSourceId() uint32 {
	if x != nil {
		return x.RelOsintDataSourceId
	}
	return 0
}

func (x *OsintDetectWord) GetWord() string {
	if x != nil {
		return x.Word
	}
	return ""
}

func (x *OsintDetectWord) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *OsintDetectWord) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *OsintDetectWord) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

type OsintDetectWordForUpsert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OsintDetectWordId    uint32 `protobuf:"varint,1,opt,name=osint_detect_word_id,json=osintDetectWordId,proto3" json:"osint_detect_word_id,omitempty"`
	RelOsintDataSourceId uint32 `protobuf:"varint,2,opt,name=rel_osint_data_source_id,json=relOsintDataSourceId,proto3" json:"rel_osint_data_source_id,omitempty"`
	Word                 string `protobuf:"bytes,3,opt,name=word,proto3" json:"word,omitempty"`
	ProjectId            uint32 `protobuf:"varint,4,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
}

func (x *OsintDetectWordForUpsert) Reset() {
	*x = OsintDetectWordForUpsert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_osint_entities_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OsintDetectWordForUpsert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OsintDetectWordForUpsert) ProtoMessage() {}

func (x *OsintDetectWordForUpsert) ProtoReflect() protoreflect.Message {
	mi := &file_osint_entities_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OsintDetectWordForUpsert.ProtoReflect.Descriptor instead.
func (*OsintDetectWordForUpsert) Descriptor() ([]byte, []int) {
	return file_osint_entities_proto_rawDescGZIP(), []int{7}
}

func (x *OsintDetectWordForUpsert) GetOsintDetectWordId() uint32 {
	if x != nil {
		return x.OsintDetectWordId
	}
	return 0
}

func (x *OsintDetectWordForUpsert) GetRelOsintDataSourceId() uint32 {
	if x != nil {
		return x.RelOsintDataSourceId
	}
	return 0
}

func (x *OsintDetectWordForUpsert) GetWord() string {
	if x != nil {
		return x.Word
	}
	return ""
}

func (x *OsintDetectWordForUpsert) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

var File_osint_entities_proto protoreflect.FileDescriptor

var file_osint_entities_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x2e, 0x6f, 0x73,
	0x69, 0x6e, 0x74, 0x22, 0xc9, 0x01, 0x0a, 0x05, 0x4f, 0x73, 0x69, 0x6e, 0x74, 0x12, 0x19, 0x0a,
	0x08, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x07, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x23, 0x0a, 0x0d,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22,
	0x94, 0x01, 0x0a, 0x0e, 0x4f, 0x73, 0x69, 0x6e, 0x74, 0x46, 0x6f, 0x72, 0x55, 0x70, 0x73, 0x65,
	0x72, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1d, 0x0a,
	0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xd3, 0x01, 0x0a, 0x0f, 0x4f, 0x73, 0x69, 0x6e, 0x74,
	0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x2f, 0x0a, 0x14, 0x6f, 0x73,
	0x69, 0x6e, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x11, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x44,
	0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x08, 0x6d, 0x61, 0x78, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x1d,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x9e, 0x01, 0x0a,
	0x18, 0x4f, 0x73, 0x69, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x46, 0x6f, 0x72, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x12, 0x2f, 0x0a, 0x14, 0x6f, 0x73, 0x69,
	0x6e, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x11, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x08, 0x6d, 0x61, 0x78, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x22, 0xe0, 0x02,
	0x0a, 0x12, 0x52, 0x65, 0x6c, 0x4f, 0x73, 0x69, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x12, 0x36, 0x0a, 0x18, 0x72, 0x65, 0x6c, 0x5f, 0x6f, 0x73, 0x69, 0x6e,
	0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x14, 0x72, 0x65, 0x6c, 0x4f, 0x73, 0x69, 0x6e, 0x74,
	0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x2f, 0x0a, 0x14,
	0x6f, 0x73, 0x69, 0x6e, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x11, 0x6f, 0x73, 0x69, 0x6e,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x19, 0x0a,
	0x08, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x07, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x2b, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x2e,
	0x6f, 0x73, 0x69, 0x6e, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x64,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x63, 0x61,
	0x6e, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x63, 0x61, 0x6e,
	0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x22, 0xab, 0x02, 0x0a, 0x1b, 0x52, 0x65, 0x6c, 0x4f, 0x73, 0x69, 0x6e, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x46, 0x6f, 0x72, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74,
	0x12, 0x36, 0x0a, 0x18, 0x72, 0x65, 0x6c, 0x5f, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x5f, 0x64, 0x61,
	0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x14, 0x72, 0x65, 0x6c, 0x4f, 0x73, 0x69, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x2f, 0x0a, 0x14, 0x6f, 0x73, 0x69, 0x6e,
	0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x11, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x73, 0x69,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6f, 0x73, 0x69,
	0x6e, 0x74, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x49, 0x64, 0x12, 0x2b, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x2e, 0x6f, 0x73, 0x69, 0x6e,
	0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x23, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x64, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x44,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x63, 0x61, 0x6e, 0x5f, 0x61, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x63, 0x61, 0x6e, 0x41, 0x74, 0x22, 0xeb,
	0x01, 0x0a, 0x0f, 0x4f, 0x73, 0x69, 0x6e, 0x74, 0x44, 0x65, 0x74, 0x65, 0x63, 0x74, 0x57, 0x6f,
	0x72, 0x64, 0x12, 0x2f, 0x0a, 0x14, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x5f, 0x64, 0x65, 0x74, 0x65,
	0x63, 0x74, 0x5f, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x11, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x44, 0x65, 0x74, 0x65, 0x63, 0x74, 0x57, 0x6f, 0x72,
	0x64, 0x49, 0x64, 0x12, 0x36, 0x0a, 0x18, 0x72, 0x65, 0x6c, 0x5f, 0x6f, 0x73, 0x69, 0x6e, 0x74,
	0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x14, 0x72, 0x65, 0x6c, 0x4f, 0x73, 0x69, 0x6e, 0x74, 0x44,
	0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x77,
	0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x77, 0x6f, 0x72, 0x64, 0x12,
	0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1d,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0xb6, 0x01, 0x0a,
	0x18, 0x4f, 0x73, 0x69, 0x6e, 0x74, 0x44, 0x65, 0x74, 0x65, 0x63, 0x74, 0x57, 0x6f, 0x72, 0x64,
	0x46, 0x6f, 0x72, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x12, 0x2f, 0x0a, 0x14, 0x6f, 0x73, 0x69,
	0x6e, 0x74, 0x5f, 0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x5f, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x11, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x44, 0x65,
	0x74, 0x65, 0x63, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x49, 0x64, 0x12, 0x36, 0x0a, 0x18, 0x72, 0x65,
	0x6c, 0x5f, 0x6f, 0x73, 0x69, 0x6e, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x14, 0x72, 0x65,
	0x6c, 0x4f, 0x73, 0x69, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x49, 0x64, 0x2a, 0x49, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x06, 0x0a, 0x02,
	0x4f, 0x4b, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x55, 0x52,
	0x45, 0x44, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x47, 0x52,
	0x45, 0x53, 0x53, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x04,
	0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63,
	0x61, 0x2d, 0x72, 0x69, 0x73, 0x6b, 0x65, 0x6e, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x73,
	0x69, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_osint_entities_proto_rawDescOnce sync.Once
	file_osint_entities_proto_rawDescData = file_osint_entities_proto_rawDesc
)

func file_osint_entities_proto_rawDescGZIP() []byte {
	file_osint_entities_proto_rawDescOnce.Do(func() {
		file_osint_entities_proto_rawDescData = protoimpl.X.CompressGZIP(file_osint_entities_proto_rawDescData)
	})
	return file_osint_entities_proto_rawDescData
}

var file_osint_entities_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_osint_entities_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_osint_entities_proto_goTypes = []interface{}{
	(Status)(0),                         // 0: osint.osint.Status
	(*Osint)(nil),                       // 1: osint.osint.Osint
	(*OsintForUpsert)(nil),              // 2: osint.osint.OsintForUpsert
	(*OsintDataSource)(nil),             // 3: osint.osint.OsintDataSource
	(*OsintDataSourceForUpsert)(nil),    // 4: osint.osint.OsintDataSourceForUpsert
	(*RelOsintDataSource)(nil),          // 5: osint.osint.RelOsintDataSource
	(*RelOsintDataSourceForUpsert)(nil), // 6: osint.osint.RelOsintDataSourceForUpsert
	(*OsintDetectWord)(nil),             // 7: osint.osint.OsintDetectWord
	(*OsintDetectWordForUpsert)(nil),    // 8: osint.osint.OsintDetectWordForUpsert
}
var file_osint_entities_proto_depIdxs = []int32{
	0, // 0: osint.osint.RelOsintDataSource.status:type_name -> osint.osint.Status
	0, // 1: osint.osint.RelOsintDataSourceForUpsert.status:type_name -> osint.osint.Status
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_osint_entities_proto_init() }
func file_osint_entities_proto_init() {
	if File_osint_entities_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_osint_entities_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Osint); i {
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
		file_osint_entities_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OsintForUpsert); i {
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
		file_osint_entities_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OsintDataSource); i {
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
		file_osint_entities_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OsintDataSourceForUpsert); i {
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
		file_osint_entities_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RelOsintDataSource); i {
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
		file_osint_entities_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RelOsintDataSourceForUpsert); i {
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
		file_osint_entities_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OsintDetectWord); i {
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
		file_osint_entities_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OsintDetectWordForUpsert); i {
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
			RawDescriptor: file_osint_entities_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_osint_entities_proto_goTypes,
		DependencyIndexes: file_osint_entities_proto_depIdxs,
		EnumInfos:         file_osint_entities_proto_enumTypes,
		MessageInfos:      file_osint_entities_proto_msgTypes,
	}.Build()
	File_osint_entities_proto = out.File
	file_osint_entities_proto_rawDesc = nil
	file_osint_entities_proto_goTypes = nil
	file_osint_entities_proto_depIdxs = nil
}
