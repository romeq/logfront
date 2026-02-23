# Contributing
### Initial setup

Set up your local environment before you start development (see below).

1. Install and update Golang and [`golangci-lint`](https://golangci-lint.run)
    ```shell
   # Make sure both your applications run at least 1.26
   $ go version
   go version go1.26.0 #..
   $ golangci-lint version # has to be built with go1.26 or newer
   golangci-lint has version 2.10.1 built with go1.26.0 #...
    ```
2. Set up git hooks
    ```shell
   git config --local core.hooksPath .githooks
    ```

### Before creating a PR

```shell
# run linters
golangci-lint run
# make sure the project builds and works
go run ./cmd/core -f config.yaml
# make sure tests pass
go test ./...
```
