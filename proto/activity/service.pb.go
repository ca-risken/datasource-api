// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: activity/service.proto

package activity

import (
	context "context"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type DescribeARNRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Arn string `protobuf:"bytes,1,opt,name=arn,proto3" json:"arn,omitempty"`
}

func (x *DescribeARNRequest) Reset() {
	*x = DescribeARNRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_activity_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DescribeARNRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeARNRequest) ProtoMessage() {}

func (x *DescribeARNRequest) ProtoReflect() protoreflect.Message {
	mi := &file_activity_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescribeARNRequest.ProtoReflect.Descriptor instead.
func (*DescribeARNRequest) Descriptor() ([]byte, []int) {
	return file_activity_service_proto_rawDescGZIP(), []int{0}
}

func (x *DescribeARNRequest) GetArn() string {
	if x != nil {
		return x.Arn
	}
	return ""
}

type DescribeARNResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Arn *ARN `protobuf:"bytes,1,opt,name=arn,proto3" json:"arn,omitempty"`
}

func (x *DescribeARNResponse) Reset() {
	*x = DescribeARNResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_activity_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DescribeARNResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeARNResponse) ProtoMessage() {}

func (x *DescribeARNResponse) ProtoReflect() protoreflect.Message {
	mi := &file_activity_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescribeARNResponse.ProtoReflect.Descriptor instead.
func (*DescribeARNResponse) Descriptor() ([]byte, []int) {
	return file_activity_service_proto_rawDescGZIP(), []int{1}
}

func (x *DescribeARNResponse) GetArn() *ARN {
	if x != nil {
		return x.Arn
	}
	return nil
}

type ListCloudTrailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId      uint32       `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	AwsId          uint32       `protobuf:"varint,2,opt,name=aws_id,json=awsId,proto3" json:"aws_id,omitempty"`
	Region         string       `protobuf:"bytes,3,opt,name=region,proto3" json:"region,omitempty"`
	StartTime      int64        `protobuf:"varint,4,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime        int64        `protobuf:"varint,5,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	AttributeKey   AttributeKey `protobuf:"varint,6,opt,name=attribute_key,json=attributeKey,proto3,enum=datasource.activity.AttributeKey" json:"attribute_key,omitempty"`
	AttributeValue string       `protobuf:"bytes,7,opt,name=attribute_value,json=attributeValue,proto3" json:"attribute_value,omitempty"`
	NextToken      string       `protobuf:"bytes,8,opt,name=next_token,json=nextToken,proto3" json:"next_token,omitempty"`
}

func (x *ListCloudTrailRequest) Reset() {
	*x = ListCloudTrailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_activity_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCloudTrailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCloudTrailRequest) ProtoMessage() {}

func (x *ListCloudTrailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_activity_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCloudTrailRequest.ProtoReflect.Descriptor instead.
func (*ListCloudTrailRequest) Descriptor() ([]byte, []int) {
	return file_activity_service_proto_rawDescGZIP(), []int{2}
}

func (x *ListCloudTrailRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *ListCloudTrailRequest) GetAwsId() uint32 {
	if x != nil {
		return x.AwsId
	}
	return 0
}

func (x *ListCloudTrailRequest) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *ListCloudTrailRequest) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *ListCloudTrailRequest) GetEndTime() int64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

func (x *ListCloudTrailRequest) GetAttributeKey() AttributeKey {
	if x != nil {
		return x.AttributeKey
	}
	return AttributeKey_UNKNOWN
}

func (x *ListCloudTrailRequest) GetAttributeValue() string {
	if x != nil {
		return x.AttributeValue
	}
	return ""
}

func (x *ListCloudTrailRequest) GetNextToken() string {
	if x != nil {
		return x.NextToken
	}
	return ""
}

type ListCloudTrailResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cloudtrail []*CloudTrail `protobuf:"bytes,1,rep,name=cloudtrail,proto3" json:"cloudtrail,omitempty"`
	NextToken  string        `protobuf:"bytes,2,opt,name=next_token,json=nextToken,proto3" json:"next_token,omitempty"`
}

func (x *ListCloudTrailResponse) Reset() {
	*x = ListCloudTrailResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_activity_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCloudTrailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCloudTrailResponse) ProtoMessage() {}

