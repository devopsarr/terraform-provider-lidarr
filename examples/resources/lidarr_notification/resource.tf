resource "lidarr_notification" "example" {
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

  implementation  = "CustomScript"
  config_contract = "CustomScriptSettings"

  path = "/scripts/lidarr.sh"
}