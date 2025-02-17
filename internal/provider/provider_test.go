package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"uberspace": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	host := os.Getenv("UBERSPACE_HOST")
	if host == "" {
		t.Fatal("UBERSPACE_HOST must be set for acceptance tests")
	}

	user := os.Getenv("UBERSPACE_USER")
	if user == "" {
		t.Fatal("UBERSPACE_USER must be set for acceptance tests")
	}

	password := os.Getenv("UBERSPACE_PASSWORD")
	privateKey := os.Getenv("UBERSPACE_PRIVATE_KEY")

	if password == "" && privateKey == "" {
		t.Fatal("either UBERSPACE_PASSWORD or UBERSPACE_PRIVATE_KEY must be set for acceptance tests")
	}
}