func (x *ListCloudTrailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_activity_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCloudTrailResponse.ProtoReflect.Descriptor instead.
func (*ListCloudTrailResponse) Descriptor() ([]byte, []int) {
	return file_activity_service_proto_rawDescGZIP(), []int{3}
}

func (x *ListCloudTrailResponse) GetCloudtrail() []*CloudTrail {
	if x != nil {
		return x.Cloudtrail
	}
	return nil
}

func (x *ListCloudTrailResponse) GetNextToken() string {
	if x != nil {
		return x.NextToken
	}
	return ""
}

type ListConfigHistoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	AwsId     uint32 `protobuf:"varint,2,opt,name=aws_id,json=awsId,proto3" json:"aws_id,omitempty"`
	Region    string `protobuf:"bytes,3,opt,name=region,proto3" json:"region,omitempty"`
	// https://docs.aws.amazon.com/cli/latest/reference/configservice/get-resource-config-history.html#options
	ResourceType       string `protobuf:"bytes,4,opt,name=resource_type,json=resourceType,proto3" json:"resource_type,omitempty"`
	ResourceId         string `protobuf:"bytes,5,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	LaterTime          int64  `protobuf:"varint,6,opt,name=later_time,json=laterTime,proto3" json:"later_time,omitempty"`
	EarlierTime        int64  `protobuf:"varint,7,opt,name=earlier_time,json=earlierTime,proto3" json:"earlier_time,omitempty"`
	ChronologicalOrder string `protobuf:"bytes,8,opt,name=chronological_order,json=chronologicalOrder,proto3" json:"chronological_order,omitempty"` // default: Reverse
	StartingToken      string `protobuf:"bytes,9,opt,name=starting_token,json=startingToken,proto3" json:"starting_token,omitempty"`                // A token to specify where to start paginating.
}

func (x *ListConfigHistoryRequest) Reset() {
	*x = ListConfigHistoryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_activity_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListConfigHistoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListConfigHistoryRequest) ProtoMessage() {}

func (x *ListConfigHistoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_activity_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListConfigHistoryRequest.ProtoReflect.Descriptor instead.
func (*ListConfigHistoryRequest) Descriptor() ([]byte, []int) {
	return file_activity_service_proto_rawDescGZIP(), []int{4}
}

func (x *ListConfigHistoryRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *ListConfigHistoryRequest) GetAwsId() uint32 {
	if x != nil {
		return x.AwsId
	}
	return 0
}

func (x *ListConfigHistoryRequest) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *ListConfigHistoryRequest) GetResourceType() string {
	if x != nil {
		return x.ResourceType
	}
	return ""
}

func (x *ListConfigHistoryRequest) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *ListConfigHistoryRequest) GetLaterTime() int64 {
	if x != nil {
		return x.LaterTime
	}
	return 0
}

func (x *ListConfigHistoryRequest) GetEarlierTime() int64 {
	if x != nil {
		return x.EarlierTime
	}
	return 0
}

func (x *ListConfigHistoryRequest) GetChronologicalOrder() string {
	if x != nil {
		return x.ChronologicalOrder
	}
	return ""
}

func (x *ListConfigHistoryRequest) GetStartingToken() string {
	if x != nil {
		return x.StartingToken
	}
	return ""
}

type ListConfigHistoryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Configuration []*Configuration `protobuf:"bytes,1,rep,name=configuration,proto3" json:"configuration,omitempty"`
	NextToken     string           `protobuf:"bytes,2,opt,name=next_token,json=nextToken,proto3" json:"next_token,omitempty"`
}

func (x *ListConfigHistoryResponse) Reset() {
	*x = ListConfigHistoryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_activity_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListConfigHistoryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListConfigHistoryResponse) ProtoMessage() {}

func (x *ListConfigHistoryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_activity_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListConfigHistoryResponse.ProtoReflect.Descriptor instead.
func (*ListConfigHistoryResponse) Descriptor() ([]byte, []int) {
	return file_activity_service_proto_rawDescGZIP(), []int{5}
}

func (x *ListConfigHistoryResponse) GetConfiguration() []*Configuration {
	if x != nil {
		return x.Configuration
	}
	return nil
}

func (x *ListConfigHistoryResponse) GetNextToken() string {
	if x != nil {
		return x.NextToken
	}
	return ""
}

var File_activity_service_proto protoreflect.FileDescriptor

var file_activity_service_proto_rawDesc = []byte{
	0x0a, 0x16, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x1a, 0x15, 0x61,
	0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x37, 0x0a,
	0x12, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x41, 0x52, 0x4e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x03, 0x61, 0x72, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x0f, 0xfa, 0x42, 0x0c, 0x72, 0x0a, 0x32, 0x08, 0x5e, 0x61, 0x72, 0x6e, 0x3a, 0x2e, 0x2a,
	0x24, 0x52, 0x03, 0x61, 0x72, 0x6e, 0x22, 0x41, 0x0a, 0x13, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x41, 0x52, 0x4e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a,
	0x03, 0x61, 0x72, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x64, 0x61, 0x74,
	0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79,
	0x2e, 0x41, 0x52, 0x4e, 0x52, 0x03, 0x61, 0x72, 0x6e, 0x22, 0xea, 0x02, 0x0a, 0x15, 0x4c, 0x69,
	0x73, 0x74, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x54, 0x72, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x2a, 0x02, 0x20, 0x00,
	0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x06, 0x61,
	0x77, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x2a, 0x02, 0x20, 0x00, 0x52, 0x05, 0x61, 0x77, 0x73, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x06, 0x72,
	0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x72, 0x02, 0x10, 0x01, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x2d, 0x0a, 0x0a,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x42, 0x0e, 0xfa, 0x42, 0x0b, 0x22, 0x09, 0x18, 0xef, 0x85, 0xcf, 0xff, 0xaf, 0x07, 0x28, 0x00,
	0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x08, 0x65,
	0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x42, 0x0e, 0xfa,
	0x42, 0x0b, 0x22, 0x09, 0x18, 0xef, 0x85, 0xcf, 0xff, 0xaf, 0x07, 0x28, 0x00, 0x52, 0x07, 0x65,
	0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x46, 0x0a, 0x0d, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x21, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x76,
	0x69, 0x74, 0x79, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x4b, 0x65, 0x79,
	0x52, 0x0c, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x27,
	0x0a, 0x0f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x5f, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75,
	0x74, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x65, 0x78, 0x74, 0x5f,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x65, 0x78,
	0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x78, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c,
	0x6f, 0x75, 0x64, 0x54, 0x72, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x3f, 0x0a, 0x0a, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x74, 0x72, 0x61, 0x69, 0x6c, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x43, 0x6c, 0x6f, 0x75, 0x64,
	0x54, 0x72, 0x61, 0x69, 0x6c, 0x52, 0x0a, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x74, 0x72, 0x61, 0x69,
	0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x65, 0x78, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x22, 0xb0, 0x03, 0x0a, 0x18, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x48,
	0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a,
	0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x2a, 0x02, 0x20, 0x00, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x06, 0x61, 0x77, 0x73, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x2a, 0x02, 0x20, 0x00, 0x52, 0x05,
	0x61, 0x77, 0x73, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x06,
	0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x28, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02,
	0x10, 0x01, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x2d,
	0x0a, 0x0a, 0x6c, 0x61, 0x74, 0x65, 0x72, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x03, 0x42, 0x0e, 0xfa, 0x42, 0x0b, 0x22, 0x09, 0x18, 0xef, 0x85, 0xcf, 0xff, 0xaf, 0x07,
	0x28, 0x00, 0x52, 0x09, 0x6c, 0x61, 0x74, 0x65, 0x72, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x31, 0x0a,
	0x0c, 0x65, 0x61, 0x72, 0x6c, 0x69, 0x65, 0x72, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x03, 0x42, 0x0e, 0xfa, 0x42, 0x0b, 0x22, 0x09, 0x18, 0xef, 0x85, 0xcf, 0xff, 0xaf,
	0x07, 0x28, 0x00, 0x52, 0x0b, 0x65, 0x61, 0x72, 0x6c, 0x69, 0x65, 0x72, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x4a, 0x0a, 0x13, 0x63, 0x68, 0x72, 0x6f, 0x6e, 0x6f, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x61,
	0x6c, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x42, 0x19, 0xfa,
	0x42, 0x16, 0x72, 0x14, 0x52, 0x07, 0x52, 0x65, 0x76, 0x65, 0x72, 0x73, 0x65, 0x52, 0x07, 0x46,
	0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x52, 0x00, 0x52, 0x12, 0x63, 0x68, 0x72, 0x6f, 0x6e, 0x6f,
	0x6c, 0x6f, 0x67, 0x69, 0x63, 0x61, 0x6c, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x0e,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x22, 0x84, 0x01, 0x0a, 0x19, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x48, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x6e,
	0x65, 0x78, 0x74, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x6e, 0x65, 0x78, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x32, 0xd2, 0x02, 0x0a, 0x0f, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x60,
	0x0a, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x41, 0x52, 0x4e, 0x12, 0x27, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x76,
	0x69, 0x74, 0x79, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x41, 0x52, 0x4e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x44, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x41, 0x52, 0x4e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x69, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x54, 0x72, 0x61,
	0x69, 0x6c, 0x12, 0x2a, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e,
	0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c, 0x6f,
	0x75, 0x64, 0x54, 0x72, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b,
	0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x61, 0x63, 0x74, 0x69,
	0x76, 0x69, 0x74, 0x79, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x54, 0x72,
	0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x72, 0x0a, 0x11, 0x4c,
	0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79,
	0x12, 0x2d, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x61, 0x63,
	0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x69, 0x74, 0x79, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x61,
	0x2d, 0x72, 0x69, 0x73, 0x6b, 0x65, 0x6e, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x69, 0x74, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_activity_service_proto_rawDescOnce sync.Once
	file_activity_service_proto_rawDescData = file_activity_service_proto_rawDesc
)

func file_activity_service_proto_rawDescGZIP() []byte {
	file_activity_service_proto_rawDescOnce.Do(func() {
		file_activity_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_activity_service_proto_rawDescData)
	})
	return file_activity_service_proto_rawDescData
}

var file_activity_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_activity_service_proto_goTypes = []interface{}{
	(*DescribeARNRequest)(nil),        // 0: datasource.activity.DescribeARNRequest
	(*DescribeARNResponse)(nil),       // 1: datasource.activity.DescribeARNResponse
	(*ListCloudTrailRequest)(nil),     // 2: datasource.activity.ListCloudTrailRequest
	(*ListCloudTrailResponse)(nil),    // 3: datasource.activity.ListCloudTrailResponse
	(*ListConfigHistoryRequest)(nil),  // 4: datasource.activity.ListConfigHistoryRequest
	(*ListConfigHistoryResponse)(nil), // 5: datasource.activity.ListConfigHistoryResponse
	(*ARN)(nil),                       // 6: datasource.activity.ARN
	(AttributeKey)(0),                 // 7: datasource.activity.AttributeKey
	(*CloudTrail)(nil),                // 8: datasource.activity.CloudTrail
	(*Configuration)(nil),             // 9: datasource.activity.Configuration
}
var file_activity_service_proto_depIdxs = []int32{
	6, // 0: datasource.activity.DescribeARNResponse.arn:type_name -> datasource.activity.ARN
	7, // 1: datasource.activity.ListCloudTrailRequest.attribute_key:type_name -> datasource.activity.AttributeKey
	8, // 2: datasource.activity.ListCloudTrailResponse.cloudtrail:type_name -> datasource.activity.CloudTrail
	9, // 3: datasource.activity.ListConfigHistoryResponse.configuration:type_name -> datasource.activity.Configuration
	0, // 4: datasource.activity.ActivityService.DescribeARN:input_type -> datasource.activity.DescribeARNRequest
	2, // 5: datasource.activity.ActivityService.ListCloudTrail:input_type -> datasource.activity.ListCloudTrailRequest
	4, // 6: datasource.activity.ActivityService.ListConfigHistory:input_type -> datasource.activity.ListConfigHistoryRequest
	1, // 7: datasource.activity.ActivityService.DescribeARN:output_type -> datasource.activity.DescribeARNResponse
	3, // 8: datasource.activity.ActivityService.ListCloudTrail:output_type -> datasource.activity.ListCloudTrailResponse
	5, // 9: datasource.activity.ActivityService.ListConfigHistory:output_type -> datasource.activity.ListConfigHistoryResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_activity_service_proto_init() }
func file_activity_service_proto_init() {
	if File_activity_service_proto != nil {
		return
	}
	file_activity_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_activity_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DescribeARNRequest); i {
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
		file_activity_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DescribeARNResponse); i {
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
		file_activity_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListCloudTrailRequest); i {
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
		file_activity_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListCloudTrailResponse); i {
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
		file_activity_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListConfigHistoryRequest); i {
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
		file_activity_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListConfigHistoryResponse); i {
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
			RawDescriptor: file_activity_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_activity_service_proto_goTypes,
		DependencyIndexes: file_activity_service_proto_depIdxs,
		MessageInfos:      file_activity_service_proto_msgTypes,
	}.Build()
	File_activity_service_proto = out.File
	file_activity_service_proto_rawDesc = nil
	file_activity_service_proto_goTypes = nil
	file_activity_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ActivityServiceClient is the client API for ActivityService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ActivityServiceClient interface {
	DescribeARN(ctx context.Context, in *DescribeARNRequest, opts ...grpc.CallOption) (*DescribeARNResponse, error)
	ListCloudTrail(ctx context.Context, in *ListCloudTrailRequest, opts ...grpc.CallOption) (*ListCloudTrailResponse, error)
	ListConfigHistory(ctx context.Context, in *ListConfigHistoryRequest, opts ...grpc.CallOption) (*ListConfigHistoryResponse, error)
}

type activityServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewActivityServiceClient(cc grpc.ClientConnInterface) ActivityServiceClient {
	return &activityServiceClient{cc}
}

func (c *activityServiceClient) DescribeARN(ctx context.Context, in *DescribeARNRequest, opts ...grpc.CallOption) (*DescribeARNResponse, error) {
	out := new(DescribeARNResponse)
	err := c.cc.Invoke(ctx, "/datasource.activity.ActivityService/DescribeARN", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *activityServiceClient) ListCloudTrail(ctx context.Context, in *ListCloudTrailRequest, opts ...grpc.CallOption) (*ListCloudTrailResponse, error) {
	out := new(ListCloudTrailResponse)
	err := c.cc.Invoke(ctx, "/datasource.activity.ActivityService/ListCloudTrail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *activityServiceClient) ListConfigHistory(ctx context.Context, in *ListConfigHistoryRequest, opts ...grpc.CallOption) (*ListConfigHistoryResponse, error) {
	out := new(ListConfigHistoryResponse)
	err := c.cc.Invoke(ctx, "/datasource.activity.ActivityService/ListConfigHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityServiceServer is the server API for ActivityService service.
type ActivityServiceServer interface {
	DescribeARN(context.Context, *DescribeARNRequest) (*DescribeARNResponse, error)
	ListCloudTrail(context.Context, *ListCloudTrailRequest) (*ListCloudTrailResponse, error)
	ListConfigHistory(context.Context, *ListConfigHistoryRequest) (*ListConfigHistoryResponse, error)
}

// UnimplementedActivityServiceServer can be embedded to have forward compatible implementations.
type UnimplementedActivityServiceServer struct {
}

func (*UnimplementedActivityServiceServer) DescribeARN(context.Context, *DescribeARNRequest) (*DescribeARNResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeARN not implemented")
}
func (*UnimplementedActivityServiceServer) ListCloudTrail(context.Context, *ListCloudTrailRequest) (*ListCloudTrailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCloudTrail not implemented")
}
func (*UnimplementedActivityServiceServer) ListConfigHistory(context.Context, *ListConfigHistoryRequest) (*ListConfigHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListConfigHistory not implemented")
}

func RegisterActivityServiceServer(s *grpc.Server, srv ActivityServiceServer) {
	s.RegisterService(&_ActivityService_serviceDesc, srv)
}

func _ActivityService_DescribeARN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeARNRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActivityServiceServer).DescribeARN(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datasource.activity.ActivityService/DescribeARN",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActivityServiceServer).DescribeARN(ctx, req.(*DescribeARNRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ActivityService_ListCloudTrail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCloudTrailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActivityServiceServer).ListCloudTrail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datasource.activity.ActivityService/ListCloudTrail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActivityServiceServer).ListCloudTrail(ctx, req.(*ListCloudTrailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ActivityService_ListConfigHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListConfigHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActivityServiceServer).ListConfigHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datasource.activity.ActivityService/ListConfigHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActivityServiceServer).ListConfigHistory(ctx, req.(*ListConfigHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ActivityService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "datasource.activity.ActivityService",
	HandlerType: (*ActivityServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DescribeARN",
			Handler:    _ActivityService_DescribeARN_Handler,
		},
		{
			MethodName: "ListCloudTrail",
			Handler:    _ActivityService_ListCloudTrail_Handler,
		},
		{
			MethodName: "ListConfigHistory",
			Handler:    _ActivityService_ListConfigHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "activity/service.proto",
}
