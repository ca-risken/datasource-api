// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.0
// source: aws/entity.proto

package aws

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
	return file_aws_entity_proto_enumTypes[0].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_aws_entity_proto_enumTypes[0]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_aws_entity_proto_rawDescGZIP(), []int{0}
}

// AWS
type AWS struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AwsId        uint32 `protobuf:"varint,1,opt,name=aws_id,json=awsId,proto3" json:"aws_id,omitempty"`
	Name         string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	ProjectId    uint32 `protobuf:"varint,3,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	AwsAccountId string `protobuf:"bytes,4,opt,name=aws_account_id,json=awsAccountId,proto3" json:"aws_account_id,omitempty"` // 12桁のAWSアカウントID
	CreatedAt    int64  `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt    int64  `protobuf:"varint,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *AWS) Reset() {
	*x = AWS{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aws_entity_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AWS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AWS) ProtoMessage() {}

func (x *AWS) ProtoReflect() protoreflect.Message {
	mi := &file_aws_entity_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AWS.ProtoReflect.Descriptor instead.
func (*AWS) Descriptor() ([]byte, []int) {
	return file_aws_entity_proto_rawDescGZIP(), []int{0}
}

func (x *AWS) GetAwsId() uint32 {
	if x != nil {
		return x.AwsId
	}
	return 0
}

func (x *AWS) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AWS) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *AWS) GetAwsAccountId() string {
	if x != nil {
		return x.AwsAccountId
	}
	return ""
}

func (x *AWS) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *AWS) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

// AWSRelDataSource
type AWSRelDataSource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AwsId           uint32 `protobuf:"varint,1,opt,name=aws_id,json=awsId,proto3" json:"aws_id,omitempty"`
	AwsDataSourceId uint32 `protobuf:"varint,2,opt,name=aws_data_source_id,json=awsDataSourceId,proto3" json:"aws_data_source_id,omitempty"`
	ProjectId       uint32 `protobuf:"varint,3,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	AssumeRoleArn   string `protobuf:"bytes,4,opt,name=assume_role_arn,json=assumeRoleArn,proto3" json:"assume_role_arn,omitempty"`
	ExternalId      string `protobuf:"bytes,5,opt,name=external_id,json=externalId,proto3" json:"external_id,omitempty"`
	Status          Status `protobuf:"varint,6,opt,name=status,proto3,enum=datasource.aws.Status" json:"status,omitempty"`
	StatusDetail    string `protobuf:"bytes,7,opt,name=status_detail,json=statusDetail,proto3" json:"status_detail,omitempty"`
	ScanAt          int64  `protobuf:"varint,8,opt,name=scan_at,json=scanAt,proto3" json:"scan_at,omitempty"`
	CreatedAt       int64  `protobuf:"varint,9,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt       int64  `protobuf:"varint,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	SpecificVersion string `protobuf:"bytes,11,opt,name=specific_version,json=specificVersion,proto3" json:"specific_version,omitempty"`
}

func (x *AWSRelDataSource) Reset() {
	*x = AWSRelDataSource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aws_entity_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AWSRelDataSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AWSRelDataSource) ProtoMessage() {}

func (x *AWSRelDataSource) ProtoReflect() protoreflect.Message {
	mi := &file_aws_entity_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AWSRelDataSource.ProtoReflect.Descriptor instead.
func (*AWSRelDataSource) Descriptor() ([]byte, []int) {
	return file_aws_entity_proto_rawDescGZIP(), []int{1}
}

func (x *AWSRelDataSource) GetAwsId() uint32 {
	if x != nil {
		return x.AwsId
	}
	return 0
}

func (x *AWSRelDataSource) GetAwsDataSourceId() uint32 {
	if x != nil {
		return x.AwsDataSourceId
	}
	return 0
}

func (x *AWSRelDataSource) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *AWSRelDataSource) GetAssumeRoleArn() string {
	if x != nil {
		return x.AssumeRoleArn
	}
	return ""
}

func (x *AWSRelDataSource) GetExternalId() string {
	if x != nil {
		return x.ExternalId
	}
	return ""
}

func (x *AWSRelDataSource) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_UNKNOWN
}

func (x *AWSRelDataSource) GetStatusDetail() string {
	if x != nil {
		return x.StatusDetail
	}
	return ""
}

func (x *AWSRelDataSource) GetScanAt() int64 {
	if x != nil {
		return x.ScanAt
	}
	return 0
}

func (x *AWSRelDataSource) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *AWSRelDataSource) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *AWSRelDataSource) GetSpecificVersion() string {
	if x != nil {
		return x.SpecificVersion
	}
	return ""
}

