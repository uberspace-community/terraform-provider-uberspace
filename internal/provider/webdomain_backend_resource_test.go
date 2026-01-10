package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccWebdomainBackendResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWebdomainBackendResourceConfig("terra", "test-backend.terra.uber.space", 1024, "/terra-backend", false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_backend.test",
						tfjsonpath.New("asteroid"),
						knownvalue.StringExact("terra"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_backend.test",
						tfjsonpath.New("domain"),
						knownvalue.StringExact("test-backend.terra.uber.space"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_backend.test",
						tfjsonpath.New("port"),
						knownvalue.Int64Exact(1024),
					),
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_backend.test",
						tfjsonpath.New("path"),
						knownvalue.StringExact("/terra-backend"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_backend.test",
						tfjsonpath.New("remove_prefix"),
						knownvalue.Bool(false),
					),
				},
			},
			{
				Config: testAccWebdomainBackendResourceConfig("terra", "test-backend.terra.uber.space", 1024, "/terra-backend-updated", true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_backend.test",
						tfjsonpath.New("path"),
						knownvalue.StringExact("/terra-backend-updated"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_backend.test",
						tfjsonpath.New("remove_prefix"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func testAccWebdomainBackendResourceConfig(asteroid, domain string, port int, path string, removePrefix bool) string {
	return fmt.Sprintf(`
resource "uberspace_webdomain" "test" {
  asteroid = %[1]q
  name     = %[2]q
}

resource "uberspace_webdomain_backend" "test" {
  asteroid      = %[1]q
  domain        = uberspace_webdomain.test.name
  destination   = "PORT"
  port          = %[3]d
  path          = %[4]q
  remove_prefix = %[5]t
}
`, asteroid, domain, port, path, removePrefix)
}
