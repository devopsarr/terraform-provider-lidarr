resource "sonarr_metadata_roksbox" "example" {
  enable         = true
  name           = "Example"
  track_metadata = true
  artist_images  = false
  album_images   = true
}