resource "uberspace_supervisor_service" "app" {
  name    = "minio"
  command = "/home/isabell/bin/minio server /home/isabell/minio --address 0.0.0.0:9000 --console-address 0.0.0.0:9001"
  environment = {
    "MINIO_ACCESS_KEY"           = "minio"
    "MINIO_SECRET_KEY"           = "minio123"
    "MINIO_BROWSER_REDIRECT_URL" = "https://console.example.com/"
  }

  lifecycle {
    replace_triggered_by = [
      // Add a trigger here to restart the service when the binary changes
      // uberspace_remote_file.minio
    ]
  }
}