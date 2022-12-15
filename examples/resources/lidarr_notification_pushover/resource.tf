resource "lidarr_notification_join" "example" {
  on_grab               = false
  on_import_failure     = false
  on_upgrade            = false
  on_download_failure   = false
  on_release_import     = true
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  api_key  = "Key"
  priority = 2
}