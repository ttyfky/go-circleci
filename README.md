# go-circleci

[![godoc](https://godoc.org/github.com/ttyfky/go-circleci?status.svg)](https://pkg.go.dev/github.com/ttyfky/go-circleci)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

This is a Go client of [CircleCI API v2](https://circleci.com/docs/api/v2/). 

This API is implemented by following the public docs by CircleCI.
* [API Reference](https://circleci.com/docs/2.0/api-intro/) 
* [CircleCI API v2](https://circleci.com/docs/api/v2/)
* [CircleCI API Developer's Guide](https://circleci.com/docs/2.0/api-developers-guide/)

## Install

```console
$ go get github.com/ttyfky/go-circleci
```

## Use

```go
import "github.com/ttyfky/go-circleci"
```

### Authentication
CircleCI supports two types of authentication 1. api_key_header and 2. basic_auth.

This Client currently supports api_key_header authentication.

#### Get API token
Get API token by following [instruction](https://circleci.com/docs/2.0/api-developers-guide/#authentication-and-authorization) before using this client.


### Create Client
Give API token to `NewClient` function.
```go
token := "API_TOKEN"
client := circleci.NewClient(token)
```

### API call
Use resource service in the client to call API of each resources in CircleCI.

```go
client := circleci.NewClient(token)
workflowID := "ID"
workflow, _ := client.Workflow.Get(workflowID)
```

More examples are availablein [example_test.go](./example_test.go).

# API availability

Not all of the APIs are implemented yet. It's more based on demand of the actions. 
The table below shows the API's availability as of Jan 2021.

`Preview` means it's preview phase in CircleCI side.

| API               | Availability |
|-------------------|--------------|
| Context (Preview) |  Available |
| Insights          |  Not Implemented |
| User (Preview)    |  Not Implemented |
| Pipeline          |  Not Implemented |
| Job (Preview)     |  Available |
| Workflow          |  Available |
| Project           |  Partially Available |

Note: Environment variable handling is part of Project API, but extracted as `ProjectEnvVar` it for convenience. 
