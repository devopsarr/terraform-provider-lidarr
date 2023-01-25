resource "lidarr_import_list_lidarr_list" "example" {
  enable_automatic_add = false
  should_monitor       = "specificAlbum"
  should_search        = false
  root_folder_path     = "/config"
  monitor_new_items    = "all"
  quality_profile_id   = 1
  metadata_profile_id  = 1
  name                 = "Example"
  list_id              = "itunes/album/new"
  tags                 = [1, 2, 3]
}
