resource "lidarr_notification_twitter" "example" {
  on_grab               = false
  on_import_failure     = true
  on_upgrade            = true
  on_download_failure   = false
  on_release_import     = true
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  access_token        = "Token"
  access_token_secret = "TokenSecret"
  consumer_key        = "Key"
  consumer_secret     = "Secret"
  mention             = "someone"
}