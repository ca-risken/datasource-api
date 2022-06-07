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

.PHONY: list-data-source
list-data-source:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "aws_id":1}' \
		localhost:8081 datasource.aws.AWSService.ListDataSource

.PHONY: attach-data-source
attach-data-source:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "attach_data_source": {"aws_id":1, "aws_data_source_id":1004, "project_id":1, "assume_role_arn":"arn:aws:iam::123456789012:role/role-name", "external_id":"test", "status":"CONFIGURED"}}' \
		localhost:8081 datasource.aws.AWSService.AttachDataSource

.PHONY: detach-data-source
detach-data-source:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "aws_id":1, "aws_data_source_id":1004}' \
		localhost:8081 datasource.aws.AWSService.DetachDataSource

.PHONY: invoke-scan
invoke-scan:
	grpcurl \
		-plaintext \
		-d '{"project_id":1, "aws_id":1, "aws_data_source_id":1001}' \
		localhost:8081 datasource.aws.AWSService.InvokeScan

.PHONY: invoke-scan-all
invoke-scan-all:
	grpcurl \
		-plaintext \
		localhost:8081 datasource.aws.AWSService.InvokeScanAll
