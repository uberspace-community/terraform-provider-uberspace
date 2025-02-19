package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccWebBackendResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWebBackendResourceConfig("example.terra.uber.space", 9090),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_web_backend.test",
						tfjsonpath.New("uri"),
						knownvalue.StringExact("example.terra.uber.space/"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_web_backend.test",
						tfjsonpath.New("port"),
						knownvalue.Int32Exact(9090),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "uberspace_web_backend.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccWebBackendResourceConfig("example.terra.uber.space", 9092),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_web_backend.test",
						tfjsonpath.New("uri"),
						knownvalue.StringExact("example.terra.uber.space/"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_web_backend.test",
						tfjsonpath.New("port"),
						knownvalue.Int32Exact(9092),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWebBackendResourceConfig(domain string, port int) string {
	return fmt.Sprintf(`
resource "uberspace_web_domain" "test" {
  domain = %[1]q
}

resource "uberspace_web_backend" "test" {
  depends_on = [uberspace_web_domain.test]

  uri = %[2]q
  port = %[3]d
}
`, domain, domain+"/", port)
}
