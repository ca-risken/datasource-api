.PHONY: all
all: help

.PHONY: help
help:
	@echo "Usage: make <sub-command>"
	@echo "\n---------------- sub-command list ----------------"
	@cat Makefile | grep -e "^.PHONY:" | grep -v "all" | cut -f2 -d' '

.PHONY: list-service
list-service:
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

.PHONY: invoke-scan
invoke-scan:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "rel_osint_data_source_id":1}' \
		localhost:8081 datasource.osint.OsintService.InvokeScan

.PHONY: invoke-scan_all
invoke-scan_all:
	grpcurl \
		-plaintext \
		localhost:8081 datasource.osint.OsintService.InvokeScanAll
