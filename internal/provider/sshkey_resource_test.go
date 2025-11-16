package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const (
	testAccSshkeyValueOne = "aaaab3nza1lzdgixnte5aaaaibqykvvdmu9pq/9jv3uqlcq8b/pt6hdcctmn2v4rh"
	testAccSshkeyValueTwo = "aaaab3nza1lzdgixnte5aaaaibxw4t9ytv4p8k2dqtv4rdc3p13quk6r6v64b6gsh"
)

func TestAccSshkeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSshkeyResourceConfig("tf", "ssh-ed25519", testAccSshkeyValueOne, "terraform@example.com"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_sshkey.test",
						tfjsonpath.New("asteroid"),
						knownvalue.StringExact("tf"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_sshkey.test",
						tfjsonpath.New("key_type"),
						knownvalue.StringExact("ssh-ed25519"),
					),
					statecheck.ExpectKnownValue(
						"uberspace_sshkey.test",
						tfjsonpath.New("key"),
						knownvalue.StringExact(testAccSshkeyValueOne),
					),
					statecheck.ExpectKnownValue(
						"uberspace_sshkey.test",
						tfjsonpath.New("key_comment"),
						knownvalue.StringExact("terraform@example.com"),
					),
				},
			},
			{
				Config: testAccSshkeyResourceConfig("tf", "ssh-ed25519", testAccSshkeyValueTwo, "terraform+updated@example.com"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"uberspace_sshkey.test",
						tfjsonpath.New("key"),
						knownvalue.StringExact(testAccSshkeyValueTwo),
					),
					statecheck.ExpectKnownValue(
						"uberspace_sshkey.test",
						tfjsonpath.New("key_comment"),
						knownvalue.StringExact("terraform+updated@example.com"),
					),
				},
			},
		},
	})
}

func testAccSshkeyResourceConfig(asteroid, keyType, key, keyComment string) string {
	comment := ""
	if keyComment != "" {
		comment = fmt.Sprintf("  key_comment = %q\n", keyComment)
	}

	return fmt.Sprintf(`
resource "uberspace_sshkey" "test" {
  asteroid = %q
  key      = %q
  key_type = %q
%s}
`, asteroid, key, keyType, comment)
}
