package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccMyCnfDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccMyCnfDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.uberspace_mycnf.test",
						tfjsonpath.New("client"),
						knownvalue.MapExact(
							map[string]knownvalue.Check{
								"user":     knownvalue.StringExact("terra"),
								"password": knownvalue.NotNull(),
							},
						),
					),
					statecheck.ExpectKnownValue(
						"data.uberspace_mycnf.test",
						tfjsonpath.New("clientreadonly"),
						knownvalue.MapExact(
							map[string]knownvalue.Check{
								"user":     knownvalue.StringExact("terra_ro"),
								"password": knownvalue.NotNull(),
							},
						),
					),
				},
			},
		},
	})
}

const testAccMyCnfDataSourceConfig = `
data "uberspace_mycnf" "test" {}
`
