---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_notification_apprise Resource - terraform-provider-lidarr"
subcategory: "Notifications"
description: |-
  Notification Apprise resource.
  For more information refer to Notification https://wiki.servarr.com/lidarr/settings#connect and Apprise https://wiki.servarr.com/lidarr/supported#apprise.
---

# lidarr_notification_apprise (Resource)

<!-- subcategory:Notifications -->Notification Apprise resource.
For more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Apprise](https://wiki.servarr.com/lidarr/supported#apprise).

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) NotificationApprise name.
- `server_url` (String) Server URL.

### Optional

- `auth_password` (String, Sensitive) Password.
- `auth_username` (String) Username.
- `configuration_key` (String, Sensitive) Configuration key.
- `field_tags` (Set of String) Tags and emojis.
- `include_health_warnings` (Boolean) Include health warnings.
- `notification_type` (Number) Notification type. `0` Info, `1` Success, `2` Warning, `3` Failure.
- `on_album_delete` (Boolean) On album delete flag.
- `on_application_update` (Boolean) On application update flag.
- `on_artist_delete` (Boolean) On artist delete flag.
- `on_download_failure` (Boolean) On download failure flag.
- `on_grab` (Boolean) On grab flag.
- `on_health_issue` (Boolean) On health issue flag.
- `on_health_restored` (Boolean) On health restored flag.
- `on_import_failure` (Boolean) On download flag.
- `on_release_import` (Boolean) On release import flag.
- `on_upgrade` (Boolean) On upgrade flag.
- `stateless_urls` (String) Stateless URLs.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_notification_apprise.example 1
```