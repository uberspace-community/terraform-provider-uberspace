# Terraform Provider for uberspace (U8)

This repository contains the Terraform provider for managing resources on [uberspace](https://uberspace.de/) U8 accounts.

## Usage

```terraform
terraform {
  required_providers {
    uberspace = {
      source  = "uberspace-community/uberspace"
      version = "0.2.0-alpha.2"
    }
  }

  required_version = ">= 1.5.0"
}

provider "uberspace" {
  # Alternatively set UBERSPACE_APIKEY in the environment
  apikey = "example-api-key"
}

resource "uberspace_sshkey" "example" {
  asteroid = "isabell"
  key = filebase64("~/.ssh/id_ed25519.pub")
  key_type = "ssh-ed25519"
}

resource "uberspace_webdomain" "minio" {
  asteroid = "isabell"
  domain   = "minio.isabell.uber.space"
}

resource "uberspace_webdomain_backend" "minio" {
  // a web backend usually depends on a web domain
  depends_on = [uberspace_webdomain.minio]

  asteroid    = "isabell"
  destination = "STATIC"
  domain      = "minio.isabell.uber.space"
  path        = "/foo"
}

resource "uberspace_webdomain_header" "cors" {
  // a web backend usually depends on a web domain
  depends_on = [uberspace_webdomain.minio]

  asteroid   = "isabell"
  domain     = "minio.isabell.uber.space"
  path       = "/"
  name       = "X-Custom-Header"
  value      = "custom"
}
```

See the [uberspace Provider Documentation](https://registry.terraform.io/providers/uberspace-community/uberspace/latest/docs)
for detailed resource and data source information.

## Retrieving the API Key

To use the provider, you need to retrieve your API key from your uberspace account. 

1. Connect to your uberspace account via SSH.
2. Run `cat /readonly/$USER/marvin_client.toml`
3. Copy the value of `api_key` from the output.

## Developing & Contributing to the Provider

The [DEVELOPER.md](DEVELOPER.md) file is a basic outline on how to build and develop the provider locally.
