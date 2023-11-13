# EC2 Info Application

This Go application retrieves information about EC2 instances and their associated Amazon Machine Images (AMIs) using the AWS SDK for Go.

## Prerequisites

- [Go](https://golang.org/dl/) installed on your machine.
- AWS credentials configured with the necessary permissions.
- Ensure that the desired AWS region is set by exporting the variable, for example: `export AWS_REGION=us-east-1`.
- Ensure that folangci-lint is installer `go get -u github.com/golangci/golangci-lint/cmd/golangci-lint`

## Getting Started

Clone the repository:

```bash
git clone https://github.com/your-username/ec2-info.git
cd ec2-info
```

## Makefile Commands

The project includes a Makefile to simplify common development tasks:

### Install Dependencies

```bash
make deps
```

Downloads project dependencies.

### Install Dependencies

```bash
make lint
```

Runs linters to ensure code quality.

### Run Unit Tests

```bash
make test
```

Executes unit tests.

### Build Application

```bash
make build
```

Builds the application.

### Build and Run

```bash
make run
```

Builds and runs the application.

### Clean

```bash
make clean
```

Removes built binaries.

## Usage

To run the application:

```bash
make run
```

alternative usage:

```go
go run main.go
```

This will fetch information about EC2 instances and AMIs, printing the result in a pretty JSON format.

## License

This project is licensed under the MIT License.
