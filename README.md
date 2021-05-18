# Manifestor

Small utility to interpolate environment variables into Kubernetes manifests.
Unlike many other templating mechanisms, this utility will never output an invalid YAML.

Under the hood, it parses the input yaml files, finds string values, interpolates environment values into them using bash-like interpolation syntax and then re-assembles the resulting YAML(s) that are printed on stdout.

## Use

```bash
$ VAR=VALUE manifestor [<yaml file> ...] | kubectl apply -f -
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
