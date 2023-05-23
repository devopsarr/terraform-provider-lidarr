resource "lidarr_notification_apprise" "example" {
  on_grab               = false
  on_import_failure     = true
  on_upgrade            = true
  on_download_failure   = false
  on_release_import     = true
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  notification_type = 1
  server_url        = "https://apprise.go"
  auth_username     = "User"
  auth_password     = "Password"
  field_tags        = ["warning", "skull"]
}