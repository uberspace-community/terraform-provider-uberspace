package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccWebdomainHeaderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWebdomainHeaderResourceConfig(
					"tf",
					"test-header.tf.uber8.space",
					"X-Custom-Header",
					"initial",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_header.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("X-Custom-Header"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_header.test",
						tfjsonpath.New("value"),
						knownvalue.StringExact("initial"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_header.test",
						tfjsonpath.New("path"),
						knownvalue.StringExact("/"),
					),
				},
			},
			{
				Config: testAccWebdomainHeaderResourceConfig(
					"tf",
					"test-header.tf.uber8.space",
					"X-Custom-Header",
					"updated",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_webdomain_header.test",
						tfjsonpath.New("value"),
						knownvalue.StringExact("updated"),
					),
				},
			},
		},
	})
}

func testAccWebdomainHeaderResourceConfig(asteroid, domain, headerName, headerValue string) string {
	return fmt.Sprintf(`
resource "uberspace_webdomain" "test" {
  asteroid = %[1]q
  name     = %[2]q
}

resource "uberspace_webdomain_header" "test" {
  depends_on = [uberspace_webdomain.test]

  asteroid = %[1]q
  domain   = %[2]q
  path     = "/"
  name     = %[3]q
  value    = %[4]q
}
`, asteroid, domain, headerName, headerValue)
}
