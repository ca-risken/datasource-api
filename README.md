# RISKEN DataSource API

# RISKEN AWS

![Build Status](https://codebuild.ap-northeast-1.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiM09KWEhQRHZEQTZoeDN5dnNxaU9MZDI5TCtLckRvd2dMQVRMUVpBeWd3WVg2S0lFNUVvd0ZsN3l1U014WC9nUW9RWlozcHlpN2FLeFl4ZjZEQm9keGZNPSIsIml2UGFyYW1ldGVyU3BlYyI6IlYzdEFGZFdDdkRsV1QyL04iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

`RISKEN` is a monitoring tool for your cloud platforms, web-site, source-code... 
`RISKEN DataSource API` provides APIs for controlling various datasources(e.g. aws, google, osint..). For example, to reference and register scan settings, and to invoke scan execution.

Please check [RISKEN Documentation](https://docs.security-hub.jp/).

## Installation

### Requirements

This module requires the following modules:

- [Go](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Protocol Buffer](https://grpc.io/docs/protoc-installation/)

### Install packages

This module is developed in the `Go language`, please run the following command after installing the `Go`.

```bash
$ make install
```

### Building

Build the containers on your machine with the following command

```bash
$ make build
```

### Running Apps

Deploy the pre-built containers to the Kubernetes environment on your local machine.

- Follow the [documentation](https://docs.security-hub.jp/admin/infra_local/#risken) to download the Kubernetes manifest sample.
- Fix the Kubernetes object specs of the manifest file as follows and deploy it.

TODO: `k8s-sample/overlays/local/datasource-api.yaml`

| service        | spec                                | before (public images) | after (pre-build images on your machine) |
| -------------- | ----------------------------------- | ---------------------- | ---------------------------------------- |
| datasource-api | spec.template.spec.containers.image | `TBD`                  | `risken-datasource-api:latest`           |

## Community

Info on reporting bugs, getting help, finding roadmaps,
and more can be found in the [RISKEN Community](https://github.com/ca-risken/community).

## License

[MIT](LICENSE).
