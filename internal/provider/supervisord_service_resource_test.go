package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccSupervisorServiceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSupervisorServiceResourceConfig("example", "go run golang.org/x/tools/cmd/godoc@latest"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_supervisor_service.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("example"),
					),
				},
			},
			// ImportState testing
			// {
			// 	ResourceName:      "uberspace_supervisor_service.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
			// Update and Read testing
			{
				Config: testAccSupervisorServiceResourceConfig("example", "go run golang.org/x/tools/cmd/godoc@latest"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_supervisor_service.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("example"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSupervisorServiceResourceConfig(name, command string) string {
	return fmt.Sprintf(`
resource "uberspace_supervisor_service" "test" {
  name = %[1]q
  command = %[2]q
}
`, name, command)
}