// AWSForUpsert
// (Unique keys: aws_account_id)
type AWSForUpsert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ProjectId    uint32 `protobuf:"varint,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	AwsAccountId string `protobuf:"bytes,3,opt,name=aws_account_id,json=awsAccountId,proto3" json:"aws_account_id,omitempty"` // 12桁のAWSアカウントID
}

func (x *AWSForUpsert) Reset() {
	*x = AWSForUpsert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aws_entity_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AWSForUpsert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AWSForUpsert) ProtoMessage() {}

func (x *AWSForUpsert) ProtoReflect() protoreflect.Message {
	mi := &file_aws_entity_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AWSForUpsert.ProtoReflect.Descriptor instead.
func (*AWSForUpsert) Descriptor() ([]byte, []int) {
	return file_aws_entity_proto_rawDescGZIP(), []int{2}
}

func (x *AWSForUpsert) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AWSForUpsert) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *AWSForUpsert) GetAwsAccountId() string {
	if x != nil {
		return x.AwsAccountId
	}
	return ""
}

// DataSource(data_sourceと紐づくaws_rel_data_sourceの状態)
type DataSource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AwsDataSourceId uint32  `protobuf:"varint,1,opt,name=aws_data_source_id,json=awsDataSourceId,proto3" json:"aws_data_source_id,omitempty"` // aws_data_source.aws_data_source_id
	DataSource      string  `protobuf:"bytes,2,opt,name=data_source,json=dataSource,proto3" json:"data_source,omitempty"`                     // aws_data_source.data_source
	MaxScore        float32 `protobuf:"fixed32,3,opt,name=max_score,json=maxScore,proto3" json:"max_score,omitempty"`                         // aws_data_source.max_score
	AwsId           uint32  `protobuf:"varint,4,opt,name=aws_id,json=awsId,proto3" json:"aws_id,omitempty"`                                   // aws_rel_data_source.aws_id
	ProjectId       uint32  `protobuf:"varint,5,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`                       // aws_rel_data_source.project_id
	AssumeRoleArn   string  `protobuf:"bytes,7,opt,name=assume_role_arn,json=assumeRoleArn,proto3" json:"assume_role_arn,omitempty"`          // aws_rel_data_source.assume_role_arn
	ExternalId      string  `protobuf:"bytes,8,opt,name=external_id,json=externalId,proto3" json:"external_id,omitempty"`                     // aws_rel_data_source.external_id
	Status          Status  `protobuf:"varint,9,opt,name=status,proto3,enum=datasource.aws.Status" json:"status,omitempty"`                   // aws_rel_data_source.status
	StatusDetail    string  `protobuf:"bytes,10,opt,name=status_detail,json=statusDetail,proto3" json:"status_detail,omitempty"`              // aws_rel_data_source.status_detail
	ScanAt          int64   `protobuf:"varint,11,opt,name=scan_at,json=scanAt,proto3" json:"scan_at,omitempty"`                               // aws_rel_data_source.scan_at
	SpecificVersion string  `protobuf:"bytes,12,opt,name=specific_version,json=specificVersion,proto3" json:"specific_version,omitempty"`     // aws_rel_data_source.specific_version
}

func (x *DataSource) Reset() {
	*x = DataSource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aws_entity_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataSource) ProtoMessage() {}

func (x *DataSource) ProtoReflect() protoreflect.Message {
	mi := &file_aws_entity_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataSource.ProtoReflect.Descriptor instead.
func (*DataSource) Descriptor() ([]byte, []int) {
	return file_aws_entity_proto_rawDescGZIP(), []int{3}
}

func (x *DataSource) GetAwsDataSourceId() uint32 {
	if x != nil {
		return x.AwsDataSourceId
	}
	return 0
}

func (x *DataSource) GetDataSource() string {
	if x != nil {
		return x.DataSource
	}
	return ""
}

func (x *DataSource) GetMaxScore() float32 {
	if x != nil {
		return x.MaxScore
	}
	return 0
}

func (x *DataSource) GetAwsId() uint32 {
	if x != nil {
		return x.AwsId
	}
	return 0
}

func (x *DataSource) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *DataSource) GetAssumeRoleArn() string {
	if x != nil {
		return x.AssumeRoleArn
	}
	return ""
}

func (x *DataSource) GetExternalId() string {
	if x != nil {
		return x.ExternalId
	}
	return ""
}

func (x *DataSource) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_UNKNOWN
}

func (x *DataSource) GetStatusDetail() string {
	if x != nil {
		return x.StatusDetail
	}
	return ""
}

func (x *DataSource) GetScanAt() int64 {
	if x != nil {
		return x.ScanAt
	}
	return 0
}

func (x *DataSource) GetSpecificVersion() string {
	if x != nil {
		return x.SpecificVersion
	}
	return ""
}

