resource "uberspace_webdomain" "minio" {
  asteroid = "isabell"
  name     = "minio.isabell.uber.space"
}

resource "uberspace_webdomain_header" "cors" {
  // a web backend usually depends on a web domain
  depends_on = [uberspace_webdomain.minio]

  asteroid = "isabell"
  domain   = "minio.isabell.uber.space"
  path     = "/"
  name     = "X-Custom-Header"
  value    = "custom"
}
