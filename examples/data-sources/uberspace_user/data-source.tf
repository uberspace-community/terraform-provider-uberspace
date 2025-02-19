// access the uberspace user
data "uberspace_user" "user" {}

output "user" {
  value = data.uberspace_user.user.name
}