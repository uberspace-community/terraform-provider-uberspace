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
				Config: testAccWebdomainBackendResourceConfig("tf", "test.tf.uber8.space", "/tf-backend", false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_backend.test",
						tfjsonpath.New("asteroid"),
						knownvalue.StringExact("tf"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_backend.test",
						tfjsonpath.New("domain"),
						knownvalue.StringExact("test.tf.uber8.space"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_backend.test",
						tfjsonpath.New("destination"),
						knownvalue.StringExact("STATIC"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_backend.test",
						tfjsonpath.New("path"),
						knownvalue.StringExact("/tf-backend"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_backend.test",
						tfjsonpath.New("remove_prefix"),
						knownvalue.Bool(false),
					),
				},
			},
			{
				Config: testAccWebdomainBackendResourceConfig("tf", "test.tf.uber8.space", "/tf-backend-updated", true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_backend.test",
						tfjsonpath.New("path"),
						knownvalue.StringExact("/tf-backend-updated"),
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

func testAccWebdomainBackendResourceConfig(asteroid, domain, path string, removePrefix bool) string {
	return fmt.Sprintf(`
resource "uberspace_webdomain" "test" {
  asteroid = %[1]q
  name     = %[2]q
}

resource "uberspace_webdomain_backend" "test" {
  asteroid      = %[1]q
  destination   = "STATIC"
  domain        = uberspace_webdomain.test.name
  path          = %[3]q
  remove_prefix = %[4]t
}
`, asteroid, domain, path, removePrefix)
}
