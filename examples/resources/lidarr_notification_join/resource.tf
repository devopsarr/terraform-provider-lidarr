resource "lidarr_notification_join" "example" {
  on_grab               = false
  on_upgrade            = true
  on_release_import     = true
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  device_names = "device1,device2"
  api_key      = "Key"
  priority     = 2
}