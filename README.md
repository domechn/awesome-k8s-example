# k8s-admission-webhook-example

Admission webhook controller example with less code

## Deploy

1. `/bin/sh -c .deployment/setup-certificates.sh` to generate the required secret for admission webhook, you need to modify `namespace`, `service`, and `secret` in `setup-ceritificates.sh`
2. `kubectl apply -f .deployment/deploy.yaml`
3. `kubectl apply -f .deployment/webhook.yaml`
