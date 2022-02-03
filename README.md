# Manifestor

Small utility to interpolate environment variables into Kubernetes manifests.
Unlike many other templating mechanisms, this utility will never output an invalid YAML.

Under the hood, it parses the input yaml files, finds string values, interpolates environment values into them using bash-like interpolation syntax and then re-assembles the resulting YAML(s) that are printed on stdout.

Also, there is an option of executing pre and post interpolation processors in ECMAScript 5.1(+).

## Use

```bash
$ VAR=VALUE manifestor [--processors=] [<yaml file> ...] | kubectl apply -f -
```

## Interpolation format

We are using [envsubst](https://github.com/drone/envsubst) to interpolate environment variables into the string values.

TL;DR;

running
```bash
HOSTNAME="localhost" manifestor ingress.yaml
```

where `ingress.yaml` template looks like this:
```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
    annotations:
        kubernetes.io/ingress.class: nginx
    name: api
    namespace: api
spec:
    rules:
        - host: ${HOSTNAME}
          http:
            paths:
                - backend:
                    serviceName: api
                    servicePort: 3030
                  path: /api
```

will result in the following output:

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
    annotations:
        kubernetes.io/ingress.class: nginx
    name: api
    namespace: api
spec:
    rules:
        - host: localhost
          http:
            paths:
                - backend:
                    serviceName: api
                    servicePort: 3030
                  path: /api
```

## Pre and post interpolation processors

String interpolation of environment variables is quite limited.
It won't cater for the cases where a number, boolean value or a whole sub-object needs to be changed, inserted or deleted.
In order to cover those cases, we support executing so-called _processors_.

A _processor_ is a JavaScript function that takes an object representing one entity from the manifest and can change it in place.

### Naming convention

TODO

### Available functions

TODO

### Available objects

TODO
