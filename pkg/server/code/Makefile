####################################################
## grpcurl example
####################################################
.PHONY: all
all: help

.PHONY: help
help:
	@echo "Usage: make <sub-command>"
	@echo "\n---------------- sub-command list ----------------"
	@cat Makefile | grep -e "^.PHONY:" | grep -v "all" | cut -f2 -d' '

.PHONY: list-service
list-service:
	grpcurl -plaintext localhost:8081 list datasource.code.CodeService

.PHONY: list-datasource
list-datasource:
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
