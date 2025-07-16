resource "uberspace_sshkey" "example" {
  asteroid = "isabell"
  key      = filebase64("~/.ssh/id_ed25519.pub")
  key_type = "ssh-ed25519"
}
