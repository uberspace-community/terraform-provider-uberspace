// access the .my.cnf file of the user
data "uberspace_mycnf" "mycnf" {}

output "user" {
  value = data.uberspace_mycnf.mycnf.client.user
}

output "password" {
  value     = data.uberspace_mycnf.mycnf.client.password
  sensitive = true
}

output "ro_user" {
  value = data.uberspace_mycnf.mycnf.clientreadonly.user
}

output "ro_password" {
  value     = data.uberspace_mycnf.mycnf.clientreadonly.password
  sensitive = true
}