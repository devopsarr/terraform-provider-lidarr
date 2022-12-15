resource "lidarr_notification_discord" "example" {
  on_grab               = false
  on_import_failure     = true
  on_upgrade            = true
  on_rename             = false
  on_download_failure   = false
  on_track_retag        = false
  on_release_import     = true
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  web_hook_url  = "http://discord-web-hook.com"
  username      = "User"
  avatar        = "https://i.imgur.com/oBPXx0D.png"
  grab_fields   = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
  import_fields = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
}