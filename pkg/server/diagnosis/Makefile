.PHONY: all
all: help

.PHONY: help
help:
	@echo "Usage: make <sub-command>"
	@echo "\n---------------- sub-command list ----------------"
	@cat Makefile | grep -e "^.PHONY:" | grep -v "all" | cut -f2 -d' '

.PHONY: list-service
list-service:
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

.PHONY: invoke-scan-all
invoke-scan-all:
	grpcurl \
		-plaintext \
		localhost:8081 datasource.diagnosis.DiagnosisService.InvokeScanAll