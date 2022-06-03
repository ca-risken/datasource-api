# RISKEN DataSource API

# RISKEN AWS

![Build Status]()

`RISKEN` is a monitoring tool for your cloud platforms, web-site, source-code... 
`RISKEN DataSource API` is a security monitoring system for datasouces(e.g. aws, google, osint..) that searches, analyzes, evaluate, and alerts on discovered threat information.

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

`k8s-sample/overlays/local/datasource-api.yaml`

| service        | spec                                | before (public images) | after (pre-build images on your machine) |
| -------------- | ----------------------------------- | ---------------------- | ---------------------------------------- |
| datasource-api | spec.template.spec.containers.image | `TBD`                  | `risken-datasource-api:latest`           |

## Community

Info on reporting bugs, getting help, finding roadmaps,
and more can be found in the [RISKEN Community](https://github.com/ca-risken/community).

## License

[MIT](LICENSE).
