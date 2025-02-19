package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccMySQLDatabaseResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMySQLDatabaseResourceConfig("example"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_mysql_database.test",
						tfjsonpath.New("suffix"),
						knownvalue.StringExact("example"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "uberspace_mysql_database.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccMySQLDatabaseResourceConfig("example"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_mysql_database.test",
						tfjsonpath.New("suffix"),
						knownvalue.StringExact("example"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMySQLDatabaseResourceConfig(suffix string) string {
	return fmt.Sprintf(`
resource "uberspace_mysql_database" "test" {
  suffix = %[1]q
}
`, suffix)
}
