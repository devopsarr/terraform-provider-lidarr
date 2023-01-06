---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_notification_subsonic Resource - terraform-provider-lidarr"
subcategory: "Notifications"
description: |-
  Notification Subsonic resource.
  For more information refer to Notification https://wiki.servarr.com/lidarr/settings#connect and Subsonic https://wiki.servarr.com/lidarr/supported#xbmc.
---

# lidarr_notification_subsonic (Resource)

<!-- subcategory:Notifications -->Notification Subsonic resource.
For more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Subsonic](https://wiki.servarr.com/lidarr/supported#xbmc).

## Example Usage

```terraform
resource "lidarr_notification_subsonic" "example" {
  on_grab           = false
  on_upgrade        = false
  on_rename         = false
  on_track_retag    = false
  on_release_import = true
  on_health_issue   = false

  include_health_warnings = false
  name                    = "Example"

  host     = "http://subsonic.com"
  port     = 8080
  username = "User"
  password = "MyPass"
  notify   = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `host` (String) Host.
- `include_health_warnings` (Boolean) Include health warnings.
- `name` (String) NotificationSubsonic name.
- `on_grab` (Boolean) On grab flag.
- `on_health_issue` (Boolean) On health issue flag.
- `on_release_import` (Boolean) On movie file delete for upgrade flag.
- `on_rename` (Boolean) On rename flag.
- `on_track_retag` (Boolean) On movie file delete flag.
- `on_upgrade` (Boolean) On upgrade flag.
- `port` (Number) Port.

### Optional

- `notify` (Boolean) Notification flag.
- `password` (String, Sensitive) Password.
- `tags` (Set of Number) List of associated tags.
- `update_library` (Boolean) Update library flag.
- `url_base` (String) URL base.
- `use_ssl` (Boolean) Use SSL flag.
- `username` (String) Username.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_notification_subsonic.example 1
```