package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccWebdomainResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWebdomainResourceConfig("tf", "test.tf.uber8.space"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_webdomain.test",
						tfjsonpath.New("domain"),
						knownvalue.StringExact("test.tf.uber8.space"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_webdomain.test",
						tfjsonpath.New("asteroid"),
						knownvalue.StringExact("tf"),
					),
				},
			},
			// ImportState testing
			// {
			// 	ResourceName:      "uberspace_webdomain.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
			// Update and Read testing
			{
				Config: testAccWebdomainResourceConfig("tf", "test.tf.uber8.space"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_webdomain.test",
						tfjsonpath.New("domain"),
						knownvalue.StringExact("test.tf.uber8.space"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_webdomain.test",
						tfjsonpath.New("asteroid"),
						knownvalue.StringExact("tf"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWebdomainResourceConfig(asteroid, domain string) string {
	return fmt.Sprintf(`
resource "uberspace_webdomain" "test" {
  asteroid = %[1]q
  domain = %[2]q
}
`, asteroid, domain)
}
