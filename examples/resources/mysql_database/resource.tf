// create a new mysql database, the name will be "{username}_{suffix}", e.g. "isabell_test"
resource "uberspace_mysql_database" "test" {
  suffix = "test"
}