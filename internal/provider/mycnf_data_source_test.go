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
						tfjsonpath.New("client.user"),
						knownvalue.StringExact("fx"),
					),
					statecheck.ExpectKnownValue(
						"data.uberspace_mycnf.test",
						tfjsonpath.New("client.password"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.uberspace_mycnf.test",
						tfjsonpath.New("clientreadonly.user"),
						knownvalue.StringExact("fx_ro"),
					),
					statecheck.ExpectKnownValue(
						"data.uberspace_mycnf.test",
						tfjsonpath.New("clientreadonly.password"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testAccMyCnfDataSourceConfig = `
data "uberspace_mycnf" "test" {}
`
