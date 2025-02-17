package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccRemoteFileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRemoteFileResourceConfig("example", "echo 'Hello, World!'"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_remote_file.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("example"),
					),
				},
			},
			// ImportState testing
			/*{
				ResourceName:      "uberspace_remote_file.test",
				ImportState:       true,
				ImportStateVerify: true,
			},*/
			// Update and Read testing
			{
				Config: testAccRemoteFileResourceConfig("example", "echo 'Hello, World!'"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_remote_file.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("example"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRemoteFileResourceConfig(name, command string) string {
	return fmt.Sprintf(`
resource "uberspace_remote_file" "test" {
  name = %[1]q
  command = %[1]q
}
`, name, command)
}
