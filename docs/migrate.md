## Usage

create oam ApplicationConfiguration yaml from exist k8s resource, e.g: deployment, sts, service....

```
Usage:
  oamctl migrate [flags] <name>

Flags:
      --deployment string   name of deployment
  -h, --help                help for migrate
      --ingress string      name of ingress
      --service string      name of service

Params:
      name string   name of oam app

Global Flags:
      --config string      config file (default is $HOME/.oamctl.yaml)
  -n, --namespace string   operate namespace (default "default")

```

## Examples
### Migrate to OAM from k8s deployment

Assume we have a deployment and apply it:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-2migrate
  labels:
    app: test
spec:
  replicas: 3
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
        - name: test
          image: somefive/hello-world:dev
          env:
            - name: k1
              value: v1
            - name: k2
              value: v2

```

We can use oamctl migrate this deployment to oam app:

```
oamctl migrate --deployment=deployment-2migrate  oam-app

```

We will get:

```
apiVersion: core.oam.dev/v1alpha1
kind: ComponentSchematic
metadata:
  name: deployment-2migrate
spec:
  workloadType: core.oam.dev/v1alpha1.Server
  containers:
    - name: test
      image: somefive/hello-world:dev
      env:
        - name: k1
          fromParam: test-k1
        - name: k2
          fromParam: test-k2
      ports:
  parameters:
    - name: test-k1
      value: v1
    - name: test-k2
      value: v2
---
apiVersion: core.oam.dev/v1alpha1
kind: ApplicationConfiguration
metadata:
  name: oam-app
spec:
  components:
    - name: deployment-2migrate
      instanceName: deployment-2migrate_instance
      parameterValues:
        - name: test-k1
          value: v1
        - name: test-k2
          value: v2
      traits:
        - name: manual-scaler
          parameterValues:
            - name: replicaCount
              value: 3
```
