resource "lidarr_metadata_kodi" "example" {
  enable          = true
  name            = "Example"
  artist_metadata = true
  album_images    = true
  artist_images   = true
  album_metadata  = false
}