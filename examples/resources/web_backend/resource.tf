resource "uberspace_web_domain" "minio" {
  domain = "minio.isabell.uber.space"
}

resource "uberspace_web_backend" "minio" {
  // a web backend usually depends on a web domain
  depends_on = [uberspace_web_domain.minio]

  uri  = "minio.isabell.uber.space/"
  port = 9001
}