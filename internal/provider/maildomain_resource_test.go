package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccMaildomainResource(t *testing.T) {
	maildomain := fmt.Sprintf("%s.terra.uber.space", acctest.RandomWithPrefix("mail"))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMaildomainResourceConfig("terra", maildomain),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_maildomain.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(maildomain),
					),
					statecheck.ExpectKnownValue(
						"uberspace_maildomain.test",
						tfjsonpath.New("asteroid"),
						knownvalue.StringExact("terra"),
					),
				},
			},
			{
				Config: testAccMaildomainResourceConfig("terra", maildomain),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_maildomain.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(maildomain),
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
