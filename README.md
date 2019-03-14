# validate-resource-limits

This project is a sample Kubernetes [validating admission controller](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/) that verifies that all pods being deployed have resources limits (for CPU and memory) specified.

## Steps to deploy

### Build the image

```bash
docker build -t stuartleeks/validate-resource-limits .
```

NOTE: The image needs to be pushed ;-)

### Generate certs

Run `./scripts/create-certs.sh` to generate the certs

### Create the secrets

```bash
kubectl create secret generic validateresourcelimits -n default \
        --from-file=key.pem=generated/app.key \
        --from-file=cert.pem=generated/app.crt
```

### Susbtitute the CA bundle in the deployment manifest

```bash
./scripts/update-manifest.sh
```

### Deploy the generated manifest

```bash
 kubectl apply -f generated/webhook.yaml
 ```

## Misc notes

The service name `validateresourcelimits.default.svc` is specified in the `app.config` for cert generation