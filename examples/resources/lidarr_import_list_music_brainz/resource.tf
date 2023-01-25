resource "lidarr_import_list_music_brainz" "example" {
  enable_automatic_add = false
  should_monitor       = "specificAlbum"
  should_search        = false
  root_folder_path     = "/config"
  monitor_new_items    = "all"
  quality_profile_id   = 1
  metadata_profile_id  = 1
  name                 = "Example"
  series_id            = "0383dadf-2a4e-4d10-a46a-e9e041da8eb3"
  tags                 = [1, 2, 3]
}
