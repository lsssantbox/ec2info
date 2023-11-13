# EC2 Info Application

This Go application retrieves information about EC2 instances and their associated Amazon Machine Images (AMIs) using the AWS SDK for Go printing the result in a JSON format.


## Prerequisites

- [Go](https://golang.org/dl/) installed on your machine.
- AWS credentials configured with the necessary permissions.
- Ensure that the desired AWS region is set by exporting the variable, for example: `export AWS_REGION=us-east-1`.


## Getting Started

```bash
git clone https://github.com/lsssantbox/ec2info
cd ec2info
make run 
```

## Makefile Commands

The project includes a Makefile to simplify common development tasks:

```bash
$ make
Usage:
  make deps       - Download dependencies
  make lint       - Run linters
  make test       - Run unit tests
  make build      - Build the application
  make run        - Build and run the application
  make clean      - Remove built binaries

```


## License

This project is licensed under the MIT License.
