package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCronTabEntryResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCronTabEntryResourceConfig("* * * * * ls > /dev/null 2>&1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_crontab_entry.test",
						tfjsonpath.New("entry"),
						knownvalue.StringExact("* * * * * ls > /dev/null 2>&1"),
					),
				},
			},
			// ImportState testing
			// {
			// 	ResourceName:      "uberspace_crontab_entry.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
			// Update and Read testing
			{
				Config: testAccCronTabEntryResourceConfig("* * * * * ls -l > /dev/null 2>&1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_crontab_entry.test",
						tfjsonpath.New("entry"),
						knownvalue.StringExact("* * * * * ls -l > /dev/null 2>&1"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCronTabEntryResourceConfig(entry string) string {
	return fmt.Sprintf(`
resource "uberspace_crontab_entry" "test" {
  entry = %[1]q
}
`, entry)
}
