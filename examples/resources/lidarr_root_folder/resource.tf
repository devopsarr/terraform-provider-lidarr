resource "lidarr_root_folder" "example" {
  name                    = "Example"
  quality_profile_id      = 1
  metadata_profile_id     = 1
  monitor_option          = "future"
  new_item_monitor_option = "all"
  path                    = "/music"
  tags                    = [1]
}