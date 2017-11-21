# Slack POC

A Slack app built in Golang

Start the app by executing:

```golang
    go run server.go --queues=zendesk
```

Start the worker by

```golang
    go run worker.go --queues=zendesk
```