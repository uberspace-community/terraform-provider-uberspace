resource "crontab_entry" "example" {
  entry = "0 0 * * * /usr/bin/php /var/www/virtual/username/htdocs/cron.php"
}