// DataSourceForAttach
type DataSourceForAttach struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AwsId           uint32 `protobuf:"varint,1,opt,name=aws_id,json=awsId,proto3" json:"aws_id,omitempty"`
	AwsDataSourceId uint32 `protobuf:"varint,2,opt,name=aws_data_source_id,json=awsDataSourceId,proto3" json:"aws_data_source_id,omitempty"`
	ProjectId       uint32 `protobuf:"varint,3,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	AssumeRoleArn   string `protobuf:"bytes,4,opt,name=assume_role_arn,json=assumeRoleArn,proto3" json:"assume_role_arn,omitempty"`
	ExternalId      string `protobuf:"bytes,5,opt,name=external_id,json=externalId,proto3" json:"external_id,omitempty"` // assume_role時に指定する外部ID
	Status          Status `protobuf:"varint,6,opt,name=status,proto3,enum=datasource.aws.Status" json:"status,omitempty"`
	StatusDetail    string `protobuf:"bytes,7,opt,name=status_detail,json=statusDetail,proto3" json:"status_detail,omitempty"`
	ScanAt          int64  `protobuf:"varint,8,opt,name=scan_at,json=scanAt,proto3" json:"scan_at,omitempty"`
	SpecificVersion string `protobuf:"bytes,9,opt,name=specific_version,json=specificVersion,proto3" json:"specific_version,omitempty"`
}

func (x *DataSourceForAttach) Reset() {
	*x = DataSourceForAttach{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aws_entity_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataSourceForAttach) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataSourceForAttach) ProtoMessage() {}

func (x *DataSourceForAttach) ProtoReflect() protoreflect.Message {
	mi := &file_aws_entity_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataSourceForAttach.ProtoReflect.Descriptor instead.
func (*DataSourceForAttach) Descriptor() ([]byte, []int) {
	return file_aws_entity_proto_rawDescGZIP(), []int{4}
}

func (x *DataSourceForAttach) GetAwsId() uint32 {
	if x != nil {
		return x.AwsId
	}
	return 0
}

func (x *DataSourceForAttach) GetAwsDataSourceId() uint32 {
	if x != nil {
		return x.AwsDataSourceId
	}
	return 0
}

func (x *DataSourceForAttach) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *DataSourceForAttach) GetAssumeRoleArn() string {
	if x != nil {
		return x.AssumeRoleArn
	}
	return ""
}

func (x *DataSourceForAttach) GetExternalId() string {
	if x != nil {
		return x.ExternalId
	}
	return ""
}

func (x *DataSourceForAttach) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_UNKNOWN
}

func (x *DataSourceForAttach) GetStatusDetail() string {
	if x != nil {
		return x.StatusDetail
	}
	return ""
}

func (x *DataSourceForAttach) GetScanAt() int64 {
	if x != nil {
		return x.ScanAt
	}
	return 0
}

func (x *DataSourceForAttach) GetSpecificVersion() string {
	if x != nil {
		return x.SpecificVersion
	}
	return ""
}

var File_aws_entity_proto protoreflect.FileDescriptor

var file_aws_entity_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x77, 0x73, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x61,
	0x77, 0x73, 0x22, 0xb3, 0x01, 0x0a, 0x03, 0x41, 0x57, 0x53, 0x12, 0x15, 0x0a, 0x06, 0x61, 0x77,
	0x73, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x61, 0x77, 0x73, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x61, 0x77, 0x73, 0x5f, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61, 0x77,
	0x73, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x95, 0x03, 0x0a, 0x10, 0x41, 0x57, 0x53,
	0x52, 0x65, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x15, 0x0a,
	0x06, 0x61, 0x77, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x61,
	0x77, 0x73, 0x49, 0x64, 0x12, 0x2b, 0x0a, 0x12, 0x61, 0x77, 0x73, 0x5f, 0x64, 0x61, 0x74, 0x61,
	0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0f, 0x61, 0x77, 0x73, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64,
	0x12, 0x26, 0x0a, 0x0f, 0x61, 0x73, 0x73, 0x75, 0x6d, 0x65, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x5f,
	0x61, 0x72, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x73, 0x73, 0x75, 0x6d,
	0x65, 0x52, 0x6f, 0x6c, 0x65, 0x41, 0x72, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x78, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65,
	0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x2e, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x61, 0x77, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x5f, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x17,
	0x0a, 0x07, 0x73, 0x63, 0x61, 0x6e, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x73, 0x63, 0x61, 0x6e, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x29, 0x0a, 0x10, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69,
	0x63, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0f, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x22, 0x67, 0x0a, 0x0c, 0x41, 0x57, 0x53, 0x46, 0x6f, 0x72, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x61, 0x77, 0x73, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61, 0x77, 0x73,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x8f, 0x03, 0x0a, 0x0a, 0x44, 0x61,
	0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x2b, 0x0a, 0x12, 0x61, 0x77, 0x73, 0x5f,
	0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x0f, 0x61, 0x77, 0x73, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x61, 0x74, 0x61,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x63,
	0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x08, 0x6d, 0x61, 0x78, 0x53, 0x63,
	0x6f, 0x72, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x61, 0x77, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x05, 0x61, 0x77, 0x73, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x61, 0x73, 0x73,
	0x75, 0x6d, 0x65, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x61, 0x72, 0x6e, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x61, 0x73, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x41, 0x72,
	0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x69, 0x64,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x49, 0x64, 0x12, 0x2e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x16, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e,
	0x61, 0x77, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x64, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x63, 0x61, 0x6e, 0x5f,
	0x61, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x63, 0x61, 0x6e, 0x41, 0x74,
	0x12, 0x29, 0x0a, 0x10, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x5f, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x73, 0x70, 0x65, 0x63,
	0x69, 0x66, 0x69, 0x63, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xda, 0x02, 0x0a, 0x13,
	0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x46, 0x6f, 0x72, 0x41, 0x74, 0x74,
	0x61, 0x63, 0x68, 0x12, 0x15, 0x0a, 0x06, 0x61, 0x77, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x05, 0x61, 0x77, 0x73, 0x49, 0x64, 0x12, 0x2b, 0x0a, 0x12, 0x61, 0x77,
	0x73, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0f, 0x61, 0x77, 0x73, 0x44, 0x61, 0x74, 0x61, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x61, 0x73, 0x73, 0x75, 0x6d, 0x65,
	0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x61, 0x72, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x61, 0x73, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x41, 0x72, 0x6e, 0x12, 0x1f,
	0x0a, 0x0b, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x49, 0x64, 0x12,
	0x2e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x16, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x61, 0x77, 0x73,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x23, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x44, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x63, 0x61, 0x6e, 0x5f, 0x61, 0x74, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x63, 0x61, 0x6e, 0x41, 0x74, 0x12, 0x29, 0x0a,
	0x10, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69,
	0x63, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2a, 0x49, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12,
	0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x4f, 0x4e, 0x46, 0x49,
	0x47, 0x55, 0x52, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x4e, 0x5f, 0x50, 0x52,
	0x4f, 0x47, 0x52, 0x45, 0x53, 0x53, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x10, 0x04, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x63, 0x61, 0x2d, 0x72, 0x69, 0x73, 0x6b, 0x65, 0x6e, 0x2f, 0x64, 0x61, 0x74, 0x61,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x61, 0x77, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aws_entity_proto_rawDescOnce sync.Once
	file_aws_entity_proto_rawDescData = file_aws_entity_proto_rawDesc
)

func file_aws_entity_proto_rawDescGZIP() []byte {
	file_aws_entity_proto_rawDescOnce.Do(func() {
		file_aws_entity_proto_rawDescData = protoimpl.X.CompressGZIP(file_aws_entity_proto_rawDescData)
	})
	return file_aws_entity_proto_rawDescData
}

var file_aws_entity_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_aws_entity_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_aws_entity_proto_goTypes = []interface{}{
	(Status)(0),                 // 0: datasource.aws.Status
	(*AWS)(nil),                 // 1: datasource.aws.AWS
	(*AWSRelDataSource)(nil),    // 2: datasource.aws.AWSRelDataSource
	(*AWSForUpsert)(nil),        // 3: datasource.aws.AWSForUpsert
	(*DataSource)(nil),          // 4: datasource.aws.DataSource
	(*DataSourceForAttach)(nil), // 5: datasource.aws.DataSourceForAttach
}
var file_aws_entity_proto_depIdxs = []int32{
	0, // 0: datasource.aws.AWSRelDataSource.status:type_name -> datasource.aws.Status
	0, // 1: datasource.aws.DataSource.status:type_name -> datasource.aws.Status
	0, // 2: datasource.aws.DataSourceForAttach.status:type_name -> datasource.aws.Status
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_aws_entity_proto_init() }
func file_aws_entity_proto_init() {
	if File_aws_entity_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aws_entity_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AWS); i {
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
		file_aws_entity_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AWSRelDataSource); i {
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
		file_aws_entity_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AWSForUpsert); i {
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
		file_aws_entity_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataSource); i {
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
		file_aws_entity_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataSourceForAttach); i {
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
			RawDescriptor: file_aws_entity_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_aws_entity_proto_goTypes,
		DependencyIndexes: file_aws_entity_proto_depIdxs,
		EnumInfos:         file_aws_entity_proto_enumTypes,
		MessageInfos:      file_aws_entity_proto_msgTypes,
	}.Build()
	File_aws_entity_proto = out.File
	file_aws_entity_proto_rawDesc = nil
	file_aws_entity_proto_goTypes = nil
	file_aws_entity_proto_depIdxs = nil
}
