resource "uberspace_maildomain" "mail" {
  asteroid = "isabell"
  name     = "mail.isabell.uber.space"
}

resource "uberspace_mailuser" "isabell" {
  // a mail user usually depends on a mail domain
  depends_on = [uberspace_maildomain.mail]

  name            = "isabell"
  password_hash   = "xxx"
  asteroid_name   = "isabell"
  maildomain_name = uberspace_maildomain.mail.name
}