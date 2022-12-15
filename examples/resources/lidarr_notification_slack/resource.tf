resource "lidarr_notification_slack" "example" {
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

  web_hook_url = "http://my.slack.com/test"
  username     = "user"
  channel      = "example-channel"
}