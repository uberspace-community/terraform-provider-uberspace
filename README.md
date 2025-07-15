# Terraform Provider for uberspace (U8)

This repository contains the Terraform provider for managing resources on [uberspace](https://uberspace.de/) U8 accounts.

## Usage

```terraform
provider "uberspace" {
  # Alternatively set UBERSPACE_APIKEY in the environment
  apikey = "example-api-key"
}
```

See the [uberspace Provider Documentation](https://registry.terraform.io/providers/uberspace-community/uberspace/latest/docs)
for detailed resource and data source information.

## Developing & Contributing to the Provider

The [DEVELOPER.md](DEVELOPER.md) file is a basic outline on how to build and develop the provider locally.
