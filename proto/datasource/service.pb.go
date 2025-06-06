// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.29.3
// source: datasource/service.proto

package datasource

import (
	context "context"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AnalyzeAttackFlowRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId    uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	ResourceName string `protobuf:"bytes,2,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	CloudType    string `protobuf:"bytes,3,opt,name=cloud_type,json=cloudType,proto3" json:"cloud_type,omitempty"`
	CloudId      string `protobuf:"bytes,4,opt,name=cloud_id,json=cloudId,proto3" json:"cloud_id,omitempty"`
}

func (x *AnalyzeAttackFlowRequest) Reset() {
	*x = AnalyzeAttackFlowRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_datasource_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyzeAttackFlowRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzeAttackFlowRequest) ProtoMessage() {}

func (x *AnalyzeAttackFlowRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datasource_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyzeAttackFlowRequest.ProtoReflect.Descriptor instead.
func (*AnalyzeAttackFlowRequest) Descriptor() ([]byte, []int) {
	return file_datasource_service_proto_rawDescGZIP(), []int{0}
}

func (x *AnalyzeAttackFlowRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *AnalyzeAttackFlowRequest) GetResourceName() string {
	if x != nil {
		return x.ResourceName
	}
	return ""
}

func (x *AnalyzeAttackFlowRequest) GetCloudType() string {
	if x != nil {
		return x.CloudType
	}
	return ""
}

func (x *AnalyzeAttackFlowRequest) GetCloudId() string {
	if x != nil {
		return x.CloudId
	}
	return ""
}

type AnalyzeAttackFlowResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nodes []*Resource             `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	Edges []*ResourceRelationship `protobuf:"bytes,2,rep,name=edges,proto3" json:"edges,omitempty"`
}

func (x *AnalyzeAttackFlowResponse) Reset() {
	*x = AnalyzeAttackFlowResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_datasource_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyzeAttackFlowResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzeAttackFlowResponse) ProtoMessage() {}

func (x *AnalyzeAttackFlowResponse) ProtoReflect() protoreflect.Message {
	mi := &file_datasource_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyzeAttackFlowResponse.ProtoReflect.Descriptor instead.
func (*AnalyzeAttackFlowResponse) Descriptor() ([]byte, []int) {
	return file_datasource_service_proto_rawDescGZIP(), []int{1}
}

func (x *AnalyzeAttackFlowResponse) GetNodes() []*Resource {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *AnalyzeAttackFlowResponse) GetEdges() []*ResourceRelationship {
	if x != nil {
		return x.Edges
	}
	return nil
}

var File_datasource_service_proto protoreflect.FileDescriptor

var file_datasource_service_proto_rawDesc = []byte{
	0x0a, 0x18, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x64, 0x61, 0x74, 0x61,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc9, 0x01, 0x0a, 0x18, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a,
	0x65, 0x41, 0x74, 0x74, 0x61, 0x63, 0x6b, 0x46, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x26, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x2a, 0x02, 0x20, 0x00, 0x52,
	0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x2f, 0x0a, 0x0d, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x72, 0x05, 0x10, 0x01, 0x18, 0xff, 0x01, 0x52, 0x0c, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2e, 0x0a, 0x0a, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x0f, 0xfa, 0x42, 0x0c, 0x72, 0x0a, 0x52, 0x03, 0x61, 0x77, 0x73, 0x52, 0x03, 0x67, 0x63, 0x70,
	0x52, 0x09, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x24, 0x0a, 0x08, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xfa,
	0x42, 0x06, 0x72, 0x04, 0x10, 0x01, 0x18, 0x20, 0x52, 0x07, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x49,
	0x64, 0x22, 0x7f, 0x0a, 0x19, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x41, 0x74, 0x74, 0x61,
	0x63, 0x6b, 0x46, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a,
	0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x36, 0x0a, 0x05, 0x65, 0x64,
	0x67, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52,
	0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70, 0x52, 0x05, 0x65, 0x64, 0x67,
	0x65, 0x73, 0x32, 0xfb, 0x01, 0x0a, 0x11, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x41, 0x0a, 0x0f, 0x43, 0x6c, 0x65, 0x61,
	0x6e, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x60, 0x0a, 0x11, 0x41,
	0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x41, 0x74, 0x74, 0x61, 0x63, 0x6b, 0x46, 0x6c, 0x6f, 0x77,
	0x12, 0x24, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x41, 0x6e,
	0x61, 0x6c, 0x79, 0x7a, 0x65, 0x41, 0x74, 0x74, 0x61, 0x63, 0x6b, 0x46, 0x6c, 0x6f, 0x77, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x2e, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x41, 0x74, 0x74, 0x61, 0x63,
	0x6b, 0x46, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a,
	0x0f, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x53, 0x63, 0x61, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63,
	0x61, 0x2d, 0x72, 0x69, 0x73, 0x6b, 0x65, 0x6e, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x61,
	0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_datasource_service_proto_rawDescOnce sync.Once
	file_datasource_service_proto_rawDescData = file_datasource_service_proto_rawDesc
)

func file_datasource_service_proto_rawDescGZIP() []byte {
	file_datasource_service_proto_rawDescOnce.Do(func() {
		file_datasource_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_datasource_service_proto_rawDescData)
	})
	return file_datasource_service_proto_rawDescData
}

