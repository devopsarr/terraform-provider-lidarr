---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_notification_prowl Resource - terraform-provider-lidarr"
subcategory: "Notifications"
description: |-
  <!-- subcategory:Notifications -->
  
  Notification Prowl resource.
  For more information refer to Notification https://wiki.servarr.com/lidarr/settings#connect and Prowl https://wiki.servarr.com/lidarr/supported#prowl.
---

# lidarr_notification_prowl (Resource)

<!-- subcategory:Notifications -->
Notification Prowl resource.
For more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Prowl](https://wiki.servarr.com/lidarr/supported#prowl).

## Example Usage

```terraform
resource "lidarr_notification_prowl" "example" {
  on_grab               = false
  on_upgrade            = true
  on_release_import     = true
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  api_key  = "APIKey"
  priority = -2
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_key` (String, Sensitive) API key.
- `name` (String) NotificationProwl name.

### Optional

- `include_health_warnings` (Boolean) Include health warnings.
- `on_album_delete` (Boolean) On album delete flag.
- `on_application_update` (Boolean) On application update flag.
- `on_artist_delete` (Boolean) On artist delete flag.
- `on_grab` (Boolean) On grab flag.
- `on_health_issue` (Boolean) On health issue flag.
- `on_health_restored` (Boolean) On health restored flag.
- `on_release_import` (Boolean) On release import flag.
- `on_upgrade` (Boolean) On upgrade flag.
- `priority` (Number) Priority.`-2` Very Low, `-1` Low, `0` Normal, `1` High, `2` Emergency.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_notification_prowl.example 1
```
