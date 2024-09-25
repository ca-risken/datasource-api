package azure

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	azuremock "github.com/ca-risken/datasource-api/pkg/azure/mocks"
	"github.com/ca-risken/datasource-api/pkg/db"
	"github.com/ca-risken/datasource-api/pkg/db/mocks"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/pkg/test"
	"github.com/ca-risken/datasource-api/proto/azure"
	"gorm.io/gorm"
)

const (
	length65string = "12345678901234567890123456789012345678901234567890123456789012345"
)

func TestListAzureDataSource(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *azure.ListAzureDataSourceRequest
		want         *azure.ListAzureDataSourceResponse
		mockResponce *[]model.AzureDataSource
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &azure.ListAzureDataSourceRequest{AzureDataSourceId: 1},
			want: &azure.ListAzureDataSourceResponse{AzureDataSource: []*azure.AzureDataSource{
				{AzureDataSourceId: 1, Name: "one", Description: "desc", MaxScore: 1.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{AzureDataSourceId: 2, Name: "two", Description: "desc", MaxScore: 1.0, CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponce: &[]model.AzureDataSource{
				{AzureDataSourceID: 1, Name: "one", Description: "desc", MaxScore: 1.0, CreatedAt: now, UpdatedAt: now},
				{AzureDataSourceID: 2, Name: "two", Description: "desc", MaxScore: 1.0, CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK empty",
			input:     &azure.ListAzureDataSourceRequest{Name: "not exists name"},
			want:      &azure.ListAzureDataSourceResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &azure.ListAzureDataSourceRequest{Name: length65string},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &azure.ListAzureDataSourceRequest{AzureDataSourceId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewAzureRepoInterface(t)
			svc := AzureService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListAzureDataSource", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListAzureDataSource(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestListAzure(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *azure.ListAzureRequest
		want         *azure.ListAzureResponse
		mockResponce *[]model.Azure
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &azure.ListAzureRequest{ProjectId: 1, AzureId: 1, SubscriptionId: "pj"},
			want: &azure.ListAzureResponse{Azure: []*azure.Azure{
				{AzureId: 1, Name: "one", ProjectId: 1, SubscriptionId: "pj", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
				{AzureId: 2, Name: "two", ProjectId: 1, SubscriptionId: "pj", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			}},
			mockResponce: &[]model.Azure{
				{AzureID: 1, Name: "one", ProjectID: 1, SubscriptionID: "pj", CreatedAt: now, UpdatedAt: now},
				{AzureID: 2, Name: "two", ProjectID: 1, SubscriptionID: "pj", CreatedAt: now, UpdatedAt: now},
			},
		},
		{
			name:      "OK empty",
			input:     &azure.ListAzureRequest{ProjectId: 1, AzureId: 1, SubscriptionId: "pj"},
			want:      &azure.ListAzureResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &azure.ListAzureRequest{AzureId: 1, SubscriptionId: "pj"},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &azure.ListAzureRequest{ProjectId: 1, AzureId: 1, SubscriptionId: "pj"},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewAzureRepoInterface(t)
			svc := AzureService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListAzure", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListAzure(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetAzure(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *azure.GetAzureRequest
		want         *azure.GetAzureResponse
		mockResponce *model.Azure
		mockError    error
		wantErr      bool
	}{
		{
			name:         "OK",
			input:        &azure.GetAzureRequest{ProjectId: 1, AzureId: 1},
			want:         &azure.GetAzureResponse{Azure: &azure.Azure{AzureId: 1, Name: "one", ProjectId: 1, SubscriptionId: "pj", CreatedAt: now.Unix(), UpdatedAt: now.Unix()}},
			mockResponce: &model.Azure{AzureID: 1, Name: "one", ProjectID: 1, SubscriptionID: "pj", CreatedAt: now, UpdatedAt: now},
		},
		{
			name:      "OK empty",
			input:     &azure.GetAzureRequest{ProjectId: 1, AzureId: 1},
			want:      &azure.GetAzureResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &azure.GetAzureRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &azure.GetAzureRequest{ProjectId: 1, AzureId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewAzureRepoInterface(t)
			svc := AzureService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetAzure", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetAzure(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestPutAzure(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *azure.PutAzureRequest
		want         *azure.PutAzureResponse
		mockResponce *model.Azure
		mockError    error
		wantErr      bool
	}{
		{
			name: "OK",
			input: &azure.PutAzureRequest{ProjectId: 1, Azure: &azure.AzureForUpsert{
				AzureId: 1, Name: "one", ProjectId: 1, SubscriptionId: "pj", VerificationCode: "valid code"},
			},
			want: &azure.PutAzureResponse{Azure: &azure.Azure{
				AzureId: 1, Name: "one", ProjectId: 1, SubscriptionId: "pj", VerificationCode: "valid code", CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponce: &model.Azure{
				AzureID: 1, Name: "one", ProjectID: 1, SubscriptionID: "pj", VerificationCode: "valid code", CreatedAt: now, UpdatedAt: now,
			},
		},
		{
			name:    "NG invalid param",
			input:   &azure.PutAzureRequest{ProjectId: 1},
			wantErr: true,
		},
		{
			name: "Invalid DB error",
			input: &azure.PutAzureRequest{ProjectId: 1, Azure: &azure.AzureForUpsert{
				AzureId: 1, Name: "one", ProjectId: 1, SubscriptionId: "pj", VerificationCode: "valid code"},
			},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewAzureRepoInterface(t)
			svc := AzureService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("UpsertAzure", test.RepeatMockAnything(2)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.PutAzure(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDeleteAzure(t *testing.T) {
	cases := []struct {
		name      string
		input     *azure.DeleteAzureRequest
		wantErr   bool
		mockCall  bool
		mockError error
	}{
		{
			name:     "OK",
			input:    &azure.DeleteAzureRequest{ProjectId: 1, AzureId: 1},
			mockCall: true,
		},
		{
			name:     "NG invalid param",
			input:    &azure.DeleteAzureRequest{ProjectId: 1},
			wantErr:  true,
			mockCall: false,
		},
		{
			name:      "Invalid DB error",
			input:     &azure.DeleteAzureRequest{ProjectId: 1, AzureId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
			mockCall:  true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewAzureRepoInterface(t)
			svc := AzureService{repository: mockDB}

			if c.mockCall {
				mockDB.On("ListRelAzureDataSource", test.RepeatMockAnything(3)...).Return(&[]db.RelAzureDataSource{{AzureDataSourceID: 1}}, nil)
				mockDB.On("DeleteRelAzureDataSource", test.RepeatMockAnything(4)...).Return(nil)
				mockDB.On("DeleteAzure", test.RepeatMockAnything(3)...).Return(c.mockError).Once()
			}
			_, err := svc.DeleteAzure(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
		})
	}
}

func TestListRelAzureDataSource(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *azure.ListRelAzureDataSourceRequest
		want         *azure.ListRelAzureDataSourceResponse
		mockResponce *[]db.RelAzureDataSource
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &azure.ListRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1},
			want: &azure.ListRelAzureDataSourceResponse{RelAzureDataSource: []*azure.RelAzureDataSource{
				{AzureId: 1, AzureDataSourceId: 1, ProjectId: 1, Status: azure.Status_OK, StatusDetail: "", CreatedAt: now.Unix(), UpdatedAt: now.Unix(), Name: "name", MaxScore: 1.0, Description: "desc", SubscriptionId: "pj"},
				{AzureId: 2, AzureDataSourceId: 1, ProjectId: 1, Status: azure.Status_OK, StatusDetail: "", CreatedAt: now.Unix(), UpdatedAt: now.Unix(), Name: "name", MaxScore: 1.0, Description: "desc", SubscriptionId: "pj"},
			}},
			mockResponce: &[]db.RelAzureDataSource{
				{AzureID: 1, AzureDataSourceID: 1, ProjectID: 1, Status: azure.Status_OK.String(), StatusDetail: "", CreatedAt: now, UpdatedAt: now, Name: "name", MaxScore: 1.0, Description: "desc", SubscriptionID: "pj"},
				{AzureID: 2, AzureDataSourceID: 1, ProjectID: 1, Status: azure.Status_OK.String(), StatusDetail: "", CreatedAt: now, UpdatedAt: now, Name: "name", MaxScore: 1.0, Description: "desc", SubscriptionID: "pj"},
			},
		},
		{
			name:      "OK empty",
			input:     &azure.ListRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1},
			want:      &azure.ListRelAzureDataSourceResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &azure.ListRelAzureDataSourceRequest{AzureId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &azure.ListRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewAzureRepoInterface(t)
			svc := AzureService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("ListRelAzureDataSource", test.RepeatMockAnything(3)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.ListRelAzureDataSource(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestGetRelAzureDataSource(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *azure.GetRelAzureDataSourceRequest
		want         *azure.GetRelAzureDataSourceResponse
		mockResponce *db.RelAzureDataSource
		mockError    error
		wantErr      bool
	}{
		{
			name:  "OK",
			input: &azure.GetRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1, AzureDataSourceId: 1},
			want: &azure.GetRelAzureDataSourceResponse{RelAzureDataSource: &azure.RelAzureDataSource{
				AzureId: 1, AzureDataSourceId: 1, ProjectId: 1, Status: azure.Status_OK, StatusDetail: "", CreatedAt: now.Unix(), UpdatedAt: now.Unix(), Name: "name", MaxScore: 1.0, Description: "desc", SubscriptionId: "pj"},
			},
			mockResponce: &db.RelAzureDataSource{
				AzureID: 1, AzureDataSourceID: 1, ProjectID: 1, Status: azure.Status_OK.String(), StatusDetail: "", CreatedAt: now, UpdatedAt: now, Name: "name", MaxScore: 1.0, Description: "desc", SubscriptionID: "pj",
			},
		},
		{
			name:      "OK empty",
			input:     &azure.GetRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1, AzureDataSourceId: 1},
			want:      &azure.GetRelAzureDataSourceResponse{},
			mockError: gorm.ErrRecordNotFound,
		},
		{
			name:    "NG invalid param",
			input:   &azure.GetRelAzureDataSourceRequest{AzureId: 1},
			wantErr: true,
		},
		{
			name:      "Invalid DB error",
			input:     &azure.GetRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1, AzureDataSourceId: 1},
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewAzureRepoInterface(t)
			svc := AzureService{repository: mockDB}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("GetRelAzureDataSource", test.RepeatMockAnything(4)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.GetRelAzureDataSource(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestAttachRelAzureDataSource(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name         string
		input        *azure.AttachRelAzureDataSourceRequest
		want         *azure.AttachRelAzureDataSourceResponse
		mockResponce *db.RelAzureDataSource
		mockError    error
		wantErr      bool

		callGetAzure bool
	}{
		{
			name: "OK",
			input: &azure.AttachRelAzureDataSourceRequest{ProjectId: 1, RelAzureDataSource: &azure.RelAzureDataSourceForUpsert{
				AzureId: 1, AzureDataSourceId: 1, ProjectId: 1, Status: azure.Status_OK, StatusDetail: "", ScanAt: now.Unix()},
			},
			want: &azure.AttachRelAzureDataSourceResponse{RelAzureDataSource: &azure.RelAzureDataSource{
				AzureId: 1, AzureDataSourceId: 1, ProjectId: 1, Status: azure.Status_OK, StatusDetail: "", ScanAt: now.Unix(), CreatedAt: now.Unix(), UpdatedAt: now.Unix()},
			},
			mockResponce: &db.RelAzureDataSource{
				AzureID: 1, AzureDataSourceID: 1, ProjectID: 1, Status: azure.Status_OK.String(), ScanAt: now, CreatedAt: now, UpdatedAt: now,
			},
			callGetAzure: true,
		},
		{
			name:         "NG invalid param",
			input:        &azure.AttachRelAzureDataSourceRequest{ProjectId: 1},
			wantErr:      true,
			callGetAzure: false,
		},
		{
			name: "Invalid DB error",
			input: &azure.AttachRelAzureDataSourceRequest{ProjectId: 1, RelAzureDataSource: &azure.RelAzureDataSourceForUpsert{
				AzureId: 1, AzureDataSourceId: 1, ProjectId: 1, Status: azure.Status_OK, StatusDetail: "", ScanAt: now.Unix()},
			},
			mockError:    gorm.ErrInvalidDB,
			wantErr:      true,
			callGetAzure: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewAzureRepoInterface(t)
			mockAzure := azuremock.NewAzureServiceClient(t)
			logger := logging.NewLogger()

			svc := AzureService{repository: mockDB, azureClient: mockAzure, logger: logger}
			if c.callGetAzure {
				mockDB.On("GetAzure", test.RepeatMockAnything(3)...).Return(&model.Azure{}, nil)
				mockAzure.On("VerifyCode", test.RepeatMockAnything(3)...).Return(true, nil)
			}

			if c.mockResponce != nil || c.mockError != nil {
				mockDB.On("UpsertRelAzureDataSource", test.RepeatMockAnything(2)...).Return(c.mockResponce, c.mockError).Once()
			}
			got, err := svc.AttachRelAzureDataSource(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
			if c.wantErr && err == nil {
				t.Fatalf("Unexpected no error")
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected mapping: want=%+v, got=%+v", c.want, got)
			}
		})
	}
}

func TestDetachRelAzureDataSource(t *testing.T) {
	cases := []struct {
		name      string
		input     *azure.DetachRelAzureDataSourceRequest
		mockCall  bool
		mockError error
		wantErr   bool
	}{
		{
			name:     "OK",
			input:    &azure.DetachRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1, AzureDataSourceId: 1},
			mockCall: true,
		},
		{
			name:     "NG invalid param",
			input:    &azure.DetachRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1},
			mockCall: false,
			wantErr:  true,
		},
		{
			name:      "Invalid DB error",
			input:     &azure.DetachRelAzureDataSourceRequest{ProjectId: 1, AzureId: 1, AzureDataSourceId: 1},
			mockCall:  true,
			mockError: gorm.ErrInvalidDB,
			wantErr:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ctx context.Context
			mockDB := mocks.NewAzureRepoInterface(t)
			svc := AzureService{repository: mockDB}

			if c.mockCall {
				mockDB.On("DeleteRelAzureDataSource", test.RepeatMockAnything(4)...).Return(c.mockError).Once()
			}
			_, err := svc.DetachRelAzureDataSource(ctx, c.input)
			if !c.wantErr && err != nil {
				t.Fatalf("Unexpected error occured: %+v", err)
			}
		})
	}
}
