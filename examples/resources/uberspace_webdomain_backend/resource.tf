resource "uberspace_webdomain" "minio" {
  asteroid = "isabell"
  domain   = "minio.isabell.uber.space"
}

resource "uberspace_webdomain_backend" "minio" {
  // a web backend usually depends on a web domain
  depends_on = [uberspace_webdomain.minio]

  asteroid = "isabell"
  domain   = "minio.isabell.uber.space"
  path     = "/"
  port     = 9001
}