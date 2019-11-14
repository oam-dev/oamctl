# oamctl

oamctl is a tiny tool help oam user manage oam app.


## Install

### Download the repo

```shell script
$ mkdir -p ${GOPATH}/src/github.com/oam-dev/
$ cd ${GOPATH}/src/github.com/oam-dev/
$ git clone git@github.com:oam-dev/oamctl.git
```

### Build the binary

Run the following command to install _oamctl_

```
go install
```


## Usage

### [oamctl migrate](docs/migrate.md)

This command can create oam ApplicationConfiguration yaml from exist k8s resource

### Trait-list

This command will list trait with version and workload type that could be applied.