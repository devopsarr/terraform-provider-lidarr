resource "lidarr_import_list_headphones" "example" {
  enable_automatic_add = false
  should_monitor       = "specificAlbum"
  should_search        = false
  root_folder_path     = "/config"
  monitor_new_items    = "all"
  quality_profile_id   = 1
  metadata_profile_id  = 1
  name                 = "Example"
  base_url             = "http://127.0.0.1:8181"
  api_key              = "APIKey"
  tags                 = [1, 2, 3]
}
