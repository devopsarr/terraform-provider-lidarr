---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_notification_ntfy Resource - terraform-provider-lidarr"
subcategory: "Notifications"
description: |-
  Notification Ntfy resource.
  For more information refer to Notification https://wiki.servarr.com/lidarr/settings#connect and Ntfy https://wiki.servarr.com/lidarr/supported#ntfy.
---

# lidarr_notification_ntfy (Resource)

<!-- subcategory:Notifications -->Notification Ntfy resource.
For more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Ntfy](https://wiki.servarr.com/lidarr/supported#ntfy).

## Example Usage

```terraform
resource "lidarr_notification_ntfy" "example" {
  on_grab               = false
  on_import_failure     = true
  on_upgrade            = true
  on_download_failure   = false
  on_release_import     = true
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  priority   = 1
  server_url = "https://ntfy.sh"
  username   = "User"
  password   = "%s"
  topics     = ["Topic1234", "Topic4321"]
  field_tags = ["warning", "skull"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) NotificationNtfy name.
- `topics` (Set of String) Topics.

### Optional

- `click_url` (String) Click URL.
- `field_tags` (Set of String) Tags and emojis.
- `include_health_warnings` (Boolean) Include health warnings.
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
- `password` (String, Sensitive) Password.
- `priority` (Number) Priority. `1` Min, `2` Low, `3` Default, `4` High, `5` Max.
- `server_url` (String) Server URL.
- `tags` (Set of Number) List of associated tags.
- `username` (String) Username.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_notification_ntfy.example 1
```