resource "lidarr_artist" "example" {
  monitored           = true
  artist_name         = "Queen"
  path                = "/music/Queen"
  quality_profile_id  = 1
  metadata_profile_id = 1
  foreign_artist_id   = "0383dadf-2a4e-4d10-a46a-e9e041da8eb3"
}
