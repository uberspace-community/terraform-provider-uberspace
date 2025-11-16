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

func TestAccMailuserResource(t *testing.T) {
	t.Parallel()

	asteroid := "tf"
	maildomain := fmt.Sprintf("%s.tf.uber8.space", acctest.RandomWithPrefix("mail"))
	username := acctest.RandomWithPrefix("acctest")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMailuserResourceConfig(asteroid, maildomain, username),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_mailuser.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(username),
					),
					statecheck.ExpectKnownValue(
						"uberspace_mailuser.test",
						tfjsonpath.New("mailaddr"),
						knownvalue.StringExact(fmt.Sprintf("%s@%s", username, maildomain)),
					),
					statecheck.ExpectKnownValue(
						"uberspace_mailuser.test",
						tfjsonpath.New("asteroid"),
						knownvalue.StringExact(asteroid),
					),
				},
			},
			{
				Config: testAccMailuserResourceConfig(asteroid, maildomain, username),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_mailuser.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(username),
					),
					statecheck.ExpectKnownValue(
						"uberspace_mailuser.test",
						tfjsonpath.New("mailaddr"),
						knownvalue.StringExact(fmt.Sprintf("%s@%s", username, maildomain)),
					),
					statecheck.ExpectKnownValue(
						"uberspace_mailuser.test",
						tfjsonpath.New("asteroid"),
						knownvalue.StringExact(asteroid),
					),
				},
			},
		},
	})
}

func testAccMailuserResourceConfig(asteroid, maildomain, username string) string {
	return fmt.Sprintf(`
resource "uberspace_maildomain" "test" {
  asteroid = %[1]q
  name     = %[2]q
}

resource "uberspace_mailuser" "test" {
  depends_on = [uberspace_maildomain.test]

  asteroid_name   = %[1]q
  maildomain_name = uberspace_maildomain.test.name
  name            = %[3]q
  password_hash   = "xxx"
}
`, asteroid, maildomain, username)
}
