package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccWebDomainResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccWebDomainResourceConfig("test.terra.uber.space"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_web_domain.test",
						tfjsonpath.New("domain"),
						knownvalue.StringExact("test.terra.uber.space"),
					),
				},
			},
			// ImportState testing
			// {
			// 	ResourceName:      "uberspace_web_domain.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
			// Update and Read testing
			{
				Config: testAccWebDomainResourceConfig("test.terra.uber.space"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_web_domain.test",
						tfjsonpath.New("domain"),
						knownvalue.StringExact("test.terra.uber.space"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccWebDomainResourceConfig(domain string) string {
	return fmt.Sprintf(`
resource "uberspace_web_domain" "test" {
  domain = %[1]q
}
`, domain)
}
