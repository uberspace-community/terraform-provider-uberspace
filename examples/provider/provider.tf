// manage crontab on the local host
provider "crontab" {}

// manage crontab on a remote host via SSH with a password
provider "crontab" {
  host     = "example.com"
  user     = "root"
  password = "password"
}

// manage crontab on a remote host via SSH with a private key
provider "crontab" {
  host        = "1.2.3.4"
  user        = "root"
  private_key = file("~/.ssh/id_ed25519")
}