var file_datasource_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_datasource_service_proto_goTypes = []interface{}{
	(*AnalyzeAttackFlowRequest)(nil),  // 0: datasource.AnalyzeAttackFlowRequest
	(*AnalyzeAttackFlowResponse)(nil), // 1: datasource.AnalyzeAttackFlowResponse
	(*Resource)(nil),                  // 2: datasource.Resource
	(*ResourceRelationship)(nil),      // 3: datasource.ResourceRelationship
	(*emptypb.Empty)(nil),             // 4: google.protobuf.Empty
}
var file_datasource_service_proto_depIdxs = []int32{
	2, // 0: datasource.AnalyzeAttackFlowResponse.nodes:type_name -> datasource.Resource
	3, // 1: datasource.AnalyzeAttackFlowResponse.edges:type_name -> datasource.ResourceRelationship
	4, // 2: datasource.DataSourceService.CleanDataSource:input_type -> google.protobuf.Empty
	0, // 3: datasource.DataSourceService.AnalyzeAttackFlow:input_type -> datasource.AnalyzeAttackFlowRequest
	4, // 4: datasource.DataSourceService.NotifyScanError:input_type -> google.protobuf.Empty
	4, // 5: datasource.DataSourceService.CleanDataSource:output_type -> google.protobuf.Empty
	1, // 6: datasource.DataSourceService.AnalyzeAttackFlow:output_type -> datasource.AnalyzeAttackFlowResponse
	4, // 7: datasource.DataSourceService.NotifyScanError:output_type -> google.protobuf.Empty
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_datasource_service_proto_init() }
func file_datasource_service_proto_init() {
	if File_datasource_service_proto != nil {
		return
	}
	file_datasource_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_datasource_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyzeAttackFlowRequest); i {
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
		file_datasource_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyzeAttackFlowResponse); i {
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
			RawDescriptor: file_datasource_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_datasource_service_proto_goTypes,
		DependencyIndexes: file_datasource_service_proto_depIdxs,
		MessageInfos:      file_datasource_service_proto_msgTypes,
	}.Build()
	File_datasource_service_proto = out.File
	file_datasource_service_proto_rawDesc = nil
	file_datasource_service_proto_goTypes = nil
	file_datasource_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DataSourceServiceClient is the client API for DataSourceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DataSourceServiceClient interface {
	CleanDataSource(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	AnalyzeAttackFlow(ctx context.Context, in *AnalyzeAttackFlowRequest, opts ...grpc.CallOption) (*AnalyzeAttackFlowResponse, error)
	NotifyScanError(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type dataSourceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDataSourceServiceClient(cc grpc.ClientConnInterface) DataSourceServiceClient {
	return &dataSourceServiceClient{cc}
}

func (c *dataSourceServiceClient) CleanDataSource(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/datasource.DataSourceService/CleanDataSource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceServiceClient) AnalyzeAttackFlow(ctx context.Context, in *AnalyzeAttackFlowRequest, opts ...grpc.CallOption) (*AnalyzeAttackFlowResponse, error) {
	out := new(AnalyzeAttackFlowResponse)
	err := c.cc.Invoke(ctx, "/datasource.DataSourceService/AnalyzeAttackFlow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataSourceServiceClient) NotifyScanError(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/datasource.DataSourceService/NotifyScanError", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataSourceServiceServer is the server API for DataSourceService service.
type DataSourceServiceServer interface {
	CleanDataSource(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	AnalyzeAttackFlow(context.Context, *AnalyzeAttackFlowRequest) (*AnalyzeAttackFlowResponse, error)
	NotifyScanError(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
}

// UnimplementedDataSourceServiceServer can be embedded to have forward compatible implementations.
type UnimplementedDataSourceServiceServer struct {
}

func (*UnimplementedDataSourceServiceServer) CleanDataSource(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CleanDataSource not implemented")
}
func (*UnimplementedDataSourceServiceServer) AnalyzeAttackFlow(context.Context, *AnalyzeAttackFlowRequest) (*AnalyzeAttackFlowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AnalyzeAttackFlow not implemented")
}
func (*UnimplementedDataSourceServiceServer) NotifyScanError(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyScanError not implemented")
}

func RegisterDataSourceServiceServer(s *grpc.Server, srv DataSourceServiceServer) {
	s.RegisterService(&_DataSourceService_serviceDesc, srv)
}

func _DataSourceService_CleanDataSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataSourceServiceServer).CleanDataSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datasource.DataSourceService/CleanDataSource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataSourceServiceServer).CleanDataSource(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataSourceService_AnalyzeAttackFlow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AnalyzeAttackFlowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataSourceServiceServer).AnalyzeAttackFlow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datasource.DataSourceService/AnalyzeAttackFlow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataSourceServiceServer).AnalyzeAttackFlow(ctx, req.(*AnalyzeAttackFlowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataSourceService_NotifyScanError_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataSourceServiceServer).NotifyScanError(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/datasource.DataSourceService/NotifyScanError",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataSourceServiceServer).NotifyScanError(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _DataSourceService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "datasource.DataSourceService",
	HandlerType: (*DataSourceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CleanDataSource",
			Handler:    _DataSourceService_CleanDataSource_Handler,
		},
		{
			MethodName: "AnalyzeAttackFlow",
			Handler:    _DataSourceService_AnalyzeAttackFlow_Handler,
		},
		{
			MethodName: "NotifyScanError",
			Handler:    _DataSourceService_NotifyScanError_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "datasource/service.proto",
}
