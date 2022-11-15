resource "lidarr_notification_custom_script" "example" {
  on_grab               = false
  on_upgrade            = true
  on_rename             = false
  on_release_import     = false
  on_download_failure   = false
  on_import_failure     = true
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  path = "/scripts/lidarr.sh"
}