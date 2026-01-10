package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccMaildomainResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMaildomainResourceConfig("terra", "mail.terra.uber.space"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_maildomain.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("mail.terra.uber.space"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_maildomain.test",
						tfjsonpath.New("asteroid"),
						knownvalue.StringExact("terra"),
					),
				},
			},
			{
				Config: testAccMaildomainResourceConfig("terra", "mail.terra.uber.space"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_maildomain.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("mail.terra.uber.space"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_maildomain.test",
						tfjsonpath.New("asteroid"),
						knownvalue.StringExact("terra"),
					),
				},
			},
		},
	})
}

func testAccMaildomainResourceConfig(asteroid, name string) string {
	return fmt.Sprintf(`
resource "uberspace_maildomain" "test" {
  asteroid = %[1]q
  name     = %[2]q
}
`, asteroid, name)
}
