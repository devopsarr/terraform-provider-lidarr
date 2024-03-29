---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_notification_plex Resource - terraform-provider-lidarr"
subcategory: "Notifications"
description: |-
  <!-- subcategory:Notifications -->
  
  Notification Plex resource.
  For more information refer to Notification https://wiki.servarr.com/lidarr/settings#connect and Plex https://wiki.servarr.com/lidarr/supported#plexserver.
---

# lidarr_notification_plex (Resource)

<!-- subcategory:Notifications -->
Notification Plex resource.
For more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Plex](https://wiki.servarr.com/lidarr/supported#plexserver).

## Example Usage

```terraform
resource "lidarr_notification_plex" "example" {
  on_upgrade        = true
  on_rename         = false
  on_track_retag    = false
  on_release_import = true

  name = "Example"

  host       = "plex.lcl"
  port       = 32400
  auth_token = "AuthTOKEN"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `auth_token` (String, Sensitive) Auth Token.
- `host` (String) Host.
- `name` (String) NotificationPlex name.

### Optional

- `on_album_delete` (Boolean) On album delete flag.
- `on_artist_delete` (Boolean) On artist delete flag.
- `on_release_import` (Boolean) On release import flag.
- `on_rename` (Boolean) On rename flag.
- `on_track_retag` (Boolean) On track retag flag.
- `on_upgrade` (Boolean) On upgrade flag.
- `port` (Number) Port.
- `tags` (Set of Number) List of associated tags.
- `update_library` (Boolean) Update library flag.
- `use_ssl` (Boolean) Use SSL flag.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_notification_plex.example 1
```
