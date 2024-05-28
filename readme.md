# Go Multithreading Challenge

This challenge applies the Go Multithreading feature to APIs and seeks the fastest result between different integrations.

### Prerequisites

- [Golang v1.22+](https://golang.org/) 

## Deliverables

You will need to deliver:
- `main.go`

## Requirements

The two requests will be made simultaneously to the following APIs:

- [BrasilAPI](https://brasilapi.com.br/api/cep/v1/01153000 + cep)
- [ViaCEP](http://viacep.com.br/ws/" + cep + "/json/)

The requirements for this challenge are:

1. **Accept the Faster Response:**
   - Accept the API that delivers the fastest response and discard the slower response.

2. **Command Line Output:**
   - The result of the request must be displayed on the command line with the address data, as well as which API sent it.

3. **Response Time Limit:**
   - Limit the response time to 1 second. Otherwise, a timeout error should be displayed.

### Running the aplication

`go run main.go`