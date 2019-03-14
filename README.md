# validate-resource-limits

## Steps to deploy

### Build the image

*** TODO (build, update manifest)

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