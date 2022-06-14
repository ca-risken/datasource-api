TARGETS = aws google diagnosis osint code
MOCK_TARGETS = $(TARGETS:=.mock)
BUILD_OPT=""
IMAGE_TAG=latest
MANIFEST_TAG=latest
IMAGE_NAME=datasource-api
IMAGE_REGISTRY=local

.PHONY: all
all: run

.PHONY: install
install:
	go get \
		google.golang.org/grpc \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/envoyproxy/protoc-gen-validate@v0.6.7 \
		github.com/grpc-ecosystem/go-grpc-middleware

.PHONY: clean
clean:
	rm -f proto/*/*.pb.go

.PHONY: fmt
fmt: proto/**/*.proto
	@clang-format -i proto/**/*.proto

.PHONY: doc
doc: fmt
	protoc \
		--proto_path=proto \
		--error_format=gcc \
		-I $(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
		--doc_out=markdown,README.md:doc \
		proto/**/*.proto;

# build without protoc-gen-validate
.PHONY: proto-without-validation
proto-without-validate: fmt
	for svc in "aws" "google" "code" "diagnosis" "osint"; do \
		protoc \
			--proto_path=proto \
			--error_format=gcc \
			--go_out=plugins=grpc,paths=source_relative:proto \
			proto/$$svc/*.proto; \
	done

.PHONY: proto
# proto : proto-without-validate proto-mock 
proto : proto-without-validate # TODO add proto-mock

PHONY: build
build: test
	IMAGE_TAG=$(IMAGE_TAG) IMAGE_NAME=$(IMAGE_NAME) BUILD_OPT="$(BUILD_OPT)" . hack/docker-build.sh

PHONY: build-ci
build-ci:
	IMAGE_TAG=$(IMAGE_TAG) IMAGE_NAME=$(IMAGE_NAME) BUILD_OPT="$(BUILD_OPT)" . hack/docker-build.sh
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

PHONY: push-image
push-image:
	docker push $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

PHONY: pull-image
pull-image:
	docker pull $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

PHONY: tag-image
tag-image:
	docker tag $(SOURCE_IMAGE_NAME):$(SOURCE_IMAGE_TAG) $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

PHONY: create-manifest
create-manifest:
	docker manifest create $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG_BASE)_linux_amd64 $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG_BASE)_linux_arm64
	docker manifest annotate --arch amd64 $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG_BASE)_linux_amd64
	docker manifest annotate --arch arm64 $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG_BASE)_linux_arm64

PHONY: push-manifest
push-manifest:
	docker manifest push $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(MANIFEST_TAG)
	docker manifest inspect $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(MANIFEST_TAG)

PHONY: test
test:
	GO111MODULE=on go test ./...

.PHONY: lint
lint: FAKE
	GO111MODULE=on golangci-lint run --timeout 5m

.PHONY: generate-mock
generate-mock: proto-mock
proto-mock: $(MOCK_TARGETS)
%.mock: FAKE
	sh hack/generate-mock.sh proto/$(*)

FAKE:

####################################################
## grpcurl example
####################################################
.PHONY: help
help:
	@echo "Usage: make <sub-command>"
	@echo "\n---------------- sub-command list ----------------"
	@cat Makefile | grep -e "^.PHONY:" | grep -v "all" | cut -f2 -d' '

####################################################
## AWS
####################################################
.PHONY: list-aws-service
list-aws-service:
	grpcurl -plaintext localhost:8081 list datasource.aws.AWSService

.PHONY: list-aws
list-aws:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "aws_id":1}' \
		localhost:8081 datasource.aws.AWSService.ListAWS

.PHONY: put-aws
put-aws:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "aws":{"name":"account-01", "project_id":1, "aws_account_id":"123456789001"}}' \
		localhost:8081 datasource.aws.AWSService.PutAWS

.PHONY: delete-aws
delete-aws:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "aws_id":2}' \
		localhost:8081 datasource.aws.AWSService.DeleteAWS

.PHONY: list-aws-data-source
list-aws-data-source:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "aws_id":1}' \
		localhost:8081 datasource.aws.AWSService.ListDataSource

.PHONY: attach-aws-data-source
attach-aws-data-source:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "attach_data_source": {"aws_id":1, "aws_data_source_id":1004, "project_id":1, "assume_role_arn":"arn:aws:iam::123456789012:role/role-name", "external_id":"test", "status":"CONFIGURED"}}' \
		localhost:8081 datasource.aws.AWSService.AttachDataSource

.PHONY: detach-aws-data-source
detach-aws-data-source:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "aws_id":1, "aws_data_source_id":1004}' \
		localhost:8081 datasource.aws.AWSService.DetachDataSource

.PHONY: invoke-aws-scan
invoke-aws-scan:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "aws_id":1, "aws_data_source_id":1001}' \
		localhost:8081 datasource.aws.AWSService.InvokeScan

.PHONY: invoke-aws-scan-all
invoke-aws-scan-all:
	grpcurl \
		-plaintext \
		localhost:8081 datasource.aws.AWSService.InvokeScanAll

####################################################
## Code
####################################################
.PHONY: list-code-service
list-code-service:
	grpcurl -plaintext localhost:8081 list datasource.code.CodeService

.PHONY: list-code-datasource
list-code-datasource:
	grpcurl \
		-plaintext \
		-d '{"code_data_source_id":1001, "name":"code:gitleaks"}' \
		localhost:8081 datasource.code.CodeService.ListDataSource

.PHONY: list-gitleaks
list-gitleaks:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "code_data_source_id":1001, "gitleaks_id":1}' \
		localhost:8081 datasource.code.CodeService.ListGitleaks

.PHONY: put-gitleaks
put-gitleaks:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "gitleaks": {"gitleaks_id":1, "code_data_source_id":1001, "name":"test-gitleaks", "project_id":1, "type":2, "target_resource":"pipe-cd"}}' \
		localhost:8081 datasource.code.CodeService.PutGitleaks

.PHONY: delete-gitleaks
delete-gitleaks:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "gitleaks_id":1}' \
		localhost:8081 datasource.code.CodeService.DeleteGitleaks

.PHONY: list-enterprise-org
list-enterprise-org:
	grpcurl \
		-plaintext \
		-d '{"project_id": 1001, "gitleaks_id":1}' \
		localhost:8081 datasource.code.CodeService.ListEnterpriseOrg

.PHONY: put-enterprise-org
put-enterprise-org:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "enterprise_org": {"gitleaks_id":1, "login":"login", "project_id":1}}' \
		localhost:8081 datasource.code.CodeService.PutEnterpriseOrg

.PHONY: delete-enterprise-org
delete-enterprise-org:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "gitleaks_id":1, "login": "login"}' \
		localhost:8081 datasource.code.CodeService.DeleteEnterpriseOrg

.PHONY: invoke-scan-gitleaks
invoke-scan-gitleaks:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "gitleaks_id":4}' \
		localhost:8081 datasource.code.CodeService.InvokeScanGitleaks

.PHONY: invoke-scan-all-gitleaks
invoke-scan-all-gitleaks:
	grpcurl -plaintext localhost:8081 datasource.code.CodeService.InvokeScanAllGitleaks

####################################################
## Diagnosis
####################################################
.PHONY: list-diagnosis-service
list-diagnosis-service:
	grpcurl -plaintext localhost:8081 list datasource.diagnosis.DiagnosisService

.PHONY: list-diagnosis_datasource
list-diagnosis_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.ListDiagnosisDataSource

.PHONY: get-diagnosis_datasource
get-diagnosis_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "diagnosis_data_source_id":1001}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.GetDiagnosisDataSource

.PHONY: put-diagnosis_datasource
put-diagnosis_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1 ,"diagnosis_data_source":{"name":"test_ds","description":"for_test","max_score":10}}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.PutDiagnosisDataSource

.PHONY: delete-diagnosis_datasource
delete-diagnosis_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "diagnosis_data_source_id":1002}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.DeleteDiagnosisDataSource

.PHONY: list-wpscan_setting
list-wpscan_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.ListWpscanSetting

.PHONY: get-wpscan_setting
get-wpscan_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "wpscan_setting_id":1}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.GetWpscanSetting

.PHONY: put-wpscan_setting
put-wpscan_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "wpscan_setting": {"diagnosis_data_source_id":1002, "project_id":1, "target_url":"http://example.com", "status":"CONFIGURED", "options":"{}"}}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.PutWpscanSetting

.PHONY: delete-wpscan_setting
delete-wpscan_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "wpscan_setting_id":1}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.DeleteWpscanSetting

.PHONY: list-portscan_setting
list-portscan_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.ListPortscanSetting

.PHONY: get-portscan_setting
get-portscan_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "portscan_setting_id":1}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.GetPortscanSetting

.PHONY: put-portscan_setting
put-portscan_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "portscan_setting": {"diagnosis_data_source_id":1003, "project_id":1, "name":"test_portscan"}}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.PutPortscanSetting

.PHONY: delete-portscan_setting
delete-portscan_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "portscan_setting_id":1}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.DeletePortscanSetting

.PHONY: list-portscan_target
list-portscan_target:
	grpcurl \
		-plaintext \
		-d '{"project_id":1}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.ListPortscanTarget

.PHONY: get-portscan_target
get-portscan_target:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "portscan_target_id":1}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.GetPortscanTarget

.PHONY: put-portscan_target
put-portscan_target:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "portscan_target": {"portscan_setting_id":1, "project_id":1, "target":"127.0.0.1", "status":"CONFIGURED"}}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.PutPortscanTarget

.PHONY: delete-portscan_target
delete-portscan_target:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "portscan_target_id":1001}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.DeletePortscanTarget

.PHONY: list-application_scan
list-application_scan:
	grpcurl \
		-plaintext \
		-d '{"project_id":1}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.ListApplicationScan

.PHONY: get-application_scan
get-application_scan:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "application_scan_id":1001}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.GetApplicationScan

.PHONY: put-application_scan
put-application_scan:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "application_scan": {"diagnosis_data_source_id":1004, "project_id":1, "name":"test_target","scan_type":"BASIC","status":"CONFIGURED"}}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.PutApplicationScan

.PHONY: delete-application_scan
delete-application_scan:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "application_scan_id":1002}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.DeleteApplicationScan

.PHONY: list-application_scan_basic_setting
list-application_scan_basic_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.ListApplicationScanBasicSetting

.PHONY: get-application_scan_basic_setting
get-application_scan_basic_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "application_scan_id":1001}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.GetApplicationScanBasicSetting

.PHONY: put-application_scan_basic_setting
put-application_scan_basic_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "application_scan_basic_setting": {"application_scan_id":1, "project_id":1, "target":"http://localhost:8080", "max_depth":10, "max_children": 10}}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.PutApplicationScanBasicSetting

.PHONY: delete-application_scan_basic_setting
delete-application_scan_basic_setting:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "application_scan_basic_setting_id":1}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.DeleteApplicationScanBasicSetting

.PHONY: invoke-scan-wpscan
invoke-scan-wpscan:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "setting_id":1,"diagnosis_data_source_id":1002}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.InvokeScan

.PHONY: invoke-scan-portscan
invoke-scan-portscan:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "setting_id":1,"diagnosis_data_source_id":1003}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.InvokeScan

.PHONY: invoke-scan-application-scan
invoke-scan-application-scan:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "setting_id":1,"diagnosis_data_source_id":1004}' \
		localhost:8081 datasource.diagnosis.DiagnosisService.InvokeScan

.PHONY: invoke-diagnosis-scan-all
invoke-diagnosis-scan-all:
	grpcurl \
		-plaintext \
		localhost:8081 datasource.diagnosis.DiagnosisService.InvokeScanAll

####################################################
## Google
####################################################
.PHONY: list-google-service
list-google-service:
	grpcurl -plaintext localhost:8081 list datasource.google.GoogleService

.PHONY: list-google-datasource
list-google-datasource:
	grpcurl \
		-plaintext \
		-d '{"google_data_source_id":1001}' \
		localhost:8081 datasource.google.GoogleService.ListGoogleDataSource

.PHONY: list-gcp
list-gcp:
	grpcurl \
		-plaintext \
		-d '{"project_id":1}' \
		localhost:8081 datasource.google.GoogleService.ListGCP

.PHONY: get-gcp
get-gcp:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "gcp_id":1}' \
		localhost:8081 datasource.google.GoogleService.GetGCP

.PHONY: put-gcp
put-gcp:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "gcp": {"name":"1", "project_id":1, "gcp_project_id":"my-project", "verification_code":"xxxxxxxx"}}' \
		localhost:8081 datasource.google.GoogleService.PutGCP

.PHONY: delete-gcp
delete-gcp:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "gcp_id":1}' \
		localhost:8081 datasource.google.GoogleService.DeleteGCP

.PHONY: list-gcp-datasource
list-gcp-datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "gcp_id":1}' \
		localhost:8081 datasource.google.GoogleService.ListGCPDataSource

.PHONY: get-gcp-datasource
get-gcp-datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "gcp_id":1, "google_data_source_id":1001}' \
		localhost:8081 datasource.google.GoogleService.GetGCPDataSource

.PHONY: attach-gcp-datasource
attach-gcp-datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "gcp_data_source": {"gcp_id":1, "google_data_source_id":1001 "project_id":1}}' \
		localhost:8081 datasource.google.GoogleService.AttachGCPDataSource

.PHONY: detach-gcp-datasource
detach-gcp-datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "gcp_id":1, "google_data_source_id":1001}' \
		localhost:8081 datasource.google.GoogleService.DetachGCPDataSource

.PHONY: invoke-scan-gcp
invoke-scan-gcp:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "gcp_id":1, "google_data_source_id":1001}' \
		localhost:8081 datasource.google.GoogleService.InvokeScanGCP

.PHONY: invoke-google-scan-all
invoke-google-scan-all:
	grpcurl -plaintext localhost:8081 datasource.google.GoogleService.InvokeScanAll

####################################################
## OSINT
####################################################
.PHONY: list-osint-service
list-osint-service:
	grpcurl -plaintext localhost:8081 list datasource.osint.OsintService

.PHONY: list-osint
list-osint:
	grpcurl \
		-plaintext \
		-d '{"project_id":1}' \
		localhost:8081 datasource.osint.OsintService.ListOsint

.PHONY: get-osint
get-osint:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "osint_id":1}' \
		localhost:8081 datasource.osint.OsintService.GetOsint

.PHONY: put-osint
put-osint:
	grpcurl \
		-plaintext \
		-d '{"project_id":1 ,"osint":{"resource_type":"Domain","resource_name":"cyberagent.co.jp","project_id":1}}' \
		localhost:8081 datasource.osint.OsintService.PutOsint

.PHONY: delete-osint
delete-osint:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "osint_id":2}' \
		localhost:8081 datasource.osint.OsintService.DeleteOsint

.PHONY: list-osint_datasource
list-osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1}' \
		localhost:8081 datasource.osint.OsintService.ListOsintDataSource

.PHONY: get-osint_datasource
get-osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "osint_data_source_id":1001}' \
		localhost:8081 datasource.osint.OsintService.GetOsintDataSource

.PHONY: put-osint_datasource
put-osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1 ,"osint_data_source":{"name":"test_ds","description":"for_test","max_score":10}}' \
		localhost:8081 datasource.osint.OsintService.PutOsintDataSource

.PHONY: delete-osint_datasource
delete-osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "osint_data_source_id":1002}' \
		localhost:8081 datasource.osint.OsintService.DeleteOsintDataSource

.PHONY: list-rel_osint_datasource
list-rel_osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1}' \
		localhost:8081 datasource.osint.OsintService.ListRelOsintDataSource

.PHONY: get-rel_osint_datasource
get-rel_osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "rel_osint_data_source_id":1001}' \
		localhost:8081 datasource.osint.OsintService.GetRelOsintDataSource

.PHONY: put-rel_osint_datasource
put-rel_osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "rel_osint_data_source": {"osint_id":1, "osint_data_source_id":1001, "project_id":1, "status":"CONFIGURED"}}' \
		localhost:8081 datasource.osint.OsintService.PutRelOsintDataSource

.PHONY: delete-rel_osint_datasource
delete-rel_osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "rel_osint_data_source_id":1}' \
		localhost:8081 datasource.osint.OsintService.DeleteRelOsintDataSource

.PHONY: list-osint_detect_word
list-osint_detect_word:
	grpcurl \
		-plaintext \
		-d '{"project_id":1}' \
		localhost:8081 datasource.osint.OsintService.ListOsintDetectWord

.PHONY: get-osint_detect_word
get-osint_detect_word:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "osint_detect_word_id":1}' \
		localhost:8081 datasource.osint.OsintService.GetOsintDetectWord

.PHONY: put-osint_detect_word
put-osint_detect_word:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "osint_detect_word": {"word":"fuga","rel_osint_data_source_id":1, "project_id":1}}' \
		localhost:8081 datasource.osint.OsintService.PutOsintDetectWord

.PHONY: delete-osint_detect_word
delete-osint_detect_word:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "osint_detect_word_id":1}' \
		localhost:8081 datasource.osint.OsintService.DeleteOsintDetectWord

.PHONY: invoke-osint-scan
invoke-osint-scan:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "rel_osint_data_source_id":1}' \
		localhost:8081 datasource.osint.OsintService.InvokeScan

.PHONY: invoke-osint-scan_all
invoke-osint-scan_all:
	grpcurl \
		-plaintext \
		localhost:8081 datasource.osint.OsintService.InvokeScanAll

