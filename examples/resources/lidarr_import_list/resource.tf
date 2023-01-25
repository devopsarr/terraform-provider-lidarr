resource "lidarr_import_list" "example" {
  enable_automatic_add = false
  should_monitor       = "entireArtist"
  should_search        = false
  list_type            = "program"
  monitor_new_items    = "all"
  root_folder_path     = lidarr_root_folder.example.path
  quality_profile_id   = lidarr_quality_profile.example.id
  metadata_profile_id  = lidarr_metadata_profile.example.id
  name                 = "Example"
  implementation       = "LidarrImport"
  config_contract      = "LidarrSettings"
  tags                 = [1, 2]

  tag_ids     = [1, 2]
  profile_ids = [1]
  base_url    = "http://127.0.0.1:8686"
  api_key     = "APIKey"
}