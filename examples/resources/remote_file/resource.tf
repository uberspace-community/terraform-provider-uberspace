// create a remote file resource with the upload of a Python script
resource "uberspace_remote_file" "examplepy" {
  src        = "example.py"
  dst        = "/home/isabell/bin/example.py"
  src_hash   = filesha256("example.py")
  executable = true
}

// crate a remote file from a download from the internet, check its hash and make it executable
data "http" "minio" {
  url = "https://dl.min.io/server/minio/release/linux-amd64/minio.sha256sum"
}

resource "uberspace_remote_file" "minio" {
  src        = "https://dl.min.io/server/minio/release/linux-amd64/minio"
  dst        = "/home/isabell/bin/minio"
  src_hash   = data.http.minio.response_body
  executable = true
}


// create a remote file resource from a string
resource "uberspace_remote_file" "exampletxt" {
  content = "example content"
  dst     = "/home/isabell/example.txt"
}
