resource "lidarr_notification_plex" "example" {
  on_upgrade        = true
  on_rename         = false
  on_track_retag    = false
  on_release_import = true

  name = "Example"

  host       = "plex.lcl"
  port       = 32400
  auth_token = "AuthTOKEN"
}