data "uberspace_user" "user" {}

// create a crontab entry to run a Python script every 5 minutes
resource "uberspace_crontab_entry" "example" {
  entry = "*/5 * * * * python /home/${data.uberspace_user.user.name}/bin/example.py"
}
