resource "lidarr_notification_signal" "example" {
  on_grab               = false
  on_import_failure     = true
  on_upgrade            = true
  on_download_failure   = false
  on_release_import     = true
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  auth_username = "User"
  auth_password = "Token"

  host          = "localhost"
  port          = 8080
  use_ssl       = true
  sender_number = "1234"
  receiver_id   = "4321"
}