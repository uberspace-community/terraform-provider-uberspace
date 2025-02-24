# Developing & Contributing to the Provider

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24
- [uberspace](https://uberspace.de/) account in version U7

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

If you are building the provider, follow the instructions to install it as a plugin.
After placing it into your plugins directory, run terraform init to initialize it.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine
(see [Requirements](#requirements) above).

To compile the provider, run `go install`.
This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```