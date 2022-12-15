resource "lidarr_notification_synology_indexer" "example" {
  on_upgrade        = true
  on_rename         = false
  on_track_retag    = false
  on_release_import = true

  name = "Example"

  update_library = true
}