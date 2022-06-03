FROM golang:1.17 as builder
WORKDIR /src/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/datasource-api

FROM public.ecr.aws/risken/base/risken-base:v0.0.1
COPY --from=builder /bin/datasource-api /usr/local/bin/
ENV PORT= \
    PROFILE_EXPORTER= \
    PROFILE_TYPES= \
    DB_MASTER_HOST= \
    DB_MASTER_USER= \
    DB_MASTER_PASSWORD= \
    DB_SLAVE_HOST= \
    DB_SLAVE_USER= \
    DB_SLAVE_PASSWORD= \
    DB_SCHEMA=mimosa \
    DB_PORT=3306 \
    DB_LOG_MODE=false \
    AWS_REGION= \
    AWS_ACCESS_KEY_ID= \
    AWS_SECRET_ACCESS_KEY= \
    AWS_SESSION_TOKEN= \
    SQS_ENDPOINT= \
    GUARD_DUTY_QUEUE_URL= \
    ACCESS_ANALYZER_QUEUE_URL= \
    ADMIN_CHECKER_QUEUE_URL= \
    CLOUDSPLOIT_QUEUE_URL= \
    PORTSCAN_QUEUE_URL= \
    ASSET_QUEUE_URL= \
    CLOUD_SPLOIT_QUEUE_URL= \
    PORTSCAN_QUEUE_URL= \
    GOOGLE_CREDENTIAL_PATH= \
    GOOGLE_SERVICE_ACCOUNT_JSON= \
    GITLEAKS_QUEUE_URL= \
    GITLEAKS_FULL_SCAN_QUEUE_URL= \
    SUBDOMAIN_QUEUE_URL= \
    WEBSITE_QUEUE_URL= \
    DIAGNOSIS_WPSCAN_QUEUE_URL= \
    DIAGNOSIS_APPLICATIONSCAN_QUEUE_URL= \
    DIAGNOSIS_PORTSCAN_QUEUE_URL= \
    DATA_KEY= \
    TZ=Asia/Tokyo
WORKDIR /usr/local/
ENTRYPOINT ["/usr/local/bin/env-injector", "docker-entrypoint.sh"]
CMD ["bin/datasource-api"]
