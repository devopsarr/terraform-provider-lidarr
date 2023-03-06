resource "lidarr_metadata_profile" "example" {
  name                  = "Example"
  primary_album_types   = [1, 2]
  secondary_album_types = [1]
  release_statuses      = [3]
}