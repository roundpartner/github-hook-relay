# GitHook SQS Relay
Relays SQS into webhooks

# Building
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w'
docker build -t github-hook-relay .
```

```bash
docker run \
    --name github-hook-relay \
    --rm \
    --env AWS_REGION=eu-west-2 \
    github-hook-relay
```
