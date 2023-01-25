resource "lidarr_import_list_spotify_albums" "example" {
  enable_automatic_add = false
  should_monitor       = "specificAlbum"
  should_search        = false
  root_folder_path     = "/config"
  monitor_new_items    = "all"
  quality_profile_id   = 1
  metadata_profile_id  = 1
  name                 = "Example"
  access_token         = "accessToken"
  refresh_token        = "refreshToken"
  expires              = "0001-01-01T00:01:00Z"
}
